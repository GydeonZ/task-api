package repository

import (
	"errors"

	"github.com/GydeonZ/task-api/internal/domain"
	"github.com/GydeonZ/task-api/pkg/apperrors"
	"gorm.io/gorm"
)

type listRepository struct {
	db *gorm.DB
}

func NewListRepository(db *gorm.DB) ListRepository {
	return &listRepository{db: db}
}

func (r *listRepository) Create(list *domain.List) error {
	return r.db.Create(list).Error
}

func (r *listRepository) FindByBoardID(boardID uint) ([]domain.List, error) {
	var lists []domain.List
	result := r.db.Where("board_id = ?", boardID).Order("position asc").Find(&lists)
	return lists, result.Error
}

func (r *listRepository) FindByID(id uint) (*domain.List, error) {
	var list domain.List
	result := r.db.First(&list, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, result.Error
	}
	return &list, nil
}

func (r *listRepository) Update(list *domain.List) error {
	return r.db.Save(list).Error
}

func (r *listRepository) Delete(id uint) error {
	result := r.db.Delete(&domain.List{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}
