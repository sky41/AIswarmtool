import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from '../views/Dashboard.vue'
import Alerts from '../views/Alerts.vue'
import Policies from '../views/Policies.vue'
const routes = [
  { path: '/', name: 'Dashboard', component: Dashboard },
  { path: '/alerts', name: 'Alerts', component: Alerts },
  { path: '/policies', name: 'Policies', component: Policies }
]
const router = createRouter({
  history: createWebHistory(),
  routes
})
export default router
