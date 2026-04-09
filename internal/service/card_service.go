package service

import (
	"github.com/GydeonZ/task-api/internal/domain"
	"github.com/GydeonZ/task-api/internal/repository"
	"github.com/GydeonZ/task-api/pkg/apperrors"
)

type CardService interface {
	CreateCard(userID, boardID, listID uint, req *domain.CreateCardRequest) (*domain.Card, error)
	GetCards(userID, boardID, listID uint) ([]domain.Card, error)
	GetCard(userID, boardID, listID, cardID uint) (*domain.Card, error)
	UpdateCard(userID, boardID, listID, cardID uint, req *domain.UpdateCardRequest) (*domain.Card, error)
	DeleteCard(userID, boardID, listID, cardID uint) error
	MoveCard(userID, cardID uint, req *domain.MoveCardRequest) (*domain.Card, error)
}

type cardService struct {
	cardRepo  repository.CardRepository
	listRepo  repository.ListRepository
	boardRepo repository.BoardRepository
}

func NewCardService(cardRepo repository.CardRepository, listRepo repository.ListRepository, boardRepo repository.BoardRepository) CardService {
	return &cardService{cardRepo: cardRepo, listRepo: listRepo, boardRepo: boardRepo}
}

func (s *cardService) verifyBoardOwnership(userID, boardID uint) error {
	board, err := s.boardRepo.FindByID(boardID)
	if err != nil {
		return err
	}
	if board.UserID != userID {
		return apperrors.ErrUnauthorized
	}
	return nil
}

func (s *cardService) CreateCard(userID, boardID, listID uint, req *domain.CreateCardRequest) (*domain.Card, error) {
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
	card := &domain.Card{
		Title:       req.Title,
		Description: req.Description,
		Position:    req.Position,
		DueDate:     req.DueDate,
		ListID:      listID,
	}
	if err := s.cardRepo.Create(card); err != nil {
		return nil, err
	}
	return card, nil
}

func (s *cardService) GetCards(userID, boardID, listID uint) ([]domain.Card, error) {
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
	return s.cardRepo.FindByListID(listID)
}

func (s *cardService) GetCard(userID, boardID, listID, cardID uint) (*domain.Card, error) {
	if err := s.verifyBoardOwnership(userID, boardID); err != nil {
		return nil, err
	}
	card, err := s.cardRepo.FindByID(cardID)
	if err != nil {
		return nil, err
	}
	if card.ListID != listID {
		return nil, apperrors.ErrNotFound
	}
	return card, nil
}

func (s *cardService) UpdateCard(userID, boardID, listID, cardID uint, req *domain.UpdateCardRequest) (*domain.Card, error) {
	if err := s.verifyBoardOwnership(userID, boardID); err != nil {
		return nil, err
	}
	card, err := s.cardRepo.FindByID(cardID)
	if err != nil {
		return nil, err
	}
	if card.ListID != listID {
		return nil, apperrors.ErrNotFound
	}
	if req.Title != "" {
		card.Title = req.Title
	}
	if req.Description != "" {
		card.Description = req.Description
	}
	card.Position = req.Position
	card.DueDate = req.DueDate
	if err := s.cardRepo.Update(card); err != nil {
		return nil, err
	}
	return card, nil
}

func (s *cardService) DeleteCard(userID, boardID, listID, cardID uint) error {
	if err := s.verifyBoardOwnership(userID, boardID); err != nil {
		return err
	}
	card, err := s.cardRepo.FindByID(cardID)
	if err != nil {
		return err
	}
	if card.ListID != listID {
		return apperrors.ErrNotFound
	}
	return s.cardRepo.Delete(cardID)
}

func (s *cardService) MoveCard(userID, cardID uint, req *domain.MoveCardRequest) (*domain.Card, error) {
	card, err := s.cardRepo.FindByID(cardID)
	if err != nil {
		return nil, err
	}
	list, err := s.listRepo.FindByID(card.ListID)
	if err != nil {
		return nil, err
	}
	board, err := s.boardRepo.FindByID(list.BoardID)
	if err != nil {
		return nil, err
	}
	if board.UserID != userID {
		return nil, apperrors.ErrUnauthorized
	}
	if err := s.cardRepo.MoveToList(cardID, req.ListID, req.Position); err != nil {
		return nil, err
	}
	card.ListID = req.ListID
	card.Position = req.Position
	return card, nil
}
