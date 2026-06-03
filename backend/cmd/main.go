package main

import (
	"fmt"
	"log"
	"github.com/aiswarmtool/backend/internal/handler"
	"github.com/aiswarmtool/backend/internal/repository"
	"github.com/aiswarmtool/backend/internal/service"
	"github.com/aiswarmtool/backend/pkg/config"
	"github.com/aiswarmtool/backend/pkg/docker"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.DBName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	dockerCli, err := docker.NewClient(config.AppConfig.Docker.Host)
	if err != nil {
		log.Printf("Warning: Failed to connect to Docker: %v", err)
	}
	defer dockerCli.Close()
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)
	monitor := service.NewMonitor(svc, dockerCli)
	go monitor.Start(config.AppConfig.Monitoring.Interval)
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.GET("/health", h.HealthCheck)
	api := r.Group("/api/v1")
	{
		api.GET("/health", h.HealthCheck)
		metrics := api.Group("/metrics")
		{
			metrics.GET("", h.GetMetrics)
			metrics.GET("/latest", h.GetLatestMetrics)
		}
		alerts := api.Group("/alerts")
		{
			alerts.GET("", h.GetIncidents)
			alerts.GET("/rules", h.GetAlertRules)
			alerts.GET("/rules/:id", h.GetAlertRule)
			alerts.POST("/rules", h.CreateAlertRule)
			alerts.PUT("/rules/:id", h.UpdateAlertRule)
			alerts.DELETE("/rules/:id", h.DeleteAlertRule)
			alerts.GET("/:id", h.GetIncident)
			alerts.POST("/:id/acknowledge", h.AcknowledgeIncident)
			alerts.POST("/:id/resolve", h.ResolveIncident)
		}
		healing := api.Group("/healing")
		{
			healing.GET("/policies", h.GetHealingPolicies)
			healing.GET("/policies/:id", h.GetHealingPolicy)
			healing.POST("/policies", h.CreateHealingPolicy)
			healing.PUT("/policies/:id", h.UpdateHealingPolicy)
			healing.DELETE("/policies/:id", h.DeleteHealingPolicy)
			healing.GET("/executions", h.GetHealingExecutions)
		}
		cluster := api.Group("/cluster")
		{
			cluster.GET("", h.GetClusterInfo)
			cluster.GET("/nodes", h.GetNodes)
			cluster.GET("/services", h.GetServices)
		}
	}
	addr := fmt.Sprintf(":%s", config.AppConfig.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:5173",
		}
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
