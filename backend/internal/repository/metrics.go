package repository

import (
	"time"
	"github.com/aiswarmtool/backend/internal/model"
)

func (r *Repository) CreateMetricsHistory(metrics *model.MetricsHistory) error {
	return r.db.Create(metrics).Error
}

func (r *Repository) GetMetricsByResource(resourceType, resourceID string, startTime, endTime time.Time, limit int) ([]*model.MetricsHistory, error) {
	var metrics []*model.MetricsHistory
	query := r.db.Where("metric_type = ? AND resource_id = ?", resourceType, resourceID)
	if !startTime.IsZero() {
		query = query.Where("collected_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		query = query.Where("collected_at <= ?", endTime)
	}
	err := query.Order("collected_at DESC").Limit(limit).Find(&metrics).Error
	return metrics, err
}

func (r *Repository) GetLatestMetrics(resourceType string) ([]*model.MetricsHistory, error) {
	var metrics []*model.MetricsHistory
	subQuery := r.db.Model(&model.MetricsHistory{}).
		Select("MAX(id) as id").
		Where("metric_type = ?", resourceType).
		Group("resource_id")
	err := r.db.Where("id IN (?)", subQuery).Find(&metrics).Error
	return metrics, err
}
