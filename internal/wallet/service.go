package wallet

import (
	"context"
	"errors"

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

func (s *WalletService) AddCollaborator(ctx context.Context, walletCode string, request WalletCollaboratorRequest) error {
	user, err := s.userRepository.FindByCode(request.UserCode)
	if err != nil {
		return err
	}

	userEmail := ctx.Value("user_email").(string)
	currentUser, err := s.userRepository.FindByEmail(userEmail)
	if err != nil {
		return err
	}

	wallet, err := s.repository.FindByCode(walletCode)
	if err != nil {
		return err
	}

	if wallet.UserID != currentUser.ID {
		return errors.New("invalid wallet owner")
	}

	wallet.Collaborators = append(wallet.Collaborators, *user)

	_, err = s.repository.Store(wallet)
	if err != nil {
		return err
	}

	return nil
}

func (s *WalletService) GetByCode(ctx context.Context, code string) (*WalletResponse, error) {
	wallet, err := s.repository.FindByCode(code)
	if err != nil {
		return nil, err
	}

	userEmail := ctx.Value("user_email").(string)
	currentUser, err := s.userRepository.FindByEmail(userEmail)
	if err != nil {
		return nil, err
	}

	if wallet.UserID != currentUser.ID {
		return nil, errors.New("invalid wallet owner")
	}

	return wallet.ToResponse(), nil
}
