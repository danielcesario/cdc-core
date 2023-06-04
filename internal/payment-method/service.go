package paymentmethod

import (
	"context"

	"github.com/danielcesario/cdc-core/internal/user"
	"github.com/google/uuid"
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

func (s *TransactionService) CreatePaymentMethod(ctx context.Context, request PaymentMethodRequest) (*PaymentMethodResponse, error) {
	userEmail := ctx.Value("user_email").(string)
	currentUser, err := s.userRepository.FindByEmail(userEmail)
	if err != nil {
		return nil, err
	}

	paymentMethod := request.toPaymentMethod()
	paymentMethod.Code = uuid.NewString()
	paymentMethod.UserID = currentUser.ID
	paymentMethod.User = *currentUser

	_, err = s.repository.Store(paymentMethod)
	if err != nil {
		return nil, err
	}

	return paymentMethod.toResponse(), nil
}
