package model

import (
	"time"
	"gorm.io/gorm"
)

type MetricsHistory struct {
	ID                int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	MetricType        string         `gorm:"size:50;not null;index" json:"metric_type"`
	ResourceID        string         `gorm:"size:100;not null;index" json:"resource_id"`
	ResourceName      string         `gorm:"size:100" json:"resource_name"`
	CPUUsage          *float64       `gorm:"type:decimal(5,2)" json:"cpu_usage"`
	MemoryUsage       *float64       `gorm:"type:decimal(5,2)" json:"memory_usage"`
	MemoryUsedBytes   *int64         `json:"memory_used_bytes"`
	MemoryTotalBytes  *int64         `json:"memory_total_bytes"`
	DiskUsage         *float64       `gorm:"type:decimal(5,2)" json:"disk_usage"`
	NetworkRxBytes    *int64         `json:"network_rx_bytes"`
	NetworkTxBytes    *int64         `json:"network_tx_bytes"`
	Status            string         `gorm:"size:50" json:"status"`
	CollectedAt       time.Time      `gorm:"not null;index" json:"collected_at"`
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`
}

func (MetricsHistory) TableName() string {
	return "metrics_history"
}
