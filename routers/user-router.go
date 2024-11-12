package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/janhaans/recipe-api/handlers"
)

func AddUserRoutes(router *gin.Engine, userHandler *handlers.UserHandler) {
	// Add the routes for the RecipeHandler
	router.POST("/users/signin", userHandler.SigninUserHandler)
}