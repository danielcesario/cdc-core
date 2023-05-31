package wallet

import (
	"context"

	"github.com/danielcesario/cdc-core/internal/user"
	"github.com/google/uuid"
)

type WalletService struct {
	repository     Repository
	userRepository user.Repository
}

func NewWalletService(repository Repository, userRepository user.Repository) *WalletService {
	return &WalletService{
		repository:     repository,
		userRepository: userRepository,
	}
}

func (s *WalletService) Create(ctx context.Context, request WalletRequest) (*WalletResponse, error) {
	userEmail := ctx.Value("user_email").(string)
	currentUser, err := s.userRepository.FindByEmail(userEmail)
	if err != nil {
		return nil, err
	}

	var wallet Wallet = *request.toWallet()
	wallet.User = *currentUser
	wallet.UserID = currentUser.ID
	wallet.Code = uuid.NewString()
	wallet.Active = true

	_, err = s.repository.Store(&wallet)
	if err != nil {
		return nil, err
	}

	return wallet.ToResponse(), nil
}

func (s *WalletService) List(ctx context.Context) ([]*WalletResponse, error) {
	userEmail := ctx.Value("user_email").(string)
	currentUser, err := s.userRepository.FindByEmail(userEmail)
	if err != nil {
		return nil, err
	}

	wallets, err := s.repository.ListByUser(currentUser.Code)
	if err != nil {
		return nil, err
	}

	var response []*WalletResponse
	for _, wallet := range wallets {
		response = append(response, wallet.ToResponse())
	}

	return response, nil
}
