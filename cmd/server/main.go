package main

import (
	"log"

	_ "clothes-shop-api/docs"
	"clothes-shop-api/internal/config"
	"clothes-shop-api/internal/db"
	"clothes-shop-api/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Clothes Shop API
// @version 1.0
// @description A RESTful API for a clothes shop built with Golang and Gin.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	_ = godotenv.Load()

	cfg := config.LoadConfig()
	db.Connect(cfg.DBUrl)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	// Setup routes
	routes.SetupRoutes(r, db.DB, cfg.JWTSecret)

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("🚀 Server running at :8080")
	log.Println("📖 Swagger UI available at: http://localhost:8080/swagger/index.html")
	r.Run(":8080")
}
