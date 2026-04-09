package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string  `gorm:"not null"`
	Email    string  `gorm:"uniqueIndex;not null"`
	Password string  `gorm:"not null"`
	Boards   []Board
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
