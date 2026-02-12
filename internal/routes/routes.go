package routes

import (
	"clothes-shop-api/internal/config"
	"clothes-shop-api/internal/handlers"
	"clothes-shop-api/internal/repositories"
	"clothes-shop-api/internal/services"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return handlers.JWTAuthMiddleware(jwtSecret)
}

func AdminOnlyMiddleware() gin.HandlerFunc {
	return handlers.AdminOnlyMiddleware()
}

func SetupRoutes(r *gin.Engine, jwtSecret string) {
	cfg := config.LoadConfig()

	// Initialize repositories
	productRepo := repositories.NewProductRepository(config.DB)
	userRepo := repositories.NewUserRepository(config.DB)
	categoryRepo := repositories.NewCategoryRepository(config.DB)
	brandRepo := repositories.NewBrandRepository(config.DB)

	// Initialize services
	emailService := services.NewEmailService(cfg)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productRepo)
	authHandler := handlers.NewAuthHandler(userRepo, emailService, jwtSecret, cfg)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	brandHandler := handlers.NewBrandHandler(brandRepo)

	// Auth routes
	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)
	r.GET("/auth/verify-email", authHandler.VerifyEmail)
	r.PUT("/auth/update-email", JWTAuthMiddleware(jwtSecret), authHandler.UpdateEmail)
	r.GET("/auth/confirm-email-update", authHandler.ConfirmEmailUpdate)

	// Admin routes
	admin := r.Group("/admin")
	admin.Use(JWTAuthMiddleware(jwtSecret), AdminOnlyMiddleware())
	{
		admin.POST("/users", authHandler.CreateUser)
		admin.GET("/users", authHandler.GetAllUsers)
		admin.PUT("/users/role", authHandler.UpdateUserRole)
		admin.PATCH("/users/toggle-active", authHandler.ToggleUserActive)
		admin.DELETE("/users/soft-delete", authHandler.SoftDeleteUser)
	}

	// Product routes
	r.GET("/products", productHandler.GetAllProducts)
	r.POST("/products", productHandler.CreateProduct)
	r.PUT("/products/:id", productHandler.UpdateProduct)
	r.PATCH("/products/:id/toggle-active", productHandler.ToggleActive)
	r.DELETE("/products/:id/soft-delete", productHandler.SoftDelete)

	// Product variant routes
	r.PATCH("/product-variants/:id/toggle-active", productHandler.ToggleVariantActive)
	r.DELETE("/product-variants/:id/soft-delete", productHandler.SoftDeleteVariant)

	// Category routes
	r.GET("/categories", categoryHandler.GetAllCategories)
	r.POST("/categories", categoryHandler.CreateCategory)
	r.PUT("/categories/:id", categoryHandler.UpdateCategory)
	r.PATCH("/categories/:id/toggle-active", categoryHandler.ToggleCategoryActive)
	r.DELETE("/categories/:id/soft-delete", categoryHandler.SoftDeleteCategory)

	// Brand routes
	r.GET("/brands", brandHandler.GetAllBrands)
	r.POST("/brands", brandHandler.CreateBrand)
	r.PUT("/brands/:id", brandHandler.UpdateBrand)
	r.PATCH("/brands/:id/toggle-active", brandHandler.ToggleBrandActive)
	r.DELETE("/brands/:id/soft-delete", brandHandler.SoftDeleteBrand)
}
