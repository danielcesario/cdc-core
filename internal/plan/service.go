package plan

import (
	"context"
	"errors"
)

type PlanService struct {
	repository Repository
}

func NewPlanService(repository Repository) *PlanService {
	return &PlanService{
		repository: repository,
	}
}

func (s *PlanService) Create(ctx context.Context, request PlanRequest) (*PlanResponse, error) {
	var plan Plan = *request.toPlan()
	plan.Active = true

	currentPlan, err := s.repository.FindBySlug(plan.Slug)
	if err == nil {
		if currentPlan.ID > 0 {
			return nil, errors.New("slug is already being used")
		}
	}

	savedIdPlan, err := s.repository.Store(&plan)
	if err != nil {
		return nil, err
	}

	plan.ID = savedIdPlan
	return plan.toResponse(), err
}

func (s *PlanService) ListAll(ctx context.Context) (result []*PlanResponse, err error) {
	plans, err := s.repository.ListAll()
	if err != nil {
		return
	}

	for _, plan := range plans {
		result = append(result, plan.toResponse())
	}

	return
}

func (s *PlanService) GetByCode(ctx context.Context, code string) (result *PlanResponse, err error) {
	plan, err := s.repository.FindBySlug(code)
	if err != nil {
		return
	}

	result = plan.toResponse()
	return
}

func (s *PlanService) Update(ctx context.Context, planCode string, request PlanRequest) (result *PlanResponse, err error) {
	plan, err := s.repository.FindBySlug(planCode)
	if err != nil {
		return
	}

	if plan.ID == 0 {
		return nil, errors.New("plan not found")
	}

	plan.Name = request.Name
	plan.Value = request.Value

	err = s.repository.Update(plan)
	if err != nil {
		return
	}

	result = plan.toResponse()
	return
}
