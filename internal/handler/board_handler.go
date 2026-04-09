package handler

import (
	"net/http"

	"github.com/GydeonZ/task-api/internal/domain"
	"github.com/GydeonZ/task-api/internal/service"
	"github.com/GydeonZ/task-api/pkg/apperrors"
	"github.com/GydeonZ/task-api/pkg/response"
	"github.com/gin-gonic/gin"
)

type BoardHandler struct {
	boardSvc service.BoardService
}

func NewBoardHandler(boardSvc service.BoardService) *BoardHandler {
	return &BoardHandler{boardSvc: boardSvc}
}

func (h *BoardHandler) GetBoards(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	boards, err := h.boardSvc.GetBoards(userID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "boards retrieved", boards)
}

func (h *BoardHandler) CreateBoard(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var req domain.CreateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	board, err := h.boardSvc.CreateBoard(userID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusCreated, "board created", board)
}

func (h *BoardHandler) GetBoard(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	boardID, err := parseID(c.Param("id"))
	if err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	board, err := h.boardSvc.GetBoard(userID, boardID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "board retrieved", board)
}

func (h *BoardHandler) UpdateBoard(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	boardID, err := parseID(c.Param("id"))
	if err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	var req domain.UpdateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	board, err := h.boardSvc.UpdateBoard(userID, boardID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "board updated", board)
}

func (h *BoardHandler) DeleteBoard(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	boardID, err := parseID(c.Param("id"))
	if err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	if err := h.boardSvc.DeleteBoard(userID, boardID); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "board deleted", nil)
}
