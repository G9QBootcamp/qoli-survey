package middlewares

import (
	"fmt"
	"net/http"

	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"github.com/labstack/echo/v4"
)

func RecoveryErrors(logger logging.Logger) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {
			defer func() {
				if rec := recover(); rec != nil {
					if err, ok := rec.(error); ok {
						logger.Error(logging.General, logging.RecoverError, "recover error in request occurred", map[logging.ExtraKey]interface{}{logging.ErrorMessage: fmt.Sprintf("%v", err)})
						c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
					} else {
						logger.Error(logging.General, logging.RecoverError, "recover error in request occurred", map[logging.ExtraKey]interface{}{logging.ErrorMessage: fmt.Sprintf("%v", err)})
						c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": rec})
					}
				}
			}()

			return next(c)
		}
	}
}
