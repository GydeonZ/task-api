package domain

import (
	"time"

	"github.com/google/uuid"
)

// BoardVisibility controls who can see the board.
type BoardVisibility string

const (
	BoardVisibilityPrivate BoardVisibility = "private"
	BoardVisibilityPublic  BoardVisibility = "public"
)

// Board represents a Trello-like project board.
type Board struct {
	ID          uuid.UUID       `json:"id"          gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OwnerID     uuid.UUID       `json:"owner_id"    gorm:"type:uuid;not null;index"`
	Title       string          `json:"title"       gorm:"not null"`
	Description string          `json:"description"`
	Visibility  BoardVisibility `json:"visibility"  gorm:"type:varchar(10);not null;default:'private'"`
	Background  string          `json:"background"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   *time.Time      `json:"deleted_at,omitempty" gorm:"index"`

	Owner   User          `json:"owner,omitempty"   gorm:"foreignKey:OwnerID"`
	Lists   []List        `json:"lists,omitempty"   gorm:"foreignKey:BoardID"`
	Members []BoardMember `json:"members,omitempty" gorm:"foreignKey:BoardID"`
}

// BoardMemberRole defines the role a member has on a board.
type BoardMemberRole string

const (
	BoardMemberRoleAdmin  BoardMemberRole = "admin"
	BoardMemberRoleMember BoardMemberRole = "member"
	BoardMemberRoleViewer BoardMemberRole = "viewer"
)

// BoardMember associates a User with a Board with a specific role.
type BoardMember struct {
	ID        uuid.UUID       `json:"id"         gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	BoardID   uuid.UUID       `json:"board_id"   gorm:"type:uuid;not null;index"`
	UserID    uuid.UUID       `json:"user_id"    gorm:"type:uuid;not null;index"`
	Role      BoardMemberRole `json:"role"       gorm:"type:varchar(10);not null;default:'member'"`
	CreatedAt time.Time       `json:"created_at"`

	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// CreateBoardRequest holds the input for creating a board.
type CreateBoardRequest struct {
	Title       string          `json:"title"       validate:"required,min=1,max=200"`
	Description string          `json:"description" validate:"omitempty,max=1000"`
	Visibility  BoardVisibility `json:"visibility"  validate:"omitempty,oneof=private public"`
	Background  string          `json:"background"  validate:"omitempty"`
}

// UpdateBoardRequest holds the input for updating a board.
type UpdateBoardRequest struct {
	Title       string          `json:"title"       validate:"omitempty,min=1,max=200"`
	Description string          `json:"description" validate:"omitempty,max=1000"`
	Visibility  BoardVisibility `json:"visibility"  validate:"omitempty,oneof=private public"`
	Background  string          `json:"background"  validate:"omitempty"`
}

// AddBoardMemberRequest holds the input for adding a member to a board.
type AddBoardMemberRequest struct {
	Email string          `json:"email" validate:"required,email"`
	Role  BoardMemberRole `json:"role"  validate:"required,oneof=admin member viewer"`
}
