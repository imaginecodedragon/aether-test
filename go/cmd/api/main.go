package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func newDatabase() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
}

type shipsResponse struct {
	Ships   []string `json:"ships"`
	Count   int      `json:"count"`
	Message string   `json:"message,omitempty"`
}

func newRouter(db *gorm.DB, shipService ShipService) *gin.Engine {
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

	router.GET("/ships", func(c *gin.Context) {
		ships := shipService.ListShips()
		response := shipsResponse{
			Ships: ships,
			Count: len(ships),
		}

		if len(ships) == 0 {
			response.Message = "No ships available."
		}

		c.JSON(http.StatusOK, response)
	})

	return router
}

func main() {
	db, err := newDatabase()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	router := newRouter(db, newShipService())

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
