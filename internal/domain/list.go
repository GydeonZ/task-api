package domain

import (
	"time"

	"github.com/google/uuid"
)

// List represents a column inside a board (e.g., "To Do", "In Progress", "Done").
type List struct {
	ID        uuid.UUID  `json:"id"         gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	BoardID   uuid.UUID  `json:"board_id"   gorm:"type:uuid;not null;index"`
	Title     string     `json:"title"      gorm:"not null"`
	Position  int        `json:"position"   gorm:"not null;default:0"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	Cards []Card `json:"cards,omitempty" gorm:"foreignKey:ListID"`
}

// CreateListRequest holds the input for creating a list.
type CreateListRequest struct {
	Title    string `json:"title"    validate:"required,min=1,max=200"`
	Position int    `json:"position" validate:"omitempty,min=0"`
}

// UpdateListRequest holds the input for updating a list.
type UpdateListRequest struct {
	Title    string `json:"title"    validate:"omitempty,min=1,max=200"`
	Position int    `json:"position" validate:"omitempty,min=0"`
}
