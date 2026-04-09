package domain

import (
	"time"

	"gorm.io/gorm"
)

type Card struct {
	gorm.Model
	Title       string     `gorm:"not null"`
	Description string
	Position    int        `gorm:"default:0"`
	DueDate     *time.Time
	ListID      uint
}

type CreateCardRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Position    int        `json:"position"`
	DueDate     *time.Time `json:"due_date"`
}

type UpdateCardRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Position    int        `json:"position"`
	DueDate     *time.Time `json:"due_date"`
}

type MoveCardRequest struct {
	ListID   uint `json:"list_id" binding:"required"`
	Position int  `json:"position"`
}
