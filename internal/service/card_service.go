package service

import (
"context"

"github.com/google/uuid"

"github.com/GydeonZ/task-api/internal/domain"
"github.com/GydeonZ/task-api/internal/repository"
"github.com/GydeonZ/task-api/pkg/apperror"
)

// CardService defines business-logic operations for cards.
type CardService interface {
CreateCard(ctx context.Context, userID, listID uuid.UUID, req *domain.CreateCardRequest) (*domain.Card, error)
GetCard(ctx context.Context, userID, cardID uuid.UUID) (*domain.Card, error)
CardsForList(ctx context.Context, userID, listID uuid.UUID) ([]domain.Card, error)
UpdateCard(ctx context.Context, userID, cardID uuid.UUID, req *domain.UpdateCardRequest) (*domain.Card, error)
MoveCard(ctx context.Context, userID, cardID uuid.UUID, req *domain.MoveCardRequest) error
DeleteCard(ctx context.Context, userID, cardID uuid.UUID) error

AddAssignee(ctx context.Context, requesterID, cardID, userID uuid.UUID) error
RemoveAssignee(ctx context.Context, requesterID, cardID, userID uuid.UUID) error

CreateComment(ctx context.Context, userID, cardID uuid.UUID, req *domain.CreateCommentRequest) (*domain.Comment, error)
UpdateComment(ctx context.Context, userID, commentID uuid.UUID, req *domain.CreateCommentRequest) (*domain.Comment, error)
DeleteComment(ctx context.Context, userID, commentID uuid.UUID) error
ListComments(ctx context.Context, userID, cardID uuid.UUID) ([]domain.Comment, error)
}

type cardService struct {
cardRepo  repository.CardRepository
listRepo  repository.ListRepository
boardRepo repository.BoardRepository
}

// NewCardService creates a new CardService.
func NewCardService(
cardRepo repository.CardRepository,
listRepo repository.ListRepository,
boardRepo repository.BoardRepository,
) CardService {
return &cardService{
cardRepo:  cardRepo,
listRepo:  listRepo,
boardRepo: boardRepo,
}
}

// boardIDForList resolves the boardID from a listID.
func (s *cardService) boardIDForList(ctx context.Context, listID uuid.UUID) (uuid.UUID, error) {
list, err := s.listRepo.FindByID(ctx, listID)
if err != nil {
return uuid.Nil, err
}
return list.BoardID, nil
}

// boardIDForCard resolves the boardID from a cardID.
func (s *cardService) boardIDForCard(ctx context.Context, cardID uuid.UUID) (uuid.UUID, error) {
card, err := s.cardRepo.FindByID(ctx, cardID)
if err != nil {
return uuid.Nil, err
}
return s.boardIDForList(ctx, card.ListID)
}

func (s *cardService) requireBoardMember(ctx context.Context, boardID, userID uuid.UUID) error {
_, err := s.boardRepo.FindMember(ctx, boardID, userID)
if err != nil {
return apperror.Forbidden("you do not have access to this board")
}
return nil
}

func (s *cardService) CreateCard(ctx context.Context, userID, listID uuid.UUID, req *domain.CreateCardRequest) (*domain.Card, error) {
boardID, err := s.boardIDForList(ctx, listID)
if err != nil {
return nil, err
}

if err := s.requireBoardMember(ctx, boardID, userID); err != nil {
return nil, err
}

priority := req.Priority
if priority == "" {
priority = domain.CardPriorityMedium
}

card := &domain.Card{
ListID:      listID,
Title:       req.Title,
Description: req.Description,
Position:    req.Position,
Priority:    priority,
DueDate:     req.DueDate,
CreatedBy:   userID,
}

if err := s.cardRepo.Create(ctx, card); err != nil {
return nil, err
}
return card, nil
}

func (s *cardService) GetCard(ctx context.Context, userID, cardID uuid.UUID) (*domain.Card, error) {
card, err := s.cardRepo.FindByID(ctx, cardID)
if err != nil {
return nil, err
}

boardID, err := s.boardIDForList(ctx, card.ListID)
if err != nil {
return nil, err
}

if err := s.requireBoardMember(ctx, boardID, userID); err != nil {
return nil, err
}

return card, nil
}

func (s *cardService) CardsForList(ctx context.Context, userID, listID uuid.UUID) ([]domain.Card, error) {
boardID, err := s.boardIDForList(ctx, listID)
if err != nil {
return nil, err
}

if err := s.requireBoardMember(ctx, boardID, userID); err != nil {
return nil, err
}

return s.cardRepo.FindByList(ctx, listID)
}

