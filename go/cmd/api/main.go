package main

import (
	"log"
	"math/rand/v2"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var groguQuotes = []string{
	"Fear is the path to the dark side.",
	"Do. Or do not. There is no try.",
	"Luminous beings are we, not this crude matter.",
}

func newDatabase() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
}

func randomGroguQuote() (string, bool) {
	if len(groguQuotes) == 0 {
		return "", false
	}

	return groguQuotes[rand.IntN(len(groguQuotes))], true
}

func newRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database unavailable"})
			return
		}

		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database ping failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	router.GET("/grogu", func(c *gin.Context) {
		quote, ok := randomGroguQuote()
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no grogu quotes configured"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"quote": quote,
		})
	})

	return router
}

func main() {
	db, err := newDatabase()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	router := newRouter(db)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
