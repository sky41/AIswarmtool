-- =============================================
-- Docker Swarm 监控与自愈系统 - 数据库 schema
-- 数据库版本: MySQL 8.0
-- 创建日期: 2026-06-03
-- =============================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- 表结构: metrics_history (监控指标历史表)
-- ----------------------------
DROP TABLE IF EXISTS `metrics_history`;
CREATE TABLE `metrics_history` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `metric_type` varchar(50) NOT NULL COMMENT '指标类型: node/container',
  `resource_id` varchar(100) NOT NULL COMMENT '资源ID: node_id或container_id',
  `resource_name` varchar(100) DEFAULT NULL COMMENT '资源名称',
  `cpu_usage` decimal(5,2) DEFAULT NULL COMMENT 'CPU使用率 (%)',
  `memory_usage` decimal(5,2) DEFAULT NULL COMMENT '内存使用率 (%)',
  `memory_used_bytes` bigint DEFAULT NULL COMMENT '内存使用量 (字节)',
  `memory_total_bytes` bigint DEFAULT NULL COMMENT '内存总量 (字节)',
  `disk_usage` decimal(5,2) DEFAULT NULL COMMENT '磁盘使用率 (%)',
  `network_rx_bytes` bigint DEFAULT NULL COMMENT '网络接收字节数',
  `network_tx_bytes` bigint DEFAULT NULL COMMENT '网络发送字节数',
  `status` varchar(50) DEFAULT NULL COMMENT '状态: healthy/warning/critical/offline',
  `collected_at` datetime NOT NULL COMMENT '采集时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_metric_type` (`metric_type`),
  KEY `idx_resource_id` (`resource_id`),
  KEY `idx_collected_at` (`collected_at`),
  KEY `idx_type_time` (`metric_type`, `collected_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='监控指标历史表';

-- ----------------------------
-- 表结构: alert_rules (告警规则表)
-- ----------------------------
DROP TABLE IF EXISTS `alert_rules`;
CREATE TABLE `alert_rules` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `rule_name` varchar(100) NOT NULL COMMENT '规则名称',
  `rule_desc` varchar(500) DEFAULT NULL COMMENT '规则描述',
  `metric_type` varchar(50) NOT NULL COMMENT '指标类型: cpu/memory/disk/status',
  `resource_type` varchar(50) NOT NULL COMMENT '资源类型: node/container',
  `condition` varchar(20) NOT NULL COMMENT '条件: gt/lt/eq/neq',
  `threshold` decimal(10,2) NOT NULL COMMENT '阈值',
  `duration` int NOT NULL DEFAULT 60 COMMENT '持续时间 (秒)',
  `severity` varchar(20) NOT NULL DEFAULT 'warning' COMMENT '严重程度: info/warning/critical',
  `enabled` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否启用: 0-禁用, 1-启用',
  `notification_type` varchar(50) DEFAULT NULL COMMENT '通知类型: email/webhook/slack',
  `notification_config` json DEFAULT NULL COMMENT '通知配置',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_enabled` (`enabled`),
  KEY `idx_metric_type` (`metric_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='告警规则表';

-- ----------------------------
-- 表结构: incident_logs (故障事件记录表)
-- ----------------------------
DROP TABLE IF EXISTS `incident_logs`;
CREATE TABLE `incident_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `incident_id` varchar(100) NOT NULL COMMENT '事件唯一ID',
  `alert_rule_id` bigint DEFAULT NULL COMMENT '关联的告警规则ID',
  `resource_type` varchar(50) NOT NULL COMMENT '资源类型: node/container',
  `resource_id` varchar(100) NOT NULL COMMENT '资源ID',
  `resource_name` varchar(100) DEFAULT NULL COMMENT '资源名称',
  `severity` varchar(20) NOT NULL COMMENT '严重程度: info/warning/critical',
  `description` text COMMENT '事件描述',
  `status` varchar(20) NOT NULL DEFAULT 'open' COMMENT '状态: open/acknowledged/resolving/resolved/closed',
  `triggered_at` datetime NOT NULL COMMENT '触发时间',
  `acknowledged_at` datetime DEFAULT NULL COMMENT '确认时间',
  `resolved_at` datetime DEFAULT NULL COMMENT '解决时间',
  `healing_policy_id` bigint DEFAULT NULL COMMENT '关联的自愈策略ID',
  `healing_executed` tinyint(1) DEFAULT 0 COMMENT '是否已执行自愈: 0-否, 1-是',
  `healing_result` text COMMENT '自愈执行结果',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_incident_id` (`incident_id`),
  KEY `idx_status` (`status`),
  KEY `idx_resource` (`resource_type`, `resource_id`),
  KEY `idx_triggered_at` (`triggered_at`),
  CONSTRAINT `fk_incident_alert_rule` FOREIGN KEY (`alert_rule_id`) REFERENCES `alert_rules` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='故障事件记录表';

-- ----------------------------
-- 表结构: self_healing_policies (自愈策略配置表)
-- ----------------------------
DROP TABLE IF EXISTS `self_healing_policies`;
CREATE TABLE `self_healing_policies` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `policy_name` varchar(100) NOT NULL COMMENT '策略名称',
  `policy_desc` varchar(500) DEFAULT NULL COMMENT '策略描述',
  `trigger_condition` text NOT NULL COMMENT '触发条件 (JSON格式)',
  `action_type` varchar(50) NOT NULL COMMENT '动作类型: restart/rollback/drain/scale',
  `action_params` json NOT NULL COMMENT '动作参数 (JSON格式)',
  `max_retries` int NOT NULL DEFAULT 3 COMMENT '最大重试次数',
  `retry_delay` int NOT NULL DEFAULT 30 COMMENT '重试延迟 (秒)',
  `enabled` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否启用: 0-禁用, 1-启用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_enabled` (`enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='自愈策略配置表';

-- ----------------------------
-- 表结构: healing_executions (自愈执行记录表)
-- ----------------------------
DROP TABLE IF EXISTS `healing_executions`;
CREATE TABLE `healing_executions` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `policy_id` bigint NOT NULL COMMENT '策略ID',
  `incident_id` varchar(100) NOT NULL COMMENT '关联事件ID',
  `resource_type` varchar(50) NOT NULL COMMENT '资源类型',
  `resource_id` varchar(100) NOT NULL COMMENT '资源ID',
  `action_type` varchar(50) NOT NULL COMMENT '动作类型',
  `attempt` int NOT NULL DEFAULT 1 COMMENT '尝试次数',
  `status` varchar(20) NOT NULL COMMENT '状态: pending/running/success/failed',
  `result` text COMMENT '执行结果',
  `started_at` datetime NOT NULL COMMENT '开始时间',
  `completed_at` datetime DEFAULT NULL COMMENT '完成时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_policy_id` (`policy_id`),
  KEY `idx_incident_id` (`incident_id`),
  KEY `idx_status` (`status`),
  CONSTRAINT `fk_execution_policy` FOREIGN KEY (`policy_id`) REFERENCES `self_healing_policies` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='自愈执行记录表';

-- ----------------------------
-- 初始化数据: 默认告警规则
-- ----------------------------
INSERT INTO `alert_rules` (`rule_name`, `rule_desc`, `metric_type`, `resource_type`, `condition`, `threshold`, `duration`, `severity`, `enabled`) VALUES
('CPU使用率过高-节点', '节点CPU使用率超过80%持续60秒', 'cpu', 'node', 'gt', 80.00, 60, 'warning', 1),
('CPU使用率严重-节点', '节点CPU使用率超过90%持续30秒', 'cpu', 'node', 'gt', 90.00, 30, 'critical', 1),
('内存使用率过高-节点', '节点内存使用率超过85%持续60秒', 'memory', 'node', 'gt', 85.00, 60, 'warning', 1),
('内存使用率严重-节点', '节点内存使用率超过95%持续30秒', 'memory', 'node', 'gt', 95.00, 30, 'critical', 1),
('磁盘使用率过高-节点', '节点磁盘使用率超过80%持续120秒', 'disk', 'node', 'gt', 80.00, 120, 'warning', 1),
('CPU使用率过高-容器', '容器CPU使用率超过80%持续60秒', 'cpu', 'container', 'gt', 80.00, 60, 'warning', 1),
('内存使用率过高-容器', '容器内存使用率超过90%持续60秒', 'memory', 'container', 'gt', 90.00, 60, 'warning', 1);

-- ----------------------------
-- 初始化数据: 默认自愈策略
-- ----------------------------
INSERT INTO `self_healing_policies` (`policy_name`, `policy_desc`, `trigger_condition`, `action_type`, `action_params`, `max_retries`, `retry_delay`, `enabled`) VALUES
('容器自动重启', '当容器异常退出时自动重启', '{\"resource_type\":\"container\",\"status\":\"exited\"}', 'restart', '{\"delay_seconds\":5}', 3, 30, 1),
('节点Drain', '当节点NotReady时自动Drain', '{\"resource_type\":\"node\",\"status\":\"notready\"}', 'drain', '{\"force\":false}', 1, 60, 1);

SET FOREIGN_KEY_CHECKS = 1;