func (s *cardService) UpdateCard(ctx context.Context, userID, cardID uuid.UUID, req *domain.UpdateCardRequest) (*domain.Card, error) {
card, err := s.cardRepo.FindByID(ctx, cardID)
if err != nil {
return nil, err
}

boardID, err := s.boardIDForList(ctx, card.ListID)
if err != nil {
return nil, err
}

if err := s.requireBoardMember(ctx, boardID, userID); err != nil {
return nil, err
}

if req.Title != "" {
card.Title = req.Title
}
if req.Description != "" {
card.Description = req.Description
}
if req.Priority != "" {
card.Priority = req.Priority
}
if req.DueDate != nil {
card.DueDate = req.DueDate
}
if req.Position > 0 {
card.Position = req.Position
}

if err := s.cardRepo.Update(ctx, card); err != nil {
return nil, err
}
return card, nil
}

func (s *cardService) MoveCard(ctx context.Context, userID, cardID uuid.UUID, req *domain.MoveCardRequest) error {
boardID, err := s.boardIDForCard(ctx, cardID)
if err != nil {
return err
}

if err := s.requireBoardMember(ctx, boardID, userID); err != nil {
return err
}

return s.cardRepo.Move(ctx, cardID, req.ListID, req.Position)
}

func (s *cardService) DeleteCard(ctx context.Context, userID, cardID uuid.UUID) error {
boardID, err := s.boardIDForCard(ctx, cardID)
if err != nil {
return err
}

if err := s.requireBoardMember(ctx, boardID, userID); err != nil {
return err
}

return s.cardRepo.Delete(ctx, cardID)
}

func (s *cardService) AddAssignee(ctx context.Context, requesterID, cardID, userID uuid.UUID) error {
boardID, err := s.boardIDForCard(ctx, cardID)
if err != nil {
return err
}

if err := s.requireBoardMember(ctx, boardID, requesterID); err != nil {
return err
}

return s.cardRepo.AddAssignee(ctx, cardID, userID)
}

func (s *cardService) RemoveAssignee(ctx context.Context, requesterID, cardID, userID uuid.UUID) error {
boardID, err := s.boardIDForCard(ctx, cardID)
if err != nil {
return err
}

if err := s.requireBoardMember(ctx, boardID, requesterID); err != nil {
return err
}

return s.cardRepo.RemoveAssignee(ctx, cardID, userID)
}

func (s *cardService) CreateComment(ctx context.Context, userID, cardID uuid.UUID, req *domain.CreateCommentRequest) (*domain.Comment, error) {
boardID, err := s.boardIDForCard(ctx, cardID)
if err != nil {
return nil, err
}

if err := s.requireBoardMember(ctx, boardID, userID); err != nil {
return nil, err
}

comment := &domain.Comment{
CardID: cardID,
UserID: userID,
Body:   req.Body,
}

if err := s.cardRepo.CreateComment(ctx, comment); err != nil {
return nil, err
}
return comment, nil
}

func (s *cardService) UpdateComment(ctx context.Context, userID, commentID uuid.UUID, req *domain.CreateCommentRequest) (*domain.Comment, error) {
comment, err := s.cardRepo.FindCommentByID(ctx, commentID)
if err != nil {
return nil, err
}

if comment.UserID != userID {
return nil, apperror.Forbidden("you can only edit your own comments")
}

comment.Body = req.Body
if err := s.cardRepo.UpdateComment(ctx, comment); err != nil {
return nil, err
}
return comment, nil
}

func (s *cardService) DeleteComment(ctx context.Context, userID, commentID uuid.UUID) error {
comment, err := s.cardRepo.FindCommentByID(ctx, commentID)
if err != nil {
return err
}

if comment.UserID != userID {
return apperror.Forbidden("you can only delete your own comments")
}

return s.cardRepo.DeleteComment(ctx, commentID)
}

func (s *cardService) ListComments(ctx context.Context, userID, cardID uuid.UUID) ([]domain.Comment, error) {
boardID, err := s.boardIDForCard(ctx, cardID)
if err != nil {
return nil, err
}

if err := s.requireBoardMember(ctx, boardID, userID); err != nil {
return nil, err
}

return s.cardRepo.ListComments(ctx, cardID)
}
