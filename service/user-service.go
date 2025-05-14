package service

import (
	"context"
	"net/http"

	"blockchain-scrap/dto"
	"blockchain-scrap/entity"
	"blockchain-scrap/pkg/errs"
	"blockchain-scrap/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserService defines the contract for user services
type UserService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, errs.MessageErr)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, errs.MessageErr)
	Authentication() gin.HandlerFunc
}

// userServiceImpl implements UserService
type userServiceImpl struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(repo repository.UserRepository) UserService {
	return &userServiceImpl{userRepo: repo}
}

// Register registers a new user
func (s *userServiceImpl) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, errs.MessageErr) {
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errs.NewBadRequest("Email already registered")
	}

	newUser := &entity.User{
		ID:       uuid.New(),
		Email:    req.Email,
		Password: req.Password,
	}

	if err := newUser.HashPassword(); err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, errs.NewInternalServerError("Failed to create new user")
	}

	return &dto.RegisterResponse{
		ID:    newUser.ID,
		Email: newUser.Email,
	}, nil
}

// Login authenticates a user
func (s *userServiceImpl) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, errs.MessageErr) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errs.NewBadRequest("Invalid email or password")
	}

	if err := user.ComparePassword(req.Password); err != nil {
		return nil, errs.NewBadRequest(err.Error())
	}

	token, err := user.CreateToken()
	if err != nil {
		return nil, errs.NewInternalServerError(err.Error())
	}

	return &dto.LoginResponse{
		Token: token,
	}, nil
}

// Authentication middleware to validate JWT token
func (s *userServiceImpl) Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		var user entity.User
		if err := user.ValidateToken(authHeader); err != nil {
			c.AbortWithStatusJSON(err.StatusCode(), err)
			return
		}

		authenticatedUser, err := s.userRepo.FindByEmail(c.Request.Context(), user.Email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}

		c.Set("userData", authenticatedUser)
		c.Next()
	}
}
