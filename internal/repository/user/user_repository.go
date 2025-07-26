package user

import (
    "gorm.io/gorm"
    "github.com/ipxsandbox/internal/entity"
)

type gormRepository struct {
    db *gorm.DB
}

func New(db *gorm.DB) Repository {
    return &gormRepository{db: db}
}

func (r *gormRepository) FindAll() ([]entity.User, error) {
    var users []entity.User
    err := r.db.Find(&users).Error
    return users, err
}

func (r *gormRepository) Create(user entity.User) (entity.User, error) {
    err := r.db.Create(&user).Error
    return user, err
}

func (r *gormRepository) FindByEmail(email string) (entity.User, error) {
    var user entity.User
    err := r.db.Where("email = ?", email).First(&user).Error
    return user, err
}