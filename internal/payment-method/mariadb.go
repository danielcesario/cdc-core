package paymentmethod

import "gorm.io/gorm"

type mariaDBRepository struct {
	db *gorm.DB
}

func NewMariaDBRepository(db *gorm.DB) Repository {
	return &mariaDBRepository{
		db: db,
	}
}

func (r *mariaDBRepository) Store(paymentMethod *PaymentMethod) (uint64, error) {
	result := r.db.Save(&paymentMethod)
	if result.Error != nil {
		return 0, result.Error
	}
	return paymentMethod.ID, nil
}
