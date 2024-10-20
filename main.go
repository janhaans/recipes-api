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


func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", GetRecipesHandler)
	router.GET("/recipes/search", GetRecipesByTagHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.Run()
}
