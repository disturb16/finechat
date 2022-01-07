package api

import (
	"context"
	"fmt"

	"github.com/disturb16/finechat/logger"
	"github.com/labstack/echo/v4"
)

// RequestLogger logs requests
func RequestLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		r := c.Request()

		// Execute the handler
		err := next(c)

		errCode := c.Response().Status
		httpErr, ok := err.(*echo.HTTPError)
		if ok {
			errCode = httpErr.Code
		}

		message := fmt.Sprintf("method: %s, status: %d, url: %s", r.Method, errCode, r.URL)
		logger.Println(c.Request().Context(), message)

		return err
	}
}

// EnrichContext enriches the context and sets context to request
func EnrichContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		r := c.Request()
		ctx := r.Context()

		// Set request ID in context
		id := r.Header.Get(logger.RequestIDKey)
		if id != "" {
			ctx = context.WithValue(ctx, logger.RequestIDKey, id)
		}

		// Set real IP in context
		addr := r.Header.Get(logger.RealIPKey)
		if addr == "" && r.RemoteAddr != "" {
			addr = r.RemoteAddr
		}

		if addr != "" {
			ctx = context.WithValue(ctx, logger.RealIPKey, addr)
		}

		// Set context to request
		c.SetRequest(r.WithContext(ctx))

		return next(c)
	}
}
