package handler

import (
	"net/http"
	"strconv"

	"github.com/GydeonZ/task-api/internal/domain"
	"github.com/GydeonZ/task-api/internal/service"
	"github.com/GydeonZ/task-api/pkg/apperrors"
	"github.com/GydeonZ/task-api/pkg/response"
	"github.com/gin-gonic/gin"
)

type CardHandler struct {
	cardSvc service.CardService
}

func NewCardHandler(cardSvc service.CardService) *CardHandler {
	return &CardHandler{cardSvc: cardSvc}
}

func (h *CardHandler) GetCards(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	boardID, err := strconv.ParseUint(c.Param("boardID"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	listID, err := strconv.ParseUint(c.Param("listID"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	cards, err := h.cardSvc.GetCards(userID, uint(boardID), uint(listID))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "cards retrieved", cards)
}

func (h *CardHandler) CreateCard(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	boardID, err := strconv.ParseUint(c.Param("boardID"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	listID, err := strconv.ParseUint(c.Param("listID"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	var req domain.CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	card, err := h.cardSvc.CreateCard(userID, uint(boardID), uint(listID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusCreated, "card created", card)
}

func (h *CardHandler) GetCard(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	boardID, err := strconv.ParseUint(c.Param("boardID"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	listID, err := strconv.ParseUint(c.Param("listID"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	cardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	card, err := h.cardSvc.GetCard(userID, uint(boardID), uint(listID), uint(cardID))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "card retrieved", card)
}

func (h *CardHandler) UpdateCard(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	boardID, err := strconv.ParseUint(c.Param("boardID"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	listID, err := strconv.ParseUint(c.Param("listID"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	cardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	var req domain.UpdateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	card, err := h.cardSvc.UpdateCard(userID, uint(boardID), uint(listID), uint(cardID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "card updated", card)
}

func (h *CardHandler) DeleteCard(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	boardID, err := strconv.ParseUint(c.Param("boardID"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	listID, err := strconv.ParseUint(c.Param("listID"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	cardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	if err := h.cardSvc.DeleteCard(userID, uint(boardID), uint(listID), uint(cardID)); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "card deleted", nil)
}

func (h *CardHandler) MoveCard(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	cardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	var req domain.MoveCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	card, err := h.cardSvc.MoveCard(userID, uint(cardID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "card moved", card)
}
