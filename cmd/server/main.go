package main

import (
	"log"
	"os"

	"clothes-shop-api/docs"
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
// @BasePath /
func main() {

	// Load .env (local only)
	_ = godotenv.Load()

	// Init database
	config.InitDB()
	config.RunMigration()

	cfg := config.LoadConfig()

	// Create Gin server
	r := gin.Default()

	// CORS (open for demo)
	r.Use(cors.Default())

	// Root → redirect to swagger
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})

	// Health check (Render dùng)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	// Setup API routes
	routes.SetupRoutes(r, cfg.JWTSecret)

	// =============================
	// 🔥 DYNAMIC PORT (Render)
	// =============================
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// =============================
	// 🔥 FIX SWAGGER HOST + SCHEME
	// =============================
	swaggerHost := os.Getenv("SWAGGER_HOST")
	if swaggerHost == "" {
		swaggerHost = "localhost:" + port
	}

	swaggerScheme := "http"
	if swaggerHost != "localhost:"+port {
		swaggerScheme = "https"
	}

	// 👇 QUAN TRỌNG
	docs.SwaggerInfo.Host = swaggerHost
	docs.SwaggerInfo.Schemes = []string{swaggerScheme}

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("====================================")
	log.Println("🚀 Server running on port:", port)
	log.Println("🌐 Swagger Host:", swaggerHost)
	log.Println("📖 Swagger UI: /swagger/index.html")
	log.Println("====================================")

	// Start server
	if err := r.Run(":" + port); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
