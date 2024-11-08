package routers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/janhaans/recipe-api/db"
	"github.com/janhaans/recipe-api/handlers"
	"github.com/stretchr/testify/assert"
)

// Declare a global variable for the shared router
var router *gin.Engine

func TestMain(m *testing.M) {
	// Set up the Gin router and the RecipeHandler once
	ctx := context.TODO()
	repo := db.NewMockDBRecipesRepository()
	recipeHandler := handlers.NewRecipeHandler(ctx, repo)
	router = gin.Default()
    AddRecipeRoutes(router, recipeHandler)

    // Run tests
    code := m.Run()

    // Exit with the result code
    os.Exit(code)
}

//Test: Create a new recipe
func TestCreateRecipe(t *testing.T) {
	// Create a new recipe JSON payload
	newRecipe := `{
		"name": "Test Recipe",
		"ingredients": ["Ingredient 1", "Ingredient 2"],
		"steps": ["Step 1", "Step 2"]
	}`

	// Create a new HTTP request
	req, _ := http.NewRequest("POST", "/recipes", bytes.NewBuffer([]byte(newRecipe)))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Test Recipe")
}

//Test: Get all recipes
func TestGetRecipes(t *testing.T) {
	// Create a new HTTP request
	req, _ := http.NewRequest("GET", "/recipes", nil)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Spaghetti Carbonara")
	assert.Contains(t, w.Body.String(), "Beef Stroganoff")
}

//Test: Get recipes by tag
func TestGetRecipeByTag(t *testing.T) {
	// Create a new HTTP request
	req, _ := http.NewRequest("GET", "/recipes?tag=pasta", nil)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Spaghetti Carbonara")
}

//Test: Update a recipe
func TestUpdateRecipe(t *testing.T) {
	// Create a new recipe JSON payload
	updatedRecipe := `{
		"name": "Spaghetti Carbonara",
		"tags": ["pasta", "italian"],
		"ingredients": ["spaghetti", "eggs", "pancetta", "pecorino cheese"],
		"steps": ["Cook spaghetti", "Fry pancetta", "Mix eggs and cheese", "Combine all ingredients"]
	}`

	// Create a new HTTP request
	req, _ := http.NewRequest("PUT", "/recipes/1", bytes.NewBuffer([]byte(updatedRecipe)))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Spaghetti Carbonara")
	assert.Contains(t, w.Body.String(), `"id": "1"`)
}

//Test: Delete a recipe
func TestDeleteRecipe(t *testing.T) {
	// Create a new HTTP request
	req, _ := http.NewRequest("DELETE", "/recipes/1", nil)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusNoContent, w.Code)
}