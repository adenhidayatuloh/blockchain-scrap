package service

import (
	"blockchain-scrap/dto"
	httprequest "blockchain-scrap/pkg/http-request"
	"encoding/json"
	"sync"
	"time"
)

type CoinService interface {
	GetCoinDetail(id, contractAddress string, timeSkip time.Duration) (*dto.ContractAddressResponse, error)
	GetAllCoins() ([]map[string]interface{}, error)
}

type coinService struct{}

func NewCoinService() *coinService {
	return &coinService{}
}

func (s *coinService) GetCoinDetail(id, contractAddress string, timeSkip time.Duration) (*dto.ContractAddressResponse, error) {
	var (
		output          = &dto.ContractAddressResponse{}
		prices          = &dto.GetPricesRequest{}
		liquidity       = []dto.GetLiquidityRequest{}
		errs            = make(chan error, 3)
		wg              sync.WaitGroup
		contractBody    []byte
		marketChartBody []byte
		liquidityBody   []byte
	)

	wg.Add(3)
	go func() {
		defer wg.Done()
		url := "https://api.coingecko.com/api/v3/coins/" + id + "/contract/" + contractAddress
		body, err := httprequest.ProcessRequest(url)
		if err != nil {
			errs <- err
			return
		}
		contractBody = body
	}()

	go func() {
		defer wg.Done()
		url := "https://api.coingecko.com/api/v3/coins/" + id + "/market_chart?vs_currency=usd&days=1"
		body, err := httprequest.ProcessRequest(url)
		if err != nil {
			errs <- err
			return
		}
		marketChartBody = body
	}()

	go func() {
		defer wg.Done()
		url := "https://api.dexscreener.com/tokens/v1/" + id + "/" + contractAddress
		body, err := httprequest.ProcessRequest(url)
		if err != nil {
			errs <- err
			return
		}
		liquidityBody = body
	}()

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			return nil, err
		}
	}

	output.Platform = contractAddress
	if err := json.Unmarshal(contractBody, output); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(marketChartBody, prices); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(liquidityBody, &liquidity); err != nil {
		return nil, err
	}

	if len(liquidity) != 0 {
		output.MarketData.Liquidity.USD = liquidity[0].Liquidity.USD
	}

	var filteredOneHour []dto.PricePoint
	var lastTime time.Time

	for i, item := range prices.Prices {
		timestampMs := int64(item[0])
		price := item[1]
		t := time.UnixMilli(timestampMs)

		if i == 0 || t.Sub(lastTime) >= timeSkip {
			filteredOneHour = append(filteredOneHour, dto.PricePoint{
				Timestamp: t.Format(time.RFC3339),
				Price:     price,
			})
			lastTime = t
		}
	}

	output.TimePrices = filteredOneHour
	return output, nil
}

func (s *coinService) GetAllCoins() ([]map[string]interface{}, error) {
	url := "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=100&page=1"
	body, err := httprequest.ProcessRequest(url)
	if err != nil {
		return nil, err
	}
	var coins []map[string]interface{}
	if err := json.Unmarshal(body, &coins); err != nil {
		return nil, err
	}
	return coins, nil
}
