package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/janhaans/recipe-api/models"
	"github.com/rs/xid"
)

type RecipeHandler struct {
	ctx context.Context
	repo models.RecipeRepository
}

func NewRecipeHandler(ctx context.Context, repo models.RecipeRepository) *RecipeHandler {
	return &RecipeHandler{
		ctx:    ctx,
		repo:  	repo,
	}
}

func (h *RecipeHandler) NewRecipeHandler(c *gin.Context) {
	var newRecipe models.Recipe
	if err := c.ShouldBindJSON(&newRecipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newRecipe.ID = xid.New().String()
	newRecipe.PublishedAt = time.Now()

	// Save the newRecipe to MongoDB
	if err := h.repo.CreateRecipe(h.ctx, newRecipe); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save recipe to database"})
		return
	}

	c.IndentedJSON(201, newRecipe)
}

func (h *RecipeHandler) GetRecipesHandler(c *gin.Context) {
	allRecipes, err := h.repo.GetRecipes(h.ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipes from database"})
		return
	}

	c.IndentedJSON(200, allRecipes)
}

func (h *RecipeHandler) GetRecipesByTag(c *gin.Context) {
	tag := c.Query("tag")
	if tag == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tag query parameter is required"})
		return
	}

	filteredRecipes, err := h.repo.GetRecipesByTag(h.ctx, tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipes from database"})
		return
	}

	c.IndentedJSON(200, filteredRecipes)
}

func (h RecipeHandler) UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var updatedRecipe models.Recipe
	if err := c.ShouldBindJSON(&updatedRecipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipe, err := h.repo.UpdateRecipe(h.ctx, id, updatedRecipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(200, recipe)
}

func (h * RecipeHandler) DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	err := h.repo.DeleteRecipe(h.ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(204, gin.H{"message": "Recipe deleted"})
}