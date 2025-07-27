package entity

type User struct {
    ID       uint   `json:"id" gorm:"primaryKey"`
    Name     string `json:"name" gorm:"not null" validate:"required,max=20"`
    Email    string `json:"email" gorm:"unique;not null" validate:"required,email"`
    Password string `json:"password" gorm:"not null" validate:"required,min=8,password"`
}

type UserResponse struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}