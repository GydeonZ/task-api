package repository

import (
	"errors"

	"github.com/GydeonZ/task-api/internal/domain"
	"github.com/GydeonZ/task-api/pkg/apperrors"
	"gorm.io/gorm"
)

type cardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) CardRepository {
	return &cardRepository{db: db}
}

func (r *cardRepository) Create(card *domain.Card) error {
	return r.db.Create(card).Error
}

func (r *cardRepository) FindByListID(listID uint) ([]domain.Card, error) {
	var cards []domain.Card
	result := r.db.Where("list_id = ?", listID).Order("position asc").Find(&cards)
	return cards, result.Error
}

func (r *cardRepository) FindByID(id uint) (*domain.Card, error) {
	var card domain.Card
	result := r.db.First(&card, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, result.Error
	}
	return &card, nil
}

func (r *cardRepository) Update(card *domain.Card) error {
	return r.db.Save(card).Error
}

func (r *cardRepository) Delete(id uint) error {
	result := r.db.Delete(&domain.Card{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *cardRepository) MoveToList(cardID, listID uint, position int) error {
	result := r.db.Model(&domain.Card{}).Where("id = ?", cardID).Updates(map[string]interface{}{
		"list_id":  listID,
		"position": position,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}
