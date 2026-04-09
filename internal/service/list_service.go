package service

import (
	"github.com/GydeonZ/task-api/internal/domain"
	"github.com/GydeonZ/task-api/internal/repository"
	"github.com/GydeonZ/task-api/pkg/apperrors"
)

type ListService interface {
	CreateList(userID, boardID uint, req *domain.CreateListRequest) (*domain.List, error)
	GetLists(userID, boardID uint) ([]domain.List, error)
	UpdateList(userID, boardID, listID uint, req *domain.UpdateListRequest) (*domain.List, error)
	DeleteList(userID, boardID, listID uint) error
}

type listService struct {
	listRepo  repository.ListRepository
	boardRepo repository.BoardRepository
}

func NewListService(listRepo repository.ListRepository, boardRepo repository.BoardRepository) ListService {
	return &listService{listRepo: listRepo, boardRepo: boardRepo}
}

func (s *listService) verifyBoardOwnership(userID, boardID uint) error {
	board, err := s.boardRepo.FindByID(boardID)
	if err != nil {
		return err
	}
	if board.UserID != userID {
		return apperrors.ErrUnauthorized
	}
	return nil
}

func (s *listService) CreateList(userID, boardID uint, req *domain.CreateListRequest) (*domain.List, error) {
	if err := s.verifyBoardOwnership(userID, boardID); err != nil {
		return nil, err
	}
	list := &domain.List{
		Title:    req.Title,
		Position: req.Position,
		BoardID:  boardID,
	}
	if err := s.listRepo.Create(list); err != nil {
		return nil, err
	}
	return list, nil
}

func (s *listService) GetLists(userID, boardID uint) ([]domain.List, error) {
	if err := s.verifyBoardOwnership(userID, boardID); err != nil {
		return nil, err
	}
	return s.listRepo.FindByBoardID(boardID)
}

func (s *listService) UpdateList(userID, boardID, listID uint, req *domain.UpdateListRequest) (*domain.List, error) {
	if err := s.verifyBoardOwnership(userID, boardID); err != nil {
		return nil, err
	}
	list, err := s.listRepo.FindByID(listID)
	if err != nil {
		return nil, err
	}
	if list.BoardID != boardID {
		return nil, apperrors.ErrNotFound
	}
	if req.Title != "" {
		list.Title = req.Title
	}
	list.Position = req.Position
	if err := s.listRepo.Update(list); err != nil {
		return nil, err
	}
	return list, nil
}

func (s *listService) DeleteList(userID, boardID, listID uint) error {
	if err := s.verifyBoardOwnership(userID, boardID); err != nil {
		return err
	}
	list, err := s.listRepo.FindByID(listID)
	if err != nil {
		return err
	}
	if list.BoardID != boardID {
		return apperrors.ErrNotFound
	}
	return s.listRepo.Delete(listID)
}
