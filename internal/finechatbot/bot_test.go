package finechatbot

import (
	"io"
	"os"
	"testing"
)

func TestGetSymbol(t *testing.T) {
	testCases := []struct {
		Name          string
		Value         string
		ExpectedError error
	}{
		{
			Name:          "Get symbol with valid input",
			Value:         "/stock=aapl.us",
			ExpectedError: nil,
		},
		{
			Name:          "Should fail with invalid input",
			Value:         "/some",
			ExpectedError: ErrInvalidSymbol,
		},
		{
			Name:          "Should fail without symbol",
			Value:         "/stock=",
			ExpectedError: ErrInvalidSymbol,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			_, err := GetSymbol(tc.Value)

			if err != tc.ExpectedError {
				t.Errorf("Expected error %v, got: %v", tc.ExpectedError, err)
			}
		})
	}
}

func TestParseStockData(t *testing.T) {

	f, err := os.Open("googl.us.csv")
	if err != nil {
		t.Fail()
	}

	testCases := []struct {
		Name          string
		Data          io.Reader
		ExpectedValue string
	}{
		{
			Name:          "Should parse data",
			Data:          f,
			ExpectedValue: "2740.34",
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			value, err := parseStockData(tc.Data)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if value != tc.ExpectedValue {
				t.Errorf("Expected value: %s, got %s", tc.ExpectedValue, value)
			}
		})
	}
}
