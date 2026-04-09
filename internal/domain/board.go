package domain

import "gorm.io/gorm"

type Board struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string
	UserID      uint
	User        User
	Lists       []List
}

type CreateBoardRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type UpdateBoardRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
