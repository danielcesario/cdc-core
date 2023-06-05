package paymentmethod

import (
	"context"
	"errors"

	"github.com/danielcesario/cdc-core/internal/user"
	"github.com/google/uuid"
)

type PaymentMethodService struct {
	repository     Repository
	userRepository user.Repository
}

func NewPaymentMethodService(repository Repository, userRepository user.Repository) *PaymentMethodService {
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

func (s *PaymentMethodService) Update(ctx context.Context, code string, request PaymentMethodRequest) (*PaymentMethodResponse, error) {
	savedPaymentMethod, err := s.repository.FindByCode(code)
	if err != nil {
		return nil, err
	}

	userEmail := ctx.Value("user_email").(string)
	currentUser, err := s.userRepository.FindByEmail(userEmail)
	if err != nil {
		return nil, err
	}

	if savedPaymentMethod.UserID != currentUser.ID {
		return nil, errors.New("invalid payment method owner")
	}

	newPaymentMethod := request.toPaymentMethod()
	savedPaymentMethod.PaymentType = newPaymentMethod.PaymentType
	savedPaymentMethod.Description = newPaymentMethod.Description
	savedPaymentMethod.CloseDay = newPaymentMethod.CloseDay
	savedPaymentMethod.DueDay = newPaymentMethod.DueDay

	_, err = s.repository.Store(savedPaymentMethod)
	if err != nil {
		return nil, err
	}

	return savedPaymentMethod.toResponse(), nil
}
