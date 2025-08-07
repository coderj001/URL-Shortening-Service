package auth

import (
	"fmt"
	"net/http"

	"github.com/coderj001/URL-shortener/config"
	"github.com/coderj001/URL-shortener/database"
	"github.com/coderj001/URL-shortener/helpers"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	UnAuthorized = "unauthorized"
	Authorized   = "authorized"
)

func AuthMiddleware(db *database.MySQLStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.Set("auth_status", UnAuthorized)
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.GetConfig().JWTSecret), nil
		})

		if err != nil || !token.Valid {
			helpers.HandleError(c, http.StatusUnauthorized, fmt.Errorf("invalid token"))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			helpers.HandleError(c, http.StatusUnauthorized, fmt.Errorf("invalid token claims"))
			c.Abort()
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			helpers.HandleError(c, http.StatusUnauthorized, fmt.Errorf("username not found in token"))
			c.Abort()
			return
		}

		user, err := db.GetUser(username)
		if err != nil {
			helpers.HandleError(c, http.StatusUnauthorized, fmt.Errorf("user verification failed: %v", err))
			c.Abort()
			return
		}

		// Set authorized status and user
		c.Set("auth_status", Authorized)
		c.Set("user", user)

		c.Next()

	}
}
