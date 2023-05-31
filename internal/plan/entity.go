package plan

import (
	"gorm.io/gorm"
)

type Plan struct {
	gorm.Model
	ID     uint64 `gorm:"autoIncrement"`
	Name   string
	Value  int64
	Slug   string
	Active bool
}

func (p *Plan) toResponse() *PlanResponse {
	return &PlanResponse{
		Name:   p.Name,
		Value:  p.Value,
		Slug:   p.Slug,
		Active: p.Active,
	}
}

type PlanRequest struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
	Slug  string `json:"slug"`
}

func (p *PlanRequest) toPlan() *Plan {
	return &Plan{
		Name:  p.Name,
		Value: p.Value,
		Slug:  p.Slug,
	}
}

type PlanResponse struct {
	Name   string `json:"name"`
	Value  int64  `json:"value"`
	Slug   string `json:"slug"`
	Active bool   `json:"active"`
}
