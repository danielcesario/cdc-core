package category

import "gorm.io/gorm"

type mariaDBRepository struct {
	db *gorm.DB
}

func NewMariaDBRepository(db *gorm.DB) Repository {
	return &mariaDBRepository{
		db: db,
	}
}

func (r *mariaDBRepository) Store(category *Category) (uint64, error) {
	result := r.db.Save(&category)
	if result.Error != nil {
		return 0, result.Error
	}
	return category.ID, nil
}

func (r *mariaDBRepository) ListByUser(userCode string) ([]*Category, error) {
	var categories []*Category

	result := r.db.Joins("JOIN user ON user.id = category.user_id").
		Where("user.code = ?", userCode).
		Find(&categories)

	if result.Error != nil {
		return nil, result.Error
	}

	return categories, nil
}

func (r *mariaDBRepository) FindByCode(code string) (*Category, error) {
	var category Category
	record := r.db.Model(&Category{}).
		Preload("User").
		Where("code = ?", code).
		First(&category)

	if record.Error != nil {
		return nil, record.Error
	}

	return &category, nil
}
