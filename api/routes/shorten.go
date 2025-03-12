package routes

import (
	"math"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/coderj001/URL-shortener/database"
	"github.com/coderj001/URL-shortener/helpers"
	"github.com/gin-gonic/gin"
)

type request struct {
	URL    string        `json:"url"`
	Expiry time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	ShortID         string        `json:"short_id"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset float64       `json:"rate_limit_reset"`
}

func ShortenURL(c *gin.Context, db *database.MySQLStore) {
	var body request
	if err := helpers.ParseRequest(c, &body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Rate limiting
	remaining, resetAt, err := db.CheckRateLimit(c.ClientIP())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "rate limit check failed"})
		return
	}

	if remaining < 0 {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error":            "rate limit exceeded",
			"rate_limit_reset": time.Until(resetAt).Minutes(),
		})
		return
	}

	// URL validation
	if !govalidator.IsURL(body.URL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid URL"})
		return
	}

	if !helpers.RemoveDomainError(body.URL) {
		c.JSON(http.StatusForbidden, gin.H{"error": "domain not allowed"})
		return
	}

	body.URL = helpers.EnforceHTTP(body.URL)

	// TODO: update for premium users - can have 4 digits shortner
	var short string
	short, err = helpers.GenerateID(9)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unable to generate shortID"})
		return
	}

	// Check existing short URL
	existing, err := db.GetURL(short)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if existing != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "short URL already exists"})
		return
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	if err := db.SaveURL(short, body.URL, body.Expiry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save URL"})
		return
	}

	c.JSON(http.StatusOK, response{
		URL:             body.URL,
		ShortID:         short,
		Expiry:          body.Expiry,
		XRateRemaining:  remaining,
		XRateLimitReset: math.Ceil(time.Until(resetAt).Minutes()),
	})
}
