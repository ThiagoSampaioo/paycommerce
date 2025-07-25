package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Token ausente ou mal formatado"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse o token com verificação de método de assinatura
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de assinatura inválido")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Token inválido"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["company_id"] == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Token malformado"})
		}

		companyIDFloat, ok := claims["company_id"].(float64)
		if !ok {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "ID da empresa inválido no token"})
		}

		c.Set("company_id", uint(companyIDFloat))

		return next(c)
	}
}
