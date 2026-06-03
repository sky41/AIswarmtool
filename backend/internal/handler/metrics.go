package handler

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetMetrics(c *gin.Context) {
	resourceType := c.Query("type")
	resourceID := c.Query("id")
	startTime := c.Query("start")
	endTime := c.Query("end")
	limitStr := c.Query("limit")
	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}
	metrics, err := h.service.GetMetrics(resourceType, resourceID, startTime, endTime, limit)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}
	h.respondJSON(c, http.StatusOK, metrics)
}

func (h *Handler) GetLatestMetrics(c *gin.Context) {
	resourceType := c.Query("type")
	if resourceType == "" {
		resourceType = "node"
	}
	metrics, err := h.service.GetLatestMetrics(resourceType)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}
	h.respondJSON(c, http.StatusOK, metrics)
}
