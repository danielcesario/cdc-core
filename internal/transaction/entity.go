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

type PaymentType int

const (
	BANKSLIP PaymentType = iota // 0
	TRANSFER
	CREDIT_CARD
	DEBIT_CARD
	MONEY
)

type InstalmentStatus int

const (
	SCHEDULLED InstalmentStatus = iota // 0
	PRESENTED
)

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
	PaymentType      PaymentType
	PaymentMethodID  uint64
	PaymentMethod    paymentmethod.PaymentMethod `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CategoryID       uint64
	Category         category.Category `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Entries          []Entry
}

type Entry struct {
	gorm.Model
	ID               uint64 `gorm:"autoIncrement"`
	Code             string
	TransactionID    uint64
	Transaction      user.User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Amount           int
	DueDate          time.Time
	InstalmentStatus InstalmentStatus
}
