package service

import (
	"blockchain-scrap/dto"
	"blockchain-scrap/pkg/errs"
	httprequest "blockchain-scrap/pkg/http-request"
	"blockchain-scrap/repository"
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/google/uuid"
)

type BlockchainService interface {
	GetBlockchainDetailByContractAddress(ctx context.Context, contractAddress string, timeSkip time.Duration) (*dto.ContractAddressResponse, errs.MessageErr)
	GetAllBlockchains() ([]map[string]interface{}, errs.MessageErr)
	GetBlockchainDetailByContractAddressAndID(id, contractAddress string, timeSkip time.Duration) (*dto.ContractAddressResponse, errs.MessageErr)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.BlockchainSearchResponse, errs.MessageErr)
	FindByID(ctx context.Context, ID uuid.UUID) (*dto.ContractAddressResponse, errs.MessageErr)
}

type blockchainService struct {
	searchRepo repository.BlockchainSearchRepository
	tokenRepo  repository.TokenRepository
}

func NewBlockchainService(searchRepo repository.BlockchainSearchRepository, tokenRepo repository.TokenRepository) BlockchainService {
	return &blockchainService{searchRepo: searchRepo, tokenRepo: tokenRepo}
}

// FindByID implements BlockchainService.
func (s *blockchainService) FindByID(ctx context.Context, ID uuid.UUID) (*dto.ContractAddressResponse, errs.MessageErr) {
	search, err := s.searchRepo.FindByID(ctx, ID)
	if err != nil {
		return nil, errs.NewNotFound("Data not found")
	}

	response := &dto.ContractAddressResponse{}
	if err := json.Unmarshal(search.ResponseData, response); err != nil {
		return nil, errs.NewInternalServerError("Failed to process data")
	}

	return response, nil
}

func (s *blockchainService) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.BlockchainSearchResponse, errs.MessageErr) {
	searches, err := s.searchRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, errs.NewNotFound("Data not found")
	}

	responses := make([]*dto.BlockchainSearchResponse, 0, len(searches))
	for _, search := range searches {
		responses = append(responses, &dto.BlockchainSearchResponse{
			ID:              search.ID,
			ContractAddress: search.ContractAddress,
		})
	}

	return responses, nil
}

