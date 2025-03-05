package routes

import (
	"net/http"

	"github.com/coderj001/URL-shortener/api/database"
	"github.com/gin-gonic/gin"
)

func ResolveURL(c *gin.Context, db *database.MySQLStore) {
	short := c.Param("url")

	original, err := db.GetURL(short)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if original == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "short URL not found"})
		return
	}

	err = db.ClickCount(short)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error while updating click count"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, original)
}
