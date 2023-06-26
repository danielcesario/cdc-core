package transaction

type Reader interface {
	FindByCode(transactionCode string) (*Transaction, error)
	Count(search SearchTransactionRequest) (int64, error)
	Search(search SearchTransactionRequest) ([]TransactionSearchResult, error)
}

type Writer interface {
	Store(transaction Transaction) (uint64, error)
}

type Repository interface {
	Reader
	Writer
}
