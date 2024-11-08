package db

import (
	"context"
	"fmt"
	"time"

	"github.com/janhaans/recipe-api/models"
)

type MockDBRecipesRepository struct {}

func NewMockDBRecipesRepository() *MockDBRecipesRepository {
	return &MockDBRecipesRepository{}
}

func (r *MockDBRecipesRepository) GetRecipes(ctx context.Context) ([]models.Recipe, error) {
	return []models.Recipe{
		{
			ID:           "1",
			Name:         "Spaghetti Carbonara",
			Tags:         []string{"pasta", "italian"},
			Ingredients:  []string{"spaghetti", "eggs", "pancetta", "pecorino cheese"},
			Instructions: []string{"Cook spaghetti", "Fry pancetta", "Mix eggs and cheese", "Combine all ingredients"},
			PublishedAt:  time.Now(),
		},
		{
			ID:           "2",
			Name:         "Beef Stroganoff",
			Tags:         []string{"beef", "russian"},
			Ingredients:  []string{"beef", "onion", "mushrooms", "sour cream"},
			Instructions: []string{"Fry beef", "Fry onion and mushrooms", "Add sour cream", "Combine all ingredients"},
			PublishedAt:  time.Now(),
		},
	}, nil
}

func (r *MockDBRecipesRepository) GetRecipe(ctx context.Context, id string) (*models.Recipe, error) {
	if id != "1" {
		return nil, fmt.Errorf("recipe not found")
	}
	return &models.Recipe{
		ID:           "1",
		Name:         "Spaghetti Carbonara",
		Tags:         []string{"pasta", "italian"},
		Ingredients:  []string{"spaghetti", "eggs", "pancetta", "pecorino cheese"},
		Instructions: []string{"Cook spaghetti", "Fry pancetta", "Mix eggs and cheese", "Combine all ingredients"},
		PublishedAt:  time.Now(),
	}, nil
}

func (r *MockDBRecipesRepository) GetRecipesByTag(ctx context.Context, tag string) ([]models.Recipe, error) {
	return []models.Recipe{
		{
			ID:           "1",
			Name:         "Spaghetti Carbonara",
			Tags:         []string{"pasta", "italian"},
			Ingredients:  []string{"spaghetti", "eggs", "pancetta", "pecorino cheese"},
			Instructions: []string{"Cook spaghetti", "Fry pancetta", "Mix eggs and cheese", "Combine all ingredients"},
			PublishedAt:  time.Now(),
		},
	}, nil
}

func (r *MockDBRecipesRepository) CreateRecipe(ctx context.Context, recipe models.Recipe) error {
	return nil
}

func (r *MockDBRecipesRepository) UpdateRecipe(ctx context.Context, id string, recipe models.Recipe) (*models.Recipe, error) {
	if id != "1" {
		return nil, fmt.Errorf("recipe not found")
	}
	return &models.Recipe{
		ID:           "1",
		Name:         "Spaghetti Carbonara",
		Tags:         []string{"pasta", "italian"},
		Ingredients:  []string{"spaghetti", "eggs", "pancetta", "pecorino cheese"},
		Instructions: []string{"Cook spaghetti", "Fry pancetta", "Mix eggs and cheese", "Combine all ingredients"},
		PublishedAt:  time.Now(),
	}, nil
}

func (r *MockDBRecipesRepository) DeleteRecipe(ctx context.Context, id string) error {
	if id != "1" {
		return fmt.Errorf("recipe not found")
	}
	return nil
}
