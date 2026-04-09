package repository

import (
	"errors"

	"github.com/GydeonZ/task-api/internal/domain"
	"github.com/GydeonZ/task-api/pkg/apperrors"
	"gorm.io/gorm"
)

type boardRepository struct {
	db *gorm.DB
}

func NewBoardRepository(db *gorm.DB) BoardRepository {
	return &boardRepository{db: db}
}

func (r *boardRepository) Create(board *domain.Board) error {
	return r.db.Create(board).Error
}

func (r *boardRepository) FindByUserID(userID uint) ([]domain.Board, error) {
	var boards []domain.Board
	result := r.db.Where("user_id = ?", userID).Find(&boards)
	return boards, result.Error
}

func (r *boardRepository) FindByID(id uint) (*domain.Board, error) {
	var board domain.Board
	result := r.db.Preload("Lists").First(&board, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, result.Error
	}
	return &board, nil
}

func (r *boardRepository) Update(board *domain.Board) error {
	return r.db.Save(board).Error
}

func (r *boardRepository) Delete(id uint) error {
	result := r.db.Delete(&domain.Board{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}
