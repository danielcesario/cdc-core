package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/danielcesario/cdc-core/internal/category"
	paymentmethod "github.com/danielcesario/cdc-core/internal/payment-method"
	"github.com/danielcesario/cdc-core/internal/user"
	"github.com/danielcesario/cdc-core/internal/wallet"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type TransactionService struct {
	repository              Repository
	userRepository          user.Repository
	walletRepository        wallet.Repository
	paymentmethodRepository paymentmethod.Repository
	categoryRepository      category.Repository
}

type Processor interface {
	Process(transaction *Transaction) error
}

func GetPaymentProcessor(paymentType paymentmethod.PaymentType) Processor {
	switch paymentType {
	case paymentmethod.CREDIT_CARD:
		return &ProcessCreditCard{}
	case paymentmethod.DEBIT_CARD:
		return &ProcessDebitCard{}
	case paymentmethod.BANKSLIP:
		return &ProcessBankSlip{}
	case paymentmethod.TRANSFER:
		return &ProcessTransfer{}
	case paymentmethod.MONEY:
		return &ProcessMoney{}
	default:
		// Retornar um processador padrão ou lidar com o pagamento desconhecido
		return nil
	}
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

func (s *TransactionService) Create(ctx context.Context, request TransactionRequest) (response *TransactionResponse, err error) {
	userEmail := ctx.Value("user_email").(string)
	currentUser, err := s.userRepository.FindByEmail(userEmail)
	if err != nil {
		return
	}

	transaction := request.toTransaction()
	transaction.Code = uuid.NewString()
	transaction.User = *currentUser
	transaction.UserID = currentUser.ID

	g := new(errgroup.Group)

	g.Go(func() error {
		wallet, err := s.validateWallet(request.Wallet, currentUser.ID)
		if err != nil {
			return err
		}

		transaction.Wallet = *wallet
		return nil
	})

	g.Go(func() error {
		if request.PaymentMethod != "" {
			pm, err := s.validatePaymentMethod(request.PaymentMethod, currentUser.ID)
			if err != nil {
				return err
			}

			transaction.PaymentMethod = *pm
			return nil
		}
		return nil
	})

	g.Go(func() error {
		category, err := s.validateCategory(request.Category, currentUser.ID)
		if err != nil {
			return err
		}

		transaction.Category = *category
		return nil
	})

	if err = g.Wait(); err != nil {
		return
	}

	if transaction.TransactionType == DEBIT {
		processor := GetPaymentProcessor(transaction.PaymentMethod.PaymentType)
		processor.Process(transaction)
	} else {
		processor := &ProcessCreditTransaction{}
		processor.Process(transaction)
	}

	savedId, err := s.repository.Store(*transaction)
	if err != nil {
		return
	}

	transaction.ID = savedId
	response = transaction.toResponse()
	return
}

func (s *TransactionService) GetByCode(ctx context.Context, code string) (*TransactionResponse, error) {
	transaction, err := s.repository.FindByCode(code)
	if err != nil {
		return nil, err
	}

	userEmail := ctx.Value("user_email").(string)
	currentUser, err := s.userRepository.FindByEmail(userEmail)
	if err != nil {
		return nil, err
	}

	if transaction.UserID != currentUser.ID {
		return nil, errors.New("invalid transaction owner")
	}

	g := new(errgroup.Group)

	g.Go(func() error {
		wallet, err := s.validateWallet(transaction.Wallet.Code, currentUser.ID)
		if err != nil {
			return err
		}

		transaction.Wallet = *wallet
		return nil
	})

	g.Go(func() error {
		if transaction.PaymentMethodID > 0 {
			pm, err := s.validatePaymentMethod(transaction.PaymentMethod.Code, currentUser.ID)
			if err != nil {
				return err
			}

			transaction.PaymentMethod = *pm
			return nil
		}
		return nil
	})

	g.Go(func() error {
		category, err := s.validateCategory(transaction.Category.Code, currentUser.ID)
		if err != nil {
			return err
		}

		transaction.Category = *category
		return nil
	})

	if err = g.Wait(); err != nil {
		return nil, err
	}

	return transaction.toResponse(), nil
}

func (s *TransactionService) Search(ctx context.Context, request SearchTransactionRequest) (*SearchPageResponse, error) {
	userEmail := ctx.Value("user_email").(string)
	currentUser, err := s.userRepository.FindByEmail(userEmail)
	if err != nil {
		return nil, err
	}

	g := new(errgroup.Group)

	g.Go(func() error {
		_, err := s.validateWallet(request.Wallet, currentUser.ID)
		if err != nil {
			return err
		}

		return nil
	})

	if err = g.Wait(); err != nil {
		return nil, err
	}

	totalItems, err := s.repository.Count(request)
	if err != nil {
		return nil, err
	}

	result, err := s.repository.Search(request)
	if err != nil {
		return nil, err
	}

	return &SearchPageResponse{
		Content:       result,
		TotalElements: int(totalItems),
		TotalPages:    int(totalItems) / request.Params.Limit.Size,
		Page:          request.Params.Limit.Page,
		Size:          request.Params.Limit.Size,
		Sort:          fmt.Sprintf("%s,%s", request.Params.Sort.Field, request.Params.Sort.Direction),
		IsLast:        request.Params.Limit.Page == (int(totalItems) / request.Params.Limit.Size),
	}, nil
}

func (s *TransactionService) validateWallet(walletCode string, userID uint64) (wallet *wallet.Wallet, err error) {
	wallet, err = s.walletRepository.FindByCode(walletCode)
	if err != nil {
		return
	}

	if wallet.UserID != userID {
		err = errors.New("invalid wallet owner")
		return
	}

	return
}

func (s *TransactionService) validatePaymentMethod(code string, userID uint64) (pm *paymentmethod.PaymentMethod, err error) {
	pm, err = s.paymentmethodRepository.FindByCode(code)
	if err != nil {
		return
	}

	if pm.UserID != userID {
		err = errors.New("invalid paymentmethod owner")
		return
	}

	return
}

func (s *TransactionService) validateCategory(code string, userID uint64) (cat *category.Category, err error) {
	cat, err = s.categoryRepository.FindByCode(code)
	if err != nil {
		return
	}

	if cat.UserID != userID {
		err = errors.New("invalid category owner")
		return
	}

	return
}
