package transaction

import (
	"time"

	"github.com/google/uuid"
)

type ProcessCreditCard struct {
}

func (p *ProcessCreditCard) Process(transaction *Transaction) error {
	paymentMethod := transaction.PaymentMethod
	txDate := transaction.Date
	currentMonthCloseDay := time.Date(txDate.Year(), txDate.Month(), paymentMethod.CloseDay, 0, 0, 0, 0, time.UTC)

	var firstInstalment time.Time
	if txDate.After(currentMonthCloseDay) {
		firstInstalment = transaction.Date.AddDate(0, 1, 0)
	} else {
		firstInstalment = transaction.Date
	}

	instalmentValue := transaction.TotalAmount / transaction.TotalInstalments
	roundDiff := transaction.TotalAmount % transaction.TotalInstalments

	for i := 0; i < transaction.TotalInstalments; i++ {
		instalmentDueDate := firstInstalment.AddDate(0, i, 0)

		var instalmentStatus InstalmentStatus
		if instalmentDueDate.After(time.Now()) {
			instalmentStatus = SCHEDULLED
		} else {
			instalmentStatus = PRESENTED
		}

		entry := Entry{
			Code:             uuid.NewString(),
			TransactionID:    transaction.ID,
			Amount:           instalmentValue,
			DueDate:          instalmentDueDate,
			InstalmentStatus: instalmentStatus,
		}

		if i == (transaction.TotalInstalments-1) && roundDiff > 0 {
			entry.Amount += roundDiff
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
	var instalmentStatus InstalmentStatus
	if transaction.Date.After(time.Now()) {
		instalmentStatus = SCHEDULLED
	} else {
		instalmentStatus = PRESENTED
	}

	entry := Entry{
		Code:             uuid.NewString(),
		TransactionID:    transaction.ID,
		Amount:           transaction.TotalAmount,
		DueDate:          transaction.Date,
		InstalmentStatus: instalmentStatus,
	}

	transaction.Entries = append(transaction.Entries, entry)
}
