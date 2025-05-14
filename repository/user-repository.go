package repository

import (
	"context"

	"blockchain-scrap/entity"
	"blockchain-scrap/pkg/errs"

	"gorm.io/gorm"
)

// UserRepository defines the contract for user-related database operations
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) errs.MessageErr
	FindByEmail(ctx context.Context, email string) (*entity.User, errs.MessageErr)
}

// userRepositoryImpl implements UserRepository
type userRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

// Create saves new user data to the database
func (r *userRepositoryImpl) Create(ctx context.Context, user *entity.User) errs.MessageErr {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return errs.NewInternalServerError(err.Error())
	}
	return nil
}

// FindByEmail searches for a user by email address
func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.User, errs.MessageErr) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, errs.NewNotFound("email not found")
	}
	return &user, nil
}
