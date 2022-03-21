package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Bl1tz23/coinconv/internal/conversion"
	"github.com/Bl1tz23/coinconv/pkg/source/coinmarketcap"
	xhttp "github.com/Bl1tz23/coinconv/pkg/util/http"

	"github.com/caarlos0/env"
	"github.com/urfave/cli/v2"
)

type Source uint

const (
	SourceCoinmarketcap Source = iota

	flagFromFiat   = "from-fiat"
	flagFromCrypto = "from-crypto"
	flagToFiat     = "to-fiat"
	flagToCrypto   = "to-crypto"
	flagAmount     = "amount"
)

var (
	fromSymbol string
	toSymbol   string
	amountStr  string
	config     *Config

	rootCmd = &cli.App{
		Name: "coinconv",
		Description: `Coinconv is a converter for currencies from specified source.
Required environment variables: 
1. COINCONV_SOURCE (available values: coinmarketcap, default: coinmarketcap)
2. COINCONV_API_KEY`,
		UsageText: "--amount AMOUNT [ --from-fiat | --from-crypto ] (required for coinmarketcap) SYMBOL [ --to-fiat | --to-crypto ] (required for coinmarketcap) SYMBOL",
		Flags: []cli.Flag{
			// TODO: create an issue
			// It's not working without flag, e.g. 100 BTC USD.
			// cli package doesn't parse flags if args[1] is not an app flag.
			&cli.StringFlag{
				Destination: &amountStr,
				Required:    true,
				Name:        flagAmount,
				Usage:       "amount of currency to convert",
			},
			&cli.StringFlag{
				Destination: &fromSymbol,
				Name:        flagFromFiat,
				Usage:       "indicates 'from' symbol which refers to fiat currency",
			},
			&cli.StringFlag{
				Destination: &fromSymbol,
				Name:        flagFromCrypto,
				Usage:       "indicates 'from' symbol which refers to cryptocurrency",
			},
			&cli.StringFlag{
				Destination: &toSymbol,
				Name:        flagToFiat,
				Usage:       "indicates 'to' symbol which refers to fiat currency",
			},
			&cli.StringFlag{
				Destination: &toSymbol,
				Name:        flagToCrypto,
				Usage:       "indicates 'to' symbol which refers to cryptocurrency",
			},
		},
		Before: func(c *cli.Context) error {

			config = new(Config)
			err := env.Parse(config)
			if err != nil {
				panic(err)
			}

			if config.Source != SourceCoinmarketcap.String() {
				return nil
			}

			if !(c.IsSet(flagFromFiat) || c.IsSet(flagFromCrypto)) {
				log.Fatal("--from-fiat or --from-crypto flag is required for coinmarketcap source")
			}

			if !(c.IsSet(flagToFiat) || c.IsSet(flagToCrypto)) {
				log.Fatal("--to-fiat or --to-crypto flag is required for coinmarketcap source")
			}

			if c.IsSet(flagFromFiat) && c.IsSet(flagFromCrypto) {
				log.Fatal("ambigious flags used to specify 'from' currency type")
			}

			if c.IsSet(flagToFiat) && c.IsSet(flagToCrypto) {
				log.Fatal("ambigious flags used to specify 'to' currency type")
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			source := coinmarketcap.New(xhttp.NewClient(0), config.APIKey, false)

			conversion, err := conversion.New(amountStr, fromSymbol, toSymbol,
				conversion.WithIsFiatOption(c.IsSet(flagFromFiat), c.IsSet(flagToFiat)))
			if err != nil {
				log.Fatal(err)
			}

			fromCode, err := conversion.From.GetCode(ctx, source)
			if err != nil {
				log.Fatalf("failed to get 'from' currency code: %s", err)
			}

			toCode, err := conversion.To.GetCode(ctx, source)
			if err != nil {
				log.Fatalf("failed to get 'to' currency code: %s", err)
			}

			result, err := source.ExchangeRate(ctx, conversion.Amount, fromCode, toCode)
			if err != nil {
				log.Fatalf("failed to convert: %s", err)
			}

			fmt.Println(result)

			return nil
		},
	}
)

func (s Source) String() string {
	return [...]string{"coinmarketcap"}[s]
}

type Config struct {
	APIKey string `env:"COINCONV_API_KEY"`
	Source string `env:"COINCONV_SOURCE" envDefault:"coinmarketcap"`
	// Sandbox bool   `env:"COINCONV_SANDBOX" envDefault:"true"`
}

func main() {
	rootCmd.Run(os.Args)
}
