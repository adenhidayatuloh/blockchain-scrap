package dto

import "github.com/google/uuid"

type ContractAddressResponse struct {
	ID              string           `json:"id"`
	Symbol          string           `json:"symbol"`
	Platform        string           `json:"contract_address"`
	WebSlug         string           `json:"web_slug"`
	ListingDay      int              `json:"listing_day"`
	MarketData      MarketData       `json:"market_data"`
	TimePrices      []PricePoint     `json:"timestamp_prices"`
	LiquidityInfo   DexLiquidityInfo `json:"dex_liquidity_info"`
	TokenAnalytics  TokenAnalytics   `json:"token_analytics"`
	Image           Image            `json:"image"`
	SummaryAnalysis string           `json:"summary_analysis"`
}

type MarketData struct {
	CurrentPrice              CurrencyValue `json:"current_price"`
	PriceChangePercentage1h   CurrencyValue `json:"price_change_percentage_1h_in_currency"`
	MarketCap                 CurrencyValue `json:"market_cap"`
	TotalVolume               CurrencyValue `json:"total_volume"`
	MarketCapChangePercentage CurrencyValue `json:"market_cap_change_percentage_24h_in_currency"`
	FullyDilutedValuation     CurrencyValue `json:"fully_diluted_valuation"`
	Liquidity                 CurrencyValue `json:"liquidity"`

	CirculatingSupply float64 `json:"circulating_supply"`
	TotalSupply       float64 `json:"total_supply"`
	MaxSupply         float64 `json:"max_supply"`
}

type DexLiquidityInfo struct {
	LiquidityPoolSize float64 `json:"liquidity_pool_size"`
	TopDex            string  `json:"top_dex"`
	Volume24h         float64 `json:"volume_24h"`
	SlippageNote      string  `json:"slippage_note"`
	DexLiquidityRatio float64 `json:"dex_liquidity_ratio"`
	LiquidityTrend7D  float64 `json:"liquidity_trend_7d"`
}
type TokenAnalytics struct {
	TopHolder   float64 `json:"top_holder"`
	TopWallets  float64 `json:"top_wallets"`
	TokenViewer float64 `json:"token_viewer"`
	SniperBot   float64 `json:"sniper_bot"`
	DevSold     bool    `json:"dev_sold"`
	DevBuyback  bool    `json:"dev_buyback"`
}

// General purpose for USD values
type CurrencyValue struct {
	USD float64 `json:"usd"`
}

type PricePoint struct {
	Timestamp string  `json:"timestamp"` // formatted timestamp
	Price     float64 `json:"price"`
}

type GetPricesRequest struct {
	Prices [][]float64 `json:"prices"`
}

type GetLiquidityRequest struct {
	Liquidity CurrencyValue `json:"liquidity"`
}

type Image struct {
	Small string `json:"small"`
}

type AIRequest struct {
	Prompt      string                  `json:"prompt"`
	Collections ContractAddressResponse `json:"collections"`
}

type AIResponse struct {
	AssistantMessage string `json:"assistant_message"`
}

type SearchRecordResponse struct {
	ID              uuid.UUID `json:"id"`
	ContractAddress string    `json:"contract_address"`
}
