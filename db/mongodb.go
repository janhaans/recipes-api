package db

import (
	"context"
	"fmt"
	"log"

	"github.com/janhaans/recipe-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBRecipesRepository struct {
	client *mongo.Client
	dbName string
	recipes *mongo.Collection
}

func NewMongoDBRecipesRepository(client *mongo.Client, dbName string) *MongoDBRecipesRepository {
	return &MongoDBRecipesRepository{
		client: client,
		dbName: dbName,
		recipes: client.Database(dbName).Collection("recipes"),
	}
}

func (r *MongoDBRecipesRepository) GetRecipes(ctx context.Context) ([]models.Recipe, error) {
	cursor, err := r.recipes.Find(ctx, bson.D{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var recipes []models.Recipe
	if err := cursor.All(ctx, &recipes); err != nil {
		log.Println(err)
		return nil, err
	}

	return recipes, nil
}

func (r *MongoDBRecipesRepository) GetRecipe(ctx context.Context, id string) (*models.Recipe, error) {
	filter := bson.M{"id": id}
	var recipe models.Recipe
	err := r.recipes.FindOne(ctx, filter).Decode(&recipe)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("recipe not found")
	} else if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to fetch recipe from database")
	}

	return &recipe, nil
}

func (r *MongoDBRecipesRepository) GetRecipesByTag(ctx context.Context, tag string) ([]models.Recipe, error) {
	filter := bson.M{"tags": tag}
	cursor, err := r.recipes.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var recipes []models.Recipe
	if err := cursor.All(ctx, &recipes); err != nil {
		log.Println(err)
		return nil, err
	}

	return recipes, nil
}

func (r *MongoDBRecipesRepository) CreateRecipe(ctx context.Context, recipe models.Recipe) error {
	_, err := r.recipes.InsertOne(ctx, recipe)
	if err!=nil {
		log.Println(err)
	}
	return err
}

func (r *MongoDBRecipesRepository) UpdateRecipe(ctx context.Context, id string, recipe models.Recipe) (*models.Recipe, error) {
	filter := bson.M{"id": id}
	update := bson.M{
		"$set": bson.M{
			"name":         recipe.Name,
			"tags":         recipe.Tags,
			"ingredients":  recipe.Ingredients,
			"instructions": recipe.Instructions,
		},
	}
	result, err := r.recipes.UpdateOne(ctx, filter, update)
	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("recipe not found")
	}
	if err!=nil {
		return nil, fmt.Errorf("failed to update recipe in database")
	}

	var resultRecipe models.Recipe
	err = r.recipes.FindOne(ctx, filter).Decode(&resultRecipe)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated recipe from database")
	}

	return &resultRecipe, nil
}

func (r *MongoDBRecipesRepository) DeleteRecipe(ctx context.Context, id string) error {
	filter := bson.M{"id": id}
	result, err := r.recipes.DeleteOne(ctx, filter)
	if err == mongo.ErrNoDocuments {
		return fmt.Errorf("recipe not found")
	} else if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to delete recipe from database")
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("recipe not found")
	}
	
	return err
}