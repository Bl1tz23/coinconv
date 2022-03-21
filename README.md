# COINCONV

## Description

### Binary usage:
```
CLI currency converter.

USAGE:
   --amount AMOUNT [ --from-fiat | --from-crypto ] (required for coinmarketcap) SYMBOL [ --to-fiat | --to-crypto ] (required for coinmarketcap) SYMBOL

DESCRIPTION:
   Coinconv is a converter for currencies from specified source.
   Required environment variables: 
   1. COINCONV_SOURCE (available values: coinmarketcap, default: coinmarketcap)
   2. COINCONV_API_KEY
   3. SANDBOX (default: true)

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --amount value       indicates that 'from' symbol refers to fiat currency
   --from-fiat value    indicates that 'from' symbol refers to fiat currency
   --from-crypto value  indicates that 'from' symbol refers to cryptocurrency
   --to-fiat value      indicates that 'to' symbol refers to fiat currency
   --to-crypto value    indicates that 'to' symbol refers to cryptocurrency
   --help, -h           show help (default: false)
```

### Docker

Example: `make run-docker command="--amount 1 --from-crypto BTC --to-fiat USD sandbox=false api_key=YOUR_API_KEY`