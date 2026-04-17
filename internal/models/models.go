package models

import "time"

// User represents a user in the system
type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password 	string    `json:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Board represents a task board (like a project in Trello)
type Board struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	Slug      string    `json:"slug" db:"slug"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// List represents a list on a board (like columns in Trello)
type List struct {
	ID        int       `json:"id" db:"id"`
	BoardID   int       `json:"board_id" db:"board_id"`
	Title     string    `json:"title" db:"title"`
	Position  int       `json:"position" db:"position"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Task represents a task card (like cards in Trello)
type Task struct {
	ID          int        `json:"id" db:"id"`
	ListID      int        `json:"list_id" db:"list_id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	Position    int        `json:"position" db:"position"`
	IsCompleted bool       `json:"is_completed" db:"is_completed"`
	DueDate     *time.Time `json:"due_date" db:"due_date"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// LoginRequest represents the request to log in a user
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=255"`
}

// CreateUserRequest represents the request to create a new user
type RegisterUserRequest struct {
	Username string `json:"username" validate:"required, min=5,max=15"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=255,regexp=^(?=.*[0-9])(?=.*[^a-zA-Z0-9]).+$"`
}

// CreateBoardRequest represents the request to create a board
type CreateBoardRequest struct {
	Title string `json:"title" validate:"required,min=1,max=255"`
}

// UpdateBoardRequest represents the request to update a board
type UpdateBoardRequest struct {
	Title string `json:"title" validate:"min=1,max=255"`
}

// CreateListRequest represents the request to create a list
type CreateListRequest struct {
	Title string `json:"title" validate:"required,min=1,max=255"`
}

// UpdateListRequest represents the request to update a list
type UpdateListRequest struct {
	Title    string `json:"title" validate:"min=1,max=255"`
	Position int    `json:"position"`
}

// CreateTaskRequest represents the request to create a task
type CreateTaskRequest struct {
	Title       string     `json:"title" validate:"required,min=1,max=500"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
}

// UpdateTaskRequest represents the request to update a task
type UpdateTaskRequest struct {
	Title       string     `json:"title" validate:"min=1,max=500"`
	Description string     `json:"description"`
	IsCompleted bool       `json:"is_completed"`
	DueDate     *time.Time `json:"due_date"`
	Position    int        `json:"position"`
}

// APIResponse is a standard API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Page     int `query:"page" validate:"min=1"`
	PageSize int `query:"page_size" validate:"min=1,max=100"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int64       `json:"total_pages"`
	Data       interface{} `json:"data"`
}