package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/janhaans/recipe-api/handlers"
)

func AddRecipeRoutes(router *gin.Engine, recipeHandler *handlers.RecipeHandler) {
	// Add the routes for the RecipeHandler
	routerGroup := router.Group("/recipes")
	{
		routerGroup.POST("/", recipeHandler.NewRecipeHandler)
		routerGroup.GET("/", recipeHandler.GetRecipesHandler)
		routerGroup.GET("/search", recipeHandler.GetRecipesByTagHandler)
		routerGroup.PUT("/:id", recipeHandler.UpdateRecipeHandler)
		routerGroup.DELETE(":id", recipeHandler.DeleteRecipeHandler)
	}
}