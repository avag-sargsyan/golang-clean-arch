package usecase

import (
	"fmt"
	"github.com/avag-sargsyan/golang-clean-arch/internal/domain/model"
	"github.com/avag-sargsyan/golang-clean-arch/internal/usecase/repository"
)

// Primary Port
type UserService interface {
	GetUsers() ([]model.User, error)
	CreateUser(name string) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUsers() ([]model.User, error) {
	return s.repo.FindAll()
}

func (s *userService) CreateUser(name string) error {
	user := model.User{ID: fmt.Sprintf("%d", len(name)), Name: name}
	return s.repo.Save(user)
}
