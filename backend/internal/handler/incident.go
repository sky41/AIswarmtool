package handler

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetIncidents(c *gin.Context) {
	status := c.Query("status")
	startTime := c.Query("start")
	endTime := c.Query("end")
	limitStr := c.Query("limit")
	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}
	incidents, err := h.service.GetIncidents(status, startTime, endTime, limit)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}
	h.respondJSON(c, http.StatusOK, incidents)
}

func (h *Handler) GetIncident(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}
	incident, err := h.service.GetIncident(id)
	if err != nil {
		h.respondError(c, http.StatusNotFound, err)
		return
	}
	h.respondJSON(c, http.StatusOK, incident)
}

func (h *Handler) AcknowledgeIncident(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}
	if err := h.service.AcknowledgeIncident(id); err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}
	h.respondJSON(c, http.StatusOK, nil)
}

func (h *Handler) ResolveIncident(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}
	if err := h.service.ResolveIncident(id); err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}
	h.respondJSON(c, http.StatusOK, nil)
}
