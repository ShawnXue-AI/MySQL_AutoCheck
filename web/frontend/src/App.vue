<template>
  <div class="app-container">
    <aside class="sidebar">
      <div class="logo-section">
        <div class="logo-icon">
          <svg viewBox="0 0 40 40" fill="none"><rect width="40" height="40" rx="10" fill="var(--accent-blue)"/><path d="M12 20l6 6 10-10" stroke="var(--bg-primary)" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/></svg>
        </div>
        <div class="logo-text">
          <span class="logo-title">Inspection</span>
          <span class="logo-sub">// MySQL巡检系统</span>
        </div>
      </div>
      <nav class="nav-menu">
        <router-link to="/" class="nav-item" :class="{ active: $route.path === '/' }">
          <svg viewBox="0 0 24 24" class="nav-icon"><path d="M3 13h8V3H3v10zm0 8h8v-6H3v6zm10 0h8V11h-8v10zm0-18v6h8V3h-8z" fill="currentColor"/></svg>
          <span>报告生成</span>
          <kbd class="nav-shortcut">H</kbd>
        </router-link>
        <router-link to="/tasks" class="nav-item" :class="{ active: $route.path === '/tasks' }">
          <svg viewBox="0 0 24 24" class="nav-icon"><path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-7 14l-5-5 1.41-1.41L12 14.17l7.59-7.59L21 8l-9 9z" fill="currentColor"/></svg>
          <span>任务记录</span>
          <kbd class="nav-shortcut">T</kbd>
        </router-link>
        <router-link to="/admin" class="nav-item" :class="{ active: $route.path === '/admin' }">
          <svg viewBox="0 0 24 24" class="nav-icon"><path d="M19.43 12.98c.04-.32.07-.64.07-.98s-.03-.66-.07-.98l2.11-1.65c.19-.15.24-.42.12-.64l-2-3.46c-.12-.22-.39-.3-.61-.22l-2.49 1c-.52-.4-1.08-.73-1.69-.98l-.38-2.65C14.46 2.18 14.25 2 14 2h-4c-.25 0-.46.18-.49.42l-.38 2.65c-.61.25-1.17.59-1.69.98l-2.49-1c-.23-.09-.49 0-.61.22l-2 3.46c-.13.22-.07.49.12.64l2.11 1.65c-.04.32-.07.65-.07.98s.03.66.07.98l-2.11 1.65c-.19.15-.24.42-.12.64l2 3.46c.12.22.39.3.61.22l2.49-1c.52.4 1.08.73 1.69.98l.38 2.65c.03.24.24.42.49.42h4c.25 0 .46-.18.49-.42l.38-2.65c.61-.25 1.17-.59 1.69-.98l2.49 1c.23.09.49 0 .61-.22l2-3.46c.12-.22.07-.49-.12-.64l-2.11-1.65zM12 15.5c-1.93 0-3.5-1.57-3.5-3.5s1.57-3.5 3.5-3.5 3.5 1.57 3.5 3.5-1.57 3.5-3.5 3.5z" fill="currentColor"/></svg>
          <span>后端管理</span>
          <kbd class="nav-shortcut">A</kbd>
        </router-link>
        <router-link to="/benchmarks" class="nav-item" :class="{ active: $route.path === '/benchmarks' }">
          <svg viewBox="0 0 24 24" class="nav-icon"><path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zM9 17H7v-7h2v7zm4 0h-2V7h2v10zm4 0h-2v-4h2v4z" fill="currentColor"/></svg>
          <span>巡检基准值</span>
          <kbd class="nav-shortcut">B</kbd>
        </router-link>
        <router-link to="/persondays" class="nav-item" :class="{ active: $route.path === '/persondays' }">
          <svg viewBox="0 0 24 24" class="nav-icon"><path d="M19 4h-1V2h-2v2H8V2H6v2H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm0 16H5V10h14v10zm0-12H5V6h14v2z" fill="currentColor"/></svg>
          <span>人天记录</span>
          <kbd class="nav-shortcut">P</kbd>
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
          <span class="time-clock">{{ currentTime }}</span>
          <span class="time-date">{{ currentDate }}</span>
          <span class="time-city">{{ currentCity }}</span>
        </div>
      </div>
    </aside>
    <main class="main-content">
      <header class="top-bar">
        <h1 class="page-title">{{ $route.path === '/' ? '报告生成' : $route.path === '/admin' ? '后端管理' : $route.path === '/benchmarks' ? '巡检基准值' : $route.path === '/persondays' ? '人天记录' : '任务记录' }}</h1>
        <div class="top-bar-right">
          <span class="top-version">v1.0.0</span>
        </div>
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
      currentCity: 'Beijing',
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
/* ============================================================
   Design Tokens — Developer Tool / Terminal-Inspired Dark Theme
   ============================================================ */
