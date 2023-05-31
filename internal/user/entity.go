package user

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	ID        uint64 `gorm:"autoIncrement"`
	Role      string `gorm:"unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	gorm.Model
	ID        uint64 `gorm:"autoIncrement"`
	Name      string
	Email     string
	Code      string
	Password  string
	Active    bool
	Roles     []Role `gorm:"many2many:user_role;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		Code:  u.Code,
		Name:  u.Name,
		Email: u.Email,
	}
}

func (u *User) GetRoles() []string {
	var roles []string
	for _, role := range u.Roles {
		roles = append(roles, role.Role)
	}
	return roles
}

type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserRequest) toUser() *User {
	return &User{
		Name:     u.Name,
		Email:    u.Email,
		Active:   false,
		Code:     uuid.NewString(),
		Password: u.Password,
	}
}

func (u *UserRequest) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

type UserResponse struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
