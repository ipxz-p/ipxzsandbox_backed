package auth_usercase

import (
	"github.com/ipxsandbox/internal/entity"
	"github.com/ipxsandbox/internal/pkg/jwtutil"
	userRepository "github.com/ipxsandbox/internal/repository/user"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecaseInterface interface {
	Register(user entity.User) (entity.UserResponse, error)
	Login(email string, password string) (accessToken string, refreshToken string, err error)
}

type authUsecase struct {
	userRepo userRepository.Repository
}

func NewAuthUsecase(repo userRepository.Repository) AuthUsecaseInterface {
	return &authUsecase{userRepo: repo}
}

func (uc *authUsecase) Register(user entity.User) (entity.UserResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return entity.UserResponse{}, err
	}
	user.Password = string(hashed)

	createdUser, err := uc.userRepo.Create(user)
	if err != nil {
		return entity.UserResponse{}, err
	}

	return entity.UserResponse{
		ID:    createdUser.ID,
		Name:  createdUser.Name,
		Email: createdUser.Email,
	}, nil
}

func (uc *authUsecase) Login(email, password string) (string, string, error) {
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err := jwtutil.GenerateTokens(user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}