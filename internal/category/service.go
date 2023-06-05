package category

import (
	"context"
	"errors"

	"github.com/danielcesario/cdc-core/internal/user"
	"github.com/google/uuid"
)

type CategoryService struct {
	repository     Repository
	userRepository user.Repository
}

func NewCategoryService(repository Repository, userRepository user.Repository) *CategoryService {
	return &CategoryService{
		repository:     repository,
		userRepository: userRepository,
	}
}

func (s *CategoryService) Create(ctx context.Context, request CategoryRequest) (*CategoryResponse, error) {
	userEmail := ctx.Value("user_email").(string)
	currentUser, err := s.userRepository.FindByEmail(userEmail)
	if err != nil {
		return nil, err
	}

	category := request.toCategory()
	category.Code = uuid.NewString()
	category.UserID = currentUser.ID
	category.User = *currentUser

	_, err = s.repository.Store(category)
	if err != nil {
		return nil, err
	}

	return category.toResponse(), nil
}

func (s *CategoryService) List(ctx context.Context) (response []*CategoryResponse, err error) {
	userEmail := ctx.Value("user_email").(string)
	currentUser, err := s.userRepository.FindByEmail(userEmail)
	if err != nil {
		return
	}

	result, err := s.repository.ListByUser(currentUser.Code)
	if err != nil {
		return
	}

	for _, category := range result {
		response = append(response, category.toResponse())
	}

	return
}

func (s *CategoryService) Update(ctx context.Context, code string, request CategoryRequest) (*CategoryResponse, error) {
	savedCategory, err := s.repository.FindByCode(code)
	if err != nil {
		return nil, err
	}

	userEmail := ctx.Value("user_email").(string)
	currentUser, err := s.userRepository.FindByEmail(userEmail)
	if err != nil {
		return nil, err
	}

	if savedCategory.UserID != currentUser.ID {
		return nil, errors.New("invalid category owner")
	}

	newCategory := request.toCategory()
	savedCategory.Description = newCategory.Description
	savedCategory.Color = newCategory.Color

	_, err = s.repository.Store(savedCategory)
	if err != nil {
		return nil, err
	}

	return savedCategory.toResponse(), nil
}
