package shortener

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	apitypes "github.com/coderj001/URL-shortener/api"
	"github.com/coderj001/URL-shortener/auth"
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
	value, _ := c.Get("auth_status")

	var body request
	if err := helpers.ParseRequest(c, &body); err != nil {
		helpers.HandleError(
			c,
			http.StatusBadRequest,
			fmt.Errorf("invalid request"),
		)
		return
	}

	// Rate limiting
	remaining, resetAt, err := db.CheckRateLimit(c.ClientIP())
	if err != nil {
		helpers.HandleError(
			c,
			http.StatusInternalServerError,
			fmt.Errorf("rate limit check failed"),
		)
		return
	}

	if remaining < 0 {
		helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf(
			"rate limit exceeded, wait till %v",
			time.Until(resetAt).Minutes(),
		))
		return
	}

	// URL validation
	if !govalidator.IsURL(body.URL) {
		helpers.HandleError(
			c,
			http.StatusBadRequest,
			fmt.Errorf("falid to validate"),
		)
		return
	}

	if !helpers.RemoveDomainError(body.URL) {
		helpers.HandleError(
			c,
			http.StatusForbidden,
			fmt.Errorf("domain not allowed"),
		)
		return
	}

	body.URL = helpers.EnforceHTTP(body.URL)

	length := 9
	if value == auth.Authorized {
		length = 4
	}
	short, err := helpers.GenerateID(length)
	if err != nil {
		helpers.HandleError(
			c,
			http.StatusForbidden,
			fmt.Errorf("unable to generate shortid"),
		)
		return
	}

	// Check existing short URL
	existing, err := db.GetURL(short)
	if err != nil {
		helpers.HandleError(
			c,
			http.StatusInternalServerError,
			fmt.Errorf("database error"),
		)
		return
	}
	if existing != "" {
		helpers.HandleError(
			c,
			http.StatusConflict,
			fmt.Errorf("short URL already exists"),
		)
		return
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	var userID *uint
	if value == auth.Authorized {
		if v, ok := c.Get("user"); ok {
			if user, ok := v.(apitypes.User); ok {
				userID = &user.ID
			}
		}
	}
	if err := db.SaveURL(short, body.URL, body.Expiry, userID); err != nil {
		helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("failed to save URL"))
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
