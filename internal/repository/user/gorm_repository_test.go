package user

import (
	"testing"

	"github.com/ipxsandbox/internal/entity"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&entity.User{})
	assert.NoError(t, err)

	return db
}

func TestCreateAndFindAll(t *testing.T) {
	db := setupTestDB(t)
	repo := New(db)

	user := entity.User{Name: "Test User", Email: "test@example.com"}
	createdUser, err := repo.Create(user)
	assert.NoError(t, err)
	assert.NotZero(t, createdUser.ID)

	users, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "Test User", users[0].Name)
}