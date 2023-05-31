package plan

import "gorm.io/gorm"

type mariaDBRepository struct {
	db *gorm.DB
}

func NewMariaDBRepository(db *gorm.DB) Repository {
	return &mariaDBRepository{
		db: db,
	}
}

func (r *mariaDBRepository) FindBySlug(slug string) (*Plan, error) {
	var plan Plan
	record := r.db.Model(&Plan{}).Where("slug = ?", slug).First(&plan)
	if record.Error != nil {
		return nil, record.Error
	}

	return &plan, nil
}

func (r *mariaDBRepository) Store(plan *Plan) (uint64, error) {
	result := r.db.Create(&plan)
	if result.Error != nil {
		return 0, result.Error
	}
	return plan.ID, nil
}

func (r *mariaDBRepository) ListAll() ([]*Plan, error) {
	var plans []*Plan
	result := r.db.Order("name asc").Find(&plans)

	if result.Error != nil {
		return nil, result.Error
	}

	return plans, nil
}

func (r *mariaDBRepository) Update(plan *Plan) error {
	result := r.db.Save(&plan)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
