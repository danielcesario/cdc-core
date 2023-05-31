package user

import (
	"context"
	"errors"
)

type UserService struct {
	repository Repository
}

func NewUserService(repository Repository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) Register(ctx context.Context, request UserRequest) (*UserResponse, error) {
	var user User = *request.toUser()
	_, err := s.repository.FindByEmail(request.Email)
	if err == nil {
		return nil, errors.New("email is already being used")
	}

	roleUser, err := s.repository.GetRoleByName("USER_DEFAULT")
	if err != nil {
		return nil, errors.New("error on load user role")
	}

	user.Roles = []Role{*roleUser}

	savedIdUser, err := s.repository.Store(&user)
	if err != nil {
		return nil, err
	}

	user.ID = savedIdUser
	return user.ToResponse(), nil
}

func (s *UserService) CheckEmailAndPassword(ctx context.Context, email, password string) (*User, error) {
	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if err := user.CheckPassword(password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *UserService) ListUsers(ctx context.Context) (response []*UserResponse, err error) {
	users, err := s.repository.FindAll()
	if err != nil {
		return
	}

	for _, user := range users {
		response = append(response, user.ToResponse())
	}

	return response, nil
}
