package middlewares

import (
	"net/http"
	"strings"

	"github.com/G9QBootcamp/qoli-survey/pkg/jwtutils"
	"github.com/labstack/echo/v4"
)

func JWTAuth(secretKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid token")
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := jwtutils.ValidateToken(token, secretKey)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			c.Set("userID", claims.UserID)
			c.Set("role", claims.Role)
			return next(c)
		}
	}
}
