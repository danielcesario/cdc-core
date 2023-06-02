package wallet

import "gorm.io/gorm"

type mariaDBRepository struct {
	db *gorm.DB
}

func NewMariaDBRepository(db *gorm.DB) Repository {
	return &mariaDBRepository{
		db: db,
	}
}

func (r *mariaDBRepository) Store(wallet *Wallet) (uint64, error) {
	result := r.db.Save(&wallet)
	if result.Error != nil {
		return 0, result.Error
	}
	return wallet.ID, nil
}

func (r *mariaDBRepository) ListByUser(userCode string) ([]*Wallet, error) {
	var wallets []*Wallet

	result := r.db.Joins("JOIN user ON user.id = wallet.user_id").
		Preload("User").
		Preload("Collaborators").
		Where("user.code = ?", userCode).
		Find(&wallets)

	if result.Error != nil {
		return nil, result.Error
	}

	return wallets, nil
}

func (r *mariaDBRepository) FindByCode(code string) (*Wallet, error) {
	var wallet Wallet
	record := r.db.Model(&Wallet{}).
		Preload("User").
		Preload("Collaborators").
		Where("code = ?", code).
		First(&wallet)

	if record.Error != nil {
		return nil, record.Error
	}

	return &wallet, nil
}
