package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/janhaans/recipe-api/recipe"
	"github.com/rs/xid"
)

var recipes = recipe.Recipes

func NewRecipeHandler(c *gin.Context) {
	var recipe recipe.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)
}

func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}
