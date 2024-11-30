package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"github.com/labstack/echo/v4"
)

type bodyLogWriter struct {
	io.Writer
	http.ResponseWriter
	status int
	size   int
}

func (w *bodyLogWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.size += len(b)
	w.Writer.Write(b) // Capture the response body
	return w.ResponseWriter.Write(b)
}

func DefaultStructuredLogger(cfg *config.Config, logger logging.Logger) echo.MiddlewareFunc {

	return structuredLogger(logger)
}

func structuredLogger(l logging.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()
			res := c.Response()

			// Read and save the request body
			bodyBytes, _ := io.ReadAll(req.Body)
			req.Body.Close()
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			// Capture the response body
			resBody := new(bytes.Buffer)
			blw := &bodyLogWriter{
				Writer:         resBody,
				ResponseWriter: res.Writer,
				status:         http.StatusOK, // Default status
			}
			res.Writer = blw

			// Process the request
			err := next(c)

			// Log details
			path := req.URL.Path
			raw := req.URL.RawQuery
			if raw != "" {
				path += "?" + raw
			}
			byteHeaders, _ := json.Marshal(req.Header)

			keys := map[logging.ExtraKey]interface{}{}
			keys[logging.Path] = path
			keys[logging.ClientIp] = c.RealIP()
			keys[logging.Method] = req.Method
			keys[logging.Latency] = time.Since(start)
			keys[logging.StatusCode] = blw.status
			keys[logging.ErrorMessage] = err
			keys[logging.BodySize] = blw.size
			keys[logging.RequestBody] = string(bodyBytes)
			keys[logging.ResponseBody] = resBody.String()
			keys[logging.RequestHeader] = string(byteHeaders)

			l.Info(logging.RequestResponse, logging.Api, "", keys)

			return err
		}
	}
}
