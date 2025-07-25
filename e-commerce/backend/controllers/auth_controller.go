package controllers

import (
	"e-commerce/config"
	"e-commerce/models"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Register(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Dados inválidos"})
	}

	// Default para cliente
	if user.Role == "" {
		user.Role = "cliente"
	}

	// Hash da senha
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao criptografar a senha"})
	}
	user.Password = string(hash)

	// Salvar no banco
	if err := config.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao salvar usuário"})
	}

	user.Password = ""
	return c.JSON(http.StatusCreated, user)
}

func Login(c echo.Context) error {
	var req models.User
	var user models.User

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Dados inválidos"})
	}

	// Buscar usuário pelo e-mail
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Usuário ou senha incorretos"})
	}

	// Verificar senha
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Usuário ou senha incorretos"})
	}

	// Gerar token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"name":    user.Name,
		"email":   user.Email,
	})

	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Erro ao gerar token"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": tokenStr,
		"user": echo.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}
