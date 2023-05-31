package user

type Reader interface {
	FindByEmail(email string) (*User, error)
	GetRoleByName(name string) (*Role, error)
	FindAll() ([]*User, error)
}

type Writer interface {
	Store(user *User) (uint64, error)
}

type Repository interface {
	Reader
	Writer
}
