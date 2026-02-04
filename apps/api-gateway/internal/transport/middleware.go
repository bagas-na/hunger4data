package http

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JWTClaims struct {
	jwt.RegisteredClaims
}

func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"message": "missing Authorization header",
				})
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"message": "invalid Authorization header format",
				})
			}

			tokenStr := parts[1]

			token, err := jwt.ParseWithClaims(
				tokenStr,
				&JWTClaims{},
				func(token *jwt.Token) (interface{}, error) { return []byte(secret), nil },
				jwt.WithIssuedAt(),
				jwt.WithExpirationRequired(),
				jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
			)
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"message": "invalid or expired token",
				})
			}

			claims, ok := token.Claims.(*JWTClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"message": "invalid token claims",
				})
			}

			c.Set("user_id", claims.Subject)
			return next(c)
		}
	}
}
