package routes

import (
	"github.com/labstack/echo/v4"
	"pagamento-gateway/controllers"
	"pagamento-gateway/middlewares"
)

func SetupRoutes(e *echo.Echo) {
	// Rotas públicas
	e.POST("/registrar", controllers.RegisterCompany)
	e.POST("/login", controllers.LoginCompany)
	

	// Rotas protegidas por API Key
	api := e.Group("/api", middlewares.ValidateAPIKey)
	api.POST("/checkout", controllers.Checkout)
	api.POST("/confirmacao/:id", controllers.ConfirmarPagamento)
	api.POST("/callback/cancelamento", controllers.CancelarPagamento)
	// Rotas protegidas por JWT
	auth := e.Group("/empresa", middlewares.AuthMiddleware)
	auth.GET("/pagamentos", controllers.ListPagamentos)
	auth.GET("/plano", controllers.VerPlano)
	auth.POST("/ativar", controllers.ActivateCompany)
	auth.PUT("/callback", controllers.AtualizarCallbackURL)
	auth.PUT("/bancario", controllers.AtualizarDadosBancarios)
	auth.POST("/saque", controllers.RealizarSaque)
	auth.GET("/saques", controllers.ListarSaques)
	auth.GET("/saldo", controllers.VerSaldo)
}
