package repository

import (
	"github.com/aiswarmtool/backend/internal/model"
)

func (r *Repository) GetHealingPolicies(enabledOnly bool) ([]*model.SelfHealingPolicy, error) {
	var policies []*model.SelfHealingPolicy
	query := r.db
	if enabledOnly {
		query = query.Where("enabled = ?", true)
	}
	err := query.Order("created_at DESC").Find(&policies).Error
	return policies, err
}

func (r *Repository) GetHealingPolicyByID(id int64) (*model.SelfHealingPolicy, error) {
	var policy model.SelfHealingPolicy
	err := r.db.First(&policy, id).Error
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

func (r *Repository) CreateHealingPolicy(policy *model.SelfHealingPolicy) error {
	return r.db.Create(policy).Error
}

func (r *Repository) UpdateHealingPolicy(policy *model.SelfHealingPolicy) error {
	return r.db.Save(policy).Error
}

func (r *Repository) DeleteHealingPolicy(id int64) error {
	return r.db.Delete(&model.SelfHealingPolicy{}, id).Error
}

func (r *Repository) CreateHealingExecution(execution *model.HealingExecution) error {
	return r.db.Create(execution).Error
}

func (r *Repository) UpdateHealingExecution(execution *model.HealingExecution) error {
	return r.db.Save(execution).Error
}

func (r *Repository) GetHealingExecutionsByIncident(incidentID string) ([]*model.HealingExecution, error) {
	var executions []*model.HealingExecution
	err := r.db.Where("incident_id = ?", incidentID).Preload("Policy").Order("created_at DESC").Find(&executions).Error
	return executions, err
}
