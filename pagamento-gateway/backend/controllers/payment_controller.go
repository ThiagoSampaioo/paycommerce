package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	
	"time"

	"pagamento-gateway/config"
	"pagamento-gateway/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)


// Dados recebidos do e-commerce no checkout
type CheckoutInput struct {
	Amount      float64 `json:"amount"`
	ReferenceID string  `json:"reference_id"`
}

// Criação de um novo pagamento (checkout)
func Checkout(c echo.Context) error {
	var input CheckoutInput
	if err := c.Bind(&input); err != nil || input.Amount <= 0 || input.ReferenceID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Dados inválidos"})
	}

	companyID := c.Get("company_id").(uint)

	// Busca empresa pra pegar o CallbackURL
	var company models.Company
	if err := config.DB.First(&company, companyID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Empresa não encontrada"})
	}

	payment := models.Payment{
		CompanyID:   companyID,
		Amount:      input.Amount,
		Status:      "AGUARDANDO_PAGAMENTO",
		CallbackURL: company.CallbackURL,
		ReferenceID: input.ReferenceID,
		CreatedAt:   time.Now(),
		ApiKey:      company.ApiKey,
	}

	if err := config.DB.Create(&payment).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao criar pagamento"})
	}

	paymentURL := fmt.Sprintf("http://localhost:5173/confirmacao/%d", payment.ID)

	return c.JSON(http.StatusOK, echo.Map{
		"message":      "Checkout criado com sucesso",
		"payment_id":   payment.ID,
		"payment_url":  paymentURL,
		"reference_id": payment.ReferenceID,
	})
}

// Confirma o pagamento e envia callback
func ConfirmarPagamento(c echo.Context) error {
	paymentID := c.Param("id")
	var payment models.Payment

	if err := config.DB.First(&payment, paymentID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Pagamento não encontrado"})
	}

	companyID := c.Get("company_id").(uint)
	if payment.CompanyID != companyID {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Acesso não autorizado"})
	}

	if payment.Status == "PAGO" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Pagamento já confirmado"})
	}

	// Atualiza status
	payment.Status = "PAGO"
	if err := config.DB.Save(&payment).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao atualizar pagamento"})
	}

	// Atualiza saldo da empresa
	config.DB.Model(&models.Company{}).
		Where("id = ?", payment.CompanyID).
		Update("balance", gorm.Expr("balance + ?", payment.Amount))

	// Envia callback de forma assíncrona
	go func(p models.Payment) {
		body := map[string]interface{}{
			"user_id":      p.CompanyID,
			"amount":       p.Amount,
			"reference_id": p.ReferenceID,
			"status":       "PAGO",
		}
		jsonBody, _ := json.Marshal(body)

		req, err := http.NewRequest("POST", p.CallbackURL, bytes.NewBuffer(jsonBody))
		if err != nil {
			log.Println("Erro ao criar request de callback:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-KEY", p.ApiKey)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Erro ao enviar callback:", err)
			return
		}
		defer resp.Body.Close()

		log.Println("Callback enviado com status:", resp.StatusCode)
	}(payment)

	return c.JSON(http.StatusOK, echo.Map{"message": "Pagamento confirmado e callback enviado"})
}

// Cancela um pagamento recebendo o reference_id no corpo (JSON)
func CancelarPagamento(c echo.Context) error {
	type CancelRequest struct {
		ReferenceID string `json:"reference_id"`
	}

	var req CancelRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Payload inválido"})
	}

	var payment models.Payment
	if err := config.DB.Where("reference_id = ?", req.ReferenceID).First(&payment).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Pagamento não encontrado"})
	}

	
	if payment.Status == "CANCELADO" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Pagamento já cancelado"})
	}

	payment.Status = "CANCELADO"
	if err := config.DB.Save(&payment).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao cancelar pagamento"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Pagamento cancelado com sucesso"})
}


func ListPagamentos(c echo.Context) error {
	companyID := c.Get("company_id").(uint)

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

	status := c.QueryParam("status") // 👈 Pega o status da query string

	var pagamentos []models.Payment
	var total int64

	query := config.DB.Model(&models.Payment{}).Where("company_id = ?", companyID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Conta total filtrado
	if err := query.Count(&total).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao contar pagamentos"})
	}

	// Busca os registros paginados e ordenados
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&pagamentos).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao buscar pagamentos"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":       pagamentos,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	})
}

