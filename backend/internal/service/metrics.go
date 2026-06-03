package service

import (
	"time"
	"github.com/aiswarmtool/backend/internal/model"
)

func (s *Service) CreateMetrics(metrics *model.MetricsHistory) error {
	return s.repo.CreateMetricsHistory(metrics)
}

func (s *Service) GetMetrics(resourceType, resourceID string, startTime, endTime string, limit int) ([]*model.MetricsHistory, error) {
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
	return s.repo.GetMetricsByResource(resourceType, resourceID, st, et, limit)
}

func (s *Service) GetLatestMetrics(resourceType string) ([]*model.MetricsHistory, error) {
	return s.repo.GetLatestMetrics(resourceType)
}
