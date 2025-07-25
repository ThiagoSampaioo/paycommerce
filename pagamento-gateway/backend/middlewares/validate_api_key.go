package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"pagamento-gateway/config"
	"pagamento-gateway/models"
)

func ValidateAPIKey(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		apiKey := c.Request().Header.Get("x-api-key")
		if apiKey == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": "API Key ausente ou mal formatada",
			})
		}

		var company models.Company
		if err := config.DB.Where("api_key = ?", apiKey).First(&company).Error; err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": "API Key inválida ou empresa não encontrada",
			})
		}

		if !company.IsActive {
			return c.JSON(http.StatusForbidden, echo.Map{
				"error": "Empresa inativa. Contrate um plano para ativar o uso da API.",
			})
		}

		// Passa o ID da empresa no contexto
		c.Set("company_id", company.ID)
		return next(c)
	}
}