:root {
  /* Backgrounds */
  --bg-primary: #0a0e14;
  --bg-secondary: #12161e;
  --bg-tertiary: #181d27;
  --bg-hover: #1c2130;
  --bg-input: #0d1117;

  /* Borders */
  --border-default: #1e2530;
  --border-muted: #1a2030;
  --border-hover: #2a3345;
  --border-active: #3b82f6;

  /* Text */
  --text-primary: #e6edf3;
  --text-secondary: #adbac7;
  --text-tertiary: #636c7b;
  --text-muted: #444c5a;

  /* Accent Colors */
  --accent-blue: #3b82f6;
  --accent-blue-hover: #5190f7;
  --accent-blue-subtle: rgba(59, 130, 246, 0.12);

  --accent-green: #3fb950;
  --accent-green-hover: #4cc560;
  --accent-green-subtle: rgba(63, 185, 80, 0.12);

  --accent-amber: #d2991d;
  --accent-amber-subtle: rgba(210, 153, 29, 0.12);

  --accent-red: #f85149;
  --accent-red-hover: #f96e67;
  --accent-red-subtle: rgba(248, 81, 73, 0.10);

  --accent-purple: #a371f7;
  --accent-purple-subtle: rgba(163, 113, 247, 0.12);

  --accent-cyan: #39c5cf;
  --accent-pink: #db61a2;

  /* Typography */
  --font-display: 'Space Grotesk', 'PingFang SC', 'Microsoft YaHei', sans-serif;
  --font-body: 'DM Sans', 'PingFang SC', 'Microsoft YaHei', sans-serif;
  --font-code: 'JetBrains Mono', 'Fira Code', 'Cascadia Code', 'Consolas', 'Monaco', monospace;

  /* Spacing scale */
  --space-xs: 4px;
  --space-sm: 8px;
  --space-md: 16px;
  --space-lg: 24px;
  --space-xl: 32px;
  --space-2xl: 48px;

  /* Radii */
  --radius-sm: 6px;
  --radius-md: 10px;
  --radius-lg: 14px;
  --radius-xl: 18px;

  /* Shadows */
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.4);
  --shadow-md: 0 4px 12px rgba(0, 0, 0, 0.5);
  --shadow-lg: 0 12px 40px rgba(0, 0, 0, 0.6);
  --shadow-glow-blue: 0 0 0 3px rgba(59, 130, 246, 0.15);
  --shadow-glow-green: 0 0 0 3px rgba(63, 185, 80, 0.12);

  /* Transitions */
  --transition-fast: 0.15s ease;
  --transition-normal: 0.2s ease;
}

/* ============================================================
   Reset & Base
   ============================================================ */
* { margin: 0; padding: 0; box-sizing: border-box; }

body {
  font-family: var(--font-body);
  background: var(--bg-primary);
  color: var(--text-secondary);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  line-height: 1.5;
}

/* Subtle noise texture overlay for depth */
body::before {
  content: '';
  position: fixed;
  inset: 0;
  z-index: -1;
  opacity: 0.015;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.85' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)'/%3E%3C/svg%3E");
  pointer-events: none;
}

