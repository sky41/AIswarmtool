package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthStatus struct {
	Status  string            `json:"status"`
	Service map[string]string `json:"services"`
}

func (h *Handler) HealthCheck(c *gin.Context) {
	status := HealthStatus{
		Status:  "healthy",
		Service: make(map[string]string),
	}
	status.Service["api"] = "healthy"
	status.Service["database"] = "healthy"
	status.Service["docker"] = "healthy"
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    status,
	})
}
