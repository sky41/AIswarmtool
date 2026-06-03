<template>
  <div class="policies">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>自愈策略配置</span>
          <el-button type="primary" @click="showPolicyDialog">新增策略</el-button>
        </div>
      </template>
      <el-table :data="healingPolicies" style="width: 100%">
        <el-table-column prop="policy_name" label="策略名称" />
        <el-table-column prop="policy_desc" label="策略描述" />
        <el-table-column prop="action_type" label="动作类型" width="120">
          <template #default="scope">
            <el-tag>{{ getActionTypeName(scope.row.action_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="max_retries" label="最大重试" width="100" />
        <el-table-column prop="retry_delay" label="重试延迟(秒)" width="120" />
        <el-table-column prop="enabled" label="启用" width="80">
          <template #default="scope">
            <el-tag :type="scope.row.enabled ? 'success' : 'info'">
              {{ scope.row.enabled ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <el-button size="small" @click="editPolicy(scope.row)">编辑</el-button>
            <el-button size="small" type="danger" @click="deletePolicy(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <el-card style="margin-top: 20px;">
      <template #header>
        <span>自愈执行记录</span>
      </template>
      <el-table :data="healingExecutions" style="width: 100%">
        <el-table-column prop="incident_id" label="事件ID" width="180" />
        <el-table-column prop="resource_type" label="资源类型" width="120" />
        <el-table-column prop="resource_id" label="资源ID" width="120" />
        <el-table-column prop="action_type" label="动作类型" width="120">
          <template #default="scope">
            <el-tag>{{ getActionTypeName(scope.row.action_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="attempt" label="尝试次数" width="100" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'success' ? 'success' : 'danger'">
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="started_at" label="开始时间" width="180" />
        <el-table-column prop="completed_at" label="完成时间" width="180" />
      </el-table>
    </el-card>
    <el-dialog v-model="policyDialogVisible" :title="isEditPolicy ? '编辑自愈策略' : '新增自愈策略'" width="600px">
      <el-form :model="policyForm" label-width="120px">
        <el-form-item label="策略名称">
          <el-input v-model="policyForm.policy_name" />
        </el-form-item>
        <el-form-item label="策略描述">
          <el-input v-model="policyForm.policy_desc" type="textarea" />
        </el-form-item>
        <el-form-item label="触发条件">
          <el-input v-model="policyForm.trigger_condition" type="textarea" placeholder='例如: {"resource_type":"container","status":"exited"}' />
        </el-form-item>
        <el-form-item label="动作类型">
          <el-select v-model="policyForm.action_type">
            <el-option label="重启容器" value="restart" />
            <el-option label="服务回滚" value="rollback" />
            <el-option label="节点排空" value="drain" />
            <el-option label="扩容" value="scale" />
          </el-select>
        </el-form-item>
        <el-form-item label="动作参数">
          <el-input v-model="policyForm.action_params_str" type="textarea" placeholder='例如: {"delay_seconds":5}' />
        </el-form-item>
        <el-form-item label="最大重试次数">
          <el-input-number v-model="policyForm.max_retries" :min="0" :max="10" />
        </el-form-item>
        <el-form-item label="重试延迟(秒)">
          <el-input-number v-model="policyForm.retry_delay" :min="1" :max="300" />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="policyForm.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="policyDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="savePolicy">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { healingApi } from '../api'
const healingPolicies = ref([])
const healingExecutions = ref([])
const policyDialogVisible = ref(false)
const isEditPolicy = ref(false)
const policyForm = ref({
  policy_name: '',
  policy_desc: '',
  trigger_condition: '',
  action_type: 'restart',
  action_params_str: '',
  max_retries: 3,
  retry_delay: 30,
  enabled: true
})
const loadData = async () => {
  try {
    const [policiesRes, executionsRes] = await Promise.all([
      healingApi.getPolicies(),
      healingApi.getExecutions()
    ])
    healingPolicies.value = policiesRes.data.data || []
    healingExecutions.value = executionsRes.data.data || []
  } catch (error) {
    ElMessage.error('加载数据失败')
  }
}
const getActionTypeName = (type) => {
  const map = { restart: '重启容器', rollback: '服务回滚', drain: '节点排空', scale: '扩容' }
  return map[type] || type
}
const showPolicyDialog = (policy = null) => {
  isEditPolicy.value = !!policy
  if (policy) {
    policyForm.value = {
      ...policy,
      action_params_str: JSON.stringify(policy.action_params, null, 2)
    }
  } else {
    policyForm.value = {
      policy_name: '',
      policy_desc: '',
      trigger_condition: '',
      action_type: 'restart',
      action_params_str: '{}',
      max_retries: 3,
      retry_delay: 30,
      enabled: true
    }
  }
  policyDialogVisible.value = true
}
const editPolicy = (row) => showPolicyDialog(row)
const savePolicy = async () => {
  try {
    let params = { ...policyForm.value }
    try {
      params.action_params = JSON.parse(params.action_params_str)
    } catch (e) {
      ElMessage.error('动作参数JSON格式错误')
      return
    }
    delete params.action_params_str
    if (isEditPolicy.value) {
      await healingApi.updatePolicy(params.id, params)
    } else {
      await healingApi.createPolicy(params)
    }
    ElMessage.success('保存成功')
    policyDialogVisible.value = false
    loadData()
  } catch (error) {
    ElMessage.error('保存失败')
  }
}
const deletePolicy = async (row) => {
  try {
    await ElMessageBox.confirm('确认删除该策略吗？')
    await healingApi.deletePolicy(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
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
