package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// BlockchainSearch merepresentasikan riwayat pencarian kontrak blockchain
type BlockchainSearch struct {
	ID              uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID          uuid.UUID       `gorm:"type:uuid;not null;index" json:"user_id"`
	ContractAddress string          `gorm:"not null;index" json:"contract_address"`
	ResponseData    json.RawMessage `gorm:"type:jsonb" json:"response_data"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
