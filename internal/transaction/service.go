package transaction

import (
	"github.com/danielcesario/cdc-core/internal/user"
)

type TransactionService struct {
	repository     Repository
	userRepository user.Repository
}

func NewWalletService(repository Repository, userRepository user.Repository) *TransactionService {
	return &TransactionService{
		repository:     repository,
		userRepository: userRepository,
	}
}
