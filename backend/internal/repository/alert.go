package repository

import (
	"github.com/aiswarmtool/backend/internal/model"
)

func (r *Repository) GetAlertRules(enabledOnly bool) ([]*model.AlertRule, error) {
	var rules []*model.AlertRule
	query := r.db
	if enabledOnly {
		query = query.Where("enabled = ?", true)
	}
	err := query.Order("created_at DESC").Find(&rules).Error
	return rules, err
}

func (r *Repository) GetAlertRuleByID(id int64) (*model.AlertRule, error) {
	var rule model.AlertRule
	err := r.db.First(&rule, id).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *Repository) CreateAlertRule(rule *model.AlertRule) error {
	return r.db.Create(rule).Error
}

func (r *Repository) UpdateAlertRule(rule *model.AlertRule) error {
	return r.db.Save(rule).Error
}

func (r *Repository) DeleteAlertRule(id int64) error {
	return r.db.Delete(&model.AlertRule{}, id).Error
}
