package dto

import (
	"time"

	"gorm.io/datatypes"
)

type TokenDTO struct {
	Address           string         `json:"address"`
	CreatedAt         time.Time      `json:"created_at"`
	DailyVolume       float64        `json:"daily_volume"`
	Decimals          int            `json:"decimals"`
	Extensions        datatypes.JSON `json:"extensions"`
	FreezeAuthority   *string        `json:"freeze_authority"`
	LogoURI           string         `json:"logoURI"`
	MintAuthority     *string        `json:"mint_authority"`
	MintedAt          *time.Time     `json:"minted_at"`
	Name              string         `json:"name"`
	PermanentDelegate *string        `json:"permanent_delegate"`
	Symbol            string         `json:"symbol"`
	Tags              datatypes.JSON `json:"tags"`
}

type TokenResponse struct {
	Tokens []*TokenDTO `json:"tokens"`
	Total  int64       `json:"total"`
}

type TokenAccountsResponse struct {
	Address  string `json:"address"`
	LogoURI  string `json:"logoURI"`
	Amount   string `json:"amount"`
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Decimals int    `json:"decimals"`
}

type TokenAccountsApiResponse struct {
	JSONRPC string              `json:"jsonrpc"`
	Result  TokenAccountsResult `json:"result"`
	ID      string              `json:"id"`
}

type TokenAccountsSolanaNativeApiResponse struct {
	JSONRPC string                          `json:"jsonrpc"`
	Result  TokenAccountsSolanaNativeResult `json:"result"`
	ID      string                          `json:"id"`
}

type TokenAccountsSolanaNativeResult struct {
	Value float64 `json:"value"`
}

type TokenAccountsResult struct {
	Context TokenContext `json:"context"`
	Value   []TokenData  `json:"value"`
}

type TokenContext struct {
	APIVersion string `json:"apiVersion"`
	Slot       int64  `json:"slot"`
}

type TokenData struct {
	Pubkey  string       `json:"pubkey"`
	Account TokenAccount `json:"account"`
}

type TokenAccount struct {
	Data       TokenAccountData `json:"data"`
	Executable bool             `json:"executable"`
	Lamports   int64            `json:"lamports"`
	Owner      string           `json:"owner"`
	RentEpoch  uint64           `json:"rentEpoch"`
	Space      int              `json:"space"`
}

type TokenAccountData struct {
	Program string          `json:"program"`
	Parsed  ParsedTokenData `json:"parsed"`
	Space   int             `json:"space"`
}

type ParsedTokenData struct {
	Info TokenInfo `json:"info"`
	Type string    `json:"type"`
}

type TokenInfo struct {
	IsNative    bool              `json:"isNative"`
	Mint        string            `json:"mint"`
	Owner       string            `json:"owner"`
	State       string            `json:"state"`
	TokenAmount TokenAmountDetail `json:"tokenAmount"`
}

type TokenAmountDetail struct {
	Amount         string  `json:"amount"`
	Decimals       int     `json:"decimals"`
	UIAmount       float64 `json:"uiAmount"`
	UIAmountString string  `json:"uiAmountString"`
}
