package dto

type ContractAddressResponse struct {
	ID         string       `json:"id"`
	Symbol     string       `json:"symbol"`
	Platform   string       `json:"contract_address"` // disesuaikan json tag
	WebSlug    string       `json:"web_slug"`
	MarketData MarketData   `json:"market_data"`
	TimePrices []PricePoint `json:"timestamp_prices"`
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
