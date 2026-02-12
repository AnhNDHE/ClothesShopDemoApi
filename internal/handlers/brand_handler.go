package handlers

import (
	"clothes-shop-api/internal/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BrandHandler struct {
	repo *repositories.BrandRepository
}

type CreateBrandRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateBrandRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func NewBrandHandler(repo *repositories.BrandRepository) *BrandHandler {
	return &BrandHandler{repo: repo}
}

// GetAllBrands godoc
// @Summary Get all brands
// @Description Retrieve a list of all brands
// @Tags brands
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Brand
// @Router /brands [get]
func (h *BrandHandler) GetAllBrands(c *gin.Context) {
	brands, err := h.repo.GetAllBrands(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, brands)
}

// CreateBrand godoc
// @Summary Create a new brand
// @Description Create a new brand with the provided details
// @Tags brands
// @Accept  json
// @Produce  json
// @Param request body CreateBrandRequest true "Brand creation data"
// @Success 201 {object} models.Brand
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /brands [post]
func (h *BrandHandler) CreateBrand(c *gin.Context) {
	var req CreateBrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	brand, err := h.repo.CreateBrand(c.Request.Context(), req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create brand"})
		return
	}

	c.JSON(http.StatusCreated, brand)
}

// UpdateBrand godoc
// @Summary Update an existing brand
// @Description Update a brand with the provided details
// @Tags brands
// @Accept  json
// @Produce  json
// @Param id path string true "Brand ID"
// @Param request body UpdateBrandRequest true "Brand update data"
// @Success 200 {object} models.Brand
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /brands/{id} [put]
func (h *BrandHandler) UpdateBrand(c *gin.Context) {
	id := c.Param("id")
	var req UpdateBrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	brand, err := h.repo.UpdateBrand(c.Request.Context(), id, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update brand"})
		return
	}

	c.JSON(http.StatusOK, brand)
}

// ToggleBrandActive godoc
// @Summary Toggle brand active status
// @Description Toggle the active status of a brand (activate/deactivate)
// @Tags brands
// @Accept  json
// @Produce  json
// @Param id path string true "Brand ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /brands/{id}/toggle-active [patch]
func (h *BrandHandler) ToggleBrandActive(c *gin.Context) {
	id := c.Param("id")

	err := h.repo.ToggleBrandActive(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle brand active status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Brand active status toggled successfully"})
}

// SoftDeleteBrand godoc
// @Summary Soft delete a brand
// @Description Mark a brand as deleted (soft delete)
// @Tags brands
// @Accept  json
// @Produce  json
// @Param id path string true "Brand ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /brands/{id}/soft-delete [delete]
func (h *BrandHandler) SoftDeleteBrand(c *gin.Context) {
	id := c.Param("id")

	err := h.repo.SoftDeleteBrand(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to soft delete brand"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Brand soft deleted successfully"})
}
