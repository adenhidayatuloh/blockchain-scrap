package service

import (
	"blockchain-scrap/dto"
	"blockchain-scrap/entity"
	"blockchain-scrap/pkg/errs"
	httprequest "blockchain-scrap/pkg/http-request"
	"blockchain-scrap/repository"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type SwapService interface {
	GetSwapTransaction(req dto.SwapRequest) (string, errs.MessageErr)
	SubmitTransaction(req dto.SubmitRequest) (string, errs.MessageErr)
	GetCurrencySwap(req dto.SwapRequest) (*dto.GetCurrencySwapResponse, errs.MessageErr)
}

type swapServiceImpl struct {
	tokenRepo    repository.TokenRepository
	tokenService TokenService
}

func NewSwapService(tokenRepo repository.TokenRepository, tokenService TokenService) SwapService {
	return &swapServiceImpl{tokenRepo: tokenRepo, tokenService: tokenService}
}

func (s *swapServiceImpl) GetSwapTransaction(req dto.SwapRequest) (string, errs.MessageErr) {

	var (
		decimalAmount int64
	)

	tokenMetadatas, errRepo := s.tokenRepo.FindByAddress([]string{req.InputMint, req.OutputMint})
	if errRepo != nil {
		return "", errRepo
	}

	tokenMap := make(map[string]*entity.Token)
	for _, t := range tokenMetadatas {
		tokenMap[t.Address] = t
	}

	inputTokenMeta, ok := tokenMap[req.InputMint]
	if !ok {
		return "", errs.NewBadRequest("token sell not supported")
	}

	_, ok = tokenMap[req.OutputMint]
	if !ok {
		return "", errs.NewBadRequest("token buy not supported")
	}

	if inputTokenMeta.Decimals > 0 {
		decimalAmount = int64(req.Amount * math.Pow10(int(inputTokenMeta.Decimals)))
	} else {
		decimalAmount = int64(req.Amount)
	}

	quoteURL := fmt.Sprintf(
		"https://quote-api.jup.ag/v6/quote?inputMint=%s&outputMint=%s&amount=%d&slippageBps=50",
		req.InputMint, req.OutputMint, decimalAmount,
	)

	body, errRequest := httprequest.ProcessJSONRequest("GET", quoteURL, nil, nil)
	if errRequest != nil {
		var jupError dto.JupiterErrorResponse
		if unmarshalErr := json.Unmarshal(body, &jupError); unmarshalErr == nil && jupError.Error != "" {
			return "", errs.NewBadRequest(fmt.Sprintf("Jupiter API error: %s (Code: %s)", jupError.Error, jupError.ErrorCode))
		}
		return "", errs.NewInternalServerError(fmt.Sprintf("Failed quote from Jupiter API: %s", errRequest.Message()))
	}

	var quoteResponse map[string]interface{}
	if err := json.Unmarshal(body, &quoteResponse); err != nil {
		return "", errs.NewInternalServerError(fmt.Sprintf("failed unmarshal respons quote Jupiter: %v", err))
	}

	swapPayload := map[string]interface{}{
		"quoteResponse": quoteResponse,
		"userPublicKey": req.PublicKey,
		"wrapUnwrapSOL": true,
	}

	swapPayloadBytes, err := json.Marshal(swapPayload)
	if err != nil {
		return "", errs.NewInternalServerError(fmt.Sprintf("failed marshal swap payload: %v", err))
	}

	swapBody, errRequest := httprequest.ProcessJSONRequest("POST", "https://quote-api.jup.ag/v6/swap", swapPayloadBytes, nil)
	if errRequest != nil {
		var jupError dto.JupiterErrorResponse
		if unmarshalErr := json.Unmarshal(swapBody, &jupError); unmarshalErr == nil && jupError.Error != "" {
			return "", errs.NewBadRequest(fmt.Sprintf("Jupiter API error: %s (Code: %s)", jupError.Error, jupError.ErrorCode))
		}
		return "", errs.NewInternalServerError(fmt.Sprintf("Failed swap from Jupiter API: %s", errRequest.Message()))
	}

	var swapResponse map[string]interface{}
	if err := json.Unmarshal(swapBody, &swapResponse); err != nil {
		return "", errs.NewInternalServerError(fmt.Sprintf("failed unmarshal respons swap Jupiter: %v", err))
	}

	transaction, ok := swapResponse["swapTransaction"].(string)
	if !ok {
		return "", errs.NewInternalServerError("invalid swapTransaction format")
	}

	return transaction, nil
}

func parseFloat(s string, fieldName string) (float64, errs.MessageErr) {
	if s == "" {
		return 0, errs.NewInternalServerError(fmt.Sprintf("field %s dari Jupiter API kosong", fieldName))
	}
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, errs.NewInternalServerError(fmt.Sprintf("gagal parse %s ('%s') dari Jupiter API ke float: %v", fieldName, s, err))
	}
	return val, nil
}

