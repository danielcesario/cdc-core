package transaction

import (
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
		return "PRESENT"
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
	UserID           uint64
	User             user.User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	WalletID         uint64
	Wallet           wallet.Wallet `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PaymentType      paymentmethod.PaymentType
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
		User:             user.UserResponse{Code: t.User.Code, Name: t.User.Name},
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
		DueDate:          e.DueDate.Format("2006-01-02T15:04:05 -07:00:00"),
		InstalmentStatus: e.InstalmentStatus.String(),
	}
}

type TransactionRequest struct {
	TotalAmount      int    `json:"total_amount"`
	TotalInstalments int    `json:"total_instalments"`
	PaymentType      string `json:"payment_type"`
	Description      string `json:"description"`
	Wallet           string `json:"wallet"`
	PaymentMethod    string `json:"payment_method"`
	Category         string `json:"category"`
}

func (tr *TransactionRequest) toTransaction() *Transaction {
	paymentType, _ := paymentmethod.GetPaymentType(tr.PaymentType)

	return &Transaction{
		TotalAmount:      tr.TotalAmount,
		TotalInstalments: tr.TotalInstalments,
		Description:      tr.Description,
		PaymentType:      paymentType,
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
	User             user.UserResponse                   `json:"user"`
	Wallet           wallet.WalletResponse               `json:"wallet"`
	PaymentMethod    paymentmethod.PaymentMethodResponse `json:"payment_method"`
	Category         category.CategoryResponse           `json:"category"`
	Entries          []EntryResponse                     `json:"entries"`
}
