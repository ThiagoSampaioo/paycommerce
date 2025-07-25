package controllers

import (
	"bytes"
	"encoding/json"
	"e-commerce/config"
	"e-commerce/models"
	"net/http"
	"strconv"
	"os"
	"time"
	"log"
	"io"

	"github.com/labstack/echo/v4"
)



func PaymentCallback(c echo.Context) error {
	apiKey := c.Request().Header.Get("X-API-KEY")

	var provider models.PaymentProvider
	if err := config.DB.Where("active = ? AND api_key = ?", true, apiKey).First(&provider).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "API Key inválida ou inativa"})
	}

	var payload struct {
		UserID      uint    `json:"user_id"`
		Amount      float64 `json:"amount"`
		ReferenceID string  `json:"reference_id"`
		Status      string  `json:"status"` // "PAGO"
	}

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Dados inválidos"})
	}

	if payload.Status != "PAGO" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Pagamento não aprovado"})
	}

	var sale models.Sale
	if err := config.DB.Preload("Items").Where("reference_id = ?", payload.ReferenceID).First(&sale).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Venda não encontrada"})
	}

	if sale.Status == "PAGO" {
		return c.JSON(http.StatusConflict, echo.Map{"error": "Venda já está confirmada"})
	}

	now := time.Now()
	sale.Status = "PAGO"
	sale.PaidAt = &now

	if err := config.DB.Save(&sale).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao atualizar venda"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Venda confirmada com sucesso"})
}





func InitiatePayment(c echo.Context) error {
	var input struct {
		UserID      uint             `json:"user_id"`
		Total       float64          `json:"total"`
		ReferenceID string           `json:"reference_id"`
		Items       []models.SaleItem `json:"items"`
		ProviderApiKey string        `json:"provider_api_key"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Dados inválidos"})
	}

	var provider models.PaymentProvider
	if err := config.DB.Where("api_key = ? AND active = ?", input.ProviderApiKey, true).First(&provider).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Provedor inválido ou inativo"})
	}

	callbackURL := os.Getenv("BASE_URL") + "/callback-payment"

	sale := models.Sale{
		UserID:      input.UserID,
		Amount:      input.Total,
		Status:      "AGUARDANDO_PAGAMENTO",
		ApiKey:      provider.ApiKey,
		ReferenceID: input.ReferenceID,
		CallbackURL: callbackURL,
		Items:       input.Items,
	}

	if err := config.DB.Create(&sale).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao salvar venda"})
	}

	// Enviar para o gateway (opcional)
	payload := map[string]interface{}{
		"user_id":      input.UserID,
		"amount":       input.Total,
		"reference_id": input.ReferenceID,
		"callback_url": callbackURL,
	}
	payloadBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", provider.PaymentURL, bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", provider.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"error": "Erro ao contatar o gateway"})
	}
	defer resp.Body.Close()

	var respData map[string]string
	json.NewDecoder(resp.Body).Decode(&respData)

	return c.JSON(http.StatusOK, echo.Map{
		"payment_url":  respData["payment_url"],
		"reference_id": input.ReferenceID,
	})
}




func GetMySales(c echo.Context) error {
	userID := int(c.Get("user_id").(float64))

	// Paginação
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	status := c.QueryParam("status")

	var total int64
	var sales []models.Sale

	query := config.DB.Model(&models.Sale{}).Where("user_id = ?", userID)

	// Se o filtro de status foi informado
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Contagem total
	if err := query.Count(&total).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao contar vendas"})
	}

	// Buscar vendas paginadas
	if err := query.Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&sales).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao buscar vendas"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"total": total,
		"page":  page,
		"limit": limit,
		"data":  sales,
	})
}


func GetAllSales(c echo.Context) error {
	role := c.Get("role")
	if role != "admin" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Acesso negado"})
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	status := c.QueryParam("status")

	var total int64
	var sales []models.Sale

	query := config.DB.Model(&models.Sale{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao contar vendas"})
	}

	if err := query.Preload("User").
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&sales).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao buscar vendas"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"total": total,
		"page":  page,
		"limit": limit,
		"data":  sales,
	})
}


func CancelSale(c echo.Context) error {
	saleID := c.Param("id")
	var sale models.Sale

	// Buscar venda
	if err := config.DB.First(&sale, saleID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Venda não encontrada"})
	}

	if sale.Status != "AGUARDANDO_PAGAMENTO" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Venda não pode ser cancelada"})
	}

	// Buscar provedor
	var provider models.PaymentProvider
	if err := config.DB.Where("active = ?", true).First(&provider).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Nenhum provedor de pagamento ativo encontrado"})
	}

	// Montar payload
	payload := map[string]string{"reference_id": sale.ReferenceID}
	payloadBytes, _ := json.Marshal(payload)

	// Criar requisição POST
	client := &http.Client{}
	req, err := http.NewRequest("POST", provider.CancelationURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao criar requisição para o gateway"})
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", provider.ApiKey)

	// Enviar
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Erro ao enviar requisição:", err)
		return c.JSON(http.StatusBadGateway, echo.Map{"error": "Erro na comunicação com o gateway"})
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	// log do payload de envio
	
	log.Println("STATUS CODE GATEWAY:", resp.StatusCode)
	log.Println("RESPOSTA GATEWAY:", string(bodyBytes))

	if resp.StatusCode >= 400 {
		return c.JSON(http.StatusBadGateway, echo.Map{
			"error":  "Falha ao cancelar pagamento no gateway",
			"status": resp.StatusCode,
			"body":   string(bodyBytes),
		})
	}

	// Atualizar status local
	sale.Status = "CANCELADA"
	if err := config.DB.Save(&sale).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao atualizar venda local"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Venda cancelada com sucesso"})
}

