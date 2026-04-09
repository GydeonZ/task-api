package repository

import "github.com/GydeonZ/task-api/internal/domain"

type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
}

type BoardRepository interface {
	Create(board *domain.Board) error
	FindByUserID(userID uint) ([]domain.Board, error)
	FindByID(id uint) (*domain.Board, error)
	Update(board *domain.Board) error
	Delete(id uint) error
}

type ListRepository interface {
	Create(list *domain.List) error
	FindByBoardID(boardID uint) ([]domain.List, error)
	FindByID(id uint) (*domain.List, error)
	Update(list *domain.List) error
	Delete(id uint) error
}

type CardRepository interface {
	Create(card *domain.Card) error
	FindByListID(listID uint) ([]domain.Card, error)
	FindByID(id uint) (*domain.Card, error)
	Update(card *domain.Card) error
	Delete(id uint) error
	MoveToList(cardID, listID uint, position int) error
}
