package seed

import (
	"e-commerce/config"
	"e-commerce/models"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func SeedAdminUser() {
	var count int64
	if err := config.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&count).Error; err != nil {
		log.Println("Erro ao verificar administradores:", err)
		return
	}

	if count > 0 {
		log.Println("Administrador já existente. Nenhum usuário criado.")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("123456"), 12)
	if err != nil {
		log.Println("Erro ao gerar hash da senha:", err)
		return
	}

	admin := models.User{
		Name:     "Usuário Sistêmico",
		Email:    "admin@shop.com",
		Password: string(hash),
		Role:     "admin",
	}

	if err := config.DB.Create(&admin).Error; err != nil {
		log.Println("Erro ao criar usuário admin:", err)
		return
	}

	log.Println("Usuário administrador criado com sucesso.")
}
