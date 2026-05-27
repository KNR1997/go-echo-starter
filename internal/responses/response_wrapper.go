package responses

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

type Error struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type Data struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ValidationError struct {
	Code   int                 `json:"code"`
	Error  string              `json:"error"`
	Fields map[string][]string `json:"fields,omitempty"`
}

func Response(c echo.Context, statusCode int, data interface{}) error {
	// nolint // context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	// nolint // context.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// nolint // context.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization")
	return c.JSON(statusCode, data)
}

func MessageResponse(c echo.Context, statusCode int, message string) error {
	return Response(c, statusCode, Data{
		Code:    statusCode,
		Message: message,
	})
}

func ErrorResponse(c echo.Context, statusCode int, message string) error {
	return Response(c, statusCode, Error{
		Code:  statusCode,
		Error: message,
	})
}

func ValidationErrorResponse(
	c echo.Context,
	statusCode int,
	message string,
	fields map[string][]string,
) error {
	return Response(c, statusCode, ValidationError{
		Code:   statusCode,
		Error:  message,
		Fields: fields,
	})
}

func ParseValidationErrors(err error) map[string][]string {
	result := map[string][]string{}

	if errs, ok := err.(validation.Errors); ok {
		for field, fieldErr := range errs {
			if fieldErr != nil {
				result[field] = []string{fieldErr.Error()}
			}
		}
	}

	return result
}
