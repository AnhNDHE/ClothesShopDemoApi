package routes

import (
	"clothes-shop-api/internal/config"
	"clothes-shop-api/internal/handlers"
	"clothes-shop-api/internal/repositories"
	"clothes-shop-api/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, jwtSecret string) {
	cfg := config.LoadConfig()

	// Initialize repositories
	productRepo := repositories.NewProductRepository(config.DB)
	userRepo := repositories.NewUserRepository(config.DB)

	// Initialize services
	emailService := services.NewEmailService(cfg)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productRepo)
	authHandler := handlers.NewAuthHandler(userRepo, emailService, jwtSecret, cfg)

	// Auth routes
	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)
	r.GET("/auth/verify-email", authHandler.VerifyEmail)

	// Product routes
	r.GET("/products", productHandler.GetAllProducts)
	r.POST("/products", productHandler.CreateProduct)
	r.PUT("/products/:id", productHandler.UpdateProduct)
	r.PATCH("/products/:id/toggle-active", productHandler.ToggleActive)
	r.DELETE("/products/:id/soft-delete", productHandler.SoftDelete)

	// Product variant routes
	r.PATCH("/product-variants/:id/toggle-active", productHandler.ToggleVariantActive)
	r.DELETE("/product-variants/:id/soft-delete", productHandler.SoftDeleteVariant)

	r.GET("/categories", productHandler.GetAllCategories)
	r.GET("/brands", productHandler.GetAllBrands)
}
