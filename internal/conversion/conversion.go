package conversion

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
)

type ConversionOption func(conversion *Conversion)

type Conversion struct {
	Amount decimal.Decimal
	From   Currency
	To     Currency
}

type Currency struct {
	Symbol string

	// Required for coinmarketcap
	IsFiat *bool
}

func New(amountStr string, from, to string, opts ...ConversionOption) (*Conversion, error) {
	amount, err := decimal.NewFromString(amountStr)
	if err != nil {
		return nil, fmt.Errorf("failed to convert amount to decimal: %w", err)
	}

	c := &Conversion{
		Amount: amount,
		From: Currency{
			Symbol: from,
		},
		To: Currency{
			Symbol: to,
		},
	}

	for _, o := range opts {
		o(c)
	}

	return c, nil
}

type conversioner interface {
	GetCurrencyCode(ctx context.Context, symbol string, isFiat *bool) (uint, error)
}

func (c *Currency) GetCode(ctx context.Context, conv conversioner) (uint, error) {
	return conv.GetCurrencyCode(ctx, c.Symbol, c.IsFiat)
}

func WithIsFiatOption(fromIsFiat, toIsFiat bool) ConversionOption {
	return func(conversion *Conversion) {
		conversion.From.IsFiat = &fromIsFiat
		conversion.To.IsFiat = &toIsFiat
	}
}
