package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/shaik-aaron/fantasy-backend/controllers"
	"github.com/shaik-aaron/fantasy-backend/intializers"
	"github.com/shaik-aaron/fantasy-backend/middleware"
)

func splitTrim(s, sep string) []string {
	var result []string
	for _, v := range strings.Split(s, sep) {
		if trimmed := strings.TrimSpace(v); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func init() {
	intializers.LoadEnv()
	intializers.ConnectToDb()
	intializers.MigrateDb()
}

func main() {
	fmt.Println("Hello, World 2!")

	router := gin.Default()

	// Add CORS middleware (comma-separated: https://app.vercel.app,http://localhost:5173)
	frontendURLs := os.Getenv("FRONTEND_URL")
	if frontendURLs == "" {
		frontendURLs = "http://localhost:5173"
	}
	allowedOrigins := make(map[string]bool)
	for _, o := range splitTrim(frontendURLs, ",") {
		allowedOrigins[strings.TrimSuffix(o, "/")] = true
	}
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			origin = strings.TrimSuffix(origin, "/")
			return allowedOrigins[origin]
		},
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
