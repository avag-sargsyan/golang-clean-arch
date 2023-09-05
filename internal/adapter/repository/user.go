package repository

import (
	"github.com/avag-sargsyan/golang-clean-arch/internal/domain/model"
	"github.com/avag-sargsyan/golang-clean-arch/internal/usecase/repository"
)

// Secondary Adapter
type InMemoryUserRepository struct {
	users []model.User
}

func NewInMemoryUserRepository() repository.UserRepository {
	return &InMemoryUserRepository{users: []model.User{}}
}

func (r *InMemoryUserRepository) FindAll() ([]model.User, error) {
	return r.users, nil
}

func (r *InMemoryUserRepository) Save(u model.User) error {
	r.users = append(r.users, u)
	return nil
}
