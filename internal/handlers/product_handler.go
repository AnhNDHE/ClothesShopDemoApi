package handlers

import (
	"clothes-shop-api/internal/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	repo *repositories.ProductRepository
}

type VariantRequest struct {
	Size  string  `json:"size" binding:"required"`
	Color string  `json:"color" binding:"required"`
	Stock int     `json:"stock" binding:"required,min=0"`
	Price float64 `json:"price" binding:"required,min=0"`
	Image string  `json:"image"`
}

type CreateProductRequest struct {
	Name         string           `json:"name" binding:"required"`
	Description  string           `json:"description"`
	CategoryName string           `json:"category_name" binding:"required"`
	BrandName    string           `json:"brand_name"`
	Variants     []VariantRequest `json:"variants" binding:"required,min=1"`
}

type UpdateProductRequest struct {
	Name         string           `json:"name" binding:"required"`
	Description  string           `json:"description"`
	CategoryName string           `json:"category_name" binding:"required"`
	BrandName    string           `json:"brand_name"`
	Variants     []VariantRequest `json:"variants" binding:"required,min=1"`
}

func NewProductHandler(repo *repositories.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

// GetAllProducts godoc
// @Summary Get all products with pagination and filters
// @Description Retrieve a list of products with optional pagination, price filter, category filter, brand filter, and search
// @Tags products
// @Accept  json
// @Produce  json
// @Param page query int false "Page number (default 1)" default(1)
// @Param limit query int false "Items per page (default 10)" default(10)
// @Param min_price query number false "Minimum price filter"
// @Param max_price query number false "Maximum price filter"
// @Param category query string false "Category name filter"
// @Param brand query string false "Brand name filter"
// @Param search query string false "Product name search"
// @Success 200 {array} models.Product
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	var minPrice, maxPrice *float64
	if mp := c.Query("min_price"); mp != "" {
		if parsed, err := strconv.ParseFloat(mp, 64); err == nil {
			minPrice = &parsed
		}
	}

	if mp := c.Query("max_price"); mp != "" {
		if parsed, err := strconv.ParseFloat(mp, 64); err == nil {
			maxPrice = &parsed
		}
	}

	var categoryName, brandName, searchName *string
	if cat := c.Query("category"); cat != "" {
		categoryName = &cat
	}

	if brand := c.Query("brand"); brand != "" {
		brandName = &brand
	}

	if search := c.Query("search"); search != "" {
		searchName = &search
	}

	products, err := h.repo.GetAllProducts(c.Request.Context(), page, limit, minPrice, maxPrice, categoryName, brandName, searchName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the provided details
// @Tags products
// @Accept  json
// @Produce  json
// @Param request body CreateProductRequest true "Product creation data"
// @Success 201 {object} models.Product
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculate min_price, max_price, total_stock from variants
	minPrice := req.Variants[0].Price
	maxPrice := req.Variants[0].Price
	totalStock := req.Variants[0].Stock

	for _, variant := range req.Variants {
		if variant.Price < minPrice {
			minPrice = variant.Price
		}
		if variant.Price > maxPrice {
			maxPrice = variant.Price
		}
		totalStock += variant.Stock
	}

	product, err := h.repo.CreateProduct(c.Request.Context(), req.Name, req.Description, minPrice, maxPrice, totalStock, req.CategoryName, req.BrandName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create variants
	variants := make([]struct {
		Size  string
		Color string
		Stock int
		Price float64
		Image string
	}, len(req.Variants))

	for i, v := range req.Variants {
		variants[i] = struct {
			Size  string
			Color string
			Stock int
			Price float64
			Image string
		}{
			Size:  v.Size,
			Color: v.Color,
			Stock: v.Stock,
			Price: v.Price,
			Image: v.Image,
		}
	}

	err = h.repo.CreateProductVariants(c.Request.Context(), product.ID.String(), variants)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product variants"})
		return
	}

	// Load the variants into the product response
	productVariants, err := h.repo.GetProductVariants(c.Request.Context(), product.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load product variants"})
		return
	}
	product.Variants = productVariants

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct godoc
// @Summary Update an existing product
// @Description Update a product with the provided details
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Param request body UpdateProductRequest true "Product update data"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculate min_price, max_price, total_stock from variants
	minPrice := req.Variants[0].Price
	maxPrice := req.Variants[0].Price
	totalStock := req.Variants[0].Stock

	for _, variant := range req.Variants {
		if variant.Price < minPrice {
			minPrice = variant.Price
		}
		if variant.Price > maxPrice {
			maxPrice = variant.Price
		}
		totalStock += variant.Stock
	}

	product, err := h.repo.UpdateProduct(c.Request.Context(), id, req.Name, req.Description, minPrice, maxPrice, totalStock, req.CategoryName, req.BrandName, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	// Update variants
	variants := make([]struct {
		Size  string
		Color string
		Stock int
		Price float64
		Image string
	}, len(req.Variants))

	for i, v := range req.Variants {
		variants[i] = struct {
			Size  string
			Color string
			Stock int
			Price float64
			Image string
		}{
			Size:  v.Size,
			Color: v.Color,
			Stock: v.Stock,
			Price: v.Price,
			Image: v.Image,
		}
	}

	err = h.repo.UpdateProductVariants(c.Request.Context(), id, variants)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product variants"})
		return
	}

	// Load the variants into the product response
	productVariants, err := h.repo.GetProductVariants(c.Request.Context(), product.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load product variants"})
		return
	}
	product.Variants = productVariants

	c.JSON(http.StatusOK, product)
}

// ToggleActive godoc
// @Summary Toggle product active status
// @Description Toggle the active status of a product (activate/deactivate)
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product
// @Failure 500 {object} map[string]string
// @Router /products/{id}/toggle-active [patch]
func (h *ProductHandler) ToggleActive(c *gin.Context) {
	id := c.Param("id")

	product, err := h.repo.ToggleActive(c.Request.Context(), id, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// SoftDelete godoc
// @Summary Soft delete a product
// @Description Mark a product as deleted (soft delete)
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product
// @Failure 500 {object} map[string]string
// @Router /products/{id}/soft-delete [delete]
func (h *ProductHandler) SoftDelete(c *gin.Context) {
	id := c.Param("id")

	product, err := h.repo.SoftDelete(c.Request.Context(), id, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// ToggleVariantActive godoc
// @Summary Toggle product variant active status
// @Description Toggle the active status of a product variant (activate/deactivate)
// @Tags product-variants
// @Accept  json
// @Produce  json
// @Param id path string true "Product Variant ID"
// @Success 200 {object} models.ProductVariant
// @Failure 500 {object} map[string]string
// @Router /product-variants/{id}/toggle-active [patch]
func (h *ProductHandler) ToggleVariantActive(c *gin.Context) {
	id := c.Param("id")

	variant, err := h.repo.ToggleVariantActive(c.Request.Context(), id, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle product variant active status"})
		return
	}

	c.JSON(http.StatusOK, variant)
}

// SoftDeleteVariant godoc
// @Summary Soft delete a product variant
// @Description Mark a product variant as deleted (soft delete)
// @Tags product-variants
// @Accept  json
// @Produce  json
// @Param id path string true "Product Variant ID"
// @Success 200 {object} models.ProductVariant
// @Failure 500 {object} map[string]string
// @Router /product-variants/{id}/soft-delete [delete]
func (h *ProductHandler) SoftDeleteVariant(c *gin.Context) {
	id := c.Param("id")

	variant, err := h.repo.SoftDeleteVariant(c.Request.Context(), id, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to soft delete product variant"})
		return
	}

	c.JSON(http.StatusOK, variant)
}