.app-container {
  display: flex;
  height: 100vh;
  overflow: hidden;
}

/* ============================================================
   Sidebar
   ============================================================ */
.sidebar {
  width: 248px;
  background: var(--bg-secondary);
  border-right: 1px solid var(--border-default);
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
  padding: 0 18px;
  border-bottom: 1px solid var(--border-default);
  flex-shrink: 0;
}

.logo-icon {
  width: 30px; height: 30px; flex-shrink: 0;
}
.logo-icon svg { display: block; width: 100%; height: 100%; }

.logo-text {
  display: flex; flex-direction: column; gap: 1px;
}

.logo-title {
  font-family: var(--font-display);
  font-size: 15px; font-weight: 700; color: var(--text-primary);
  letter-spacing: -0.3px; line-height: 1.2;
}

.logo-sub {
  font-family: var(--font-code);
  font-size: 10px; color: var(--text-tertiary); line-height: 1.2;
}

.nav-menu {
  padding: 12px 10px; display: flex; flex-direction: column; gap: 2px; flex-shrink: 0;
}

.nav-item {
  display: flex; align-items: center; gap: 10px; padding: 9px 12px;
  border-radius: var(--radius-sm); color: var(--text-tertiary); text-decoration: none;
  font-size: 13.5px; font-weight: 500; transition: all var(--transition-fast);
  position: relative;
}

.nav-item:hover { background: var(--bg-tertiary); color: var(--text-secondary); }
.nav-item.active { background: var(--accent-blue-subtle); color: var(--accent-blue); }

.nav-icon { width: 18px; height: 18px; flex-shrink: 0; }

.nav-shortcut {
  margin-left: auto;
  font-family: var(--font-code);
  font-size: 10px;
  padding: 2px 5px;
  border-radius: 3px;
  background: var(--bg-primary);
  color: var(--text-muted);
  border: 1px solid var(--border-default);
  line-height: 1.2;
  opacity: 0;
  transition: opacity var(--transition-fast);
}
.nav-item:hover .nav-shortcut,
.nav-item.active .nav-shortcut { opacity: 1; }

/* Sidebar Bottom */
.sidebar-bottom {
  margin-top: auto;
  flex-shrink: 0;
}

.sidebar-footer {
  padding: 12px 18px 8px;
  border-top: 1px solid var(--border-default);
}

.stats-mini { display: flex; align-items: center; justify-content: center; gap: 14px; }

.stat-item { display: flex; flex-direction: column; align-items: center; gap: 1px; }

.stat-value {
  font-family: var(--font-code);
  font-size: 19px; font-weight: 600; color: var(--text-primary); line-height: 1.2;
  font-variant-numeric: tabular-nums;
}

.stat-value.accent { color: var(--accent-blue); }
.stat-value.pulse { color: var(--accent-blue); animation: stat-pulse 2s ease-in-out infinite; }

@keyframes stat-pulse {
  0%, 100% { opacity: 1; } 50% { opacity: 0.45; }
}

.stat-label {
  font-size: 10px; color: var(--text-tertiary); text-transform: uppercase;
  letter-spacing: 0.6px; font-weight: 500;
}

.stat-divider { width: 1px; height: 24px; background: var(--border-default); }

.sidebar-time {
  padding: 8px 18px 14px;
  border-top: 1px solid var(--border-default);
  display: flex; flex-direction: column; align-items: center; gap: 3px;
}

.time-clock {
  font-family: var(--font-code);
  font-size: 22px; font-weight: 600; color: var(--text-primary);
  line-height: 1.15; letter-spacing: 1px;
  font-variant-numeric: tabular-nums;
}

.time-date {
  font-family: var(--font-code);
  font-size: 11px; font-weight: 500; color: var(--text-tertiary);
  letter-spacing: 0.5px; font-variant-numeric: tabular-nums;
}