func (s *blockchainService) GetBlockchainDetailByContractAddressAndID(id, contractAddress string, timeSkip time.Duration) (*dto.ContractAddressResponse, errs.MessageErr) {
	var (
		response      = &dto.ContractAddressResponse{}
		prices        = &dto.GetPricesRequest{}
		liquidity     = []dto.GetLiquidityRequest{}
		errChan       = make(chan error, 3)
		wg            sync.WaitGroup
		contractResp  []byte
		marketResp    []byte
		liquidityResp []byte
		aiResp        []byte
	)

	wg.Add(3)
	go func() {
		defer wg.Done()
		url := "https://api.coingecko.com/api/v3/coins/" + id + "/contract/" + contractAddress
		body, err := httprequest.ProcessJSONRequest("GET", url, nil, nil)
		if err != nil {
			errChan <- err
			return
		}
		contractResp = body
	}()

	go func() {
		defer wg.Done()
		url := "https://api.coingecko.com/api/v3/coins/" + id + "/market_chart?vs_currency=usd&days=1"
		body, err := httprequest.ProcessJSONRequest("GET", url, nil, nil)
		if err != nil {
			errChan <- err
			return
		}
		marketResp = body
	}()

	go func() {
		defer wg.Done()
		url := "https://api.dexscreener.com/tokens/v1/" + id + "/" + contractAddress
		body, err := httprequest.ProcessJSONRequest("GET", url, nil, nil)
		if err != nil {
			errChan <- err
			return
		}
		liquidityResp = body
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return nil, errs.NewInternalServerError("Failed to fetch data from API")
		}
	}

	response.Platform = contractAddress
	if err := json.Unmarshal(contractResp, response); err != nil {
		return nil, errs.NewInternalServerError("Failed to process contract data")
	}

	if err := json.Unmarshal(marketResp, prices); err != nil {
		return nil, errs.NewInternalServerError("Failed to process price data")
	}

	if err := json.Unmarshal(liquidityResp, &liquidity); err != nil {
		return nil, errs.NewInternalServerError("Failed to process liquidity data")
	}

	if len(liquidity) > 0 {
		response.MarketData.Liquidity.USD = liquidity[0].Liquidity.USD
	}

	// Filter time stamp in chart
	var filteredPrices []dto.PricePoint
	var lastTime time.Time

	for i, item := range prices.Prices {
		timestampMs := int64(item[0])
		price := item[1]
		t := time.UnixMilli(timestampMs)

		if i == 0 || t.Sub(lastTime) >= timeSkip {
			filteredPrices = append(filteredPrices, dto.PricePoint{
				Timestamp: t.Format(time.RFC3339),
				Price:     price,
			})
			lastTime = t
		}
	}
	response.TimePrices = filteredPrices

	if response.Symbol == "" {
		return nil, errs.NewNotFound("Contract address not found")
	}

	if len(response.TimePrices) == 0 {
		return nil, errs.NewNotFound("Token ID not found")
	}

	// Default seed data
	liquidityInfo := &dto.DexLiquidityInfo{
		LiquidityPoolSize: 67.8,
		TopDex:            "Jupiter",
		Volume24h:         130.0,
		SlippageNote:      "Low (<0.5%)",
		DexLiquidityRatio: 4.6,
		LiquidityTrend7D:  8.2,
	}

	tokenAnalytics := &dto.TokenAnalytics{
		TopHolder:   8.568,
		TopWallets:  8.568,
		TokenViewer: 8.568,
		SniperBot:   1.576,
		DevSold:     true,
		DevBuyback:  true,
	}
	response.LiquidityInfo = *liquidityInfo
	response.TokenAnalytics = *tokenAnalytics
	response.ListingDay = 575

	aiReq := &dto.AIRequest{
		Prompt:      "Analyze the following crypto data and summarize it in English. Also give suggestions on what the potential of this crypto is. All currency data is also in USD. Summarize only the important data. If there is any empty or meaningless data, just ignore it, it does not need to be in the output.",
		Collections: *response,
	}

	aiReqBody, err := json.Marshal(aiReq)
	if err != nil {
		return nil, errs.NewInternalServerError("Failed to process data for AI")
	}

	url := "https://casandra-bot.athenor.id/api/preset/completions"
	header := map[string]string{
		"ATHENOR-API-KEY": os.Getenv("API_KEY"),
	}
	body, _ := httprequest.ProcessJSONRequest("POST", url, aiReqBody, header)
	aiResp = body

	aiResponse := &dto.AIResponse{}
	if err = json.Unmarshal(aiResp, aiResponse); err != nil {
		return nil, errs.NewInternalServerError("Failed to process AI response")
	}

	response.SummaryAnalysis = aiResponse.AssistantMessage

	return response, nil
}

