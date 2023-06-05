package paymentmethod

type Reader interface {
	ListByUser(userCode string) ([]*PaymentMethod, error)
	FindByCode(code string) (*PaymentMethod, error)
}

type Writer interface {
	Store(paymentMethod *PaymentMethod) (uint64, error)
}

type Repository interface {
	Reader
	Writer
}
