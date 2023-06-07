package main

import (
	"fmt"
	"os"

	"github.com/danielcesario/cdc-core/cmd/api/handler"
	"github.com/danielcesario/cdc-core/internal/category"
	"github.com/danielcesario/cdc-core/internal/database"
	paymentmethod "github.com/danielcesario/cdc-core/internal/payment-method"
	"github.com/danielcesario/cdc-core/internal/plan"
	"github.com/danielcesario/cdc-core/internal/transaction"
	"github.com/danielcesario/cdc-core/internal/user"
	"github.com/danielcesario/cdc-core/internal/wallet"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")

	dsn := DB_USER + ":" + DB_PASSWORD + "@tcp(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME + "?charset=utf8mb4&parseTime=True&loc=Local"
	fmt.Println(dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy:                           schema.NamingStrategy{SingularTable: true},
		DryRun:                                   false,
	})

	if err != nil {
		panic("failed to connect database")
	}

	// Update Entities
	database.UpdateSchema(db)

	// User Dependencies
	userRepo := user.NewMariaDBRepository(db)
	userService := user.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Plan Dependecies
	planRepo := plan.NewMariaDBRepository(db)
	planService := plan.NewPlanService(planRepo)
	planHandler := handler.NewPlanHandler(planService)

	// Wallet Dependecies
	walletRepo := wallet.NewMariaDBRepository(db)
	walletService := wallet.NewWalletService(walletRepo, userRepo)
	walletHandler := handler.NewWalletHandler(walletService)

	// Payment Methods Dependecies
	paymentMethodRepo := paymentmethod.NewMariaDBRepository(db)
	paymentMethodService := paymentmethod.NewPaymentMethodService(paymentMethodRepo, userRepo)
	paymentMethodHandler := handler.NewPaymentMethodHandler(paymentMethodService)

	// Category Dependecies
	categoryRepo := category.NewMariaDBRepository(db)
	categoryService := category.NewCategoryService(categoryRepo, userRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	// Transaction Dependecies
	transactionRepo := transaction.NewMariaDBRepository(db)
	transactionService := transaction.NewTransactionService(
		transactionRepo,
		userRepo,
		walletRepo,
		paymentMethodRepo,
		categoryRepo,
	)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// General Handler
	handler := handler.NewHandler(
		userHandler,
		planHandler,
		walletHandler,
		paymentMethodHandler,
		categoryHandler,
		transactionHandler,
	)

	router := handler.InitRouter()
	router.Run(":8080")
}
