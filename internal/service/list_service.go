package service

import (
"context"

"github.com/google/uuid"

"github.com/GydeonZ/task-api/internal/domain"
"github.com/GydeonZ/task-api/internal/repository"
"github.com/GydeonZ/task-api/pkg/apperror"
)

// ListService defines business-logic operations for lists.
type ListService interface {
CreateList(ctx context.Context, userID, boardID uuid.UUID, req *domain.CreateListRequest) (*domain.List, error)
GetList(ctx context.Context, userID, listID uuid.UUID) (*domain.List, error)
ListsForBoard(ctx context.Context, userID, boardID uuid.UUID) ([]domain.List, error)
UpdateList(ctx context.Context, userID, listID uuid.UUID, req *domain.UpdateListRequest) (*domain.List, error)
DeleteList(ctx context.Context, userID, listID uuid.UUID) error
}

type listService struct {
listRepo  repository.ListRepository
boardRepo repository.BoardRepository
}

// NewListService creates a new ListService.
func NewListService(listRepo repository.ListRepository, boardRepo repository.BoardRepository) ListService {
return &listService{listRepo: listRepo, boardRepo: boardRepo}
}

func (s *listService) CreateList(ctx context.Context, userID, boardID uuid.UUID, req *domain.CreateListRequest) (*domain.List, error) {
if err := s.requireBoardMember(ctx, boardID, userID); err != nil {
return nil, err
}

list := &domain.List{
BoardID:  boardID,
Title:    req.Title,
Position: req.Position,
}

if err := s.listRepo.Create(ctx, list); err != nil {
return nil, err
}
return list, nil
}

func (s *listService) GetList(ctx context.Context, userID, listID uuid.UUID) (*domain.List, error) {
list, err := s.listRepo.FindByID(ctx, listID)
if err != nil {
return nil, err
}

if err := s.requireBoardMember(ctx, list.BoardID, userID); err != nil {
return nil, err
}

return list, nil
}

func (s *listService) ListsForBoard(ctx context.Context, userID, boardID uuid.UUID) ([]domain.List, error) {
if err := s.requireBoardMember(ctx, boardID, userID); err != nil {
return nil, err
}
return s.listRepo.FindByBoard(ctx, boardID)
}

func (s *listService) UpdateList(ctx context.Context, userID, listID uuid.UUID, req *domain.UpdateListRequest) (*domain.List, error) {
list, err := s.listRepo.FindByID(ctx, listID)
if err != nil {
return nil, err
}

if err := s.requireBoardMember(ctx, list.BoardID, userID); err != nil {
return nil, err
}

if req.Title != "" {
list.Title = req.Title
}
if req.Position >= 0 {
list.Position = req.Position
}

if err := s.listRepo.Update(ctx, list); err != nil {
return nil, err
}
return list, nil
}

func (s *listService) DeleteList(ctx context.Context, userID, listID uuid.UUID) error {
list, err := s.listRepo.FindByID(ctx, listID)
if err != nil {
return err
}

if err := s.requireBoardMember(ctx, list.BoardID, userID); err != nil {
return err
}

return s.listRepo.Delete(ctx, listID)
}

func (s *listService) requireBoardMember(ctx context.Context, boardID, userID uuid.UUID) error {
_, err := s.boardRepo.FindMember(ctx, boardID, userID)
if err != nil {
return apperror.Forbidden("you do not have access to this board")
}
return nil
}
