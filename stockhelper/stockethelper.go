package stockhelper

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

const apiURL string = "https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv"

func GetSymbol(val string) (string, error) {
	rgx, err := regexp.Compile("(/)([A-z]*)=([A-z]*).([A-z]*)")
	if err != nil {
		return "", err
	}

	parts := rgx.FindStringSubmatch(val)
	if len(parts) < 4 {
		return "", fmt.Errorf("Invalid stock symbol")
	}

	symbol := fmt.Sprintf("%s.%s", parts[3], parts[4])

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
