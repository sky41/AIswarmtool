<template>
  <div class="alerts">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>告警与事件中心</span>
          <el-button type="primary" @click="loadData">刷新</el-button>
        </div>
      </template>
      <el-tabs v-model="activeTab">
        <el-tab-pane label="活动告警" name="active">
          <el-table :data="activeIncidents" style="width: 100%">
            <el-table-column prop="incident_id" label="事件ID" width="180" />
            <el-table-column prop="resource_type" label="资源类型" width="120" />
            <el-table-column prop="resource_name" label="资源名称" />
            <el-table-column prop="severity" label="严重程度" width="120">
              <template #default="scope">
                <el-tag :type="getSeverityType(scope.row.severity)">{{ scope.row.severity }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="120">
              <template #default="scope">
                <el-tag :type="getStatusType(scope.row.status)">{{ scope.row.status }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="triggered_at" label="触发时间" width="180" />
            <el-table-column label="操作" width="200">
              <template #default="scope">
                <el-button size="small" @click="acknowledge(scope.row)">确认</el-button>
                <el-button size="small" type="success" @click="resolve(scope.row)">解决</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="历史记录" name="history">
          <el-table :data="historyIncidents" style="width: 100%">
            <el-table-column prop="incident_id" label="事件ID" width="180" />
            <el-table-column prop="resource_type" label="资源类型" width="120" />
            <el-table-column prop="resource_name" label="资源名称" />
            <el-table-column prop="severity" label="严重程度" width="120">
              <template #default="scope">
                <el-tag :type="getSeverityType(scope.row.severity)">{{ scope.row.severity }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="120" />
            <el-table-column prop="triggered_at" label="触发时间" width="180" />
            <el-table-column prop="resolved_at" label="解决时间" width="180" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="告警规则" name="rules">
          <div style="margin-bottom: 20px;">
            <el-button type="primary" @click="showRuleDialog">新增规则</el-button>
          </div>
          <el-table :data="alertRules" style="width: 100%">
            <el-table-column prop="rule_name" label="规则名称" />
            <el-table-column prop="metric_type" label="指标类型" width="120" />
            <el-table-column prop="resource_type" label="资源类型" width="120" />
            <el-table-column prop="condition" label="条件" width="80" />
            <el-table-column prop="threshold" label="阈值" width="100" />
            <el-table-column prop="severity" label="严重程度" width="120">
              <template #default="scope">
                <el-tag :type="getSeverityType(scope.row.severity)">{{ scope.row.severity }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="enabled" label="启用" width="80">
              <template #default="scope">
                <el-tag :type="scope.row.enabled ? 'success' : 'info'">
                  {{ scope.row.enabled ? '是' : '否' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200">
              <template #default="scope">
                <el-button size="small" @click="editRule(scope.row)">编辑</el-button>
                <el-button size="small" type="danger" @click="deleteRule(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>
    <el-dialog v-model="ruleDialogVisible" :title="isEditRule ? '编辑告警规则' : '新增告警规则'" width="500px">
      <el-form :model="ruleForm" label-width="100px">
        <el-form-item label="规则名称">
          <el-input v-model="ruleForm.rule_name" />
        </el-form-item>
        <el-form-item label="规则描述">
          <el-input v-model="ruleForm.rule_desc" type="textarea" />
        </el-form-item>
        <el-form-item label="指标类型">
          <el-select v-model="ruleForm.metric_type">
            <el-option label="CPU" value="cpu" />
            <el-option label="内存" value="memory" />
            <el-option label="磁盘" value="disk" />
            <el-option label="状态" value="status" />
          </el-select>
        </el-form-item>
        <el-form-item label="资源类型">
          <el-select v-model="ruleForm.resource_type">
            <el-option label="节点" value="node" />
            <el-option label="容器" value="container" />
          </el-select>
        </el-form-item>
        <el-form-item label="条件">
          <el-select v-model="ruleForm.condition">
            <el-option label=">" value="gt" />
            <el-option label="<" value="lt" />
            <el-option label="=" value="eq" />
            <el-option label="!=" value="neq" />
          </el-select>
        </el-form-item>
        <el-form-item label="阈值">
          <el-input-number v-model="ruleForm.threshold" :min="0" :max="100" />
        </el-form-item>
        <el-form-item label="严重程度">
          <el-select v-model="ruleForm.severity">
            <el-option label="信息" value="info" />
            <el-option label="警告" value="warning" />
            <el-option label="严重" value="critical" />
          </el-select>
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="ruleForm.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="ruleDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveRule">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { alertsApi } from '../api'
const activeTab = ref('active')
const activeIncidents = ref([])
const historyIncidents = ref([])
const alertRules = ref([])
const ruleDialogVisible = ref(false)
const isEditRule = ref(false)
const ruleForm = ref({
  rule_name: '',
  rule_desc: '',
  metric_type: 'cpu',
  resource_type: 'node',
  condition: 'gt',
  threshold: 80,
  severity: 'warning',
  enabled: true
})
const loadData = async () => {
  try {
    const [incidentsRes, rulesRes] = await Promise.all([
      alertsApi.getIncidents({ limit: 100 }),
      alertsApi.getAlertRules()
    ])
    const allIncidents = incidentsRes.data.data || []
    activeIncidents.value = allIncidents.filter(i => ['open', 'acknowledged'].includes(i.status))
    historyIncidents.value = allIncidents.filter(i => ['resolved', 'closed'].includes(i.status))
    alertRules.value = rulesRes.data.data || []
  } catch (error) {
    ElMessage.error('加载数据失败')
  }
}
const acknowledge = async (row) => {
  try {
    await alertsApi.acknowledgeIncident(row.id)
    ElMessage.success('已确认')
    loadData()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}
const resolve = async (row) => {
  try {
    await alertsApi.resolveIncident(row.id)
    ElMessage.success('已标记为解决')
    loadData()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}
const showRuleDialog = (rule = null) => {
  isEditRule.value = !!rule
  if (rule) {
    ruleForm.value = { ...rule }
  } else {
    ruleForm.value = {
      rule_name: '',
      rule_desc: '',
      metric_type: 'cpu',
      resource_type: 'node',
      condition: 'gt',
      threshold: 80,
      severity: 'warning',
      enabled: true
    }
  }
  ruleDialogVisible.value = true
}
const editRule = (row) => showRuleDialog(row)
const saveRule = async () => {
  try {
    if (isEditRule.value) {
      await alertsApi.updateAlertRule(ruleForm.value.id, ruleForm.value)
    } else {
      await alertsApi.createAlertRule(ruleForm.value)
    }
    ElMessage.success('保存成功')
    ruleDialogVisible.value = false
    loadData()
  } catch (error) {
    ElMessage.error('保存失败')
  }
}
const deleteRule = async (row) => {
  try {
    await ElMessageBox.confirm('确认删除该规则吗？')
    await alertsApi.deleteAlertRule(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}
const getSeverityType = (s) => ({ info: 'info', warning: 'warning', critical: 'danger' }[s] || 'info')
const getStatusType = (s) => {
  const map = { open: 'danger', acknowledged: 'warning', resolved: 'success', closed: 'info' }
  return map[s] || 'info'
}
onMounted(loadData)
</script>
<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
