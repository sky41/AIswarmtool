package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
	"gorm.io/gorm"
)

type AlertRule struct {
	ID                 int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	RuleName           string          `gorm:"size:100;not null" json:"rule_name"`
	RuleDesc           string          `gorm:"size:500" json:"rule_desc"`
	MetricType         string          `gorm:"size:50;not null;index" json:"metric_type"`
	ResourceType       string          `gorm:"size:50;not null" json:"resource_type"`
	Condition          string          `gorm:"size:20;not null" json:"condition"`
	Threshold          float64         `gorm:"type:decimal(10,2);not null" json:"threshold"`
	Duration           int             `gorm:"default:60" json:"duration"`
	Severity           string          `gorm:"size:20;default:warning" json:"severity"`
	Enabled            bool            `gorm:"default:true;index" json:"enabled"`
	NotificationType   *string         `gorm:"size:50" json:"notification_type"`
	NotificationConfig *JSONMap        `gorm:"type:json" json:"notification_config"`
	CreatedAt          time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

func (AlertRule) TableName() string {
	return "alert_rules"
}
