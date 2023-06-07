package category

import (
	"github.com/danielcesario/cdc-core/internal/user"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	ID          uint64 `gorm:"autoIncrement"`
	Code        string
	Description string
	Color       string
	UserID      uint64
	User        user.User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (c *Category) toResponse() *CategoryResponse {
	return &CategoryResponse{
		Code:        c.Code,
		Description: c.Description,
		Color:       c.Color,
	}
}

type CategoryRequest struct {
	Description string `json:"description`
	Color       string `json:"color"`
}

func (c *CategoryRequest) toCategory() *Category {
	return &Category{
		Description: c.Description,
		Color:       c.Color,
	}
}

type CategoryResponse struct {
	Code        string `json:"code`
	Description string `json:"description`
	Color       string `json:"color`
}
