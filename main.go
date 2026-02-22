package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/shaik-aaron/fantasy-backend/controllers"
	"github.com/shaik-aaron/fantasy-backend/intializers"
	"github.com/shaik-aaron/fantasy-backend/middleware"
)

func init() {
	intializers.LoadEnv()
	intializers.ConnectToDb()
	intializers.MigrateDb()
}

func main() {
	fmt.Println("Hello, World 2!")

	router := gin.Default()

	// Add CORS middleware
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/signup", controllers.SignUp)
	router.POST("/login", controllers.Login)
	router.GET("/validate", middleware.RequireAuth, controllers.Validate)
	router.POST("/sessions", middleware.RequireAuth, controllers.CreateSession)
	router.GET("/sessions/:userId", middleware.RequireAuth, controllers.GetSessions)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
