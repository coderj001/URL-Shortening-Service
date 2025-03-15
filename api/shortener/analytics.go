package shortener

import (
	"fmt"
	"net/http"

	apitypes "github.com/coderj001/URL-shortener/api"
	"github.com/coderj001/URL-shortener/auth"
	"github.com/coderj001/URL-shortener/database"
	"github.com/coderj001/URL-shortener/helpers"
	"github.com/gin-gonic/gin"
)

func AnalyticsShortURL(c *gin.Context, db *database.MySQLStore) {
	shortID := c.Param("shortID")
	analytics, err := db.GetURLAnalytics(shortID)

	if err != nil {
		helpers.HandleError(
			c,
			http.StatusInternalServerError,
			fmt.Errorf("database error, %v", err),
		)
		return
	}

	c.JSON(http.StatusOK, analytics)
}

func AnalyticsShortIDList(c *gin.Context, db *database.MySQLStore) {
	value, _ := c.Get("auth_status")

	var userID uint
	if value == auth.Authorized {
		if v, ok := c.Get("user"); ok {
			if user, ok := v.(apitypes.User); ok {
				userID = user.ID
			}
		}
	} else {
		helpers.HandleError(c, http.StatusUnauthorized, fmt.Errorf("unauthorized access"))
		return
	}

	urlData, err := db.GetShortIDList(userID)
	if err != nil {
		helpers.HandleError(c, http.StatusConflict, fmt.Errorf("error while fetching database"))
		return
	}

	c.JSON(http.StatusOK, urlData)
}
