package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/janhaans/recipe-api/models"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecipeHandler struct {
	ctx context.Context
	collection *mongo.Collection

}

func NewRecipeHandler(ctx context.Context, collection *mongo.Collection) *RecipeHandler {
	return &RecipeHandler{
		ctx: ctx,
		collection: collection,
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
	_, err := h.collection.InsertOne(h.ctx, newRecipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save recipe to database"})
		return
	}

	c.IndentedJSON(201, newRecipe)
}

func (h *RecipeHandler) GetRecipesHandler(c *gin.Context) {
	cursor, err := h.collection.Find(context.TODO(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipes from database"})
		return
	}
	defer cursor.Close(h.ctx)

	var allRecipes []models.Recipe
	if err := cursor.All(h.ctx, &allRecipes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode recipes"})
		return
	}

	c.IndentedJSON(200, allRecipes)
}

func (h *RecipeHandler) GetRecipesByTagHandler(c *gin.Context) {
	tag := c.Query("tag")
	if tag == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tag query parameter is required"})
		return
	}

	filter := bson.M{"tags": tag}

	cursor, err := h.collection.Find(h.ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipes from database"})
		return
	}
	defer cursor.Close(h.ctx)

	var filteredRecipes []models.Recipe
	if err := cursor.All(h.ctx, &filteredRecipes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode recipes"})
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

	filter := bson.M{"id": id}
	update := bson.M{
		"$set": bson.M{
			"name":         updatedRecipe.Name,
			"tags":         updatedRecipe.Tags,
			"ingredients":  updatedRecipe.Ingredients,
			"instructions": updatedRecipe.Instructions,
		},
	}

	result, err := h.collection.UpdateOne(h.ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update recipe in database"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	var resultRecipe models.Recipe
	err = h.collection.FindOne(h.ctx, filter).Decode(&resultRecipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated recipe from database"})
		return
	}

	c.IndentedJSON(200, resultRecipe)
}

func (h * RecipeHandler) DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	filter := bson.M{"id": id}

	result, err := h.collection.DeleteOne(h.ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete recipe from database"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	c.IndentedJSON(200, gin.H{"message": "Recipe deleted"})
}