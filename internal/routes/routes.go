package routes

import (
	"clothes-shop-api/internal/config"
	"clothes-shop-api/internal/handlers"
	"clothes-shop-api/internal/repositories"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, jwtSecret string) {
	// Initialize repositories
	productRepo := repositories.NewProductRepository(config.DB)
	userRepo := repositories.NewUserRepository(config.DB)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productRepo)
	authHandler := handlers.NewAuthHandler(userRepo, jwtSecret)

	// Auth routes
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

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
}
