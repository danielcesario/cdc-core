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
				securedPaymentMethod.POST("/", h.CreatePaymentMethod)
				securedPaymentMethod.GET("/", h.ListPaymentMethod)
				securedPaymentMethod.PUT("/:paymentMethodCode", h.UpdatePaymentMethod)
			}

			// Category Area
			securedCategory := panel.Group("/categories").Use(middlewares.Auth())
			{
				securedCategory.POST("/", h.CreateCategory)
				securedCategory.GET("/", h.ListCategory)
				securedCategory.PUT("/:categoryCode", h.UpdateCategory)
			}

			// Transaction Area
			securedTransaction := panel.Group("/transactions").Use(middlewares.Auth())
			{
				securedTransaction.POST("/", h.CreateTransaction)
				securedTransaction.GET("/:transactionCode", h.GetTransaction)
				securedTransaction.POST("/search", h.SerachTransactions)
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