.time-city {
  font-family: var(--font-code);
  font-size: 10px; color: var(--text-muted); letter-spacing: 0.4px;
  text-transform: uppercase; margin-top: 2px;
}

/* ============================================================
   Main Content Area
   ============================================================ */
.main-content {
  flex: 1; display: flex; flex-direction: column; min-width: 0; height: 100vh; overflow: hidden;
}

.top-bar {
  display: flex; align-items: center; justify-content: space-between;
  height: 56px; padding: 0 32px; border-bottom: 1px solid var(--border-default);
  background: var(--bg-primary); flex-shrink: 0;
}

.page-title {
  font-family: var(--font-display);
  font-size: 17px; font-weight: 600; color: var(--text-primary);
  letter-spacing: -0.2px;
}

.top-bar-right {
  display: flex; align-items: center; gap: 12px;
}

.top-version {
  font-family: var(--font-code);
  font-size: 11px;
  color: var(--text-muted);
  background: var(--bg-tertiary);
  padding: 3px 8px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-default);
}

.content-area { flex: 1; padding: 24px 32px; overflow-y: auto; }

/* ============================================================
   Scrollbar
   ============================================================ */
::-webkit-scrollbar { width: 5px; height: 5px; }
::-webkit-scrollbar-track { background: transparent; }
::-webkit-scrollbar-thumb { background: var(--border-hover); border-radius: 3px; }
::-webkit-scrollbar-thumb:hover { background: var(--text-muted); }
::selection { background: var(--accent-blue-subtle); color: var(--text-primary); }

/* ============================================================
   Shared Component Primitives (used across all views/components)
   ============================================================ */

/* Card */
.card {
  background: var(--bg-secondary);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
  padding: 24px;
  margin-bottom: 20px;
  transition: border-color var(--transition-normal);
}
.card:hover { border-color: var(--border-hover); }

.card-header {
  display: flex; align-items: flex-start; gap: 12px; margin-bottom: 18px;
}

.card-icon { width: 20px; height: 20px; color: var(--accent-blue); flex-shrink: 0; margin-top: 2px; }
.card-title { font-family: var(--font-display); font-size: 15px; font-weight: 600; color: var(--text-primary); margin-bottom: 2px; letter-spacing: -0.2px; }
.card-desc { font-size: 12.5px; color: var(--text-tertiary); }

/* Buttons */
.btn {
  display: inline-flex; align-items: center; gap: 8px; padding: 9px 18px;
  border: none; border-radius: var(--radius-sm); font-size: 13.5px; font-weight: 500;
  cursor: pointer; transition: all var(--transition-fast); text-decoration: none;
  font-family: var(--font-body);
}
.btn:disabled { opacity: 0.4; cursor: not-allowed; pointer-events: none; }
.btn-icon { width: 16px; height: 16px; }

