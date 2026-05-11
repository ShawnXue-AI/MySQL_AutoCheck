<template>
  <div class="app-container">
    <aside class="sidebar">
      <div class="logo-section">
        <div class="logo-icon">
          <svg viewBox="0 0 40 40" fill="none"><rect width="40" height="40" rx="10" fill="#238636" opacity="0.9"/><path d="M12 20l6 6 10-10" stroke="#fff" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/></svg>
        </div>
        <div class="logo-text">
          <span class="logo-title">Inspection</span>
          <span class="logo-sub">MySQL 巡检系统</span>
        </div>
      </div>
      <nav class="nav-menu">
        <router-link to="/" class="nav-item" :class="{ active: $route.path === '/' }">
          <svg viewBox="0 0 24 24" class="nav-icon"><path d="M3 13h8V3H3v10zm0 8h8v-6H3v6zm10 0h8V11h-8v10zm0-18v6h8V3h-8z" fill="currentColor"/></svg>
          <span>报告生成</span>
        </router-link>
        <router-link to="/tasks" class="nav-item" :class="{ active: $route.path === '/tasks' }">
          <svg viewBox="0 0 24 24" class="nav-icon"><path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-7 14l-5-5 1.41-1.41L12 14.17l7.59-7.59L21 8l-9 9z" fill="currentColor"/></svg>
          <span>任务记录</span>
        </router-link>
        <router-link to="/admin" class="nav-item" :class="{ active: $route.path === '/admin' }">
          <svg viewBox="0 0 24 24" class="nav-icon"><path d="M19.43 12.98c.04-.32.07-.64.07-.98s-.03-.66-.07-.98l2.11-1.65c.19-.15.24-.42.12-.64l-2-3.46c-.12-.22-.39-.3-.61-.22l-2.49 1c-.52-.4-1.08-.73-1.69-.98l-.38-2.65C14.46 2.18 14.25 2 14 2h-4c-.25 0-.46.18-.49.42l-.38 2.65c-.61.25-1.17.59-1.69.98l-2.49-1c-.23-.09-.49 0-.61.22l-2 3.46c-.13.22-.07.49.12.64l2.11 1.65c-.04.32-.07.65-.07.98s.03.66.07.98l-2.11 1.65c-.19.15-.24.42-.12.64l2 3.46c.12.22.39.3.61.22l2.49-1c.52.4 1.08.73 1.69.98l.38 2.65c.03.24.24.42.49.42h4c.25 0 .46-.18.49-.42l.38-2.65c.61-.25 1.17-.59 1.69-.98l2.49 1c.23.09.49 0 .61-.22l2-3.46c.12-.22.07-.49-.12-.64l-2.11-1.65zM12 15.5c-1.93 0-3.5-1.57-3.5-3.5s1.57-3.5 3.5-3.5 3.5 1.57 3.5 3.5-1.57 3.5-3.5 3.5z" fill="currentColor"/></svg>
          <span>后端管理</span>
        </router-link>
        <router-link to="/benchmarks" class="nav-item" :class="{ active: $route.path === '/benchmarks' }">
          <svg viewBox="0 0 24 24" class="nav-icon"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z" fill="currentColor"/></svg>
          <span>巡检基准值</span>
        </router-link>
      </nav>
      <div class="sidebar-bottom">
        <div class="sidebar-footer">
          <div class="stats-mini">
            <div class="stat-item">
              <span class="stat-value accent">{{ stats.active_users }}</span>
              <span class="stat-label">在线</span>
            </div>
            <div class="stat-divider"></div>
            <div class="stat-item">
              <span class="stat-value">{{ stats.total_completed }}</span>
              <span class="stat-label">累计</span>
            </div>
            <div class="stat-divider" v-if="stats.running_tasks > 0"></div>
            <div class="stat-item" v-if="stats.running_tasks > 0">
              <span class="stat-value pulse">{{ stats.running_tasks }}</span>
              <span class="stat-label">运行中</span>
            </div>
          </div>
        </div>
        <div class="sidebar-time">
          <span class="time-date">{{ currentDate }}</span>
          <span class="time-clock">{{ currentTime }}</span>
          <span class="time-city">{{ currentCity }}</span>
        </div>
      </div>
    </aside>
    <main class="main-content">
      <header class="top-bar">
        <h1 class="page-title">{{ $route.path === '/' ? '报告生成' : $route.path === '/admin' ? '后端管理' : $route.path === '/benchmarks' ? '巡检基准值' : '任务记录' }}</h1>
      </header>
      <div class="content-area">
        <router-view @stats-update="updateStats" />
      </div>
    </main>
    <Toast ref="toast" />
  </div>
</template>

<script>
import { getStats } from './api/index.js'
import Toast from './components/Toast.vue'

