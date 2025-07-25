package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"e-commerce/models"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao conectar com o banco:", err)
	}

	if err := database.AutoMigrate(&models.User{}, &models.Sale{}, &models.PaymentProvider{}, &models.SaleItem{}); err != nil {
		log.Fatal("Erro ao migrar modelos:", err)
	}

	database.AutoMigrate(&models.User{}, &models.Sale{}, &models.PaymentProvider{}, &models.SaleItem{});

	DB = database
}
