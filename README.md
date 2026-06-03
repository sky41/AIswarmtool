# Docker Swarm 监控与自愈系统

基于 Docker Swarm 集群的实时监控与故障自愈智能运维系统。

## 技术栈

- **后端**: Go + Gin + GORM
- **前端**: Vue 3 + Element Plus + ECharts
- **数据库**: MySQL 8.0
- **容器技术**: Docker Engine API

## 项目结构

```
AIswarmtool/
├── backend/              # Go 后端
│   ├── cmd/              # 程序入口
│   ├── internal/         # 内部代码
│   │   ├── model/        # 数据模型
│   │   ├── repository/   # 数据访问层
│   │   ├── service/      # 业务逻辑层
│   │   └── handler/      # API 处理器
│   ├── pkg/              # 公共包
│   │   ├── config/       # 配置
│   │   └── docker/       # Docker 集成
│   └── go.mod            # Go 依赖
├── frontend/             # Vue 前端
│   ├── src/
│   │   ├── views/        # 页面组件
│   │   ├── api/          # API 封装
│   │   └── router/       # 路由
│   └── package.json
├── database/             # 数据库脚本
│   └── schema.sql
├── docs/                 # 文档
│   ├── requirements.json # 需求规格
│   └── openapi.yaml      # API 文档
└── README.md
```

## 功能特性

1. **集群监控**
   - 实时监控节点与容器状态
   - 采集 CPU、内存等资源指标
   - 定时刷新与告警通知

2. **告警系统**
   - 灵活的告警规则配置
   - 多级别严重程度管理
   - 告警确认与解决流程

3. **自愈策略**
   - 自动重启异常容器
   - 节点故障时自动排空
   - 可配置的自愈策略

4. **管理界面**
   - 监控大盘可视化
   - 告警与事件中心
   - 策略配置管理

## 快速开始

### 1. 数据库初始化

```bash
# 创建数据库
mysql -u root -p
CREATE DATABASE swarm_monitor CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 导入表结构
mysql -u root -p swarm_monitor < database/schema.sql
```

### 2. 配置后端

复制并修改配置示例：

```bash
cd backend
# 创建 config.yaml（可选）
```

后端使用环境变量或配置文件：

```yaml
server:
  port: 8080
database:
  host: localhost
  port: 3306
  user: root
  password: your_password
  dbname: swarm_monitor
docker:
  host: unix:///var/run/docker.sock
monitoring:
  interval: 10
```

### 3. 启动后端

```bash
cd backend
go mod download
go run cmd/main.go
```

### 4. 启动前端

```bash
cd frontend
npm install
npm run dev
```

前端将在 `http://localhost:3000` 启动。

## API 文档

OpenAPI 规范文档请查看 [docs/openapi.yaml](docs/openapi.yaml)。

主要 API 端点：

- `/api/v1/metrics` - 监控指标
- `/api/v1/alerts` - 告警与事件
- `/api/v1/alerts/rules` - 告警规则
- `/api/v1/healing/policies` - 自愈策略
- `/api/v1/cluster` - 集群信息

## 核心模块说明

### 监控采集器

定时从 Docker API 获取节点与容器状态，采集资源指标。

### 告警引擎

根据配置规则评估监控数据，触发告警并记录事件。

### 自愈执行器

根据策略自动执行修复动作，如重启容器、排空节点等。

## 默认配置

系统已预设以下告警规则：

- CPU > 80%（节点/容器）
- 内存 > 85%（节点）
- 磁盘 > 80%（节点）

预设自愈策略：

- 容器异常退出时自动重启
- 节点故障时自动排空

## 注意事项

1. 确保后端服务有权访问 Docker Socket
2. 生产环境建议配置高可用部署
3. 监控数据量大时需考虑数据保留策略

## License

MIT
