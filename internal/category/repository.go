package category

type Reader interface {
	ListByUser(userCode string) ([]*Category, error)
	FindByCode(code string) (*Category, error)
}

type Writer interface {
	Store(category *Category) (uint64, error)
}

type Repository interface {
	Reader
	Writer
}
