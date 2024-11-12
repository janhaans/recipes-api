package handlers

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/janhaans/recipe-api/models"
)

type UserHandler struct {
	ctx context.Context
}

func NewUserHandler(ctx context.Context) *UserHandler {
	return &UserHandler{
		ctx: ctx,
	}
}

func (u *UserHandler) SigninUserHandler(c *gin.Context) {
	//Extract User from Request
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Authenticate User
	if newUser.Username != "admin" || newUser.Password != "password" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}

	// JWT Signing Key
	SigningKey := []byte(os.Getenv("JWT_TOKEN"))

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": newUser.Username,
		"exp":      time.Now().Add(time.Minute * 10).Unix(),
	})

	tokenString, err := token.SignedString(SigningKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Return JWT
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}