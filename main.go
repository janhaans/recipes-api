package main

import (
	"github.com/gin-gonic/gin"
	"github.com/janhaans/recipe-api/controllers"
)

func main() {
	router := gin.Default()
	router.POST("/recipe", controllers.NewRecipeHandler)
	router.Run()
}
