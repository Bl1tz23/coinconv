package coinmarketcap

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Error struct {
	Status Status `json:"status"`
}

type Status struct {
	Timestamp    time.Time `json:"timestamp"`
	ErrorCode    int       `json:"error_code"`
	ErrorMessage string    `json:"error_message"`
	Elapsed      int       `json:"elapsed"`
	CreditCount  int       `json:"credit_count"`
}

// TODO: maybe different error types for each code
type CommonError struct {
	Status Status `json:"status"`
}

func (e CommonError) Error() string {
	return fmt.Sprintf("request failed: %s", e.Status.ErrorMessage)
}

type Currency struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Sign   string `json:"sign"`
	Sybmol string `json:"symbol"`
}

func checkResponse(res *http.Response) error {
	if res.StatusCode <= 399 {
		return nil
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read error response body: %w", err)
	}
	defer res.Body.Close()

	status := Status{}
	err = json.Unmarshal(b, &status)
	if err != nil {
		return fmt.Errorf("failed to unmarhsal error response body, err: %w, body: %s", err, b)
	}

	return &CommonError{
		status,
	}
}

func String(s string) *string {
	return &s
}

func setCallHeaders(getHeadersFuncs ...func() (string, string)) map[string]string {
	headers := make(map[string]string, len(getHeadersFuncs))

	for _, ah := range getHeadersFuncs {
		key, value := ah()
		headers[key] = value
	}

	return headers
}

func queryParams(params url.Values) string {
	encodedParams := params.Encode()
	if encodedParams == "" {
		return ""
	}

	return "?" + encodedParams
}
