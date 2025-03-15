package users

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/coderj001/URL-shortener/database"
	"github.com/coderj001/URL-shortener/helpers"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context, db *database.MySQLStore) {
	var userInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := helpers.ParseRequest(c, &userInput); err != nil {
		helpers.HandleError(
			c,
			http.StatusBadRequest,
			fmt.Errorf("invalid request"),
		)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(userInput.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		helpers.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	err = db.SaveUser(userInput.Username, string(hashedPassword))
	if err != nil {
		helpers.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	tokenString, err := helpers.GenerateToken(userInput.Username)
	if err != nil {
		helpers.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"message": "User created successfully",
			"data": map[string]string{
				"username": userInput.Username,
				"token":    tokenString,
			},
		},
	)
}

func LoginUser(c *gin.Context, db *database.MySQLStore) {
	var userInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := helpers.ParseRequest(c, &userInput); err != nil {
		helpers.HandleError(c, http.StatusBadRequest, fmt.Errorf("invalid request"))
		return
	}

	hashedPassword, err := db.GetHashPassward(userInput.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helpers.HandleError(c, http.StatusBadGateway, fmt.Errorf("user don't exists"))
			return
		}
		helpers.HandleError(c, http.StatusBadGateway, fmt.Errorf("database error"))
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userInput.Password)); err != nil {
		helpers.HandleError(c, http.StatusUnauthorized, fmt.Errorf("invalid credentials"))
		return
	}

	tokenString, err := helpers.GenerateToken(userInput.Username)

	if err != nil {
		helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("token generation failed, %w", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}
