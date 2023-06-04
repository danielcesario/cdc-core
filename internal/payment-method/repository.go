package paymentmethod

type Reader interface {
}

type Writer interface {
	Store(paymentMethod *PaymentMethod) (uint64, error)
}

type Repository interface {
	Reader
	Writer
}
