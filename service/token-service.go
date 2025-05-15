package service

import (
	"blockchain-scrap/dto"
	"blockchain-scrap/entity"
	"blockchain-scrap/pkg/errs"
	httprequest "blockchain-scrap/pkg/http-request"
	"blockchain-scrap/repository"
	"encoding/json"
	"strconv"

	"github.com/gagliardetto/solana-go"
)

type TokenService interface {
	GetAllTokens(limit, offset int, search string) (*dto.TokenResponse, errs.MessageErr)
	FetchAccountInfo(address string) ([]*dto.TokenAccountsResponse, errs.MessageErr)
}

type tokenService struct {
	repo repository.TokenRepository
}

func NewTokenService(r repository.TokenRepository) TokenService {
	return &tokenService{r}
}

func (s *tokenService) GetAllTokens(limit, offset int, search string) (*dto.TokenResponse, errs.MessageErr) {
	tokens, count, err := s.repo.GetAll(limit, offset, search)
	if err != nil {
		return nil, errs.NewInternalServerError(err.Error())
	}

	tokenDTOs := make([]*dto.TokenDTO, len(tokens))
	for i, token := range tokens {
		tokenDTOs[i] = &dto.TokenDTO{
			Address:           token.Address,
			CreatedAt:         token.CreatedAt,
			DailyVolume:       token.DailyVolume,
			Decimals:          token.Decimals,
			Extensions:        token.Extensions,
			FreezeAuthority:   token.FreezeAuthority,
			LogoURI:           token.LogoURI,
			MintAuthority:     token.MintAuthority,
			MintedAt:          token.MintedAt,
			Name:              token.Name,
			PermanentDelegate: token.PermanentDelegate,
			Symbol:            token.Symbol,
			Tags:              token.Tags,
		}
	}

	return &dto.TokenResponse{
		Tokens: tokenDTOs,
		Total:  count,
	}, nil
}

func (s *tokenService) FetchAccountInfo(address string) ([]*dto.TokenAccountsResponse, errs.MessageErr) {
	_, err := solana.PublicKeyFromBase58(address)
	if err != nil {
		return nil, errs.NewBadRequest("Invalid Solana address")
	}

	url := "https://api.mainnet-beta.solana.com"

	type resultTokenAccounts struct {
		data *dto.TokenAccountsApiResponse
		err  errs.MessageErr
	}
	type resultNativeSol struct {
		data *dto.TokenAccountsSolanaNativeApiResponse
		err  errs.MessageErr
	}

	tokenCh := make(chan resultTokenAccounts)
	nativeCh := make(chan resultNativeSol)

	// Goroutine untuk getTokenAccountsByOwner
	go func() {
		payload := map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      "1",
			"method":  "getTokenAccountsByOwner",
			"params": []interface{}{
				address,
				map[string]string{"programId": "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"},
				map[string]string{"encoding": "jsonParsed"},
			},
		}
		payloadBuffer, err := json.Marshal(payload)
		if err != nil {
			tokenCh <- resultTokenAccounts{nil, errs.NewInternalServerError(err.Error())}
			return
		}
		body, errRequest := httprequest.ProcessJSONRequest("POST", url, payloadBuffer, nil)
		if errRequest != nil {
			tokenCh <- resultTokenAccounts{nil, errRequest}
			return
		}
		var result dto.TokenAccountsApiResponse
		err = json.Unmarshal(body, &result)
		if err != nil {
			tokenCh <- resultTokenAccounts{nil, errs.NewInternalServerError(err.Error())}
			return
		}
		tokenCh <- resultTokenAccounts{&result, nil}
	}()

	// Goroutine untuk getBalance
	go func() {
		payload := map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      "1",
			"method":  "getBalance",
			"params":  []interface{}{address},
		}
		payloadBuffer, err := json.Marshal(payload)
		if err != nil {
			nativeCh <- resultNativeSol{nil, errs.NewInternalServerError(err.Error())}
			return
		}
		body, errRequest := httprequest.ProcessJSONRequest("POST", url, payloadBuffer, nil)
		if errRequest != nil {
			nativeCh <- resultNativeSol{nil, errRequest}
			return
		}
		var result dto.TokenAccountsSolanaNativeApiResponse
		err = json.Unmarshal(body, &result)
		if err != nil {
			nativeCh <- resultNativeSol{nil, errs.NewInternalServerError(err.Error())}
			return
		}
		nativeCh <- resultNativeSol{&result, nil}
	}()

	// Ambil hasil dari kedua goroutine
	var (
		tokenResult  resultTokenAccounts
		nativeResult resultNativeSol
	)
	for i := 0; i < 2; i++ {
		select {
		case tr := <-tokenCh:
			tokenResult = tr
		case nr := <-nativeCh:
			nativeResult = nr
		}
	}

	if tokenResult.err != nil {
		return nil, tokenResult.err
	}
	if nativeResult.err != nil {
		return nil, nativeResult.err
	}

	solanaNativeAddress := "So11111111111111111111111111111111111111112"
	tokens := []*dto.TokenAccountsResponse{}
	var mintAddresses []string

	tokens = append(tokens, &dto.TokenAccountsResponse{
		Address: solanaNativeAddress,
		Amount:  strconv.FormatFloat(nativeResult.data.Result.Value, 'f', -1, 64),
	})
	mintAddresses = append(mintAddresses, solanaNativeAddress)

	for _, acc := range tokenResult.data.Result.Value {
		mint := acc.Account.Data.Parsed.Info.Mint
		amount := acc.Account.Data.Parsed.Info.TokenAmount.Amount
		mintAddresses = append(mintAddresses, mint)
		tokens = append(tokens, &dto.TokenAccountsResponse{
			Address: mint,
			Amount:  amount,
		})
	}

	tokenEntities, err := s.repo.FindByAddress(mintAddresses)
	if err != nil {
		return nil, errs.NewInternalServerError(err.Error())
	}
	tokenMap := make(map[string]*entity.Token)
	for _, t := range tokenEntities {
		tokenMap[t.Address] = t
	}
	for _, token := range tokens {
		if ent, ok := tokenMap[token.Address]; ok {
			token.LogoURI = ent.LogoURI
			token.Name = ent.Name
			token.Symbol = ent.Symbol
			token.Decimals = ent.Decimals
		}
	}
	return tokens, nil
}
