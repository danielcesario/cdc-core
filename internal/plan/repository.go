package plan

type Reader interface {
	FindBySlug(slug string) (*Plan, error)
	ListAll() ([]*Plan, error)
}

type Writer interface {
	Store(plan *Plan) (uint64, error)
	Update(plan *Plan) error
}

type Repository interface {
	Reader
	Writer
}
