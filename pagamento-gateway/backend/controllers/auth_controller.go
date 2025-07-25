package controllers

import (
	"net/http"
	"os"
	"time"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"pagamento-gateway/config"
	"pagamento-gateway/models"
)

func LoginCompany(c echo.Context) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}



	var input LoginInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Dados inválidos"})
	}

	var company models.Company
	if err := config.DB.Where("email = ?", input.Email).First(&company).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Empresa não encontrada"})
	}

	fmt.Println("Login com:", input.Email, input.Password)
	fmt.Println("Hash salvo no banco:", company.Password)
	if err := bcrypt.CompareHashAndPassword([]byte(company.Password), []byte(input.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Senha inválida"})
	}

	claims := jwt.MapClaims{
		"company_id": company.ID,
		"exp":        time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenString, _ := token.SignedString([]byte(secret))

	return c.JSON(http.StatusOK, echo.Map{
		"token":   tokenString,
		"message": "Login realizado com sucesso",
	})
}

