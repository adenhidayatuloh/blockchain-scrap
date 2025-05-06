package entity

import "gorm.io/gorm"

type CoinDetail struct {
	gorm.Model
	CoinID                   string  `gorm:"column:coin_id"`
	Symbol                   string  `gorm:"column:symbol"`
	ContractAddress          string  `gorm:"column:contract_address;uniqueIndex"`
	WebSlug                  string  `gorm:"column:web_slug"`
	CurrentPriceUSD          float64 `gorm:"column:current_price_usd"`
	PriceChangePct1HUSD      float64 `gorm:"column:price_change_pct_1h_usd"`
	MarketCapUSD             float64 `gorm:"column:market_cap_usd"`
	TotalVolumeUSD           float64 `gorm:"column:total_volume_usd"`
	MarketCapChangePct24HUSD float64 `gorm:"column:market_cap_change_pct_24h_usd"`
	FullyDilutedValuationUSD float64 `gorm:"column:fdv_usd"`
	LiquidityUSD             float64 `gorm:"column:liquidity_usd"`
	CirculatingSupply        float64 `gorm:"column:circulating_supply"`
	TotalSupply              float64 `gorm:"column:total_supply"`
	MaxSupply                float64 `gorm:"column:max_supply"`
}
