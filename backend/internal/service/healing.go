package service

import (
	"github.com/aiswarmtool/backend/internal/model"
)

func (s *Service) GetHealingPolicies(enabledOnly bool) ([]*model.SelfHealingPolicy, error) {
	return s.repo.GetHealingPolicies(enabledOnly)
}

func (s *Service) GetHealingPolicy(id int64) (*model.SelfHealingPolicy, error) {
	return s.repo.GetHealingPolicyByID(id)
}

func (s *Service) CreateHealingPolicy(policy *model.SelfHealingPolicy) error {
	return s.repo.CreateHealingPolicy(policy)
}

func (s *Service) UpdateHealingPolicy(policy *model.SelfHealingPolicy) error {
	return s.repo.UpdateHealingPolicy(policy)
}

func (s *Service) DeleteHealingPolicy(id int64) error {
	return s.repo.DeleteHealingPolicy(id)
}

func (s *Service) GetHealingExecutions(incidentID string) ([]*model.HealingExecution, error) {
	return s.repo.GetHealingExecutionsByIncident(incidentID)
}