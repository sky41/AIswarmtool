package model

import (
	"time"
)

type SelfHealingPolicy struct {
	ID                int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	PolicyName        string     `gorm:"size:100;not null" json:"policy_name"`
	PolicyDesc        string     `gorm:"size:500" json:"policy_desc"`
	TriggerCondition  string     `gorm:"type:text;not null" json:"trigger_condition"`
	ActionType        string     `gorm:"size:50;not null" json:"action_type"`
	ActionParams      JSONMap    `gorm:"type:json;not null" json:"action_params"`
	MaxRetries        int        `gorm:"default:3" json:"max_retries"`
	RetryDelay        int        `gorm:"default:30" json:"retry_delay"`
	Enabled           bool       `gorm:"default:true;index" json:"enabled"`
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (SelfHealingPolicy) TableName() string {
	return "self_healing_policies"
}

type HealingExecution struct {
	ID           int64             `gorm:"primaryKey;autoIncrement" json:"id"`
	PolicyID     int64             `gorm:"not null;index" json:"policy_id"`
	IncidentID   string            `gorm:"size:100;not null;index" json:"incident_id"`
	ResourceType string            `gorm:"size:50;not null" json:"resource_type"`
	ResourceID   string            `gorm:"size:100;not null" json:"resource_id"`
	ActionType   string            `gorm:"size:50;not null" json:"action_type"`
	Attempt      int               `gorm:"default:1" json:"attempt"`
	Status       string            `gorm:"size:20;not null;index" json:"status"`
	Result       *string           `gorm:"type:text" json:"result"`
	StartedAt    time.Time         `gorm:"not null" json:"started_at"`
	CompletedAt  *time.Time        `json:"completed_at"`
	CreatedAt    time.Time         `gorm:"autoCreateTime" json:"created_at"`

	Policy       *SelfHealingPolicy `gorm:"foreignKey:PolicyID" json:"policy,omitempty"`
}

func (HealingExecution) TableName() string {
	return "healing_executions"
}
