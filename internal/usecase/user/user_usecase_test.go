package user

import (
	"errors"
	"testing"

	"github.com/ipxsandbox/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Repository
type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) FindAll() ([]entity.User, error) {
	args := m.Called()
	return args.Get(0).([]entity.User), args.Error(1)
}

func (m *mockUserRepo) Create(user entity.User) (entity.User, error) {
	args := m.Called(user)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *mockUserRepo) FindByEmail(email string) (entity.User, error) {
	args := m.Called(email)
	return args.Get(0).(entity.User), args.Error(1)
}

func TestGetAllUsers(t *testing.T) {
	mockRepo := new(mockUserRepo)
	mockUsers := []entity.User{{ID: 1, Name: "Alice", Email: "alice@example.com"}}
	mockRepo.On("FindAll").Return(mockUsers, nil)

	uc := NewUserUsecase(mockRepo)
	users, err := uc.GetAllUsers()
	assert.NoError(t, err)
	assert.Equal(t, mockUsers, users)

	mockRepo.AssertExpectations(t)
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(mockUserRepo)
	inputUser := entity.User{Name: "Bob", Email: "bob@example.com"}
	returnUser := entity.User{ID: 2, Name: "Bob", Email: "bob@example.com"}

	mockRepo.On("Create", inputUser).Return(returnUser, nil)

	uc := NewUserUsecase(mockRepo)
	user, err := uc.CreateUser(inputUser)
	assert.NoError(t, err)
	assert.Equal(t, returnUser, user)

	mockRepo.AssertExpectations(t)
}

func TestCreateUser_Error(t *testing.T) {
	mockRepo := new(mockUserRepo)
	inputUser := entity.User{Name: "Bob", Email: "bob@example.com"}

	mockRepo.On("Create", inputUser).Return(entity.User{}, errors.New("create error"))

	uc := NewUserUsecase(mockRepo)
	user, err := uc.CreateUser(inputUser)
	assert.Error(t, err)
	assert.Equal(t, entity.User{}, user)

	mockRepo.AssertExpectations(t)
}