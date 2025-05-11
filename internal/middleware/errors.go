package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// ErrorHandler is a middleware that handles all errors
func ErrorHandler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				// Already handled errors
				if c.Response().Committed {
					return err
				}

				// Get the error code and message
				code := http.StatusInternalServerError
				message := "Internal Server Error"

				// Check if it's an HTTPError
				if he, ok := err.(*echo.HTTPError); ok {
					code = he.Code
					if m, ok := he.Message.(string); ok {
						message = m
					} else {
						message = he.Error()
					}
				} else {
					// For regular errors, use the error message
					message = err.Error()
				}

				// Log the error
				c.Logger().Error(err)

				// Return a standardized error response
				return c.JSON(code, ErrorResponse{
					Message: message,
					Code:    code,
				})
			}
			return nil
		}
	}
}
