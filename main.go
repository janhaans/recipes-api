/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
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

 func init() {
	recipes = make([]Recipe, 0)
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


func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", GetRecipesHandler)
	router.Run()
}
