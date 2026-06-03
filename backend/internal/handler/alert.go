package handler

import (
	"net/http"
	"strconv"
	"github.com/aiswarmtool/backend/internal/model"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAlertRules(c *gin.Context) {
	enabledOnly := c.Query("enabled") == "true"
	rules, err := h.service.GetAlertRules(enabledOnly)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}
	h.respondJSON(c, http.StatusOK, rules)
}

func (h *Handler) GetAlertRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}
	rule, err := h.service.GetAlertRule(id)
	if err != nil {
		h.respondError(c, http.StatusNotFound, err)
		return
	}
	h.respondJSON(c, http.StatusOK, rule)
}

func (h *Handler) CreateAlertRule(c *gin.Context) {
	var rule model.AlertRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}
	if err := h.service.CreateAlertRule(&rule); err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}
	h.respondJSON(c, http.StatusCreated, rule)
}

func (h *Handler) UpdateAlertRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}
	var rule model.AlertRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}
	rule.ID = id
	if err := h.service.UpdateAlertRule(&rule); err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}
	h.respondJSON(c, http.StatusOK, rule)
}

func (h *Handler) DeleteAlertRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}
	if err := h.service.DeleteAlertRule(id); err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}
	h.respondJSON(c, http.StatusOK, nil)
}
