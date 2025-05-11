package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RequestLoggerConfig defines the config for the request logger middleware
type RequestLoggerConfig struct {
	// Skipper defines a function to skip middleware execution
	Skipper middleware.Skipper

	// LogRequestBody determines if the request body should be logged
	LogRequestBody bool

	// LogResponseBody determines if the response body should be logged
	LogResponseBody bool

	// LogLevel is the level at which to log
	LogLevel int
}

// DefaultRequestLoggerConfig is the default config for the request logger middleware
var DefaultRequestLoggerConfig = RequestLoggerConfig{
	Skipper:         middleware.DefaultSkipper,
	LogRequestBody:  false,
	LogResponseBody: false,
	LogLevel:        1,
}

// RequestLogger returns a middleware that logs HTTP requests
func RequestLogger(config RequestLoggerConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultRequestLoggerConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			start := time.Now()
			req := c.Request()
			res := c.Response()

			var reqBody []byte
			if config.LogRequestBody {
				if req.Body != nil {
					reqBody, _ = io.ReadAll(req.Body)
					req.Body = io.NopCloser(bytes.NewBuffer(reqBody))
				}
			}

			var resBodyBuffer *bytes.Buffer
			if config.LogResponseBody {
				resBodyBuffer = new(bytes.Buffer)
				res.Writer = &bodyDumpResponseWriter{
					ResponseWriter: res.Writer,
					body:           resBodyBuffer,
				}
			}

			err := next(c)

			stop := time.Now()
			latency := stop.Sub(start).Milliseconds()

			status := res.Status
			method := req.Method
			path := req.URL.Path
			if path == "" {
				path = "/"
			}

			byteIn := req.Header.Get(echo.HeaderContentLength)
			if byteIn == "" {
				byteIn = "0"
			}

			logLine := "REQUEST: " + method + " " + path +
				", STATUS: " + strconv.Itoa(status) +
				", LATENCY: " + strconv.FormatInt(latency, 10) + "ms" +
				", BYTES_IN: " + byteIn +
				", BYTES_OUT: " + strconv.FormatInt(res.Size, 10)

			// Log request body if enabled
			if config.LogRequestBody && len(reqBody) > 0 {
				logLine += "\nREQUEST_BODY: " + string(reqBody)
			}

			// Log response body if enabled
			if config.LogResponseBody && resBodyBuffer != nil {
				resBody := resBodyBuffer.String()
				if len(resBody) > 0 {
					// Trim large responses
					if len(resBody) > 1000 {
						resBody = resBody[:1000] + "... (truncated)"
					}
					logLine += "\nRESPONSE_BODY: " + resBody
				}
			}

			// Log based on status code
			if status >= 500 {
				c.Logger().Error(logLine)
			} else if status >= 400 {
				c.Logger().Warn(logLine)
			} else {
				c.Logger().Info(logLine)
			}

			return err
		}
	}
}

// bodyDumpResponseWriter implements http.ResponseWriter to capture the response body
type bodyDumpResponseWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	// Write to original writer
	n, err := w.ResponseWriter.Write(b)
	// Also write to buffer for logging
	if err == nil {
		w.body.Write(b)
	}
	return n, err
}
