package user

import "github.com/ipxsandbox/internal/entity"

type Repository interface {
    FindAll() ([]entity.User, error)
    Create(user entity.User) (entity.User, error)
    FindByEmail(email string) (entity.User, error)
}