package transaction_test

import (
	"testing"
	"time"

	paymentmethod "github.com/danielcesario/cdc-core/internal/payment-method"
	"github.com/danielcesario/cdc-core/internal/transaction"
	"github.com/stretchr/testify/assert"
)

func TestOneInstalmentProcessors(t *testing.T) {
	scenarios := []struct {
		Title           string
		Processor       transaction.Processor
		TransactionDate time.Time
		ExpectedStatus  transaction.InstalmentStatus
	}{
		{
			Title:           "ProcessDebitCard One Year Before - SCHEDULLED",
			Processor:       &transaction.ProcessDebitCard{},
			TransactionDate: time.Now().AddDate(1, 0, 0),
			ExpectedStatus:  transaction.SCHEDULLED,
		},
		{
			Title:           "ProcessBankSlip Current Date - PRESENTED",
			Processor:       &transaction.ProcessBankSlip{},
			TransactionDate: time.Now(), ExpectedStatus: transaction.PRESENTED,
		},
		{
			Title:           "ProcessTransfer One Month Before - SCHEDULLED",
			Processor:       &transaction.ProcessTransfer{},
			TransactionDate: time.Now().AddDate(0, 1, 0),
			ExpectedStatus:  transaction.SCHEDULLED,
		},
		{
			Title:           "ProcessMoney Current Date - PRESENTED",
			Processor:       &transaction.ProcessMoney{},
			TransactionDate: time.Now(),
			ExpectedStatus:  transaction.PRESENTED,
		},
		{
			Title:           "ProcessCreditTransaction One Day Before - SCHEDULLED",
			Processor:       &transaction.ProcessCreditTransaction{},
			TransactionDate: time.Now().AddDate(0, 0, 1),
			ExpectedStatus:  transaction.SCHEDULLED,
		},
	}

	for _, scenario := range scenarios {
		mockTransaction := &transaction.Transaction{
			ID:          1,
			TotalAmount: 10000,
			Date:        scenario.TransactionDate,
		}

		t.Run(scenario.Title, func(t *testing.T) {
			err := scenario.Processor.Process(mockTransaction)

			assert.Nil(t, err)
			assert.Equal(t, 1, len(mockTransaction.Entries))
			assert.Equal(t, scenario.ExpectedStatus, mockTransaction.Entries[0].InstalmentStatus)
		})
	}
}

func mockTime(t *testing.T, year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func TestMultipleIntalmentProcessor(t *testing.T) {
	scenarios := []struct {
		Title                     string
		CardCloseDay              int
		TotalAmount               int
		TotalInstalments          int
		TransactionDate           time.Time
		ExpectedInstalmentsAmount []int
		ExpectedInstalmentsDates  []time.Time
		ExpectedInstalmentStatus  []transaction.InstalmentStatus
	}{
		{
			Title:                     "Process Transactions with two equal instalments before card close day",
			CardCloseDay:              10,
			TotalAmount:               50000,
			TotalInstalments:          2,
			TransactionDate:           mockTime(t, 2023, 1, 5),
			ExpectedInstalmentsAmount: []int{25000, 25000},
			ExpectedInstalmentsDates:  []time.Time{mockTime(t, 2023, 1, 5), mockTime(t, 2023, 2, 5)},
		},
		{
			Title:                     "Process Transactions with two equal instalments after card close day",
			CardCloseDay:              10,
			TotalAmount:               50000,
			TotalInstalments:          2,
			TransactionDate:           mockTime(t, 2023, 1, 15),
			ExpectedInstalmentsAmount: []int{25000, 25000},
			ExpectedInstalmentsDates:  []time.Time{mockTime(t, 2023, 2, 15), mockTime(t, 2023, 3, 15)},
		},
		{
			Title:                     "Process Transactions with different status",
			CardCloseDay:              10,
			TotalAmount:               50000,
			TotalInstalments:          2,
			TransactionDate:           mockTime(t, time.Now().Year(), int(time.Now().Month()), 5),
			ExpectedInstalmentsAmount: []int{25000, 25000},
			ExpectedInstalmentsDates:  []time.Time{mockTime(t, time.Now().Year(), int(time.Now().Month()), 5), mockTime(t, time.Now().Year(), int(time.Now().Month()), 5).AddDate(0, 1, 0)},
			ExpectedInstalmentStatus:  []transaction.InstalmentStatus{transaction.PRESENTED, transaction.SCHEDULLED},
		},
		{
			Title:                     "Process Transactions with difference on round value",
			CardCloseDay:              10,
			TotalAmount:               50033,
			TotalInstalments:          4,
			TransactionDate:           mockTime(t, 2023, 1, 15),
			ExpectedInstalmentsAmount: []int{12508, 12508, 12508, 12509},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.Title, func(t *testing.T) {
			paymentmethod := &paymentmethod.PaymentMethod{
				CloseDay: scenario.CardCloseDay,
			}

			mockTransaction := &transaction.Transaction{
				ID:               1,
				TotalAmount:      scenario.TotalAmount,
				Date:             scenario.TransactionDate,
				TotalInstalments: scenario.TotalInstalments,
				PaymentMethod:    *paymentmethod,
			}

			t.Run(scenario.Title, func(t *testing.T) {
				processor := &transaction.ProcessCreditCard{}
				err := processor.Process(mockTransaction)

				assert.Nil(t, err)
				assert.Equal(t, scenario.TotalInstalments, len(mockTransaction.Entries))
				for i, entry := range mockTransaction.Entries {
					assert.Equal(t, scenario.ExpectedInstalmentsAmount[i], entry.Amount)

					if scenario.ExpectedInstalmentsDates != nil {
						assert.Equal(t, scenario.ExpectedInstalmentsDates[i], entry.DueDate)
					}

					if scenario.ExpectedInstalmentStatus != nil {
						assert.Equal(t, scenario.ExpectedInstalmentStatus[i], entry.InstalmentStatus)
					}
				}
			})

		})
	}
}
