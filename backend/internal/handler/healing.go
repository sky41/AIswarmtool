package handler

import (
	"net/http"
	"strconv"
	"github.com/aiswarmtool/backend/internal/model"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetHealingPolicies(c *gin.Context) {
	enabledOnly := c.Query("enabled") == "true"
	policies, err := h.service.GetHealingPolicies(enabledOnly)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}
	h.respondJSON(c, http.StatusOK, policies)
}

func (h *Handler) GetHealingPolicy(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}
	policy, err := h.service.GetHealingPolicy(id)
	if err != nil {
		h.respondError(c, http.StatusNotFound, err)
		return
	}
	h.respondJSON(c, http.StatusOK, policy)
}

func (h *Handler) CreateHealingPolicy(c *gin.Context) {
	var policy model.SelfHealingPolicy
	if err := c.ShouldBindJSON(&policy); err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}
	if err := h.service.CreateHealingPolicy(&policy); err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}
	h.respondJSON(c, http.StatusCreated, policy)
}

func (h *Handler) UpdateHealingPolicy(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}
	var policy model.SelfHealingPolicy
	if err := c.ShouldBindJSON(&policy); err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}
	policy.ID = id
	if err := h.service.UpdateHealingPolicy(&policy); err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}
	h.respondJSON(c, http.StatusOK, policy)
}

func (h *Handler) DeleteHealingPolicy(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}
	if err := h.service.DeleteHealingPolicy(id); err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}
	h.respondJSON(c, http.StatusOK, nil)
}

func (h *Handler) GetHealingExecutions(c *gin.Context) {
	incidentID := c.Query("incident_id")
	executions, err := h.service.GetHealingExecutions(incidentID)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}
	h.respondJSON(c, http.StatusOK, executions)
}
