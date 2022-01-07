package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

type DTOError struct {
	Field string      `json:"field"`
	Tag   string      `json:"tag"`
	Value interface{} `json:"value"`
}

func (dtoErr DTOError) Error() string {
	return fmt.Sprintf(
		"validation error on field: %s with tag: %s, received: %v",
		dtoErr.Field,
		dtoErr.Tag,
		dtoErr.Value,
	)
}

// decodeParams decodes the query params into destination interface.
func decodeParams(c echo.Context, dst interface{}) error {
	binder := &echo.DefaultBinder{}
	return binder.BindQueryParams(c, dst)
}

// decodeBody decodes the request body into destination interface.
func decodeBody(c echo.Context, dst interface{}) error {
	data := c.Request().Body
	decoder := json.NewDecoder(data)
	return decoder.Decode(dst)
}

// parseDTOError parses the validator error and returns a HTTPError.
//
// Only the first error is returned.
func parseDTOError(err error) *echo.HTTPError {

	ve, ok := err.(validator.ValidationErrors)
	if !ok || len(ve) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	dtoErr := DTOError{
		Field: ve[0].Field(),
		Tag:   ve[0].ActualTag(),
		Value: ve[0].Value(),
	}

	return echo.NewHTTPError(http.StatusBadRequest, dtoErr.Error())
}
