package transaction

import (
	"gorm.io/gorm"
)

type mariaDBRepository struct {
	db *gorm.DB
}

func NewMariaDBRepository(db *gorm.DB) Repository {
	return &mariaDBRepository{
		db: db,
	}
}

func (r *mariaDBRepository) Store(transaction Transaction) (uint64, error) {
	result := r.db.Save(&transaction)
	if result.Error != nil {
		return 0, result.Error
	}
	return transaction.ID, nil
}

func (r *mariaDBRepository) FindByCode(code string) (*Transaction, error) {
	var transaction Transaction

	record := r.db.
		//Model(&Transaction{}).
		Preload("User").
		Preload("Wallet", func(db *gorm.DB) *gorm.DB { return db.Select("ID", "Name", "Code") }).
		Preload("PaymentMethod").
		Preload("Category").
		Preload("Entries").
		Where("code = ?", code).
		First(&transaction)

	if record.Error != nil {
		return nil, record.Error
	}

	return &transaction, nil
}
