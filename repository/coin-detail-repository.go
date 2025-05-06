package repository

import (
	"blockchain-scrap/entity"

	"gorm.io/gorm"
)

type CoinDetailRepository interface {
	Save(detail *entity.CoinDetail) error
}

type coinDetailRepository struct {
	db *gorm.DB
}

func NewCoinDetailRepository(db *gorm.DB) CoinDetailRepository {
	return &coinDetailRepository{db}
}

func (r *coinDetailRepository) Save(detail *entity.CoinDetail) error {
	return r.db.Create(detail).Error
}
