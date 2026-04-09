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

type ListHandler struct {
	listSvc service.ListService
}

func NewListHandler(listSvc service.ListService) *ListHandler {
	return &ListHandler{listSvc: listSvc}
}

func (h *ListHandler) GetLists(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	boardID, err := strconv.ParseUint(c.Param("boardID"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	lists, err := h.listSvc.GetLists(userID, uint(boardID))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "lists retrieved", lists)
}

func (h *ListHandler) CreateList(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	boardID, err := strconv.ParseUint(c.Param("boardID"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	var req domain.CreateListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	list, err := h.listSvc.CreateList(userID, uint(boardID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusCreated, "list created", list)
}

func (h *ListHandler) UpdateList(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	boardID, err := strconv.ParseUint(c.Param("boardID"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	listID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	var req domain.UpdateListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	list, err := h.listSvc.UpdateList(userID, uint(boardID), uint(listID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "list updated", list)
}

func (h *ListHandler) DeleteList(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	boardID, err := strconv.ParseUint(c.Param("boardID"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	listID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.ErrBadRequest)
		return
	}
	if err := h.listSvc.DeleteList(userID, uint(boardID), uint(listID)); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "list deleted", nil)
}
