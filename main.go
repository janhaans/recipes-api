package main

import (
	"github.com/gin-gonic/gin"
	"github.com/janhaans/recipe-api/controllers"
)

func main() {
	router := gin.Default()
	router.POST("/recipes", controllers.NewRecipeHandler)
	router.GET("/recipes", controllers.ListRecipesHandler)
	router.Run()
}
