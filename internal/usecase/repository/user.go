package repository

import "github.com/avag-sargsyan/golang-clean-arch/internal/domain/model"

// Secondary Port
type UserRepository interface {
	FindAll() ([]model.User, error)
	Save(model.User) error
}
