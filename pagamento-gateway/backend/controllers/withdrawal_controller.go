package controllers

import (
	"net/http"
	"time"
	"fmt"
	

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"pagamento-gateway/config"
	"pagamento-gateway/models"
)

func RealizarSaque(c echo.Context) error {
	type Input struct {
		Amount float64 `json:"amount"`
	}

	var input Input
	if err := c.Bind(&input); err != nil || input.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Valor inválido para saque"})
	}

	companyID, ok := c.Get("company_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Não autenticado"})
	}

	var company models.Company
	if err := config.DB.First(&company, companyID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Empresa não encontrada"})
	}

	if company.Balance < input.Amount {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Saldo insuficiente"})
	}

	// Atualiza saldo da empresa
	err := config.DB.Model(&company).
		Update("balance", gorm.Expr("balance - ?", input.Amount)).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao atualizar saldo"})
	}

	// Registra o saque
	withdrawal := models.Withdrawal{
		CompanyID: company.ID,
		Amount:    input.Amount,
		CreatedAt: time.Now(),
	}

	if err := config.DB.Create(&withdrawal).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao registrar saque"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Saque realizado com sucesso",
		"saque":   withdrawal,
	})
}

func ListarSaques(c echo.Context) error {
	companyID := c.Get("company_id").(uint)

	// Paginação
	page := 1
	limit := 10
	if p := c.QueryParam("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := c.QueryParam("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	var saques []models.Withdrawal
	var total int64

	// Conta total de saques da empresa
	config.DB.Model(&models.Withdrawal{}).
		Where("company_id = ?", companyID).
		Count(&total)

	// Busca com paginação
	err := config.DB.
		Where("company_id = ?", companyID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&saques).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao buscar saques"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":       saques,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	})
}

func VerSaldo(c echo.Context) error {
	companyID := c.Get("company_id").(uint)

	var company models.Company
	if err := config.DB.First(&company, companyID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Empresa não encontrada"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"saldo": company.Balance,
	})
}