func (s *blockchainService) GetBlockchainDetailByContractAddress(ctx context.Context, contractAddress string, timeSkip time.Duration) (*dto.ContractAddressResponse, errs.MessageErr) {
	var (
		response      = &dto.ContractAddressResponse{}
		prices        = &dto.GetPricesRequest{}
		liquidity     = []dto.GetLiquidityRequest{}
		errChan       = make(chan error, 2)
		wg            sync.WaitGroup
		contractResp  []byte
		marketResp    []byte
		liquidityResp []byte
		aiResp        []byte
	)

	token, err := s.tokenRepo.FindByAddress([]string{contractAddress})

	if err != nil {
		return nil, err
	}

	if len(token) == 0 {
		return nil, errs.NewNotFound("Contract address not found, please change different contract address")
	}

	url := "https://api.coingecko.com/api/v3/coins/id/contract/" + contractAddress
	body, _ := httprequest.ProcessJSONRequest("GET", url, nil, nil)

	contractResp = body
	response.Platform = contractAddress
	if err := json.Unmarshal(contractResp, response); err != nil {
		return nil, errs.NewInternalServerError("Failed to process contract data")
	}
	if response.Symbol == "" {
		return nil, errs.NewNotFound("Contract address not found, please change different contract address")
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		url := "https://api.coingecko.com/api/v3/coins/" + response.ID + "/market_chart?vs_currency=usd&days=1"
		body, err := httprequest.ProcessJSONRequest("GET", url, nil, nil)
		if err != nil {
			errChan <- err
			return
		}
		marketResp = body
	}()

	go func() {
		defer wg.Done()
		url := "https://api.dexscreener.com/tokens/v1/" + response.ID + "/" + contractAddress
		body, err := httprequest.ProcessJSONRequest("GET", url, nil, nil)
		if err != nil {
			errChan <- err
			return
		}
		liquidityResp = body
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return nil, errs.NewInternalServerError("Failed to fetch data from API")
		}
	}

	if err := json.Unmarshal(marketResp, prices); err != nil {
		return nil, errs.NewInternalServerError("Failed to process price data")
	}

	if err := json.Unmarshal(liquidityResp, &liquidity); err != nil {
		return nil, errs.NewInternalServerError(err.Error())
	}

	if len(liquidity) > 0 {
		response.MarketData.Liquidity.USD = liquidity[0].Liquidity.USD
	}

	// Filter time stamp in chart
	var filteredPrices []dto.PricePoint
	var lastTime time.Time

	for i, item := range prices.Prices {
		timestampMs := int64(item[0])
		price := item[1]
		t := time.UnixMilli(timestampMs)

		if i == 0 || t.Sub(lastTime) >= timeSkip {
			filteredPrices = append(filteredPrices, dto.PricePoint{
				Timestamp: t.Format(time.RFC3339),
				Price:     price,
			})
			lastTime = t
		}
	}
	response.TimePrices = filteredPrices

	if len(response.TimePrices) == 0 {
		return nil, errs.NewNotFound("Contract address not found, please change different contract address")
	}

	// Default seed data
	liquidityInfo := &dto.DexLiquidityInfo{
		LiquidityPoolSize: 67.8,
		TopDex:            "Jupiter",
		Volume24h:         130.0,
		SlippageNote:      "Low (<0.5%)",
		DexLiquidityRatio: 4.6,
		LiquidityTrend7D:  8.2,
	}

	tokenAnalytics := &dto.TokenAnalytics{
		TopHolder:   8.568,
		TopWallets:  8.568,
		TokenViewer: 8.568,
		SniperBot:   1.576,
		DevSold:     true,
		DevBuyback:  true,
	}
	response.LiquidityInfo = *liquidityInfo
	response.TokenAnalytics = *tokenAnalytics
	response.ListingDay = 575

	aiReq := &dto.AIRequest{
		Prompt:      "Analyze the following crypto data and summarize it in English. Also provide suggestions on the potential of this crypto. All currency data is also in USD. Summarize only the important data. If there is empty or meaningless data, ignore it, because it does not need to be included in the output. Go straight to the core of the analysis, suggestions and recommendations, a maximum of 1 paragraph that already covers the most important.",
		Collections: *response,
	}

	aiReqBody, errMarshal := json.Marshal(aiReq)
	if errMarshal != nil {
		return nil, errs.NewInternalServerError("Failed to process data for AI")
	}

	url = "https://casandra-bot.athenor.id/api/preset/completions"
	header := map[string]string{
		"ATHENOR-API-KEY": os.Getenv("API_KEY"),
	}
	body, _ = httprequest.ProcessJSONRequest("POST", url, aiReqBody, header)
	aiResp = body

	aiResponse := &dto.AIResponse{}
	if errUnmarshal := json.Unmarshal(aiResp, aiResponse); errUnmarshal != nil {
		return nil, errs.NewInternalServerError("Failed to process AI response")
	}
	response.SummaryAnalysis = aiResponse.AssistantMessage

	// jsonData, errMarshal := json.Marshal(response)
	// if errMarshal != nil {
	// 	return nil, errs.NewInternalServerError("Failed to process data for storage")
	// }

	// searchRecord := &entity.BlockchainSearch{
	// 	ID:              uuid.New(),
	// 	UserID:          userID,
	// 	ContractAddress: response.Platform,
	// 	ResponseData:    jsonData,
	// }
	// err = s.searchRepo.SaveOrUpdate(ctx, searchRecord)
	// if err != nil {
	// 	return nil, errs.NewInternalServerError("Failed to save search data")
	// }

	return response, nil
}

func (s *blockchainService) GetAllBlockchains() ([]map[string]interface{}, errs.MessageErr) {
	url := "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=100&page=1"
	body, err := httprequest.ProcessJSONRequest("GET", url, nil, nil)
	if err != nil {
		return nil, errs.NewInternalServerError("Failed to fetch blockchain data")
	}
	var coins []map[string]interface{}
	if err := json.Unmarshal(body, &coins); err != nil {
		return nil, errs.NewInternalServerError("Failed to process blockchain data")
	}
	return coins, nil
}
