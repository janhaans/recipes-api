/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
 }

 var recipes []Recipe
 var client *mongo.Client

 func init() {
	recipes = make([]Recipe, 0)
	var err error
	// Initialize MongoDB client
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Load recipes from file
	LoadRecipesFromFile()
 }


// LoadRecipesFromFile reads recipes from an embedded JSON file and loads them into MongoDB
func LoadRecipesFromFile() {
	if client == nil {
        log.Fatal("MongoDB client is not initialized")
    }

	collectionNames, err := client.Database("recipes-db").ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		log.Fatalf("Failed to list collection names: %v", err)
	}

	collectionExists := false
	for _, name := range collectionNames {
		if name == "recipes" {
			collectionExists = true
			break
		}
	}

	if collectionExists {
		fmt.Println("Collection 'recipes' already exists in the database.")
		return
	}

	file, err := os.ReadFile("recipes.json")
	if err != nil {
		log.Fatalf("Failed to read recipes file: %v", err)
	}

	var loadedRecipes []Recipe
	if err := json.Unmarshal(file, &loadedRecipes); err != nil {
		log.Fatalf("Failed to unmarshal recipes: %v", err)
	}

	collection := client.Database("recipes-db").Collection("recipes")
	for _, recipe := range loadedRecipes {
		_, err = collection.InsertOne(context.TODO(), recipe)
		if err != nil {
			log.Printf("Failed to load recipe %s: %v", recipe.Name, err)
		}
	}

	fmt.Println("Recipes loaded into MongoDB!")
}


func NewRecipeHandler(c *gin.Context) {
	var newRecipe Recipe
	if err := c.ShouldBindJSON(&newRecipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newRecipe.ID = xid.New().String()
	newRecipe.PublishedAt = time.Now()
	recipes = append(recipes, newRecipe)

	// Save the newRecipe to MongoDB
	collection := client.Database("recipes-db").Collection("recipes")
	_, err := collection.InsertOne(context.TODO(), newRecipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save recipe to database"})
		return
	}

	c.JSON(201, newRecipe)
}

func GetRecipesHandler(c *gin.Context) {
	collection := client.Database("recipes-db").Collection("recipes")

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipes from database"})
		return
	}
	defer cursor.Close(context.TODO())

	var allRecipes []Recipe
	if err := cursor.All(context.TODO(), &allRecipes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode recipes"})
		return
	}

	c.JSON(200, allRecipes)
}

func GetRecipesByTagHandler(c *gin.Context) {
	tag := c.Query("tag")
	if tag == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tag query parameter is required"})
		return
	}

	collection := client.Database("recipes-db").Collection("recipes")
	filter := bson.M{"tags": tag}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipes from database"})
		return
	}
	defer cursor.Close(context.TODO())

	var filteredRecipes []Recipe
	if err := cursor.All(context.TODO(), &filteredRecipes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode recipes"})
		return
	}

	c.JSON(200, filteredRecipes)
}

func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var updatedRecipe Recipe
	if err := c.ShouldBindJSON(&updatedRecipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := client.Database("recipes-db").Collection("recipes")
	filter := bson.M{"id": id}
	update := bson.M{
		"$set": bson.M{
			"name":         updatedRecipe.Name,
			"tags":         updatedRecipe.Tags,
			"ingredients":  updatedRecipe.Ingredients,
			"instructions": updatedRecipe.Instructions,
		},
	}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update recipe in database"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	var resultRecipe Recipe
	err = collection.FindOne(context.TODO(), filter).Decode(&resultRecipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated recipe from database"})
		return
	}

	c.JSON(200, resultRecipe)
}

func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	collection := client.Database("recipes-db").Collection("recipes")
	filter := bson.M{"id": id}

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete recipe from database"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Recipe deleted"})
}

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", GetRecipesHandler)
	router.GET("/recipes/search", GetRecipesByTagHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.Run()
}
