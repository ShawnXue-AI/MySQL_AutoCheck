import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import HomeView from './views/HomeView.vue'
import TasksView from './views/TasksView.vue'
import AdminView from './views/AdminView.vue'
import BenchmarksView from './views/BenchmarksView.vue'
import PersonDayView from './views/PersonDayView.vue'

const routes = [
  { path: '/', name: 'home', component: HomeView },
  { path: '/tasks', name: 'tasks', component: TasksView },
  { path: '/admin', name: 'admin', component: AdminView },
  { path: '/benchmarks', name: 'benchmarks', component: BenchmarksView },
  { path: '/persondays', name: 'persondays', component: PersonDayView }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

const app = createApp(App)
app.use(router)
app.mount('#app')
