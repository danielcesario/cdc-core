package transaction

type Reader interface {
	FindByCode(transactionCode string) (*Transaction, error)
}

type Writer interface {
	Store(transaction Transaction) (uint64, error)
}

type Repository interface {
	Reader
	Writer
}
