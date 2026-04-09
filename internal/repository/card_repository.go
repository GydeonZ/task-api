package repository

import (
"context"
"errors"

"github.com/google/uuid"
"gorm.io/gorm"

"github.com/GydeonZ/task-api/internal/domain"
"github.com/GydeonZ/task-api/pkg/apperror"
)

type cardRepository struct {
db *gorm.DB
}

// NewCardRepository creates a PostgreSQL-backed CardRepository.
func NewCardRepository(db *gorm.DB) CardRepository {
return &cardRepository{db: db}
}

func (r *cardRepository) Create(ctx context.Context, card *domain.Card) error {
if err := r.db.WithContext(ctx).Create(card).Error; err != nil {
return apperror.Internal(err)
}
return nil
}

func (r *cardRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Card, error) {
var card domain.Card
err := r.db.WithContext(ctx).
Preload("Labels").
Preload("Assignees").
Preload("Comments.User").
Preload("Checklists.Items").
Preload("Attachments").
First(&card, "id = ?", id).Error
if errors.Is(err, gorm.ErrRecordNotFound) {
return nil, apperror.NotFound("card")
}
if err != nil {
return nil, apperror.Internal(err)
}
return &card, nil
}

func (r *cardRepository) FindByList(ctx context.Context, listID uuid.UUID) ([]domain.Card, error) {
var cards []domain.Card
err := r.db.WithContext(ctx).
Where("list_id = ?", listID).
Order("position ASC").
Preload("Labels").
Preload("Assignees").
Find(&cards).Error
if err != nil {
return nil, apperror.Internal(err)
}
return cards, nil
}

func (r *cardRepository) Update(ctx context.Context, card *domain.Card) error {
if err := r.db.WithContext(ctx).Save(card).Error; err != nil {
return apperror.Internal(err)
}
return nil
}

func (r *cardRepository) Delete(ctx context.Context, id uuid.UUID) error {
if err := r.db.WithContext(ctx).Delete(&domain.Card{}, "id = ?", id).Error; err != nil {
return apperror.Internal(err)
}
return nil
}

func (r *cardRepository) Move(ctx context.Context, cardID, targetListID uuid.UUID, position int) error {
err := r.db.WithContext(ctx).
Model(&domain.Card{}).
Where("id = ?", cardID).
Updates(map[string]interface{}{
"list_id":  targetListID,
"position": position,
}).Error
if err != nil {
return apperror.Internal(err)
}
return nil
}

func (r *cardRepository) AddAssignee(ctx context.Context, cardID, userID uuid.UUID) error {
card := domain.Card{}
card.ID = cardID
user := domain.User{}
user.ID = userID
if err := r.db.WithContext(ctx).Model(&card).Association("Assignees").Append(&user); err != nil {
return apperror.Internal(err)
}
return nil
}

func (r *cardRepository) RemoveAssignee(ctx context.Context, cardID, userID uuid.UUID) error {
card := domain.Card{}
card.ID = cardID
user := domain.User{}
user.ID = userID
if err := r.db.WithContext(ctx).Model(&card).Association("Assignees").Delete(&user); err != nil {
return apperror.Internal(err)
}
return nil
}

func (r *cardRepository) CreateComment(ctx context.Context, comment *domain.Comment) error {
if err := r.db.WithContext(ctx).Create(comment).Error; err != nil {
return apperror.Internal(err)
}
return nil
}

func (r *cardRepository) FindCommentByID(ctx context.Context, id uuid.UUID) (*domain.Comment, error) {
var comment domain.Comment
err := r.db.WithContext(ctx).Preload("User").First(&comment, "id = ?", id).Error
if errors.Is(err, gorm.ErrRecordNotFound) {
return nil, apperror.NotFound("comment")
}
if err != nil {
return nil, apperror.Internal(err)
}
return &comment, nil
}

func (r *cardRepository) UpdateComment(ctx context.Context, comment *domain.Comment) error {
if err := r.db.WithContext(ctx).Save(comment).Error; err != nil {
return apperror.Internal(err)
}
return nil
}

func (r *cardRepository) DeleteComment(ctx context.Context, id uuid.UUID) error {
if err := r.db.WithContext(ctx).Delete(&domain.Comment{}, "id = ?", id).Error; err != nil {
return apperror.Internal(err)
}
return nil
}

func (r *cardRepository) ListComments(ctx context.Context, cardID uuid.UUID) ([]domain.Comment, error) {
var comments []domain.Comment
err := r.db.WithContext(ctx).
Preload("User").
Where("card_id = ?", cardID).
Order("created_at ASC").
Find(&comments).Error
if err != nil {
return nil, apperror.Internal(err)
}
return comments, nil
}
