package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"pagamento-gateway/config"
	"pagamento-gateway/routes"
)

func main() {
	// Carrega variáveis do .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro ao carregar .env")
	}

	// Conecta ao banco de dados
	config.ConnectDB()

	// Cria a instância do Echo
	e := echo.New()

	// Habilita CORS para o frontend em localhost:5173
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	AllowOrigins: []string{"http://localhost:5173", "http://localhost:8084"},
AllowHeaders: []string{
	"Authorization",
	"Content-Type",
	"x-api-key",
	"X-API-Key",
},
	AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
}))


	// Registra as rotas
	routes.SetupRoutes(e)

	// Define a porta
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Inicia o servidor
	e.Logger.Fatal(e.Start(":" + port))
}
