package repository

import (
	"blockchain-scrap/entity"
	"blockchain-scrap/pkg/errs"

	"gorm.io/gorm"
)

type TokenRepository interface {
	GetAll(limit, offset int, search string) ([]*entity.Token, int64, errs.MessageErr)
	FindByAddress(addresses []string) ([]*entity.Token, errs.MessageErr)
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db}
}

func (r *tokenRepository) FindByAddress(addresses []string) ([]*entity.Token, errs.MessageErr) {
	var record []*entity.Token
	err := r.db.Debug().
		Where("address IN ?", addresses).
		Find(&record).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewNotFound("Token in that contract address not found, please change different contract address")
		}
		return nil, errs.NewInternalServerError(err.Error())
	}

	return record, nil
}

func (r *tokenRepository) GetAll(limit, offset int, search string) ([]*entity.Token, int64, errs.MessageErr) {
	var tokens []*entity.Token
	var total int64

	query := r.db.Model(&entity.Token{})

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("LOWER(name) LIKE LOWER(?) OR LOWER(symbol) LIKE LOWER(?)", searchTerm, searchTerm)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errs.NewInternalServerError(err.Error())
	}

	if err := query.Limit(limit).Offset(offset).Order("id asc").Find(&tokens).Error; err != nil {
		return nil, 0, errs.NewInternalServerError(err.Error())
	}

	return tokens, total, nil
}
