package repository

import (
	"time"
	"github.com/aiswarmtool/backend/internal/model"
)

func (r *Repository) GetIncidents(status string, startTime, endTime time.Time, limit int) ([]*model.IncidentLog, error) {
	var incidents []*model.IncidentLog
	query := r.db.Preload("AlertRule").Preload("HealingPolicy")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if !startTime.IsZero() {
		query = query.Where("triggered_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		query = query.Where("triggered_at <= ?", endTime)
	}
	err := query.Order("triggered_at DESC").Limit(limit).Find(&incidents).Error
	return incidents, err
}

func (r *Repository) GetIncidentByID(id int64) (*model.IncidentLog, error) {
	var incident model.IncidentLog
	err := r.db.Preload("AlertRule").Preload("HealingPolicy").First(&incident, id).Error
	if err != nil {
		return nil, err
	}
	return &incident, nil
}

func (r *Repository) GetIncidentByIncidentID(incidentID string) (*model.IncidentLog, error) {
	var incident model.IncidentLog
	err := r.db.Where("incident_id = ?", incidentID).Preload("AlertRule").Preload("HealingPolicy").First(&incident).Error
	if err != nil {
		return nil, err
	}
	return &incident, nil
}

func (r *Repository) CreateIncident(incident *model.IncidentLog) error {
	return r.db.Create(incident).Error
}

func (r *Repository) UpdateIncident(incident *model.IncidentLog) error {
	return r.db.Save(incident).Error
}

func (r *Repository) GetOpenIncidentsByResource(resourceType, resourceID string) ([]*model.IncidentLog, error) {
	var incidents []*model.IncidentLog
	err := r.db.Where("resource_type = ? AND resource_id = ? AND status IN ?", resourceType, resourceID, []string{"open", "acknowledged", "resolving"}).Find(&incidents).Error
	return incidents, err
}
