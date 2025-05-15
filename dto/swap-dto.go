package dto

type SwapRequest struct {
	PublicKey  string  `json:"publicKey" binding:"required"`
	InputMint  string  `json:"inputMint" binding:"required"`
	OutputMint string  `json:"outputMint" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
}

type SwapResponse struct {
	Transaction string `json:"transaction"`
}

type SubmitRequest struct {
	SignedTransaction string `json:"signedTransaction" binding:"required"`
}

type SubmitResponse struct {
	Signature string `json:"signature"`
}

type JupiterErrorResponse struct {
	Error     string `json:"error"`
	ErrorCode string `json:"errorCode"`
}

type GetCurrencySwapResponse struct {
	InAmount        float64 `json:"in_amount"`
	OutAmount       float64 `json:"out_mount"`
	SwapUsdValue    float64 `json:"swap_usd_value"`
	IsSwappable     bool    `json:"is_swappable"`
	BalanceInAmount float64 `json:"balance_in_amount"`
}

type QuoteResponse struct {
	InputMint      string       `json:"inputMint"`
	OutputMint     string       `json:"outputMint"`
	InAmount       string       `json:"inAmount"`
	OutAmount      string       `json:"outAmount"`
	SlippageBps    int          `json:"slippageBps"`
	PriceImpactPct string       `json:"priceImpactPct"`
	SwapUsdValue   string       `json:"swapUsdValue"`
	RoutePlan      []RouteEntry `json:"routePlan"`
}

type RouteEntry struct {
	Percent  float64  `json:"percent"`
	SwapInfo SwapInfo `json:"swapInfo"`
}

type SwapInfo struct {
	Label     string `json:"label"`
	AmmKey    string `json:"ammKey"`
	InAmount  string `json:"inAmount"`
	OutAmount string `json:"outAmount"`
	FeeAmount string `json:"feeAmount"`
	FeeMint   string `json:"feeMint"`
}
