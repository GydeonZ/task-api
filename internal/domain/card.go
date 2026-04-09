package domain

import (
	"time"

	"github.com/google/uuid"
)

// CardPriority defines the urgency level of a card.
type CardPriority string

const (
	CardPriorityLow    CardPriority = "low"
	CardPriorityMedium CardPriority = "medium"
	CardPriorityHigh   CardPriority = "high"
	CardPriorityUrgent CardPriority = "urgent"
)

// Card represents a task card inside a list.
type Card struct {
	ID          uuid.UUID    `json:"id"           gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ListID      uuid.UUID    `json:"list_id"      gorm:"type:uuid;not null;index"`
	Title       string       `json:"title"        gorm:"not null"`
	Description string       `json:"description"`
	Position    int          `json:"position"     gorm:"not null;default:0"`
	Priority    CardPriority `json:"priority"     gorm:"type:varchar(10);not null;default:'medium'"`
	DueDate     *time.Time   `json:"due_date,omitempty"`
	CreatedBy   uuid.UUID    `json:"created_by"   gorm:"type:uuid;not null"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   *time.Time   `json:"deleted_at,omitempty" gorm:"index"`

	Labels      []Label      `json:"labels,omitempty"      gorm:"many2many:card_labels;"`
	Assignees   []User       `json:"assignees,omitempty"   gorm:"many2many:card_assignees;"`
	Comments    []Comment    `json:"comments,omitempty"    gorm:"foreignKey:CardID"`
	Checklists  []Checklist  `json:"checklists,omitempty"  gorm:"foreignKey:CardID"`
	Attachments []Attachment `json:"attachments,omitempty" gorm:"foreignKey:CardID"`
}

// Label is a coloured tag that can be applied to cards.
type Label struct {
	ID      uuid.UUID `json:"id"      gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	BoardID uuid.UUID `json:"board_id" gorm:"type:uuid;not null;index"`
	Name    string    `json:"name"    gorm:"not null"`
	Color   string    `json:"color"   gorm:"not null;default:'#61bd4f'"`
}

// Comment is a user comment on a card.
type Comment struct {
	ID        uuid.UUID  `json:"id"         gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CardID    uuid.UUID  `json:"card_id"    gorm:"type:uuid;not null;index"`
	UserID    uuid.UUID  `json:"user_id"    gorm:"type:uuid;not null"`
	Body      string     `json:"body"       gorm:"not null"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// Checklist is a group of to-do items on a card.
type Checklist struct {
	ID        uuid.UUID       `json:"id"      gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CardID    uuid.UUID       `json:"card_id" gorm:"type:uuid;not null;index"`
	Title     string          `json:"title"   gorm:"not null"`
	Items     []ChecklistItem `json:"items,omitempty" gorm:"foreignKey:ChecklistID"`
	CreatedAt time.Time       `json:"created_at"`
}

// ChecklistItem is a single item inside a Checklist.
type ChecklistItem struct {
	ID          uuid.UUID `json:"id"           gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ChecklistID uuid.UUID `json:"checklist_id" gorm:"type:uuid;not null;index"`
	Title       string    `json:"title"        gorm:"not null"`
	Done        bool      `json:"done"         gorm:"default:false"`
	Position    int       `json:"position"     gorm:"default:0"`
}

// Attachment is a file or link attached to a card.
type Attachment struct {
	ID        uuid.UUID `json:"id"         gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CardID    uuid.UUID `json:"card_id"    gorm:"type:uuid;not null;index"`
	Name      string    `json:"name"       gorm:"not null"`
	URL       string    `json:"url"        gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateCardRequest holds the input for creating a card.
type CreateCardRequest struct {
	Title       string       `json:"title"       validate:"required,min=1,max=500"`
	Description string       `json:"description" validate:"omitempty,max=5000"`
	Position    int          `json:"position"    validate:"omitempty,min=0"`
	Priority    CardPriority `json:"priority"    validate:"omitempty,oneof=low medium high urgent"`
	DueDate     *time.Time   `json:"due_date"`
}

// UpdateCardRequest holds the input for updating a card.
type UpdateCardRequest struct {
	Title       string       `json:"title"       validate:"omitempty,min=1,max=500"`
	Description string       `json:"description" validate:"omitempty,max=5000"`
	Position    int          `json:"position"    validate:"omitempty,min=0"`
	Priority    CardPriority `json:"priority"    validate:"omitempty,oneof=low medium high urgent"`
	DueDate     *time.Time   `json:"due_date"`
	ListID      *uuid.UUID   `json:"list_id"     validate:"omitempty"`
}

// MoveCardRequest holds the input for moving a card to another list.
type MoveCardRequest struct {
	ListID   uuid.UUID `json:"list_id"  validate:"required"`
	Position int       `json:"position" validate:"min=0"`
}

// CreateCommentRequest holds the input for creating a comment.
type CreateCommentRequest struct {
	Body string `json:"body" validate:"required,min=1,max=5000"`
}
