package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a registered user in the system.
type User struct {
	ID        uuid.UUID  `json:"id"         gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string     `json:"name"       gorm:"not null"`
	Email     string     `json:"email"      gorm:"uniqueIndex;not null"`
	Password  string     `json:"-"          gorm:"not null"`
	Avatar    string     `json:"avatar"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// RegisterRequest holds the input for user registration.
type RegisterRequest struct {
	Name     string `json:"name"     validate:"required,min=2,max=100"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginRequest holds the input for user login.
type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdateProfileRequest holds the input for updating a user's profile.
type UpdateProfileRequest struct {
	Name   string `json:"name"   validate:"omitempty,min=2,max=100"`
	Avatar string `json:"avatar" validate:"omitempty,url"`
}
