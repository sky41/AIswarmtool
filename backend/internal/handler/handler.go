package handler

import (
	"github.com/aiswarmtool/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) respondJSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, gin.H{
		"success": status < 400,
		"data":    data,
	})
}

func (h *Handler) respondError(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{
		"success": false,
		"error":   err.Error(),
	})
}
