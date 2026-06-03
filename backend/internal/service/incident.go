package service

import (
	"time"
	"github.com/aiswarmtool/backend/internal/model"
)

func (s *Service) GetIncidents(status, startTime, endTime string, limit int) ([]*model.IncidentLog, error) {
	var st, et time.Time
	if startTime != "" {
		st, _ = time.Parse(time.RFC3339, startTime)
	}
	if endTime != "" {
		et, _ = time.Parse(time.RFC3339, endTime)
	}
	if limit <= 0 {
		limit = 100
	}
	return s.repo.GetIncidents(status, st, et, limit)
}

func (s *Service) GetIncident(id int64) (*model.IncidentLog, error) {
	return s.repo.GetIncidentByID(id)
}

func (s *Service) AcknowledgeIncident(id int64) error {
	incident, err := s.repo.GetIncidentByID(id)
	if err != nil {
		return err
	}
	now := time.Now()
	incident.Status = "acknowledged"
	incident.AcknowledgedAt = &now
	return s.repo.UpdateIncident(incident)
}

func (s *Service) ResolveIncident(id int64) error {
	incident, err := s.repo.GetIncidentByID(id)
	if err != nil {
		return err
	}
	now := time.Now()
	incident.Status = "resolved"
	incident.ResolvedAt = &now
	return s.repo.UpdateIncident(incident)
}
