package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/janhaans/recipe-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongoDB(ctx context.Context) *mongo.Client {
	var err error
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Check the connection
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB!")

	return mongoClient
}

//Load Recipes from JSON file into MongoDB
func LoadRecipes(ctx context.Context, mongoClient *mongo.Client) {
	if mongoClient == nil {
        log.Fatal("MongoDB client is not initialized")
    }

	collectionNames, err := mongoClient.Database("recipes-db").ListCollectionNames(ctx, bson.D{})
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

	recipesCollection := mongoClient.Database("recipes-db").Collection("recipes")
	for _, recipe := range loadedRecipes {
		_, err = recipesCollection.InsertOne(ctx, recipe)
		if err != nil {
			log.Printf("Failed to load recipe %s: %v", recipe.Name, err)
		}
	}

	fmt.Println("Recipes loaded into MongoDB!")
}