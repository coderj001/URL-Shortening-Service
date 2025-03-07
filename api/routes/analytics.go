package routes

import (
	"net/http"

	"github.com/coderj001/URL-shortener/api/database"
	"github.com/gin-gonic/gin"
)

func AnalyticsShortURL(c *gin.Context, db *database.MySQLStore) {
	shortID := c.Param("shortID")

	counts, err := db.GetClickCount(shortID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": counts})
}
