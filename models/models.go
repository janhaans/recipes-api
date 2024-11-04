package models

import "time"

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
 }

 type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}