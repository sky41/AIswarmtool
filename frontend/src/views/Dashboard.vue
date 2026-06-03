<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon :size="40" color="#409EFF"><Monitor /></el-icon>
            <div class="stat-info">
              <div class="stat-value">{{ stats.nodes }}</div>
              <div class="stat-label">集群节点</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon :size="40" color="#67C23A"><Box /></el-icon>
            <div class="stat-info">
              <div class="stat-value">{{ stats.containers }}</div>
              <div class="stat-label">运行容器</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon :size="40" color="#E6A23C"><Warning /></el-icon>
            <div class="stat-info">
              <div class="stat-value">{{ stats.activeAlerts }}</div>
              <div class="stat-label">活动告警</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon :size="40" color="#F56C6C"><Tools /></el-icon>
            <div class="stat-info">
              <div class="stat-value">{{ stats.healingActions }}</div>
              <div class="stat-label">自愈动作</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>节点状态</span>
          </template>
          <el-table :data="nodes" style="width: 100%">
            <el-table-column prop="resource_name" label="节点名称" />
            <el-table-column prop="status" label="状态">
              <template #default="scope">
                <el-tag :type="getStatusType(scope.row.status)">{{ scope.row.status }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="collected_at" label="采集时间" />
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>容器状态</span>
          </template>
          <el-table :data="containers" style="width: 100%">
            <el-table-column prop="resource_name" label="容器名称" />
            <el-table-column prop="status" label="状态">
              <template #default="scope">
                <el-tag :type="getStatusType(scope.row.status)">{{ scope.row.status }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="collected_at" label="采集时间" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>
<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { metricsApi } from '../api'
const stats = ref({
  nodes: 0,
  containers: 0,
  activeAlerts: 0,
  healingActions: 0
})
const nodes = ref([])
const containers = ref([])
let refreshTimer = null
const loadData = async () => {
  try {
    const [nodeRes, containerRes] = await Promise.all([
      metricsApi.getLatestMetrics('node'),
      metricsApi.getLatestMetrics('container')
    ])
    nodes.value = nodeRes.data.data || []
    containers.value = containerRes.data.data || []
    stats.value.nodes = nodes.value.length
    stats.value.containers = containers.value.length
    stats.value.activeAlerts = nodes.value.filter(n => n.status !== 'healthy').length
  } catch (error) {
    console.error('Load data error:', error)
  }
}
const getStatusType = (status) => {
  const map = { healthy: 'success', warning: 'warning', critical: 'danger' }
  return map[status] || 'info'
}
onMounted(() => {
  loadData()
  refreshTimer = setInterval(loadData, 10000)
})
onUnmounted(() => {
  if (refreshTimer) clearInterval(refreshTimer)
})
</script>
<style scoped>
.dashboard {
  padding: 0;
}
.stat-card {
  margin-bottom: 20px;
}
.stat-content {
  display: flex;
  align-items: center;
  gap: 20px;
}
.stat-info {
  flex: 1;
}
.stat-value {
  font-size: 32px;
  font-weight: bold;
  color: #303133;
}
.stat-label {
  font-size: 14px;
  color: #909399;
}
</style>
