package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/GydeonZ/task-api/internal/domain"
	"github.com/GydeonZ/task-api/pkg/apperror"
)

type boardRepository struct {
	db *gorm.DB
}

// NewBoardRepository creates a PostgreSQL-backed BoardRepository.
func NewBoardRepository(db *gorm.DB) BoardRepository {
	return &boardRepository{db: db}
}

func (r *boardRepository) Create(ctx context.Context, board *domain.Board) error {
	if err := r.db.WithContext(ctx).Create(board).Error; err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func (r *boardRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Board, error) {
	var board domain.Board
	err := r.db.WithContext(ctx).
		Preload("Owner").
		Preload("Lists").
		Preload("Members.User").
		First(&board, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.NotFound("board")
	}
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return &board, nil
}

func (r *boardRepository) FindByMember(ctx context.Context, userID uuid.UUID) ([]domain.Board, error) {
	var boards []domain.Board
	err := r.db.WithContext(ctx).
		Joins("JOIN board_members ON board_members.board_id = boards.id").
		Where("board_members.user_id = ? AND boards.deleted_at IS NULL", userID).
		Preload("Owner").
		Find(&boards).Error
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return boards, nil
}

func (r *boardRepository) Update(ctx context.Context, board *domain.Board) error {
	if err := r.db.WithContext(ctx).Save(board).Error; err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func (r *boardRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Board{}, "id = ?", id).Error; err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func (r *boardRepository) AddMember(ctx context.Context, member *domain.BoardMember) error {
	if err := r.db.WithContext(ctx).Create(member).Error; err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func (r *boardRepository) RemoveMember(ctx context.Context, boardID, userID uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Delete(&domain.BoardMember{}, "board_id = ? AND user_id = ?", boardID, userID).Error
	if err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func (r *boardRepository) FindMember(ctx context.Context, boardID, userID uuid.UUID) (*domain.BoardMember, error) {
	var member domain.BoardMember
	err := r.db.WithContext(ctx).
		First(&member, "board_id = ? AND user_id = ?", boardID, userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.NotFound("board member")
	}
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return &member, nil
}

func (r *boardRepository) ListMembers(ctx context.Context, boardID uuid.UUID) ([]domain.BoardMember, error) {
	var members []domain.BoardMember
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("board_id = ?", boardID).
		Find(&members).Error
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return members, nil
}
