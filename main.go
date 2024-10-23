/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/janhaans/recipe-api/handlers"
	"github.com/janhaans/recipe-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

 var client *mongo.Client
 var collection *mongo.Collection
 var ctx = context.TODO()
 var handler *handlers.RecipeHandler


 func init() {
	var err error
	// Initialize MongoDB client
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Load recipes from file
	LoadRecipesFromFile()

	// Initialize the RecipeHandler
	collection = client.Database("recipes-db").Collection("recipes")
	handler = handlers.NewRecipeHandler(ctx, collection)
 }


// LoadRecipesFromFile reads recipes from an embedded JSON file and loads them into MongoDB
func LoadRecipesFromFile() {
	if client == nil {
        log.Fatal("MongoDB client is not initialized")
    }

	collectionNames, err := client.Database("recipes-db").ListCollectionNames(ctx, bson.D{})
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

	var loadedRecipes []models.Recipe
	if err := json.Unmarshal(file, &loadedRecipes); err != nil {
		log.Fatalf("Failed to unmarshal recipes: %v", err)
	}

	collection := client.Database("recipes-db").Collection("recipes")
	for _, recipe := range loadedRecipes {
		_, err = collection.InsertOne(ctx, recipe)
		if err != nil {
			log.Printf("Failed to load recipe %s: %v", recipe.Name, err)
		}
	}

	fmt.Println("Recipes loaded into MongoDB!")
}

func main() {
	router := gin.Default()
	router.POST("/recipes", handler.NewRecipeHandler)
	router.GET("/recipes", handler.GetRecipesHandler)
	router.GET("/recipes/search", handler.GetRecipesByTagHandler)
	router.PUT("/recipes/:id", handler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", handler.DeleteRecipeHandler)
	router.Run()
}
