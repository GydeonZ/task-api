package service

import (
"context"

"github.com/google/uuid"

"github.com/GydeonZ/task-api/internal/domain"
"github.com/GydeonZ/task-api/internal/repository"
"github.com/GydeonZ/task-api/pkg/apperror"
)

// BoardService defines business-logic operations for boards.
type BoardService interface {
CreateBoard(ctx context.Context, ownerID uuid.UUID, req *domain.CreateBoardRequest) (*domain.Board, error)
GetBoard(ctx context.Context, userID, boardID uuid.UUID) (*domain.Board, error)
ListBoards(ctx context.Context, userID uuid.UUID) ([]domain.Board, error)
UpdateBoard(ctx context.Context, userID, boardID uuid.UUID, req *domain.UpdateBoardRequest) (*domain.Board, error)
DeleteBoard(ctx context.Context, userID, boardID uuid.UUID) error

AddMember(ctx context.Context, requesterID, boardID uuid.UUID, req *domain.AddBoardMemberRequest) error
RemoveMember(ctx context.Context, requesterID, boardID, targetUserID uuid.UUID) error
ListMembers(ctx context.Context, userID, boardID uuid.UUID) ([]domain.BoardMember, error)
}

type boardService struct {
boardRepo repository.BoardRepository
userRepo  repository.UserRepository
}

// NewBoardService creates a new BoardService.
func NewBoardService(boardRepo repository.BoardRepository, userRepo repository.UserRepository) BoardService {
return &boardService{boardRepo: boardRepo, userRepo: userRepo}
}

func (s *boardService) CreateBoard(ctx context.Context, ownerID uuid.UUID, req *domain.CreateBoardRequest) (*domain.Board, error) {
visibility := req.Visibility
if visibility == "" {
visibility = domain.BoardVisibilityPrivate
}

board := &domain.Board{
OwnerID:     ownerID,
Title:       req.Title,
Description: req.Description,
Visibility:  visibility,
Background:  req.Background,
}

if err := s.boardRepo.Create(ctx, board); err != nil {
return nil, err
}

// Automatically add owner as admin member.
member := &domain.BoardMember{
BoardID: board.ID,
UserID:  ownerID,
Role:    domain.BoardMemberRoleAdmin,
}
if err := s.boardRepo.AddMember(ctx, member); err != nil {
return nil, err
}

return board, nil
}

func (s *boardService) GetBoard(ctx context.Context, userID, boardID uuid.UUID) (*domain.Board, error) {
board, err := s.boardRepo.FindByID(ctx, boardID)
if err != nil {
return nil, err
}

if board.Visibility == domain.BoardVisibilityPrivate {
if _, err := s.boardRepo.FindMember(ctx, boardID, userID); err != nil {
return nil, apperror.Forbidden("you do not have access to this board")
}
}

return board, nil
}

func (s *boardService) ListBoards(ctx context.Context, userID uuid.UUID) ([]domain.Board, error) {
return s.boardRepo.FindByMember(ctx, userID)
}

func (s *boardService) UpdateBoard(ctx context.Context, userID, boardID uuid.UUID, req *domain.UpdateBoardRequest) (*domain.Board, error) {
board, err := s.boardRepo.FindByID(ctx, boardID)
if err != nil {
return nil, err
}

if err := s.requireRole(ctx, boardID, userID, domain.BoardMemberRoleAdmin); err != nil {
return nil, err
}

if req.Title != "" {
board.Title = req.Title
}
if req.Description != "" {
board.Description = req.Description
}
if req.Visibility != "" {
board.Visibility = req.Visibility
}
if req.Background != "" {
board.Background = req.Background
}

if err := s.boardRepo.Update(ctx, board); err != nil {
return nil, err
}
return board, nil
}

func (s *boardService) DeleteBoard(ctx context.Context, userID, boardID uuid.UUID) error {
if _, err := s.boardRepo.FindByID(ctx, boardID); err != nil {
return err
}

if err := s.requireRole(ctx, boardID, userID, domain.BoardMemberRoleAdmin); err != nil {
return err
}

return s.boardRepo.Delete(ctx, boardID)
}

func (s *boardService) AddMember(ctx context.Context, requesterID, boardID uuid.UUID, req *domain.AddBoardMemberRequest) error {
if err := s.requireRole(ctx, boardID, requesterID, domain.BoardMemberRoleAdmin); err != nil {
return err
}

targetUser, err := s.userRepo.FindByEmail(ctx, req.Email)
if err != nil {
return apperror.NotFound("user with that email")
}

// Check if already a member.
if existing, _ := s.boardRepo.FindMember(ctx, boardID, targetUser.ID); existing != nil {
return apperror.Conflict("user is already a member of this board")
}

member := &domain.BoardMember{
BoardID: boardID,
UserID:  targetUser.ID,
Role:    req.Role,
}
return s.boardRepo.AddMember(ctx, member)
}

func (s *boardService) RemoveMember(ctx context.Context, requesterID, boardID, targetUserID uuid.UUID) error {
if err := s.requireRole(ctx, boardID, requesterID, domain.BoardMemberRoleAdmin); err != nil {
return err
}
return s.boardRepo.RemoveMember(ctx, boardID, targetUserID)
}

func (s *boardService) ListMembers(ctx context.Context, userID, boardID uuid.UUID) ([]domain.BoardMember, error) {
if _, err := s.boardRepo.FindMember(ctx, boardID, userID); err != nil {
return nil, apperror.Forbidden("you do not have access to this board")
}
return s.boardRepo.ListMembers(ctx, boardID)
}

// requireRole checks that userID has at least the specified role on the board.
// Admin > Member > Viewer.
func (s *boardService) requireRole(ctx context.Context, boardID, userID uuid.UUID, required domain.BoardMemberRole) error {
member, err := s.boardRepo.FindMember(ctx, boardID, userID)
if err != nil {
return apperror.Forbidden("you do not have access to this board")
}

order := map[domain.BoardMemberRole]int{
domain.BoardMemberRoleViewer: 0,
domain.BoardMemberRoleMember: 1,
domain.BoardMemberRoleAdmin:  2,
}

if order[member.Role] < order[required] {
return apperror.Forbidden("insufficient permissions")
}
return nil
}
