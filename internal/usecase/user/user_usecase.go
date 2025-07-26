package user

import (
	"github.com/ipxsandbox/internal/entity"
	"github.com/ipxsandbox/internal/repository/user"
)

type usecase struct {
	repo user.Repository
}

func NewUserUsecase(repo user.Repository) Usecase {
	return &usecase{repo: repo}
}

func (u *usecase) GetAllUsers() ([]entity.User, error) {
	return u.repo.FindAll()
}

func (u *usecase) CreateUser(user entity.User) (entity.User, error) {
	return u.repo.Create(user)
}
