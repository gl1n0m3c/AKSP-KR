import { createRouter, createWebHistory } from 'vue-router'
import CalendarView from './views/CalendarView.vue'
import AdminView from './views/AdminView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: CalendarView },
    { path: '/admin', component: AdminView },
  ],
})

export default router

