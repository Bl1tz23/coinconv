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
	// v1 deprecated and no longer supports
	EndpointGetPriceConversion = "/v2/tools/price-conversion"

	MethodGetPriceConversion = http.MethodGet
)

type GetPriceConversionRequest struct {
	Amout     float64
	ID        *string
	Symbol    *string
	Time      *time.Time
	Convert   *string
	ConvertID *string
}

type GetPriceConversionResponse struct {
	Data   GetPriceConversionData `json:"data"`
	Status Status                 `json:"status"`
}

type GetPriceConversionData struct {
	Symbol      string           `json:"symbol"`
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Amount      float64          `json:"amount"`
	LastUpdated time.Time        `json:"last_updated"`
	Quote       map[string]Quote `json:"quote"`
}

type Quote struct {
	Price       float64
	LastUpdated time.Time
}

func (c *Coinmarketcap) GetPriceConversion(ctx context.Context, request *GetPriceConversionRequest) (*GetPriceConversionResponse, error) {
	req, err := xhttp.NewRequest(ctx, MethodGetPriceConversion, c.addr+EndpointGetPriceConversion+request.EncodeQuery(), nil, setCallHeaders(c.WithAuthenticationHeader()))
	if err != nil {
		return nil, fmt.Errorf("failed to create GetPriceConversion request: %w", err)
	}

	b, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make GetPriceConversion request: %w", err)
	}

	priceConversion := GetPriceConversionResponse{}
	err = json.Unmarshal(b, &priceConversion)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal GetPriceConversionResponse: %w", err)
	}

	return &priceConversion, nil
}

func (r *GetPriceConversionRequest) EncodeQuery() string {
	params := url.Values{}

	params.Add(ParamAmount, strconv.Itoa(int(r.Amout)))

	if r.Convert != nil {
		params.Add(ParamConvert, *r.Convert)
	}

	if r.ConvertID != nil {
		params.Add(ParamConvertID, *r.ConvertID)
	}

	if r.Symbol != nil {
		params.Add(ParamSymbol, *r.Symbol)
	}

	if r.ID != nil {
		params.Add(ParamID, *r.ID)
	}

	if r.Time != nil {
		params.Add(ParamTime, r.Time.Format(time.RFC3339))
	}

	return queryParams(params)
}
