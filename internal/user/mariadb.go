package user

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

func (r *mariaDBRepository) FindByEmail(email string) (*User, error) {
	var user User
	record := r.db.Model(&User{}).Preload("Roles").Where("email = ?", email).First(&user)
	if record.Error != nil {
		return nil, record.Error
	}

	return &user, nil
}

func (r *mariaDBRepository) GetRoleByName(name string) (*Role, error) {
	var role Role
	record := r.db.Model(&Role{}).Where("role = ?", name).First(&role)
	if record.Error != nil {
		return nil, record.Error
	}

	return &role, nil
}

func (r *mariaDBRepository) Store(user *User) (uint64, error) {
	result := r.db.Create(&user)
	if result.Error != nil {
		return 0, result.Error
	}
	return user.ID, nil
}

func (r *mariaDBRepository) FindAll() ([]*User, error) {
	var users []*User
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}
