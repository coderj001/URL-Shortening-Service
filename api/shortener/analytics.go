package shortener

import (
	"fmt"
	"net/http"

	"github.com/coderj001/URL-shortener/database"
	"github.com/coderj001/URL-shortener/helpers"
	"github.com/gin-gonic/gin"
)

func AnalyticsShortURL(c *gin.Context, db *database.MySQLStore) {
	shortID := c.Param("shortID")
	analytics, err := db.GetURLAnalytics(shortID)

	if err != nil {
		helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("database error, %v", err))
		return
	}

	c.JSON(http.StatusOK, analytics)
}
