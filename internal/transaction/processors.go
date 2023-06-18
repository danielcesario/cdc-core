package transaction

import (
	"time"

	"github.com/google/uuid"
)

type ProcessCreditCard struct {
}

func (p *ProcessCreditCard) Process(transaction *Transaction) error {
	paymentMethod := transaction.PaymentMethod
	currentMonthCloseDay := time.Date(time.Now().Year(), time.Now().Month(), paymentMethod.CloseDay, 0, 0, 0, 0, time.UTC)

	var firstInstalment time.Time
	if time.Now().After(currentMonthCloseDay) {
		firstInstalment = time.Now().AddDate(0, 1, 0)
	} else {
		firstInstalment = time.Now()
	}

	instalmentValue := transaction.TotalAmount / transaction.TotalInstalments
	// TODO: Verify round difference

	for i := 0; i < transaction.TotalInstalments; i++ {
		entry := Entry{
			Code:             uuid.NewString(),
			TransactionID:    transaction.ID,
			Amount:           instalmentValue,
			DueDate:          firstInstalment.AddDate(0, i, 0),
			InstalmentStatus: SCHEDULLED,
		}

		transaction.Entries = append(transaction.Entries, entry)
	}

	return nil
}

type ProcessDebitCard struct {
}

func (p *ProcessDebitCard) Process(transaction *Transaction) error {
	processOneInstalment(transaction)
	return nil
}

type ProcessBankSlip struct {
}

func (p *ProcessBankSlip) Process(transaction *Transaction) error {
	processOneInstalment(transaction)
	return nil
}

type ProcessTransfer struct {
}

func (p *ProcessTransfer) Process(transaction *Transaction) error {
	processOneInstalment(transaction)
	return nil
}

type ProcessMoney struct {
}

func (p *ProcessMoney) Process(transaction *Transaction) error {
	processOneInstalment(transaction)
	return nil
}

type ProcessCreditTransaction struct {
}

func (p *ProcessCreditTransaction) Process(transaction *Transaction) error {
	processOneInstalment(transaction)
	return nil
}

func processOneInstalment(transaction *Transaction) {
	entry := Entry{
		Code:             uuid.NewString(),
		TransactionID:    transaction.ID,
		Amount:           transaction.TotalAmount,
		DueDate:          time.Now(),
		InstalmentStatus: PRESENTED,
	}

	transaction.Entries = append(transaction.Entries, entry)
}
