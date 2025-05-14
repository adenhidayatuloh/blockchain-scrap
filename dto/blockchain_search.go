package dto

import (
	"time"

	"github.com/google/uuid"
)

// BlockchainSearchDTO merepresentasikan struct untuk pencarian blockchain
type BlockchainSearchDTO struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	ContractAddress string    `json:"contract_address"`
	ResponseData    []byte    `json:"response_data"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// BlockchainSearchResponse merepresentasikan response untuk pencarian blockchain
type BlockchainSearchResponse struct {
	ID              uuid.UUID `json:"id"`
	ContractAddress string    `json:"contract_address"`
}
