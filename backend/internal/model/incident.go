package model

import (
	"time"
)

type IncidentLog struct {
	ID               int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	IncidentID       string     `gorm:"size:100;not null;uniqueIndex" json:"incident_id"`
	AlertRuleID      *int64     `json:"alert_rule_id"`
	ResourceType     string     `gorm:"size:50;not null;index:idx_resource,priority:1" json:"resource_type"`
	ResourceID       string     `gorm:"size:100;not null;index:idx_resource,priority:2" json:"resource_id"`
	ResourceName     string     `gorm:"size:100" json:"resource_name"`
	Severity         string     `gorm:"size:20;not null" json:"severity"`
	Description      *string    `gorm:"type:text" json:"description"`
	Status           string     `gorm:"size:20;default:open;index" json:"status"`
	TriggeredAt      time.Time  `gorm:"not null;index" json:"triggered_at"`
	AcknowledgedAt   *time.Time `json:"acknowledged_at"`
	ResolvedAt       *time.Time `json:"resolved_at"`
	HealingPolicyID  *int64     `json:"healing_policy_id"`
	HealingExecuted  bool       `gorm:"default:false" json:"healing_executed"`
	HealingResult    *string    `gorm:"type:text" json:"healing_result"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	AlertRule        *AlertRule          `gorm:"foreignKey:AlertRuleID" json:"alert_rule,omitempty"`
	HealingPolicy    *SelfHealingPolicy  `gorm:"foreignKey:HealingPolicyID" json:"healing_policy,omitempty"`
}

func (IncidentLog) TableName() string {
	return "incident_logs"
}
