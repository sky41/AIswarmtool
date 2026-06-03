package service

import (
	"context"
	"fmt"
	"time"
	"github.com/aiswarmtool/backend/internal/model"
	"github.com/aiswarmtool/backend/pkg/docker"
	"github.com/google/uuid"
)

type Monitor struct {
	service     *Service
	dockerCli   *docker.Client
	alertEngine *AlertEngine
	healer      *Healer
	stopChan    chan struct{}
}

func NewMonitor(service *Service, dockerCli *docker.Client) *Monitor {
	m := &Monitor{
		service:   service,
		dockerCli: dockerCli,
		stopChan:  make(chan struct{}),
	}
	m.alertEngine = NewAlertEngine(service)
	m.healer = NewHealer(service, dockerCli)
	return m
}

func (m *Monitor) Start(interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			m.collect()
		case <-m.stopChan:
			return
		}
	}
}

func (m *Monitor) Stop() {
	close(m.stopChan)
}

func (m *Monitor) collect() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	m.collectNodes(ctx)
	m.collectContainers(ctx)
}

func (m *Monitor) collectNodes(ctx context.Context) {
	nodes, err := m.dockerCli.GetNodes(ctx)
	if err != nil {
		return
	}
	for _, node := range nodes {
		var status string
		if node.Status.State == "ready" {
			status = "healthy"
		} else {
			status = "critical"
		}
		metrics := &model.MetricsHistory{
			MetricType:  "node",
			ResourceID:  node.ID,
			ResourceName: node.Description.Hostname,
			Status:      status,
			CollectedAt: time.Now(),
		}
		m.service.CreateMetrics(metrics)
		m.alertEngine.EvaluateNode(node)
		if status == "critical" {
			m.healer.HandleNodeIssue(node.ID, node.Description.Hostname)
		}
	}
}

func (m *Monitor) collectContainers(ctx context.Context) {
	containers, err := m.dockerCli.GetContainers(ctx)
	if err != nil {
		return
	}
	for _, container := range containers {
		var status string
		switch container.State {
		case "running":
			status = "healthy"
		case "exited", "dead":
			status = "critical"
		default:
			status = "warning"
		}
		var name string
		if len(container.Names) > 0 {
			name = container.Names[0][1:]
		}
		metrics := &model.MetricsHistory{
			MetricType:   "container",
			ResourceID:   container.ID,
			ResourceName: name,
			Status:       status,
			CollectedAt:  time.Now(),
		}
		m.service.CreateMetrics(metrics)
		if status == "critical" {
			m.alertEngine.EvaluateContainer(container)
			m.healer.HandleContainerIssue(container.ID, name, container.State)
		}
	}
}

type AlertEngine struct {
	service *Service
}

func NewAlertEngine(service *Service) *AlertEngine {
	return &AlertEngine{service: service}
}

func (e *AlertEngine) EvaluateNode(node interface{}) {
	rules, _ := e.service.GetAlertRules(true)
	for _, rule := range rules {
		if rule.ResourceType != "node" {
			continue
		}
		incidentID := uuid.New().String()
		desc := fmt.Sprintf("Node alert triggered by rule: %s", rule.RuleName)
		incident := &model.IncidentLog{
			IncidentID:   incidentID,
			AlertRuleID:  &rule.ID,
			ResourceType: "node",
			Severity:     rule.Severity,
			Description:  &desc,
			Status:       "open",
			TriggeredAt:  time.Now(),
		}
		e.service.repo.CreateIncident(incident)
	}
}

func (e *AlertEngine) EvaluateContainer(container interface{}) {
	rules, _ := e.service.GetAlertRules(true)
	for _, rule := range rules {
		if rule.ResourceType != "container" {
			continue
		}
		incidentID := uuid.New().String()
		desc := fmt.Sprintf("Container alert triggered by rule: %s", rule.RuleName)
		incident := &model.IncidentLog{
			IncidentID:   incidentID,
			AlertRuleID:  &rule.ID,
			ResourceType: "container",
			Severity:     rule.Severity,
			Description:  &desc,
			Status:       "open",
			TriggeredAt:  time.Now(),
		}
		e.service.repo.CreateIncident(incident)
	}
}

type Healer struct {
	service   *Service
	dockerCli *docker.Client
}

func NewHealer(service *Service, dockerCli *docker.Client) *Healer {
	return &Healer{
		service:   service,
		dockerCli: dockerCli,
	}
}

func (h *Healer) HandleNodeIssue(nodeID, nodeName string) {
	policies, _ := h.service.GetHealingPolicies(true)
	for _, policy := range policies {
		if policy.ActionType == "drain" {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			execution := &model.HealingExecution{
				PolicyID:     policy.ID,
				IncidentID:   uuid.New().String(),
				ResourceType: "node",
				ResourceID:   nodeID,
				ActionType:   policy.ActionType,
				Status:       "running",
				StartedAt:    time.Now(),
			}
			h.service.repo.CreateHealingExecution(execution)
			if err := h.dockerCli.DrainNode(ctx, nodeID); err != nil {
				msg := err.Error()
				execution.Status = "failed"
				execution.Result = &msg
			} else {
				msg := "Node drained successfully"
				execution.Status = "success"
				execution.Result = &msg
			}
			now := time.Now()
			execution.CompletedAt = &now
			h.service.repo.UpdateHealingExecution(execution)
			break
		}
	}
}

func (h *Healer) HandleContainerIssue(containerID, containerName, state string) {
	policies, _ := h.service.GetHealingPolicies(true)
	for _, policy := range policies {
		if policy.ActionType == "restart" {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			execution := &model.HealingExecution{
				PolicyID:     policy.ID,
				IncidentID:   uuid.New().String(),
				ResourceType: "container",
				ResourceID:   containerID,
				ActionType:   policy.ActionType,
				Status:       "running",
				StartedAt:    time.Now(),
			}
			h.service.repo.CreateHealingExecution(execution)
			if err := h.dockerCli.RestartContainer(ctx, containerID); err != nil {
				msg := err.Error()
				execution.Status = "failed"
				execution.Result = &msg
			} else {
				msg := "Container restarted successfully"
				execution.Status = "success"
				execution.Result = &msg
			}
			now := time.Now()
			execution.CompletedAt = &now
			h.service.repo.UpdateHealingExecution(execution)
			break
		}
	}
}
