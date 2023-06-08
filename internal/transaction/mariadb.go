package transaction

import "gorm.io/gorm"

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
