import axios from 'axios'
const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000
})
api.interceptors.request.use(
  config => {
    return config
  },
  error => {
    return Promise.reject(error)
  }
)
api.interceptors.response.use(
  response => {
    return response
  },
  error => {
    if (error.response) {
      const { status } = error.response
      if (status === 401) {
        console.error('认证失败，请重新登录')
      } else if (status === 500) {
        console.error('服务器错误，请稍后重试')
      }
    } else if (error.request) {
      console.error('网络错误，请检查网络连接')
    } else {
      console.error('请求配置错误', error.message)
    }
    return Promise.reject(error)
  }
)
export const metricsApi = {
  getMetrics: (params) => api.get('/metrics', { params }),
  getLatestMetrics: (type) => api.get('/metrics/latest', { params: { type } })
}
export const alertsApi = {
  getIncidents: (params) => api.get('/alerts', { params }),
  getIncident: (id) => api.get(`/alerts/${id}`),
  acknowledgeIncident: (id) => api.post(`/alerts/${id}/acknowledge`),
  resolveIncident: (id) => api.post(`/alerts/${id}/resolve`),
  getAlertRules: (enabledOnly) => api.get('/alerts/rules', { params: { enabled_only: enabledOnly } }),
  getAlertRule: (id) => api.get(`/alerts/rules/${id}`),
  createAlertRule: (data) => api.post('/alerts/rules', data),
  updateAlertRule: (id, data) => api.put(`/alerts/rules/${id}`, data),
  deleteAlertRule: (id) => api.delete(`/alerts/rules/${id}`)
}
export const healingApi = {
  getPolicies: (enabledOnly) => api.get('/healing/policies', { params: { enabled_only: enabledOnly } }),
  getPolicy: (id) => api.get(`/healing/policies/${id}`),
  createPolicy: (data) => api.post('/healing/policies', data),
  updatePolicy: (id, data) => api.put(`/healing/policies/${id}`, data),
  deletePolicy: (id) => api.delete(`/healing/policies/${id}`),
  getExecutions: (incidentId) => api.get('/healing/executions', { params: { incident_id: incidentId } })
}
export const clusterApi = {
  getClusterInfo: () => api.get('/cluster'),
  getNodes: () => api.get('/cluster/nodes'),
  getServices: () => api.get('/cluster/services')
}
export default api
