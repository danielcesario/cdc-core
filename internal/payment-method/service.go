package paymentmethod

import (
	"context"

	"github.com/danielcesario/cdc-core/internal/user"
	"github.com/google/uuid"
)

type PaymentMethodService struct {
	repository     Repository
	userRepository user.Repository
}

func NewWalletService(repository Repository, userRepository user.Repository) *PaymentMethodService {
	return &PaymentMethodService{
		repository:     repository,
		userRepository: userRepository,
	}
}

func (s *PaymentMethodService) Create(ctx context.Context, request PaymentMethodRequest) (*PaymentMethodResponse, error) {
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

func (s *PaymentMethodService) List(ctx context.Context) (response []*PaymentMethodResponse, err error) {
	userEmail := ctx.Value("user_email").(string)
	currentUser, err := s.userRepository.FindByEmail(userEmail)
	if err != nil {
		return
	}

	result, err := s.repository.ListByUser(currentUser.Code)
	if err != nil {
		return
	}

	for _, paymentMethod := range result {
		response = append(response, paymentMethod.toResponse())
	}

	return
}
