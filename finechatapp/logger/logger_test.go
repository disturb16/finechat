package logger

import (
	"context"
	"testing"
)

func TestContextArgs(t *testing.T) {
	testCases := []struct {
		Name         string
		Args         map[string]interface{}
		ExpectedArgs []interface{}
	}{
		{
			Name: "Should return empty args when no request id or real ip",
			Args: map[string]interface{}{},
			ExpectedArgs: []interface{}{
				", x-request-id: -",
				", x-real-ip: -",
			},
		},
		{
			Name: "Should return request id",
			Args: map[string]interface{}{RequestIDKey: "request-id"},
			ExpectedArgs: []interface{}{
				", x-request-id: request-id",
				", x-real-ip: -",
			},
		},
		{
			Name: "Should return real ip",
			Args: map[string]interface{}{RealIPKey: "real-ip"},
			ExpectedArgs: []interface{}{
				", x-request-id: -",
				", x-real-ip: real-ip",
			},
		},
		{
			Name: "Should return request id and real ip",
			Args: map[string]interface{}{
				RequestIDKey: "request-id",
				RealIPKey:    "real-ip",
			},
			ExpectedArgs: []interface{}{
				", x-request-id: request-id",
				", x-real-ip: real-ip",
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			for k, v := range tc.Args {
				ctx = context.WithValue(ctx, k, v)
			}

			args := contextArgs(ctx)

			for i := range args {
				if args[i] != tc.ExpectedArgs[i] {
					t.Errorf("Expected %s, got %s", tc.ExpectedArgs[i], args[i])
				}
			}

		})
	}
}
