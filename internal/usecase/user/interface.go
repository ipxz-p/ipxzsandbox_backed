package user

import "github.com/ipxsandbox/internal/entity"

type Usecase interface {
    GetAllUsers() ([]entity.User, error)
    CreateUser(user entity.User) (entity.User, error)
}