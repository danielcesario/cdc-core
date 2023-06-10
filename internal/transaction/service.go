package transaction

import (
	"context"
	"errors"

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
		return s.validateWallet(request.Wallet, *transaction)
	})

	g.Go(func() error {
		return s.validatePaymentMethod(request.PaymentMethod, *transaction)
	})

	g.Go(func() error {
		return s.validateCategory(request.Category, *transaction)
	})

	if err = g.Wait(); err != nil {
		return
	}

	processor := GetPaymentProcessor(transaction.PaymentType)
	processor.Process(transaction)

	savedId, err := s.repository.Store(*transaction)
	if err != nil {
		return
	}

	transaction.ID = savedId
	response = transaction.toResponse()
	return
}

func (s *TransactionService) validateWallet(code string, transaction Transaction) error {
	wallet, err := s.walletRepository.FindByCode(code)
	if err != nil {
		return err
	}

	if wallet.UserID != transaction.UserID {
		err = errors.New("invalid wallet owner")
		return err
	}

	transaction.Wallet = *wallet
	return nil
}

func (s *TransactionService) validatePaymentMethod(code string, transaction Transaction) error {
	paymentmethod, err := s.paymentmethodRepository.FindByCode(code)
	if err != nil {
		return err
	}

	if paymentmethod.UserID != transaction.UserID {
		err = errors.New("invalid paymentmethod owner")
		return err
	}

	transaction.PaymentMethod = *paymentmethod
	return nil
}

func (s *TransactionService) validateCategory(code string, transaction Transaction) error {
	category, err := s.categoryRepository.FindByCode(code)
	if err != nil {
		return err
	}

	if category.UserID != transaction.UserID {
		err = errors.New("invalid category owner")
		return err
	}

	transaction.Category = *category
	return nil
}
