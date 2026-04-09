package domain

import "gorm.io/gorm"

type List struct {
	gorm.Model
	Title    string `gorm:"not null"`
	Position int    `gorm:"default:0"`
	BoardID  uint
	Cards    []Card
}

type CreateListRequest struct {
	Title    string `json:"title" binding:"required"`
	Position int    `json:"position"`
}

type UpdateListRequest struct {
	Title    string `json:"title"`
	Position int    `json:"position"`
}
