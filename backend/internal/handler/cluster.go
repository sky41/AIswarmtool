package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetClusterInfo(c *gin.Context) {
	data := gin.H{
		"name": "Docker Swarm Cluster",
		"status": "healthy",
	}
	h.respondJSON(c, http.StatusOK, data)
}

func (h *Handler) GetNodes(c *gin.Context) {
	data := []gin.H{}
	h.respondJSON(c, http.StatusOK, data)
}

func (h *Handler) GetServices(c *gin.Context) {
	data := []gin.H{}
	h.respondJSON(c, http.StatusOK, data)
}
