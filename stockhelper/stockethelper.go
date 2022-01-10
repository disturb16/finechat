package stockhelper

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const apiURL string = "https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv"

var ErrInvalidSymbol error = errors.New("invalid stock symbol")

func GetSymbol(val string) (string, error) {
	if !strings.HasPrefix(val, "/stock=") || len(val) == 7 {
		return "", ErrInvalidSymbol
	}

	symbol := val[7:]
	return symbol, nil
}

func GetShare(symbol string) (string, error) {
	url := fmt.Sprintf(apiURL, symbol)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	return parseStockData(resp.Body)
}

func parseStockData(data io.Reader) (string, error) {
	csvLines, err := csv.NewReader(data).ReadAll()
	if err != nil {
		return "", err
	}

	closingPrice := csvLines[1][6]

	return closingPrice, nil
}
