package service

import (
	"github.com/GydeonZ/task-api/internal/domain"
	"github.com/GydeonZ/task-api/internal/repository"
	"github.com/GydeonZ/task-api/pkg/apperrors"
)

type BoardService interface {
	CreateBoard(userID uint, req *domain.CreateBoardRequest) (*domain.Board, error)
	GetBoards(userID uint) ([]domain.Board, error)
	GetBoard(userID, boardID uint) (*domain.Board, error)
	UpdateBoard(userID, boardID uint, req *domain.UpdateBoardRequest) (*domain.Board, error)
	DeleteBoard(userID, boardID uint) error
}

type boardService struct {
	boardRepo repository.BoardRepository
}

func NewBoardService(boardRepo repository.BoardRepository) BoardService {
	return &boardService{boardRepo: boardRepo}
}

func (s *boardService) CreateBoard(userID uint, req *domain.CreateBoardRequest) (*domain.Board, error) {
	board := &domain.Board{
		Title:       req.Title,
		Description: req.Description,
		UserID:      userID,
	}
	if err := s.boardRepo.Create(board); err != nil {
		return nil, err
	}
	return board, nil
}

func (s *boardService) GetBoards(userID uint) ([]domain.Board, error) {
	return s.boardRepo.FindByUserID(userID)
}

func (s *boardService) GetBoard(userID, boardID uint) (*domain.Board, error) {
	board, err := s.boardRepo.FindByID(boardID)
	if err != nil {
		return nil, err
	}
	if board.UserID != userID {
		return nil, apperrors.ErrUnauthorized
	}
	return board, nil
}

func (s *boardService) UpdateBoard(userID, boardID uint, req *domain.UpdateBoardRequest) (*domain.Board, error) {
	board, err := s.boardRepo.FindByID(boardID)
	if err != nil {
		return nil, err
	}
	if board.UserID != userID {
		return nil, apperrors.ErrUnauthorized
	}
	if req.Title != "" {
		board.Title = req.Title
	}
	if req.Description != "" {
		board.Description = req.Description
	}
	if err := s.boardRepo.Update(board); err != nil {
		return nil, err
	}
	return board, nil
}

func (s *boardService) DeleteBoard(userID, boardID uint) error {
	board, err := s.boardRepo.FindByID(boardID)
	if err != nil {
		return err
	}
	if board.UserID != userID {
		return apperrors.ErrUnauthorized
	}
	return s.boardRepo.Delete(boardID)
}
