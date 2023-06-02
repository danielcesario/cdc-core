package wallet

type Reader interface {
	ListByUser(userCode string) ([]*Wallet, error)
	FindByCode(walletCode string) (*Wallet, error)
}

type Writer interface {
	Store(wallet *Wallet) (uint64, error)
}

type Repository interface {
	Reader
	Writer
}
