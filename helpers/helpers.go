package helpers

import (
	"crypto/rand"
	"errors"
	"strings"
	"time"

	"github.com/coderj001/URL-shortener/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// EnforceHTTP ...
func EnforceHTTP(url string) string {
	// make every url https
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}

// RemoveDomainError ...
func RemoveDomainError(url string) bool {
	// basically this functions removes all the commonly found
	// prefixes from URL such as http, https, www
	// then checks of the remaining string is the DOMAIN itself
	if url == config.GetConfig().GetDomain() {
		return false
	}
	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)
	newURL = strings.Split(newURL, "/")[0]

	if newURL == config.GetConfig().GetDomain() {
		return false
	}
	return true
}

func GenerateID(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	var id string
	for _, byte := range b {
		// Use modulo to map random bytes to our charset
		idx := int(byte) % len(charset)
		id += string(charset[idx])
	}

	return id, nil
}

func ParseRequest(c *gin.Context, body any) error {
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	return nil
}

func GenerateToken(username string) (string, error) {
	jwtKey := []byte(config.GetConfig().JWTSecret)

	// Ensure secret key is properly loaded
	if len(jwtKey) == 0 {
		return "", errors.New("JWT secret key not set in environment variables")
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": username,
		"iss":  config.GetConfig().AppName,
		"exp":  time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat":  time.Now().Unix(),                // Issued at
	})

	tokenString, err := claims.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Token generation failed
// SigningMethodHS256 (HMAC-SHA256) Type: symmetric signing algorithm
// SigningMethodES256 (ECDSA-SHA256) Type: asymmetric signing algorithm

func HandleError(c *gin.Context, httpStatus int, err error) {
	c.JSON(httpStatus, gin.H{"error": err.Error()})
}
