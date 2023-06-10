package transaction

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ProcessCreditCard struct {
}

func (p *ProcessCreditCard) Process(transaction *Transaction) error {
	fmt.Println("ProcessCreditCard")
	return nil
}

type ProcessDebitCard struct {
}

func (p *ProcessDebitCard) Process(transaction *Transaction) error {
	fmt.Println("ProcessDebitCard")
	processOneInstalment(transaction)
	return nil
}

type ProcessBankSlip struct {
}

func (p *ProcessBankSlip) Process(transaction *Transaction) error {
	fmt.Println("ProcessBankSlip")
	processOneInstalment(transaction)
	return nil
}

type ProcessTransfer struct {
}

func (p *ProcessTransfer) Process(transaction *Transaction) error {
	fmt.Println("ProcessTransfer")
	processOneInstalment(transaction)
	return nil
}

type ProcessMoney struct {
}

func (p *ProcessMoney) Process(transaction *Transaction) error {
	fmt.Println("ProcessMoney")
	processOneInstalment(transaction)
	return nil
}

func processOneInstalment(transaction *Transaction) {
	entry := Entry{
		Code:             uuid.NewString(),
		Transaction:      *transaction,
		Amount:           transaction.TotalAmount,
		DueDate:          time.Now(),
		InstalmentStatus: PRESENTED,
	}

	transaction.Entries = append(transaction.Entries, entry)
}
