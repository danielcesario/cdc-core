package paymentmethod

import (
	"errors"

	"github.com/danielcesario/cdc-core/internal/user"
	"gorm.io/gorm"
)

type PaymentType int

const (
	BANKSLIP PaymentType = iota // 0
	TRANSFER
	CREDIT_CARD
	DEBIT_CARD
	MONEY
)

var paymentTypeNames = [...]string{
	"BANKSLIP",
	"TRANSFER",
	"CREDIT_CARD",
	"DEBIT_CARD",
	"MONEY",
}

func (p PaymentType) String() string {
	if p < BANKSLIP || p > MONEY {
		return "Unknown"
	}
	return paymentTypeNames[p]
}

func GetPaymentType(s string) (PaymentType, error) {
	switch s {
	case "BANKSLIP":
		return BANKSLIP, nil
	case "TRANSFER":
		return TRANSFER, nil
	case "CREDIT_CARD":
		return CREDIT_CARD, nil
	case "DEBIT_CARD":
		return DEBIT_CARD, nil
	case "MONEY":
		return MONEY, nil
	}
	return -1, errors.New("unknown payment type: " + s)
}

type PaymentMethod struct {
	gorm.Model
	ID          uint64 `gorm:"autoIncrement"`
	Code        string
	PaymentType PaymentType
	Description string
	DueDay      int
	CloseDay    int
	UserID      uint64
	User        user.User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (p *PaymentMethod) toResponse() *PaymentMethodResponse {
	return &PaymentMethodResponse{
		Code:        p.Code,
		PaymentType: p.PaymentType.String(),
		Description: p.Description,
		DueDay:      p.DueDay,
		CloseDay:    p.CloseDay,
	}
}

type PaymentMethodRequest struct {
	PaymentType string `json:"payment_type"`
	Description string `json:"description"`
	DueDay      int    `json:"due_day"`
	CloseDay    int    `json:"close_day"`
}

func (p *PaymentMethodRequest) toPaymentMethod() *PaymentMethod {
	paymentType, err := GetPaymentType(p.PaymentType)
	if err != nil {
		panic(err)
	}

	return &PaymentMethod{
		PaymentType: paymentType,
		Description: p.Description,
		DueDay:      p.DueDay,
		CloseDay:    p.CloseDay,
	}
}

type PaymentMethodResponse struct {
	Code        string `json:"code"`
	PaymentType string `json:"payment_type"`
	Description string `json:"description"`
	DueDay      int    `json:"due_day"`
	CloseDay    int    `json:"close_day"`
}
