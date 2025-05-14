// entity/token.go

package entity

import (
	"time"

	"gorm.io/datatypes"
)

type Token struct {
	ID                uint   `gorm:"primaryKey"`
	Address           string `gorm:"uniqueIndex;size:100"`
	CreatedAt         time.Time
	DailyVolume       float64
	Decimals          int
	FreezeAuthority   *string
	LogoURI           string
	MintAuthority     *string
	MintedAt          *time.Time
	Name              string
	PermanentDelegate *string
	Symbol            string
	Tags              datatypes.JSON
	Extensions        datatypes.JSON
}
