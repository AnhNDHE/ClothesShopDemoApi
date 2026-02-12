package handlers

import (
	"clothes-shop-api/internal/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	repo *repositories.CategoryRepository
}

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func NewCategoryHandler(repo *repositories.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{repo: repo}
}

// GetAllCategories godoc
// @Summary Get all categories
// @Description Retrieve a list of all categories
// @Tags categories
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Category
// @Router /categories [get]
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.repo.GetAllCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category with the provided details
// @Tags categories
// @Accept  json
// @Produce  json
// @Param request body CreateCategoryRequest true "Category creation data"
// @Success 201 {object} models.Category
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.repo.CreateCategory(c.Request.Context(), req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// UpdateCategory godoc
// @Summary Update an existing category
// @Description Update a category with the provided details
// @Tags categories
// @Accept  json
// @Produce  json
// @Param id path string true "Category ID"
// @Param request body UpdateCategoryRequest true "Category update data"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.repo.UpdateCategory(c.Request.Context(), id, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// ToggleCategoryActive godoc
// @Summary Toggle category active status
// @Description Toggle the active status of a category (activate/deactivate)
// @Tags categories
// @Accept  json
// @Produce  json
// @Param id path string true "Category ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories/{id}/toggle-active [patch]
func (h *CategoryHandler) ToggleCategoryActive(c *gin.Context) {
	id := c.Param("id")

	err := h.repo.ToggleCategoryActive(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle category active status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category active status toggled successfully"})
}

// SoftDeleteCategory godoc
// @Summary Soft delete a category
// @Description Mark a category as deleted (soft delete)
// @Tags categories
// @Accept  json
// @Produce  json
// @Param id path string true "Category ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories/{id}/soft-delete [delete]
func (h *CategoryHandler) SoftDeleteCategory(c *gin.Context) {
	id := c.Param("id")

	err := h.repo.SoftDeleteCategory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to soft delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category soft deleted successfully"})
}
