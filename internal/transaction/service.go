package transaction

import (
	"github.com/danielcesario/cdc-core/internal/category"
	paymentmethod "github.com/danielcesario/cdc-core/internal/payment-method"
	"github.com/danielcesario/cdc-core/internal/user"
	"github.com/danielcesario/cdc-core/internal/wallet"
)

type TransactionService struct {
	repository              Repository
	userRepository          user.Repository
	walletRepository        wallet.Repository
	paymentmethodRepository paymentmethod.Repository
	categoryRepository      category.Repository
}

func NewTransactionService(repository Repository, userRepository user.Repository, walletRepository wallet.Repository,
	paymentmethodRepository paymentmethod.Repository, categoryRepository category.Repository) *TransactionService {
	return &TransactionService{
		repository:              repository,
		userRepository:          userRepository,
		walletRepository:        walletRepository,
		paymentmethodRepository: paymentmethodRepository,
		categoryRepository:      categoryRepository,
	}
}
