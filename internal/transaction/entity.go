package transaction

import (
	"errors"
	"time"

	"github.com/danielcesario/cdc-core/internal/category"
	paymentmethod "github.com/danielcesario/cdc-core/internal/payment-method"
	"github.com/danielcesario/cdc-core/internal/user"
	"github.com/danielcesario/cdc-core/internal/wallet"
	"gorm.io/gorm"
)

type TransactionType int

const (
	CREDIT TransactionType = iota // 0
	DEBIT
)

func (i TransactionType) String() string {
	switch i {
	case 0:
		return "CREDIT"
	case 1:
		return "DEBIT"
	default:
		return "Invalid Status"
	}
}

func GetTransactionType(s string) (TransactionType, error) {
	switch s {
	case "CREDIT":
		return CREDIT, nil
	case "DEBIT":
		return DEBIT, nil
	}
	return -1, errors.New("unknown Transaction Type: " + s)
}

type InstalmentStatus int

const (
	SCHEDULLED InstalmentStatus = iota // 0
	PRESENTED
)

func (i InstalmentStatus) String() string {
	switch i {
	case 0:
		return "SCHEDULLED"
	case 1:
		return "PRESENTED"
	default:
		return "Invalid Status"
	}
}

type Transaction struct {
	gorm.Model
	ID               uint64 `gorm:"autoIncrement"`
	Code             string
	TotalAmount      int
	TotalInstalments int
	TransactionType  TransactionType
	Description      string
	Date             time.Time
	UserID           uint64
	User             user.User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	WalletID         uint64
	Wallet           wallet.Wallet `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PaymentMethodID  uint64
	PaymentMethod    paymentmethod.PaymentMethod `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CategoryID       uint64
	Category         category.Category `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Entries          []Entry           `gorm:"foreignKey:TransactionID"`
}

func (t *Transaction) toResponse() *TransactionResponse {
	var entries []EntryResponse
	for _, e := range t.Entries {
		entries = append(entries, *e.toResponse())
	}

	return &TransactionResponse{
		Code:             t.Code,
		TotalAmount:      t.TotalAmount,
		TotalInstalments: t.TotalInstalments,
		TransactionType:  t.TransactionType.String(),
		Description:      t.Description,
		Date:             t.Date.Format("2006-01-02T15:04:05 -07:00:00"),
		User:             user.UserResponse{Code: t.User.Code, Name: t.User.Name, Email: t.User.Email},
		Wallet:           wallet.WalletResponse{Name: t.Wallet.Name, Code: t.Wallet.Code},
		PaymentMethod:    paymentmethod.PaymentMethodResponse{Code: t.PaymentMethod.Code, Description: t.PaymentMethod.Description, PaymentType: t.PaymentMethod.PaymentType.String()},
		Category:         category.CategoryResponse{Code: t.Category.Code, Description: t.Category.Description, Color: t.Category.Color},
		Entries:          entries,
	}
}

type Entry struct {
	gorm.Model
	ID               uint64 `gorm:"autoIncrement"`
	Code             string
	TransactionID    uint64
	Amount           int
	DueDate          time.Time
	InstalmentStatus InstalmentStatus
}

func (e *Entry) toResponse() *EntryResponse {
	return &EntryResponse{
		Code:             e.Code,
		Amount:           e.Amount,
		DueDate:          e.DueDate.Format(time.RFC3339),
		InstalmentStatus: e.InstalmentStatus.String(),
	}
}

type TransactionRequest struct {
	TotalAmount      int    `json:"total_amount"`
	TotalInstalments int    `json:"total_instalments"`
	TransactionType  string `json:"transaction_type"`
	Description      string `json:"description"`
	Date             string `json:"date"`
	Wallet           string `json:"wallet"`
	PaymentMethod    string `json:"payment_method"`
	Category         string `json:"category"`
}

func (tr *TransactionRequest) toTransaction() *Transaction {
	transactionType, _ := GetTransactionType(tr.TransactionType)
	requestDate, err := time.Parse(time.RFC3339, tr.Date)
	if err != nil {
		requestDate = time.Now()
	}

	return &Transaction{
		TotalAmount:      tr.TotalAmount,
		TotalInstalments: tr.TotalInstalments,
		Description:      tr.Description,
		TransactionType:  transactionType,
		Date:             requestDate,
		Wallet:           wallet.Wallet{Code: tr.Wallet},
		PaymentMethod:    paymentmethod.PaymentMethod{Code: tr.PaymentMethod},
		Category:         category.Category{Code: tr.Category},
	}
}

type EntryResponse struct {
	Code             string `json:"code"`
	Amount           int    `json:"amount"`
	DueDate          string `json:"due_date"`
	InstalmentStatus string `json:"instalment_status"`
}

type TransactionResponse struct {
	Code             string                              `json:"code"`
	TotalAmount      int                                 `json:"total_amount"`
	TotalInstalments int                                 `json:"total_instalments"`
	TransactionType  string                              `json:"transaction_type"`
	Description      string                              `json:"description"`
	Date             string                              `json:"date"`
	User             user.UserResponse                   `json:"user"`
	Wallet           wallet.WalletResponse               `json:"wallet"`
	PaymentMethod    paymentmethod.PaymentMethodResponse `json:"payment_method"`
	Category         category.CategoryResponse           `json:"category"`
	Entries          []EntryResponse                     `json:"entries,omitempty"`
}

type SearchLimit struct {
	Size int `json:"size"`
	Page int `json:"page"`
}

type SerchSort struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

type SearchParams struct {
	Sort  SerchSort   `json:"sort"`
	Limit SearchLimit `json:"limit"`
}

type SearchTransactionRequest struct {
	From   string       `json:"from"`
	To     string       `json:"to"`
	Wallet string       `json:"wallet"`
	Params SearchParams `json:"params"`
}

type SerchItemResponse struct {
	TransactionType string                              `json:"transaction_type"`
	Amount          int                                 `json:"amount"`
	Date            string                              `json:"date"`
	Wallet          wallet.WalletResponse               `json:"wallet"`
	PaymentMethod   paymentmethod.PaymentMethodResponse `json:"payment_method"`
	Category        category.CategoryResponse           `json:"category"`
}

type SearchPageResponse struct {
	Content       []TransactionSearchResult `json:"content"`
	TotalElements int                       `json:"total_elements"`
	TotalPages    int                       `json:"total_pages"`
	Page          int                       `json:"page"`
	Size          int                       `json:"size"`
	Sort          string                    `json:"sort"`
	IsLast        bool                      `json:"is_last"`
}

type TransactionSearchResult struct {
	Code          string    `json:"code"`
	Type          string    `json:"type"`
	Amount        int       `json:"amount"`
	DueDate       time.Time `json:"due_date"`
	Description   string    `json:"description"`
	User          string    `json:"user"`
	PaymentMethod string    `json:"payment_method"`
	Category      string    `json:"category"`
}
