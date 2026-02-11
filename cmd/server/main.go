package main

import (
	"log"
	"os"

	_ "clothes-shop-api/docs"
	"clothes-shop-api/internal/config"
	"clothes-shop-api/internal/routes"

	"github.com/gin-contrib/cors"
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

// ❗ KHÔNG set @host → để swagger tự nhận domain (local / render)
// @BasePath /
func main() {
	// Load .env (local only, render sẽ ignore nếu không có)
	_ = godotenv.Load()

	// Init DB + migration
	config.InitDB()
	config.RunMigration()

	cfg := config.LoadConfig()

	// Gin
	r := gin.Default()

	// CORS (tạm open, sau này fix theo domain FE)
	r.Use(cors.Default())

	// Root → redirect swagger
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})

	// Healthcheck (Render dùng rất nhiều)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	// API routes
	routes.SetupRoutes(r, cfg.JWTSecret)

	// PORT (Render cấp động)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // local
	}

	// SWAGGER_HOST (flexible cho local / production)
	var swaggerScheme string
	swaggerHost := os.Getenv("SWAGGER_HOST")
	if swaggerHost == "" {
		swaggerHost = "localhost:" + port // local default
		swaggerScheme = "http"
	} else {
		swaggerScheme = "https" // production default to https
	}

	// SWAGGER_SCHEME (override if set)
	if envScheme := os.Getenv("SWAGGER_SCHEME"); envScheme != "" {
		swaggerScheme = envScheme
	}

	// Swagger with dynamic host and scheme
	swaggerURL := swaggerScheme + "://" + swaggerHost + "/swagger/doc.json"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL(swaggerURL)))

	log.Println("🚀 Server running on port:", port)
	log.Println("📖 Swagger UI: /swagger/index.html")
	log.Println("🌐 Swagger Host set to:", swaggerHost)

	// Start server
	if err := r.Run(":" + port); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
