package bot

import (
	"fmt"

	"github.com/jszwec/csvutil"
)

type StockDetails struct {
	Symbol     string `json:"symbol" csv:"Symbol"`
	ClosePrice string `json:"close_price" csv:"Close"`
}

func (s *StockDetails) String() string {
	return fmt.Sprintf("%s quote is $%s per share", s.Symbol, s.ClosePrice)
}

// parseStockData parses the stock data from the response body (a csv content).
func parseStockData(data []byte) (*StockDetails, error) {

	var details []StockDetails

	err := csvutil.Unmarshal(data, &details)
	if err != nil {
		return nil, err
	}

	if len(details) == 0 {
		return nil, fmt.Errorf("no data found")
	}

	return &details[0], nil
}
