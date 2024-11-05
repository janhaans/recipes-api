package db

import (
	"time"

	"github.com/janhaans/recipe-api/models"
)

type MockDBRecipesRepository struct {}

func NewMockDBRecipesRepository() *MockDBRecipesRepository {
	return &MockDBRecipesRepository{}
}

func (r *MockDBRecipesRepository) GetRecipes() ([]models.Recipe, error) {
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

func (r *MockDBRecipesRepository) GetRecipesByTag(tag string) ([]models.Recipe, error) {
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

func (r *MockDBRecipesRepository) CreateRecipe(recipe models.Recipe) error {
	return nil
}

func (r *MockDBRecipesRepository) UpdateRecipe(id string, recipe models.Recipe) (*models.Recipe, error) {
	return &models.Recipe{
		ID:           "1",
		Name:         "Spaghetti Carbonara",
		Tags:         []string{"pasta", "italian"},
		Ingredients:  []string{"spaghetti", "eggs", "pancetta", "pecorino cheese"},
		Instructions: []string{"Cook spaghetti", "Fry pancetta", "Mix eggs and cheese", "Combine all ingredients"},
		PublishedAt:  time.Now(),
	}, nil
}

func (r *MockDBRecipesRepository) DeleteRecipe(id string) error {
	return nil
}
