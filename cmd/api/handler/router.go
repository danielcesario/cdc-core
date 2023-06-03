package handler

import (
	"github.com/danielcesario/cdc-core/cmd/api/middlewares"
	"github.com/gin-gonic/gin"
)

func (h *Handler) InitRouter() *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())

	api := router.Group("/api")
	{
		api.POST("/login", h.Login)
		api.POST("/register", h.Register)
		api.POST("/recover", h.FakeResponse)
		api.GET("/validate", h.FakeResponse)

		// Super Admin Area
		securedAdmin := api.Group("/admin").Use(middlewares.Auth()).Use(middlewares.IsSuperAdmin())
		{
			securedAdmin.POST("/plans", h.CreatePlan)
			securedAdmin.GET("/plans", h.ListPlans)
			securedAdmin.GET("/plans/:planCode", h.GetPlan)
			securedAdmin.PATCH("/plans/:planCode", h.UpdatePlan)
			securedAdmin.GET("/users", h.GetUsers)
		}

		// Registered User Area
		panel := api.Group("/panel")
		{
			// Wallet Area
			securedWallet := panel.Group("/wallets").Use(middlewares.Auth())
			{
				securedWallet.POST("/", h.CreateWallet)
				securedWallet.GET("/", h.ListWallets)
				securedWallet.PUT("/:walletCode/collaborator", h.AddCollaborator)
				securedWallet.GET("/:walletCode", h.GetWallet)
			}

			// Payment Method Area
			securedPaymentMethod := panel.Group("/payment-methods").Use(middlewares.Auth())
			{
				securedPaymentMethod.POST("/", h.FakeResponse)
				securedPaymentMethod.GET("/", h.FakeResponse)
				securedPaymentMethod.GET("/:paymentMethodCode", h.FakeResponse)
				securedPaymentMethod.PUT("/:paymentMethodCode", h.FakeResponse)
			}

			// Transaction Area
			securedTransaction := panel.Group("/transactions").Use(middlewares.Auth())
			{
				securedTransaction.POST("/", h.FakeResponse)
				securedTransaction.GET("/:transactionCode", h.FakeResponse)
				securedTransaction.POST("/search", h.FakeResponse)
			}

			// User Area
			securedUser := panel.Group("/user").Use(middlewares.Auth())
			{
				securedUser.PATCH("/profile", h.FakeResponse)
			}
		}

	}

	return router
}
