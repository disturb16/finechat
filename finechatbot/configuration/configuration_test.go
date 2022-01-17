package configuration

import (
	"testing"
)

func TestGet(t *testing.T) {
	testCases := []struct {
		Name          string
		FilePath      string
		ExpectedError error
	}{
		{
			Name:          "Should fail when file does not exist",
			FilePath:      "t.go",
			ExpectedError: ErrNoFile,
		},
		{
			Name:          "Valid file",
			FilePath:      "config_example.yml",
			ExpectedError: nil,
		},
		{
			Name:          "Should fail when file is not valid",
			FilePath:      "invalid_file.yml",
			ExpectedError: ErrParsingFile,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			_, err := Get(tc.FilePath)
			if err != tc.ExpectedError {
				t.Errorf("Expected error %v, got %v", tc.ExpectedError, err)
			}
		})
	}
}
