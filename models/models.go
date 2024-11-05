package models

import (
	"context"
	"time"
)

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
 }

 type RecipeRepository interface {
	 GetRecipes(ctx context.Context) ([]Recipe, error)
	 GetRecipesByTagHandler(ctx context.Context, tag string) ([]Recipe, error)
	 CreateRecipe(ctx context.Context, recipe Recipe) error
	 UpdateRecipe(ctx context.Context, id string, recipe Recipe) (*Recipe, error)
	 DeleteRecipe(ctx context.Context, id string) error
 }