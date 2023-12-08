package recipe

import (
	"encoding/json"
	"os"
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

var Recipes []Recipe

func init() {
	bs, err := os.ReadFile("recipes.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(bs, &Recipes); err != nil {
		panic(err)
	}
}