export default {
  name: 'App',
  components: { Toast },
  data() {
    return {
      stats: { active_users: 0, total_completed: 0, running_tasks: 0 },
      currentDate: '',
      currentTime: '',
      currentCity: '北京',
      pollTimer: null,
      clockTimer: null
    }
  },
  mounted() {
    this.pollStats()
    this.updateClock()
    this.detectCity()
    this.pollTimer = setInterval(() => this.pollStats(), 5000)
    this.clockTimer = setInterval(() => this.updateClock(), 1000)
  },
  beforeUnmount() {
    if (this.pollTimer) clearInterval(this.pollTimer)
    if (this.clockTimer) clearInterval(this.clockTimer)
  },
  methods: {
    async pollStats() {
      try {
        const { data } = await getStats()
        this.stats = data
      } catch {}
    },
    updateClock() {
      const now = new Date()
      this.currentDate = now.toLocaleDateString('zh-CN', {
        year: 'numeric', month: '2-digit', day: '2-digit'
      })
      this.currentTime = now.toLocaleTimeString('zh-CN', {
        hour: '2-digit', minute: '2-digit', second: '2-digit',
        hour12: false
      })
    },
    async detectCity() {
      try {
        const res = await fetch('https://ipapi.co/json/')
        const data = await res.json()
        if (data.city) {
          this.currentCity = data.city
        }
      } catch {}
    },
    updateStats(stats) {
      if (stats) this.stats = stats
    }
  }
}
</script>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }

body {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'PingFang SC', 'Microsoft YaHei', sans-serif;
  background: #0d1117;
  color: #c9d1d9;
  -webkit-font-smoothing: antialiased;
}

.app-container {
  display: flex;
  height: 100vh;
  overflow: hidden;
}

/* ===== Sidebar ===== */
.sidebar {
  width: 240px;
  background: #161b22;
  border-right: 1px solid #21262d;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  height: 100vh;
}

.logo-section {
  display: flex;
  align-items: center;
  gap: 11px;
  height: 56px;
  padding: 0 20px;
  border-bottom: 1px solid #21262d;
  flex-shrink: 0;
}

.logo-icon {
  width: 32px; height: 32px; flex-shrink: 0;
}

.logo-text {
  display: flex; flex-direction: column;
}

.logo-title {
  font-size: 15px; font-weight: 700; color: #f0f6fc; letter-spacing: -0.3px; line-height: 1.2;
}

.logo-sub {
  font-size: 10.5px; color: #8b949e; letter-spacing: 0.4px; text-transform: uppercase; line-height: 1.2;
}

.nav-menu {
  padding: 10px 10px; display: flex; flex-direction: column; gap: 1px; flex-shrink: 0;
}

.nav-item {
  display: flex; align-items: center; gap: 11px; padding: 9px 14px;
  border-radius: 8px; color: #8b949e; text-decoration: none;
  font-size: 13.5px; font-weight: 500; transition: all 0.15s;
}

.nav-item:hover { background: #1c2128; color: #c9d1d9; }
.nav-item.active { background: rgba(35,134,54,0.15); color: #3fb950; }

.nav-icon { width: 19px; height: 19px; flex-shrink: 0; }

/* ===== Sidebar Bottom (stats + time locked together) ===== */
.sidebar-bottom {
  margin-top: auto;
  flex-shrink: 0;
}

.sidebar-footer {
  padding: 10px 20px 8px;
  border-top: 1px solid #21262d;
}

.stats-mini { display: flex; align-items: center; justify-content: center; gap: 16px; }

.stat-item { display: flex; flex-direction: column; align-items: center; gap: 1px; }

.stat-value {
  font-size: 20px; font-weight: 700; color: #f0f6fc; line-height: 1.2;
}

.stat-value.accent { color: #3fb950; }
.stat-value.pulse { color: #3fb950; animation: stat-pulse 2s ease-in-out infinite; }

@keyframes stat-pulse {
  0%, 100% { opacity: 1; } 50% { opacity: 0.5; }
}

.stat-label {
  font-size: 10.5px; color: #8b949e; text-transform: uppercase; letter-spacing: 0.4px;
}

.stat-divider { width: 1px; height: 26px; background: #21262d; }

.sidebar-time {
  padding: 8px 20px 14px;
  border-top: 1px solid #21262d;
  display: flex; flex-direction: column; align-items: center; gap: 3px;
}

.time-date {
  font-size: 13px; font-weight: 500; color: #8b949e;
  letter-spacing: 0.5px; font-variant-numeric: tabular-nums;
}

.time-clock {
  font-size: 22px; font-weight: 700; color: #f0f6fc;
  line-height: 1.15; letter-spacing: 1px;
  font-variant-numeric: tabular-nums;
}

.time-city {
  font-size: 11px; color: #484f58; letter-spacing: 0.4px;
  text-transform: uppercase; margin-top: 2px;
}

/* ===== Main Content ===== */
.main-content {
  flex: 1; display: flex; flex-direction: column; min-width: 0; height: 100vh; overflow: hidden;
}

.top-bar {
  display: flex; align-items: center; justify-content: space-between;
  height: 56px; padding: 0 32px; border-bottom: 1px solid #21262d;
  background: #0d1117; flex-shrink: 0;
}

.page-title { font-size: 18px; font-weight: 600; color: #f0f6fc; letter-spacing: -0.2px; }

.content-area { flex: 1; padding: 24px 32px; overflow-y: auto; }

::-webkit-scrollbar { width: 5px; }
::-webkit-scrollbar-track { background: transparent; }
::-webkit-scrollbar-thumb { background: #30363d; border-radius: 3px; }
::-webkit-scrollbar-thumb:hover { background: #484f58; }
::selection { background: rgba(35,134,54,0.3); }
</style>
