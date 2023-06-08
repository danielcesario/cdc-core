package transaction

type Reader interface {
}

type Writer interface {
	Store(transaction Transaction) (uint64, error)
}

type Repository interface {
	Reader
	Writer
}
