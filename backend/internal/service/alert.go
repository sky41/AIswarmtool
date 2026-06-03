package service

import (
	"github.com/aiswarmtool/backend/internal/model"
)

func (s *Service) GetAlertRules(enabledOnly bool) ([]*model.AlertRule, error) {
	return s.repo.GetAlertRules(enabledOnly)
}

func (s *Service) GetAlertRule(id int64) (*model.AlertRule, error) {
	return s.repo.GetAlertRuleByID(id)
}

func (s *Service) CreateAlertRule(rule *model.AlertRule) error {
	return s.repo.CreateAlertRule(rule)
}

func (s *Service) UpdateAlertRule(rule *model.AlertRule) error {
	return s.repo.UpdateAlertRule(rule)
}

func (s *Service) DeleteAlertRule(id int64) error {
	return s.repo.DeleteAlertRule(id)
}
