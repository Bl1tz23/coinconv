package coinmarketcap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	xhttp "interview/pkg/util/http"
)

const (
	CryptocurrencyMapSortByID CryptocurrencyMapSort = iota
	CryptocurrencyMapSortByCMC
)

const (
	AuxPlatform Aux = iota
	AuxFirstHostoricalData
	AuxLastHistoricalData
	AuxIsActive
	AuxStatus
)

const (
	EndpointCryptocurrencyMap = "/v1/cryptocurrency/map"

	MethodGetCryptocurrencyMap = http.MethodGet
)

type CryptocurrencyMapSort uint

func (c CryptocurrencyMapSort) String() string {
	return [...]string{"id", "cmc_rank"}[c]
}

type GetCryptocurrencyMapRequest struct {
	ListingStatus *string
	Start         *uint
	Limit         *uint
	Sort          *CryptocurrencyMapSort
	Symbol        *string
	Aux           *Aux
}

type GetCryptocurrencyMapResponse struct {
	Data   []GetCryptocurrencyMapData `json:"data"`
	Status Status                     `json:"status"`
}

type GetCryptocurrencyMapData struct {
	Currency
	Slug                string    `json:"slug"`
	IsActive            uint      `json:"is_active"`
	FirstHistoricalData time.Time `json:"first_historical_data"`
	LastHistoricalData  time.Time `json:"last_historical_data"`
	Platform            *Platform `json:"platform,omitempty"`
}

func (d *GetCryptocurrencyMapData) GetCurrency() Currency {
	return d.Currency
}

type Platform struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	Slug         string `json:"slug"`
	TokenAddress string `json:"token_address"`
}

type Aux uint

func (a Aux) String() string {
	return [...]string{"platform", "first_historical_data", "last_historical_data", "is_active", "status"}[a]
}

func (c *Coinmarketcap) GetCryptocurrencyMap(ctx context.Context, request *GetCryptocurrencyMapRequest) (*GetCryptocurrencyMapResponse, error) {
	req, err := xhttp.NewRequest(ctx, MethodGetCryptocurrencyMap, c.addr+EndpointCryptocurrencyMap+request.encodeQuery(), nil, setCallHeaders(c.WithAuthenticationHeader()))
	if err != nil {
		return nil, fmt.Errorf("failed to create GetCryptocurrencyMap request: %w", err)
	}

	b, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make GetCryptocurrencyMap reqeust: %w", err)
	}

	cryptocurrencyMap := GetCryptocurrencyMapResponse{}
	err = json.Unmarshal(b, &cryptocurrencyMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal GetCryptocurrencyMap response: %w", err)
	}

	return &cryptocurrencyMap, nil
}

func (r *GetCryptocurrencyMapRequest) encodeQuery() string {
	params := url.Values{}
	if r.Limit != nil {
		params.Add(ParamLimit, strconv.Itoa(int(*r.Limit)))
	}

	if r.Start != nil {
		params.Add(ParamStart, strconv.Itoa(int(*r.Start)))
	}

	if r.Sort != nil {
		params.Add(ParamSort, r.Sort.String())
	}

	if r.Symbol != nil {
		params.Add(ParamSymbol, *r.Symbol)
	}

	if r.Aux != nil {
		params.Add(ParamAux, r.Aux.String())
	}

	if r.ListingStatus != nil {
		params.Add(ParamListingStatus, *r.ListingStatus)
	}

	return queryParams(params)
}
