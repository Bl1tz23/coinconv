package coinmarketcap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	xhttp "interview/pkg/util/http"
)

const (
	FiatSortByID SortFiat = iota
	FiatSortByName

	MethodFiatMap = http.MethodGet

	EndpointMapFiat = "/v1/fiat/map"
)

type SortFiat uint

func (s SortFiat) String() string {
	return [...]string{"id", "name"}[s]
}

type GetFiatMapRequest struct {
	Start         *uint
	Limit         *uint
	Sort          *SortFiat
	IncludeMetals bool
}

type GetFiatMapResponse struct {
	Data   []GetFiatMapData `json:"data"`
	Status Status           `json:"status"`
}

type GetFiatMapData struct {
	Currency
}

func (d *GetFiatMapData) GetCurrency() Currency {
	return d.Currency
}

func (c *Coinmarketcap) GetFiatMap(ctx context.Context, request *GetFiatMapRequest) (*GetFiatMapResponse, error) {
	req, err := xhttp.NewRequest(ctx, MethodFiatMap, c.addr+EndpointMapFiat+"?"+request.encodeQuery(), nil, setCallHeaders(c.WithAuthenticationHeader()))
	if err != nil {
		return nil, fmt.Errorf("failed to create GetFiatMap request: %w", err)
	}

	b, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make GetFiatMap request: %w", err)
	}

	fiatMap := GetFiatMapResponse{}
	err = json.Unmarshal(b, &fiatMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal fiat map response: %w", err)
	}

	return &fiatMap, nil
}

func (r *GetFiatMapRequest) encodeQuery() string {
	params := url.Values{}

	// We do not set params and use API defaults if it's nil

	if r.Start != nil {
		params.Add(ParamStart, strconv.Itoa(int(*r.Start)))
	}

	if r.Limit != nil {
		params.Add(ParamLimit, strconv.Itoa(int(*r.Limit)))
	}

	if r.Sort != nil {
		params.Add(ParamSort, r.Sort.String())
	}

	if r.IncludeMetals {
		params.Add(ParamIncludeMetals, strconv.FormatBool(r.IncludeMetals))
	}

	return queryParams(params)
}
