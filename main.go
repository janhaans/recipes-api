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
	// Here you would typically save the newRecipe to a database
	c.JSON(201, newRecipe)
}

func GetRecipesHandler(c *gin.Context) {
	c.JSON(200, recipes)
}

func GetRecipesByTagHandler(c *gin.Context) {
	tag := c.Query("tag")
	if tag == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tag query parameter is required"})
		return
	}

	filteredRecipes := make([]Recipe, 0)
	for _, recipe := range recipes {
		for _, t := range recipe.Tags {
			if t == tag {
				filteredRecipes = append(filteredRecipes, recipe)
				break
			}
		}
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

	for i, recipe := range recipes {
		if recipe.ID == id {
			updatedRecipe.ID = recipe.ID
			updatedRecipe.PublishedAt = recipe.PublishedAt
			recipes[i] = updatedRecipe
			c.JSON(200, updatedRecipe)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
}

func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	for i, recipe := range recipes {
		if recipe.ID == id {
			recipes = append(recipes[:i], recipes[i+1:]...)
			c.JSON(200, gin.H{"message": "Recipe deleted"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
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