.btn-primary { background: var(--accent-blue); color: #fff; }
.btn-primary:hover:not(:disabled) { background: var(--accent-blue-hover); }

.btn-accent { background: linear-gradient(135deg, var(--accent-blue), #5190f7); color: #fff; }
.btn-accent:hover:not(:disabled) { background: linear-gradient(135deg, #3572e0, var(--accent-blue)); }

.btn-success { background: var(--accent-green); color: #fff; }
.btn-success:hover:not(:disabled) { background: var(--accent-green-hover); }

.btn-danger { background: var(--accent-red); color: #fff; }
.btn-danger:hover:not(:disabled) { background: var(--accent-red-hover); }

.btn-outline { background: transparent; border: 1px solid var(--border-hover); color: var(--text-secondary); }
.btn-outline:hover { border-color: var(--accent-blue); color: var(--accent-blue); background: var(--accent-blue-subtle); }
.btn-outline input[type="file"] { display: none; }

.btn-ghost { background: transparent; color: var(--text-tertiary); padding: 6px 12px; font-size: 13px; }
.btn-ghost:hover { color: var(--text-secondary); background: var(--bg-tertiary); }

.btn-block { width: 100%; justify-content: center; padding: 11px; font-size: 14px; }
.btn-sm { padding: 6px 14px; font-size: 12.5px; }
.btn-xs { padding: 4px 10px; font-size: 11.5px; }

/* Form Fields */
.field-group { margin-bottom: 14px; }
.field-label { display: block; font-size: 13px; font-weight: 500; color: var(--text-secondary); margin-bottom: 6px; }
.required { color: var(--accent-red); }

.field-input,
.field-select {
  width: 100%; padding: 10px 13px;
  background: var(--bg-input);
  border: 1px solid var(--border-hover);
  border-radius: var(--radius-sm);
  color: var(--text-primary);
  font-size: 13.5px; font-family: var(--font-body);
  outline: none; transition: border-color var(--transition-normal);
}
.field-input:focus,
.field-select:focus {
  border-color: var(--accent-blue);
  box-shadow: var(--shadow-glow-blue);
}
.field-input::placeholder { color: var(--text-muted); }

.field-select { cursor: pointer; appearance: auto; }
.field-select option { background: var(--bg-secondary); color: var(--text-secondary); }

.field-input.error { border-color: var(--accent-red); }
.field-input.error:focus { box-shadow: 0 0 0 3px rgba(248, 81, 73, 0.12); }

.field-error { display: block; font-size: 11.5px; color: var(--accent-red); margin-top: 4px; }

/* Empty State */
.empty-state { text-align: center; padding: 80px 20px; }
.empty-icon { width: 48px; height: 48px; margin: 0 auto 16px; color: var(--border-hover); }
.empty-title { font-family: var(--font-display); font-size: 16px; font-weight: 600; color: var(--text-secondary); margin-bottom: 8px; }
.empty-desc { font-size: 14px; color: var(--text-tertiary); margin-bottom: 24px; }

/* Badge / Status Pills */
.badge {
  font-size: 11px; font-weight: 500; padding: 2px 10px; border-radius: 20px;
  font-family: var(--font-code);
}
.badge-success { background: var(--accent-green-subtle); color: var(--accent-green); }
.badge-info { background: var(--accent-blue-subtle); color: var(--accent-blue); }
.badge-warning { background: var(--accent-amber-subtle); color: var(--accent-amber); }
.badge-error { background: var(--accent-red-subtle); color: var(--accent-red); }
.badge-muted { background: var(--bg-tertiary); color: var(--text-tertiary); }

/* Status Dot */
.status-dot {
  width: 7px; height: 7px; border-radius: 50%; flex-shrink: 0;
}
.status-dot.active { background: var(--accent-green); animation: dot-pulse 1.5s ease-in-out infinite; }
.status-dot.idle { background: var(--text-muted); }
.status-dot.error { background: var(--accent-red); }

@keyframes dot-pulse {
  0%, 100% { opacity: 1; } 50% { opacity: 0.3; }
}

/* Code inline */
code, .inline-code {
  font-family: var(--font-code);
  font-size: 12px;
  background: var(--accent-blue-subtle);
  color: var(--accent-blue);
  padding: 2px 6px;
  border-radius: 4px;
}

/* Tooltip / KBD */
kbd {
  font-family: var(--font-code);
  font-size: 10px;
  padding: 1px 5px;
  border-radius: 3px;
  background: var(--bg-primary);
  color: var(--text-muted);
  border: 1px solid var(--border-default);
  line-height: 1.4;
}

/* Responsive */
@media (max-width: 900px) {
  .sidebar { width: 64px; }
  .logo-text, .nav-item span, .nav-shortcut, .stat-label, .time-date, .time-city, .stats-mini { display: none; }
  .nav-item { justify-content: center; padding: 10px; }
  .logo-section { justify-content: center; padding: 0 12px; }
  .sidebar-footer, .sidebar-time { padding-left: 8px; padding-right: 8px; }
  .content-area { padding: 16px; }
  .top-bar { padding: 0 16px; }
}
</style>
