package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/GydeonZ/task-api/internal/domain"
	"github.com/GydeonZ/task-api/pkg/apperror"
)

type listRepository struct {
	db *gorm.DB
}

// NewListRepository creates a PostgreSQL-backed ListRepository.
func NewListRepository(db *gorm.DB) ListRepository {
	return &listRepository{db: db}
}

func (r *listRepository) Create(ctx context.Context, list *domain.List) error {
	if err := r.db.WithContext(ctx).Create(list).Error; err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func (r *listRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.List, error) {
	var list domain.List
	err := r.db.WithContext(ctx).
		Preload("Cards").
		First(&list, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.NotFound("list")
	}
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return &list, nil
}

func (r *listRepository) FindByBoard(ctx context.Context, boardID uuid.UUID) ([]domain.List, error) {
	var lists []domain.List
	err := r.db.WithContext(ctx).
		Where("board_id = ?", boardID).
		Order("position ASC").
		Find(&lists).Error
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return lists, nil
}

func (r *listRepository) Update(ctx context.Context, list *domain.List) error {
	if err := r.db.WithContext(ctx).Save(list).Error; err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func (r *listRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&domain.List{}, "id = ?", id).Error; err != nil {
		return apperror.Internal(err)
	}
	return nil
}
