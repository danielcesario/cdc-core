package transaction_test

import (
	"reflect"
	"testing"

	"github.com/danielcesario/cdc-core/internal/transaction"
	"github.com/stretchr/testify/assert"
)

func TestOneInstalmentProcessors(t *testing.T) {
	scenarios := []struct {
		Processor transaction.Processor
	}{
		{Processor: &transaction.ProcessDebitCard{}},
		{Processor: &transaction.ProcessBankSlip{}},
		{Processor: &transaction.ProcessTransfer{}},
		{Processor: &transaction.ProcessMoney{}},
		{Processor: &transaction.ProcessCreditTransaction{}},
	}

	for _, scenario := range scenarios {
		mockTransaction := &transaction.Transaction{
			ID:          1,
			TotalAmount: 10000,
		}

		t.Run(getType(scenario.Processor), func(t *testing.T) {
			err := scenario.Processor.Process(mockTransaction)
			assert.Nil(t, err)
			assert.Equal(t, 1, len(mockTransaction.Entries))
		})
	}
}

func getType(myvar interface{}) string {
	valueOf := reflect.ValueOf(myvar)

	if valueOf.Type().Kind() == reflect.Ptr {
		return reflect.Indirect(valueOf).Type().Name()
	} else {
		return valueOf.Type().Name()
	}
}