func (s *swapServiceImpl) GetCurrencySwap(req dto.SwapRequest) (*dto.GetCurrencySwapResponse, errs.MessageErr) {
	var (
		quoteResponse       dto.QuoteResponse
		getCurrencyResponse dto.GetCurrencySwapResponse
		decimalAmount       int64
	)

	tokenMetadatas, errRepo := s.tokenRepo.FindByAddress([]string{req.InputMint, req.OutputMint})
	if errRepo != nil {
		return nil, errRepo
	}

	fmt.Println(tokenMetadatas[0].Address)

	tokenMap := make(map[string]*entity.Token)
	for _, t := range tokenMetadatas {
		tokenMap[t.Address] = t
	}

	inputTokenMeta, ok := tokenMap[req.InputMint]
	if !ok {
		return nil, errs.NewBadRequest("Token sell not supported, please change.")
	}
	outputTokenMeta, ok := tokenMap[req.OutputMint]
	if !ok {
		return nil, errs.NewBadRequest("Token buy not supported, please change.")
	}

	if inputTokenMeta.Decimals > 0 {
		decimalAmount = int64(req.Amount * math.Pow10(int(inputTokenMeta.Decimals)))
	} else {
		decimalAmount = int64(req.Amount)
	}

	quoteURL := fmt.Sprintf(
		"https://quote-api.jup.ag/v6/quote?inputMint=%s&outputMint=%s&amount=%d&slippageBps=50",
		req.InputMint, req.OutputMint, decimalAmount,
	)

	body, errRequest := httprequest.ProcessJSONRequest("GET", quoteURL, nil, nil)
	if errRequest != nil {
		var jupError dto.JupiterErrorResponse
		if unmarshalErr := json.Unmarshal(body, &jupError); unmarshalErr == nil && jupError.Error != "" {
			return nil, errs.NewBadRequest(fmt.Sprintf("Jupiter API error: %s (Code: %s)", jupError.Error, jupError.ErrorCode))
		}
		return nil, errs.NewInternalServerError(fmt.Sprintf("Failed quote from Jupiter API: %s", errRequest.Message()))
	}

	if err := json.Unmarshal(body, &quoteResponse); err != nil {
		return nil, errs.NewInternalServerError(fmt.Sprintf("failed unmarshal respons quote Jupiter: %v", err))
	}

	getCurrencyResponse.InAmount = float64(req.Amount)

	outAmountSmallestUnit, errParseOut := parseFloat(quoteResponse.OutAmount, "OutAmount")
	if errParseOut != nil {
		getCurrencyResponse.OutAmount = 0
		getCurrencyResponse.IsSwappable = false
	} else {
		if outputTokenMeta.Decimals > 0 {
			getCurrencyResponse.OutAmount = outAmountSmallestUnit / math.Pow10(int(outputTokenMeta.Decimals))
		} else {
			getCurrencyResponse.OutAmount = outAmountSmallestUnit
		}
		getCurrencyResponse.IsSwappable = getCurrencyResponse.OutAmount > 0
	}

	swapUsdVal, errParseUsd := parseFloat(quoteResponse.SwapUsdValue, "SwapUsdValue")
	if errParseUsd != nil {
		getCurrencyResponse.SwapUsdValue = 0
	} else {
		getCurrencyResponse.SwapUsdValue = swapUsdVal
	}

	userTokenAccounts, errFetch := s.tokenService.FetchAccountInfo(req.PublicKey)
	if errFetch != nil {
		getCurrencyResponse.BalanceInAmount = 0
	} else {
		foundBalance := false
		for _, acc := range userTokenAccounts {
			if acc.Address == req.InputMint {
				balanceAmountSmallestUnit, errConvBal := strconv.ParseFloat(acc.Amount, 64)
				if errConvBal != nil {
					getCurrencyResponse.BalanceInAmount = 0
				} else {
					if inputTokenMeta.Decimals > 0 {
						getCurrencyResponse.BalanceInAmount = balanceAmountSmallestUnit / math.Pow10(int(inputTokenMeta.Decimals))
					} else {
						getCurrencyResponse.BalanceInAmount = balanceAmountSmallestUnit
					}
				}
				foundBalance = true
				break
			}
		}
		if !foundBalance {
			getCurrencyResponse.BalanceInAmount = 0
		}
	}

	if quoteResponse.OutAmount == "" || quoteResponse.OutAmount == "0" {
		getCurrencyResponse.IsSwappable = false
	}
	if getCurrencyResponse.BalanceInAmount == 0 {
		getCurrencyResponse.IsSwappable = false
	}

	if getCurrencyResponse.BalanceInAmount < getCurrencyResponse.InAmount {
		getCurrencyResponse.IsSwappable = false
	}

	return &getCurrencyResponse, nil
}

func (s *swapServiceImpl) SubmitTransaction(req dto.SubmitRequest) (string, errs.MessageErr) {
	client := rpc.New("https://api.mainnet-beta.solana.com")

	tx, err := solana.TransactionFromBase64(req.SignedTransaction)
	if err != nil {
		return "", errs.NewInternalServerError(fmt.Sprintf("failed to decode transaction: %v", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := client.SendTransactionWithOpts(
		ctx,
		tx,
		rpc.TransactionOpts{
			SkipPreflight:       false,
			PreflightCommitment: rpc.CommitmentFinalized,
		},
	)
	if err != nil {
		return "", errs.NewInternalServerError(fmt.Sprintf("failed to send transaction: %v", err))
	}

	return sig.String(), nil
}
