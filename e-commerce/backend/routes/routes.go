package routes

import (
	"e-commerce/controllers"
	"e-commerce/middlewares"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "API do E-commerce funcionando 🚀")
	})

	// Auth
	e.POST("/register", controllers.Register)
	e.POST("/login", controllers.Login)

	// Vendas via API Key (sem login)
	e.POST("/initiate-payment", controllers.InitiatePayment)
	e.POST("/callback-payment", controllers.PaymentCallback)

	// Rotas autenticadas
	api := e.Group("/api")
	api.Use(middlewares.JWTMiddleware)

	api.GET("/me", func(c echo.Context) error {
		userID := c.Get("user_id")
		role := c.Get("role")
		return c.JSON(200, echo.Map{
			"user_id": userID,
			"role":    role,
		})
	})

	api.GET("/sales/me", controllers.GetMySales)
	api.GET("/sales/all", controllers.GetAllSales)
	api.GET("/sales/cancel/:id", controllers.CancelSale)
	api.POST("/providers", controllers.CreateProvider)
	api.GET("/providers", controllers.GetAllProviders)
	api.GET("/providers/active", controllers.GetActiveProvider)
	api.PUT("/providers/:id", controllers.UpdateProvider)
	api.DELETE("/providers/:id", controllers.DeleteProvider)
	api.POST("/providers/:id/deactivate", controllers.DesativateProvider)

}
