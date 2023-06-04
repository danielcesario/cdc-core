package transaction

type Reader interface {
}

type Writer interface {
}

type Repository interface {
	Reader
	Writer
}
