package handlers

import (
	"clothes-shop-api/internal/config"
	"clothes-shop-api/internal/models"
	"clothes-shop-api/internal/repositories"
	"clothes-shop-api/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	userRepo     *repositories.UserRepository
	emailService *services.EmailService
	jwtSecret    string
	cfg          config.Config
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"`
}

type UpdateEmailRequest struct {
	NewEmail string `json:"new_email" binding:"required,email"`
}

type UpdateRoleRequest struct {
	UserID  string `json:"user_id" binding:"required"`
	NewRole string `json:"new_role" binding:"required"`
}

func NewAuthHandler(userRepo *repositories.UserRepository, emailService *services.EmailService, jwtSecret string, cfg config.Config) *AuthHandler {
	return &AuthHandler{
		userRepo:     userRepo,
		emailService: emailService,
		jwtSecret:    jwtSecret,
		cfg:          cfg,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account and send verification email
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body RegisterRequest true "User registration data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Role == "" {
		req.Role = "customer"
	}

	user, err := h.userRepo.CreateUser(c.Request.Context(), req.Email, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate verification token (JWT with short expiration)
	verificationToken, err := h.generateVerificationToken(user.ID.String(), user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verification token"})
		return
	}

	// Determine base URL
	baseURL := h.cfg.AppBaseURLLocal
	if baseURL == "" {
		baseURL = "http://localhost:8080" // fallback
	}

	// Send verification email
	err = h.emailService.SendVerificationEmail(user.Email, verificationToken, baseURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully. Please check your email to verify your account."})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body LoginRequest true "User login data"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account not verified. Please check your email and verify your account."})
		return
	}

	if !h.userRepo.CheckPassword(user.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := h.generateToken(user.ID.String(), user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  *user,
	})
}

func (h *AuthHandler) generateToken(userID, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hours
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}

func (h *AuthHandler) generateVerificationToken(userID, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"type":    "verification",
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hours for verification
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}

// VerifyEmail godoc
// @Summary Verify user email
// @Description Verify user email using token from email
// @Tags auth
// @Accept  json
// @Produce  json
// @Param token query string true "Verification token"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/verify-email [get]
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	if (*claims)["type"] != "verification" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type"})
		return
	}

	userID := (*claims)["user_id"].(string)

	err = h.userRepo.ActivateUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully. You can now log in."})
}

// CreateUser godoc
// @Summary Create a new user account (Admin only)
// @Description Create a new user account with active status and send email with credentials
// @Tags admin
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param request body CreateUserRequest true "User creation data"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/users [post]
func (h *AuthHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Role == "" {
		req.Role = "customer"
	}

	user, err := h.userRepo.CreateUserActive(c.Request.Context(), req.Email, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Send email with account credentials
	err = h.emailService.SendAccountCreatedEmail(user.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send account creation email"})
		return
	}

	token, err := h.generateToken(user.ID.String(), user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User:  *user,
	})
}

// GetAllUsers godoc
// @Summary Get all users (Admin only)
// @Description Retrieve a list of all users (requires admin role)
// @Tags admin
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {array} models.User
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/users [get]
func (h *AuthHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userRepo.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// UpdateEmail godoc
// @Summary Update user email
// @Description Request to update user email, sends confirmation email to new email
// @Tags auth
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param request body UpdateEmailRequest true "New email data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/update-email [put]
func (h *AuthHandler) UpdateEmail(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req UpdateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate email update confirmation token
	confirmationToken, err := h.generateEmailUpdateToken(userID.(string), req.NewEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate confirmation token"})
		return
	}

	// Determine base URL
	baseURL := h.cfg.AppBaseURLLocal
	if baseURL == "" {
		baseURL = "http://localhost:8080" // fallback
	}

	// Send confirmation email to new email address
	err = h.emailService.SendEmailUpdateConfirmation(req.NewEmail, confirmationToken, baseURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send confirmation email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email update requested. Please check your new email address to confirm the change."})
}

// ConfirmEmailUpdate godoc
// @Summary Confirm email update
// @Description Confirm email update using token from confirmation email
// @Tags auth
// @Accept  json
// @Produce  json
// @Param token query string true "Confirmation token"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/confirm-email-update [get]
func (h *AuthHandler) ConfirmEmailUpdate(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	if (*claims)["type"] != "email_update" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type"})
		return
	}

	userID := (*claims)["user_id"].(string)
	newEmail := (*claims)["new_email"].(string)

	err = h.userRepo.UpdateUserEmail(c.Request.Context(), userID, newEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email updated successfully."})
}

func (h *AuthHandler) generateEmailUpdateToken(userID, newEmail string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"new_email": newEmail,
		"type":      "email_update",
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // 24 hours for confirmation
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}

// ToggleUserActive godoc
// @Summary Toggle user active status (Admin only)
// @Description Toggle user active/inactive status. Admin cannot toggle their own status or other admins' status
// @Tags admin
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param request body UpdateRoleRequest true "User ID to toggle"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/users/toggle-active [patch]
func (h *AuthHandler) ToggleUserActive(c *gin.Context) {
	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if admin is trying to toggle their own status
	if adminID.(string) == req.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin cannot toggle their own active status"})
		return
	}

	// Check if target user is an admin (admins cannot toggle other admins)
	targetUser, err := h.userRepo.GetUserByID(c.Request.Context(), req.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if targetUser.Role == "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot toggle another admin's active status"})
		return
	}

	err = h.userRepo.ToggleUserActive(c.Request.Context(), req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle user active status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User active status toggled successfully"})
}

// SoftDeleteUser godoc
// @Summary Soft delete user (Admin only)
// @Description Soft delete user by setting is_deleted flag. Admin cannot delete themselves or other admins
// @Tags admin
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param request body UpdateRoleRequest true "User ID to soft delete"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/users/soft-delete [delete]
func (h *AuthHandler) SoftDeleteUser(c *gin.Context) {
	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if admin is trying to delete themselves
	if adminID.(string) == req.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin cannot delete themselves"})
		return
	}

	// Check if target user is an admin (admins cannot delete other admins)
	targetUser, err := h.userRepo.GetUserByID(c.Request.Context(), req.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if targetUser.Role == "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete another admin"})
		return
	}

	err = h.userRepo.SoftDeleteUser(c.Request.Context(), req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to soft delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User soft deleted successfully"})
}

// UpdateUserRole godoc
// @Summary Update user role (Admin only)
// @Description Update user role. Admin cannot update their own role or other admins' roles
// @Tags admin
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param request body UpdateRoleRequest true "Role update data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/users/role [put]
func (h *AuthHandler) UpdateUserRole(c *gin.Context) {
	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if admin is trying to update their own role
	if adminID.(string) == req.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin cannot update their own role"})
		return
	}

	// Check if target user is an admin (admins cannot update other admins)
	targetUser, err := h.userRepo.GetUserByID(c.Request.Context(), req.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if targetUser.Role == "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot update another admin's role"})
		return
	}

	// Validate role
	validRoles := map[string]bool{"Admin": true, "Customer": true, "Staff": true}
	if !validRoles[req.NewRole] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role. Must be Admin, Customer, or Staff"})
		return
	}

	err = h.userRepo.UpdateUserRole(c.Request.Context(), req.UserID, req.NewRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User role updated successfully"})
}
