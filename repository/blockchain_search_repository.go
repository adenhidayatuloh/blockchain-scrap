package repository

import (
	"blockchain-scrap/entity"
	"blockchain-scrap/pkg/errs"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BlockchainSearchRepository defines the contract for database operations related to search history
type BlockchainSearchRepository interface {
	Save(ctx context.Context, record *entity.BlockchainSearch) errs.MessageErr
	Update(ctx context.Context, record *entity.BlockchainSearch) errs.MessageErr
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.BlockchainSearch, errs.MessageErr)
	FindByUserIDAndContract(ctx context.Context, userID uint, contractAddress string) (*entity.BlockchainSearch, errs.MessageErr)
	FindByID(ctx context.Context, ID uuid.UUID) (*entity.BlockchainSearch, errs.MessageErr)
	SaveOrUpdate(ctx context.Context, record *entity.BlockchainSearch) errs.MessageErr
}

// blockchainSearchRepositoryImpl implements BlockchainSearchRepository
type blockchainSearchRepositoryImpl struct {
	db *gorm.DB
}

// NewBlockchainSearchRepository creates a new instance of BlockchainSearchRepository
func NewBlockchainSearchRepository(db *gorm.DB) BlockchainSearchRepository {
	return &blockchainSearchRepositoryImpl{db: db}
}

// FindByID searches for search history by ID
func (r *blockchainSearchRepositoryImpl) FindByID(ctx context.Context, ID uuid.UUID) (*entity.BlockchainSearch, errs.MessageErr) {
	var record entity.BlockchainSearch
	err := r.db.WithContext(ctx).Where("id = ?", ID).First(&record).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewNotFound("Data not found")
		}
		return nil, errs.NewInternalServerError("Failed to fetch data")
	}
	return &record, nil
}

// Save saves a new search history
func (r *blockchainSearchRepositoryImpl) Save(ctx context.Context, record *entity.BlockchainSearch) errs.MessageErr {
	err := r.db.WithContext(ctx).Create(record).Error
	if err != nil {
		return errs.NewInternalServerError("Failed to save data")
	}
	return nil
}

// Update updates search history data
func (r *blockchainSearchRepositoryImpl) Update(ctx context.Context, record *entity.BlockchainSearch) errs.MessageErr {
	err := r.db.WithContext(ctx).Save(record).Error
	if err != nil {
		return errs.NewInternalServerError("Failed to update data")
	}
	return nil
}

// FindByUserID searches for all search history by user ID
func (r *blockchainSearchRepositoryImpl) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.BlockchainSearch, errs.MessageErr) {
	var records []*entity.BlockchainSearch
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&records).Error
	if err != nil {
		return nil, errs.NewInternalServerError("Failed to fetch data")
	}
	return records, nil
}

// FindByUserIDAndContract searches for search history by user ID and contract address
func (r *blockchainSearchRepositoryImpl) FindByUserIDAndContract(ctx context.Context, userID uint, contractAddress string) (*entity.BlockchainSearch, errs.MessageErr) {
	var record entity.BlockchainSearch
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND contract_address = ?", userID, contractAddress).
		First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewNotFound("Data not found")
		}
		return nil, errs.NewInternalServerError("Failed to fetch data")
	}
	return &record, nil
}

// SaveOrUpdate saves or updates search history
func (r *blockchainSearchRepositoryImpl) SaveOrUpdate(ctx context.Context, record *entity.BlockchainSearch) errs.MessageErr {
	var existingRecord entity.BlockchainSearch

	if err := r.db.Where("contract_address = ? AND user_id = ?", record.ContractAddress, record.UserID).First(&existingRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := r.db.Create(record).Error; err != nil {
				return errs.NewInternalServerError("Failed to save data")
			}
		} else {
			return errs.NewInternalServerError("Failed to process data")
		}
	} else {
		record.ID = existingRecord.ID
		if err := r.db.Save(record).Error; err != nil {
			return errs.NewInternalServerError("Failed to update data")
		}
	}
	return nil
}
