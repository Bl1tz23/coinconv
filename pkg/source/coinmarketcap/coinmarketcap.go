package coinmarketcap

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"
)

const (
	ProAPI = "https://pro-api.coinmarketcap.com"

	// SandboxAPI is not supported due to false documentation
	// Response JSONs are different for some mysterical reason.
	SandboxAPI = "https://sandbox-api.coinmarketcap.com"

	// Header is the only way to authenticate.
	// It is not safe to pass token with query params.
	HeaderAuthentication = "X-CMC_PRO_API_KEY"

	// Query params
	ParamStart         = "start"
	ParamLimit         = "limit"
	ParamSort          = "sort"
	ParamSymbol        = "symbol"
	ParamListingStatus = "listing_status"
	ParamAux           = "aux"
	ParamIncludeMetals = "include_metals"
	ParamAmount        = "amount"
	ParamID            = "id"
	ParamTime          = "time"
	ParamConvert       = "convert"
	ParamConvertID     = "convert_id"
)

// TODO: separate coinmarketcap methods and auxillary funcs.
// Make another type of it e.g. "session" or "client"
// And embedd it to Coinmarketcap struct
type Coinmarketcap struct {
	client *http.Client

	addr   string
	secret string
}

func New(httpClient *http.Client, secret string, sandbox bool) *Coinmarketcap {
	c := &Coinmarketcap{
		client: httpClient,
		secret: secret,
	}

	// if sandbox {
	// 	c.addr = SandboxAPI
	// } else {
	c.addr = ProAPI
	// }

	return c
}

// Symbol is an alias for the currency.
func (c *Coinmarketcap) GetCurrencyCode(ctx context.Context, symbol string, isFiat *bool) (uint, error) {
	if isFiat == nil {
		return 0, errors.New("fiat currency option is required to convert currency with coinmarketcap")
	}

	if *isFiat {
		res, err := c.GetFiatMap(ctx, &GetFiatMapRequest{})
		if err != nil {
			return 0, err
		}

		for i := range res.Data {
			if res.Data[i].Sybmol == symbol {
				return res.Data[i].ID, nil
			}
		}

		return 0, errors.New("currency code not found")
	}

	res, err := c.GetCryptocurrencyMap(ctx, &GetCryptocurrencyMapRequest{
		Symbol: String(symbol),
	})
	if err != nil {
		return 0, err
	}

	for i := range res.Data {
		if symbol == res.Data[i].Sybmol {
			return res.Data[i].ID, nil
		}
	}

	return 0, errors.New("currency code not found")
}

func (c *Coinmarketcap) ExchangeRate(ctx context.Context, amount decimal.Decimal, from, to uint) (string, error) {
	res, err := c.GetPriceConversion(ctx, &GetPriceConversionRequest{
		Amout:     amount.Abs().InexactFloat64(),
		ID:        String(strconv.Itoa(int(from))),
		ConvertID: String(strconv.Itoa(int(to))),
	})
	if err != nil {
		return "", fmt.Errorf("failed to convert: %w", err)
	}

	if quote, ok := res.Data.Quote[strconv.Itoa(int(to))]; ok {
		return decimal.NewFromFloat(quote.Price).String(), nil
	}

	return "", errors.New("api returned unknown data")
}

func (c *Coinmarketcap) WithAuthenticationHeader() func() (string, string) {
	return func() (string, string) {
		return HeaderAuthentication, c.secret
	}
}

func (c *Coinmarketcap) do(req *http.Request) ([]byte, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	err = checkResponse(res)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	defer res.Body.Close()

	return b, nil
}
