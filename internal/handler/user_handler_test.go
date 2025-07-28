package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ipxsandbox/internal/entity"
	userUsecase "github.com/ipxsandbox/internal/usecase/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Usecase
type mockUserUsecase struct {
	mock.Mock
}

func (m *mockUserUsecase) GetAllUsers() ([]entity.User, error) {
	args := m.Called()
	return args.Get(0).([]entity.User), args.Error(1)
}

func (m *mockUserUsecase) CreateUser(user entity.User) (entity.User, error) {
	args := m.Called(user)
	return args.Get(0).(entity.User), args.Error(1)
}

func setupRouter(uc userUsecase.Usecase) *gin.Engine {
	handler := NewUserHandler(uc)
	r := gin.Default()
	r.GET("/users", handler.GetUsers)
	r.POST("/users", handler.CreateUser)
	return r
}

func TestGetUsersHandler_Success(t *testing.T) {
	mockUC := new(mockUserUsecase)
	mockUsers := []entity.User{{ID: 1, Name: "Alice", Email: "alice@example.com"}}
	mockUC.On("GetAllUsers").Return(mockUsers, nil)

	r := setupRouter(mockUC)

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var users []entity.User
	err := json.Unmarshal(w.Body.Bytes(), &users)
	assert.NoError(t, err)
	assert.Equal(t, mockUsers, users)

	mockUC.AssertExpectations(t)
}

func TestGetUsersHandler_Error(t *testing.T) {
	mockUC := new(mockUserUsecase)
	mockUC.On("GetAllUsers").Return(nil, errors.New("some error"))

	r := setupRouter(mockUC)

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}

func TestCreateUserHandler_Success(t *testing.T) {
	mockUC := new(mockUserUsecase)
	inputUser := entity.User{Name: "Bob", Email: "bob@example.com"}
	returnUser := entity.User{ID: 2, Name: "Bob", Email: "bob@example.com"}

	mockUC.On("CreateUser", inputUser).Return(returnUser, nil)

	r := setupRouter(mockUC)

	body, _ := json.Marshal(inputUser)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var createdUser entity.User
	err := json.Unmarshal(w.Body.Bytes(), &createdUser)
	assert.NoError(t, err)
	assert.Equal(t, returnUser, createdUser)

	mockUC.AssertExpectations(t)
}

func TestCreateUserHandler_BadRequest(t *testing.T) {
	mockUC := new(mockUserUsecase)

	r := setupRouter(mockUC)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer([]byte(`invalid json`)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}