package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/janhaans/recipe-api/handlers"
)

func AddRecipeRoutes(router *gin.Engine, recipeHandler *handlers.RecipeHandler) {
	// Add the routes for the RecipeHandler
	routerGroup := router.Group("/recipes")
	{
		//Create a new recipe
		routerGroup.POST("", recipeHandler.NewRecipeHandler)
		//Get all recipes
		routerGroup.GET("", recipeHandler.GetRecipesHandler)
		//Get recipes by tag
		routerGroup.GET("/search", recipeHandler.GetRecipesByTag)
		//Update a recipe
		routerGroup.PUT("/:id", recipeHandler.UpdateRecipeHandler)
		//Delete a recipe
		routerGroup.DELETE(":id", recipeHandler.DeleteRecipeHandler)
	}
}