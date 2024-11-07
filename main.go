/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/janhaans/recipe-api/db"
	"github.com/janhaans/recipe-api/handlers"
	"github.com/janhaans/recipe-api/routers"
)

func main() {
	ctx := context.TODO()

	// Connect to MongoDB
	mongoClient := db.ConnectToMongoDB(ctx)

	//LoadRecipesFromFile()
	db.LoadRecipes(ctx, mongoClient)

	// Initialize the RecipeHandler
	repo := db.NewMongoDBRecipesRepository(mongoClient, "recipes-db")
	recipeHandler := handlers.NewRecipeHandler(ctx, repo)

	// Initialize the Gin router
	router := gin.Default()
	routers.AddRecipeRoutes(router, recipeHandler)
	router.Run()
}
