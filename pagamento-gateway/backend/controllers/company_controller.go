package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"pagamento-gateway/config"
	"pagamento-gateway/models"
)

// Gera uma chave API aleatória
func generateAPIKey() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Cadastro da empresa (sem api_key ainda)
func RegisterCompany(c echo.Context) error {
	var company models.Company

	if err := c.Bind(&company); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Dados inválidos"})
	}

	if company.Email == "" || company.Password == "" || company.Name == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Todos os campos são obrigatórios"})
	}

	// Gera hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(company.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao gerar senha segura"})
	}
	company.Password = string(hashedPassword)
	company.IsActive = false
	company.ApiKey = "" // só será gerada na ativação

	// Salva no banco
	if err := config.DB.Create(&company).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao registrar empresa"})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Empresa registrada com sucesso. Ative seu plano para receber a API Key.",
	})
}

// Ativa a empresa após o pagamento (gera api_key aqui)
func ActivateCompany(c echo.Context) error {
	companyID, ok := c.Get("company_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Não autenticado"})
	}

	var company models.Company
	if err := config.DB.First(&company, companyID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Empresa não encontrada"})
	}

	if company.IsActive {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Plano já está ativo"})
	}

	company.IsActive = true
	company.ApiKey = generateAPIKey()

	if err := config.DB.Save(&company).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao ativar plano"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Plano ativado com sucesso",
		"api_key": company.ApiKey,
	})
}

// Ver informações da empresa logada
func VerPlano(c echo.Context) error {
	companyID, ok := c.Get("company_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Não autenticado"})
	}

	var company models.Company
	if err := config.DB.First(&company, companyID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Empresa não encontrada"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"id":        company.ID,
		"name":      company.Name,
		"email":     company.Email,
		"is_active": company.IsActive,
		"api_key":   company.ApiKey,
		"created":   company.CreatedAt,
		"balance":   company.Balance,
		"callback_url": company.CallbackURL,
		"bank_code":      company.BankCode,
		"bank_name":      company.BankName,
		"agency_account": company.AgencyAccount,
		"number_account": company.NumberAccount,
		"type_account":   company.TypeAccount,
		"type_key_pix":  company.TypeKeyPix,
		"key_pix":       company.KeyPix,
		
	})
}

func AtualizarCallbackURL(c echo.Context) error {
	companyID := c.Get("company_id").(uint)

	type RequestBody struct {
		CallbackURL string `json:"callback_url"`
	}

	var body RequestBody
	if err := c.Bind(&body); err != nil || body.CallbackURL == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "URL inválida"})
	}

	var company models.Company
	if err := config.DB.First(&company, companyID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Empresa não encontrada"})
	}

	company.CallbackURL = body.CallbackURL
	if err := config.DB.Save(&company).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao salvar callback URL"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Callback URL atualizada com sucesso"})
}

func AtualizarDadosBancarios(c echo.Context) error {
	companyID := c.Get("company_id").(uint)

	type DadosBancarios struct {
		BankCode      string `json:"bank_code"`
		BankName      string `json:"bank_name"`
		AgencyAccount string `json:"agency_account"`
		NumberAccount string `json:"number_account"`
		TypeAccount   string `json:"type_account"`   // CONTA_CORRENTE ou POUPANCA
		TypeKeyPix    string `json:"type_key_pix"`   // CPF, CNPJ, EMAIL, TELEFONE, CHAVE_ALEATORIA
		KeyPix        string `json:"key_pix"`
	}

	var input DadosBancarios
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Dados inválidos"})
	}

	var company models.Company
	if err := config.DB.First(&company, companyID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Empresa não encontrada"})
	}

	// Atualiza os dados
	company.BankCode = input.BankCode
	company.BankName = input.BankName
	company.AgencyAccount = input.AgencyAccount
	company.NumberAccount = input.NumberAccount
	company.TypeAccount = input.TypeAccount
	company.TypeKeyPix = input.TypeKeyPix
	company.KeyPix = input.KeyPix

	if err := config.DB.Save(&company).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao atualizar dados bancários"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Dados bancários atualizados com sucesso",
	})
}
