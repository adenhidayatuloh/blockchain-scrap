package handler

import (
	"blockchain-scrap/dto"
	"blockchain-scrap/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user authentication requests
type UserHandler struct {
	userSvc service.UserService
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userSvc service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

// Register handles new user registration requests
// Register godoc
// @Summary Register new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register Request"
// @Success 201 {object} dto.RegisterResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := h.userSvc.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(err.StatusCode(), gin.H{"error": err.Message()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login handles user login requests
// Login godoc
// @Summary User login
// @Description Authenticate user and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	token, err := h.userSvc.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(err.StatusCode(), gin.H{"error": err.Message()})
		return
	}

	c.JSON(http.StatusOK, token)
}
