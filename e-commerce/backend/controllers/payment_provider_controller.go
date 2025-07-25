package controllers

import (
	"e-commerce/config"
	"e-commerce/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateProvider(c echo.Context) error {
	role := c.Get("role")

	if role != "admin" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Acesso negado"})
	}

	var provider models.PaymentProvider
	
	if err := c.Bind(&provider); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Dados inválidos"})
	}

	if err := config.DB.Create(&provider).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao salvar provedor"})
	}

	return c.JSON(http.StatusCreated, provider)
}

func GetActiveProvider(c echo.Context) error {
	var provider models.PaymentProvider
	if err := config.DB.Where("active = ?", true).First(&provider).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Nenhum provedor ativo"})
	}
	return c.JSON(http.StatusOK, provider)
}

func GetAllProviders(c echo.Context) error {
	role := c.Get("role")
	if role != "admin" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Acesso negado"})
	}

	var providers []models.PaymentProvider
	if err := config.DB.Find(&providers).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao buscar provedores"})
	}

	return c.JSON(http.StatusOK, providers)
}
func UpdateProvider(c echo.Context) error {
	role := c.Get("role")
	if role != "admin" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Acesso negado"})
	}

	id := c.Param("id")
	var existing models.PaymentProvider
	if err := config.DB.First(&existing, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Provedor não encontrado"})
	}

	var updateData models.PaymentProvider
	if err := c.Bind(&updateData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Dados inválidos"})
	}

	existing.Name = updateData.Name
	existing.ApiKey = updateData.ApiKey
	existing.PaymentURL = updateData.PaymentURL
	existing.CancelationURL = updateData.CancelationURL
	existing.Active = updateData.Active

	if err := config.DB.Save(&existing).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao atualizar provedor"})
	}

	return c.JSON(http.StatusOK, existing)
}


func DeleteProvider(c echo.Context) error {
	role := c.Get("role")
	if role != "admin" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Acesso negado"})
	}

	id := c.Param("id")
	var provider models.PaymentProvider
	if err := config.DB.First(&provider, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Provedor não encontrado"})
	}

	if err := config.DB.Delete(&provider).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao deletar provedor"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Provedor deletado com sucesso"})
}

func DesativateProvider(c echo.Context) error {
	role := c.Get("role")
	if role != "admin" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Acesso negado"})
	}

	id := c.Param("id")
	var provider models.PaymentProvider
	if err := config.DB.First(&provider, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Provedor não encontrado"})
	}

	provider.Active = false
	if err := config.DB.Save(&provider).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao desativar provedor"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Provedor desativado com sucesso"})
}