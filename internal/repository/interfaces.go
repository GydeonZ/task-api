package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/GydeonZ/task-api/internal/domain"
)

// UserRepository defines persistence operations for User entities.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// BoardRepository defines persistence operations for Board entities.
type BoardRepository interface {
	Create(ctx context.Context, board *domain.Board) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Board, error)
	FindByMember(ctx context.Context, userID uuid.UUID) ([]domain.Board, error)
	Update(ctx context.Context, board *domain.Board) error
	Delete(ctx context.Context, id uuid.UUID) error

	AddMember(ctx context.Context, member *domain.BoardMember) error
	RemoveMember(ctx context.Context, boardID, userID uuid.UUID) error
	FindMember(ctx context.Context, boardID, userID uuid.UUID) (*domain.BoardMember, error)
	ListMembers(ctx context.Context, boardID uuid.UUID) ([]domain.BoardMember, error)
}

// ListRepository defines persistence operations for List entities.
type ListRepository interface {
	Create(ctx context.Context, list *domain.List) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.List, error)
	FindByBoard(ctx context.Context, boardID uuid.UUID) ([]domain.List, error)
	Update(ctx context.Context, list *domain.List) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// CardRepository defines persistence operations for Card entities.
type CardRepository interface {
	Create(ctx context.Context, card *domain.Card) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Card, error)
	FindByList(ctx context.Context, listID uuid.UUID) ([]domain.Card, error)
	Update(ctx context.Context, card *domain.Card) error
	Delete(ctx context.Context, id uuid.UUID) error
	Move(ctx context.Context, cardID, targetListID uuid.UUID, position int) error

	AddAssignee(ctx context.Context, cardID, userID uuid.UUID) error
	RemoveAssignee(ctx context.Context, cardID, userID uuid.UUID) error

	CreateComment(ctx context.Context, comment *domain.Comment) error
	FindCommentByID(ctx context.Context, id uuid.UUID) (*domain.Comment, error)
	UpdateComment(ctx context.Context, comment *domain.Comment) error
	DeleteComment(ctx context.Context, id uuid.UUID) error
	ListComments(ctx context.Context, cardID uuid.UUID) ([]domain.Comment, error)
}
