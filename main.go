package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/coderj001/URL-shortener/api/routes"
	"github.com/coderj001/URL-shortener/config"
	"github.com/coderj001/URL-shortener/database"
	"github.com/coderj001/URL-shortener/logger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func setupRoutes(router *gin.Engine, db *database.MySQLStore) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ping": "pong",
		})
	})

	router.GET("/:url", func(c *gin.Context) {
		routes.ResolveURL(c, db)
	})
	router.POST("/api/v1", func(c *gin.Context) {
		routes.ShortenURL(c, db)
	})
	router.GET("api/v1/analytics/:shortID", func(c *gin.Context) {
		routes.AnalyticsShortURL(c, db)
	})
}

func main() {

	db, err := database.NewMySQLStore()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	router := gin.Default()

	//? Middleware
	router.Use(logger.Logger())

	setupRoutes(router, db)

	log.Fatal(router.Run(fmt.Sprintf(":%s", config.GetConfig().Port)))
}
