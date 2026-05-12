<template>
  <div class="app-root">
    <div class="grid-dots"></div>
    <aside class="sidebar">
      <div class="sidebar-inner">
        <router-link to="/" class="logo-section">
          <span class="logo-prompt">~/</span>
          <span class="logo-text">inspection</span>
          <span class="logo-cursor"></span>
        </router-link>
        <nav class="nav-menu">
          <router-link to="/" class="nav-item" :class="{ active: $route.path === '/' }">
            <svg viewBox="0 0 24 24" class="nav-icon" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 13h8V3H3v10zm0 8h8v-6H3v6zm10 0h8V11h-8v10zm0-18v6h8V3h-8z"/></svg>
            <span>报告生成</span>
          </router-link>
          <router-link to="/tasks" class="nav-item" :class="{ active: $route.path === '/tasks' }">
            <svg viewBox="0 0 24 24" class="nav-icon" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M8 6h13M8 12h13M8 18h13M3 6h.01M3 12h.01M3 18h.01"/></svg>
            <span>任务记录</span>
          </router-link>
          <router-link to="/admin" class="nav-item" :class="{ active: $route.path === '/admin' }">
            <svg viewBox="0 0 24 24" class="nav-icon" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>
            <span>后端管理</span>
          </router-link>
          <router-link to="/benchmarks" class="nav-item" :class="{ active: $route.path === '/benchmarks' }">
            <svg viewBox="0 0 24 24" class="nav-icon" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 12h-4l-3 9L9 3l-3 9H2"/></svg>
            <span>巡检基准值</span>
          </router-link>
          <router-link to="/persondays" class="nav-item" :class="{ active: $route.path === '/persondays' }">
            <svg viewBox="0 0 24 24" class="nav-icon" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="18" x="3" y="3" rx="2"/><path d="M3 9h18M9 21V9"/></svg>
            <span>人天记录</span>
          </router-link>
        </nav>

        <div class="sidebar-footer">
          <div class="stats-row">
            <div class="stat-item">
              <span class="stat-value accent">{{ stats.active_users }}</span>
              <span class="stat-label">在线</span>
            </div>
            <span class="stat-sep"></span>
            <div class="stat-item">
              <span class="stat-value">{{ stats.total_completed }}</span>
              <span class="stat-label">累计</span>
            </div>
            <template v-if="stats.running_tasks > 0">
              <span class="stat-sep"></span>
              <div class="stat-item">
                <span class="stat-value pulse">{{ stats.running_tasks }}</span>
                <span class="stat-label">运行中</span>
              </div>
            </template>
          </div>
          <div class="time-row">
            <span class="time-date">{{ currentDate }}</span>
            <span class="time-clock">{{ currentTime }}</span>
            <span class="time-city">{{ currentCity }}</span>
          </div>
        </div>
      </div>
    </aside>

    <main class="main-area">
      <header class="top-bar">
        <h1 class="page-title">{{ pageTitle }}</h1>
      </header>
      <div class="content-scroll">
        <router-view @stats-update="updateStats" />
      </div>
    </main>
  </div>
</template>

<script>
import { getStats } from './api/index.js'

export default {
  name: 'App',
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
  computed: {
    pageTitle() {
      const m = {
        '/': '报告生成', '/tasks': '任务记录', '/admin': '后端管理',
        '/benchmarks': '巡检基准值', '/persondays': '人天记录'
      }
      return m[this.$route.path] || '任务记录'
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
    clearInterval(this.pollTimer)
    clearInterval(this.clockTimer)
  },
  methods: {
    async pollStats() {
      try { const { data } = await getStats(); this.stats = data } catch {}
    },
    updateClock() {
      const now = new Date()
      this.currentDate = now.toLocaleDateString('zh-CN', { year:'numeric', month:'2-digit', day:'2-digit' })
      this.currentTime = now.toLocaleTimeString('zh-CN', { hour:'2-digit', minute:'2-digit', second:'2-digit', hour12:false })
    },
    async detectCity() {
      try { const r = await fetch('https://ipapi.co/json/'); const d = await r.json(); if (d.city) this.currentCity = d.city } catch {}
    },
    updateStats(s) { if (s) this.stats = s }
  }
}
</script>

<style>
:root {
  --background: #0a0a0b;
  --foreground: #fafafa;
  --card: #111113;
  --card-foreground: #fafafa;
  --border: #27272a;
  --border-hover: #3f3f46;
  --primary: #3b82f6;
  --primary-foreground: #eff6ff;
  --muted: #18181b;
  --muted-foreground: #a1a1aa;
  --destructive: #ef4444;
  --radius: 0.625rem;
}

* { margin: 0; padding: 0; box-sizing: border-box; }

body {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'PingFang SC', 'Microsoft YaHei', sans-serif;
  background: var(--background);
  color: var(--foreground);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

.app-root {
  display: flex;
  height: 100vh;
  overflow: hidden;
  position: relative;
}

.grid-dots {
  position: fixed; inset: 0; pointer-events: none; z-index: 0;
  opacity: 0.025;
  background-image: linear-gradient(to right, currentColor 1px, transparent 1px),
                    linear-gradient(to bottom, currentColor 1px, transparent 1px);
  background-size: 32px 32px;
}

/* ===== Sidebar ===== */
.sidebar {
  width: 240px; flex-shrink: 0; height: 100vh; position: relative; z-index: 1;
}
.sidebar-inner {
  display: flex; flex-direction: column; height: 100%;
  border-right: 1px solid var(--border);
  background: rgba(17,17,19,0.55);
  padding: 12px;
}

.logo-section {
  display: flex; align-items: center; gap: 1px; padding: 8px 12px 16px;
  text-decoration: none; cursor: pointer;
}
.logo-prompt { font-size: 15px; font-weight: 600; color: var(--primary); font-family: 'JetBrains Mono','Consolas',monospace; }
.logo-text { font-size: 15px; font-weight: 700; color: var(--foreground); font-family: 'JetBrains Mono','Consolas',monospace; letter-spacing: -0.5px; }
.logo-cursor { width: 2px; height: 16px; background: var(--primary); display: inline-block; animation: blink 1s step-end infinite; margin-left: 1px; opacity: 0.7; }
@keyframes blink { 0%,100% { opacity: 1; } 50% { opacity: 0; } }

.nav-menu { display: flex; flex-direction: column; gap: 2px; flex: 1; }

.nav-item {
  display: flex; align-items: center; gap: 10px; padding: 8px 12px;
  border-radius: var(--radius); font-size: 13.5px; font-weight: 500;
  color: var(--muted-foreground); text-decoration: none;
  border: 1px solid transparent;
  transition: all 0.2s ease;
}
.nav-item:hover { color: var(--foreground); background: var(--muted); border-color: var(--border); }
.nav-item.active { color: var(--primary-foreground); background: rgba(59,130,246,0.15); border-color: rgba(59,130,246,0.3); }
.nav-item:active { transform: scale(0.97); }

.nav-icon { width: 18px; height: 18px; flex-shrink: 0; }
.nav-item.active .nav-icon { color: var(--primary); }

/* Sidebar bottom */
.sidebar-footer { margin-top: auto; flex-shrink: 0; }
.stats-row { display: flex; align-items: center; justify-content: center; gap: 14px; padding: 10px 0; border-top: 1px solid var(--border); }
.stat-item { display: flex; flex-direction: column; align-items: center; gap: 1px; }
.stat-value { font-size: 19px; font-weight: 700; color: var(--foreground); line-height: 1.1; }
.stat-value.accent { color: #22c55e; }
.stat-value.pulse { color: #22c55e; animation: pulse 2s ease-in-out infinite; }
@keyframes pulse { 0%,100% { opacity: 1; } 50% { opacity: 0.4; } }
.stat-label { font-size: 10px; color: var(--muted-foreground); text-transform: uppercase; letter-spacing: 0.5px; }
.stat-sep { width: 1px; height: 22px; background: var(--border); }

.time-row { display: flex; flex-direction: column; align-items: center; gap: 2px; padding: 8px 0 12px; border-top: 1px solid var(--border); }
.time-date { font-size: 12px; color: var(--muted-foreground); font-variant-numeric: tabular-nums; }
.time-clock { font-size: 20px; font-weight: 700; color: var(--foreground); line-height: 1.1; font-variant-numeric: tabular-nums; letter-spacing: 0.5px; }
.time-city { font-size: 10px; color: #52525b; text-transform: uppercase; letter-spacing: 0.4px; margin-top: 1px; }

/* ===== Main ===== */
.main-area { flex: 1; display: flex; flex-direction: column; min-width: 0; height: 100vh; position: relative; z-index: 1; }

.top-bar {
  display: flex; align-items: center; height: 52px; padding: 0 28px;
  border-bottom: 1px solid var(--border);
  background: rgba(10,10,11,0.55);
  flex-shrink: 0;
}
.page-title { font-size: 17px; font-weight: 600; color: var(--foreground); letter-spacing: -0.3px; }

.content-scroll { flex: 1; overflow-y: auto; padding: 24px 28px; }

/* ===== Scrollbar ===== */
::-webkit-scrollbar { width: 5px; }
::-webkit-scrollbar-track { background: transparent; }
::-webkit-scrollbar-thumb { background: var(--border); border-radius: 3px; }
::-webkit-scrollbar-thumb:hover { background: var(--border-hover); }
::selection { background: rgba(59,130,246,0.3); }

/* ===== Shared Components ===== */
.card {
  background: rgba(17,17,19,0.5); border: 1px solid var(--border); border-radius: var(--radius);
  padding: 20px; margin-bottom: 16px;
}
.card-header { display: flex; align-items: flex-start; gap: 10px; margin-bottom: 16px; }
.card-icon { width: 17px; height: 17px; color: var(--primary); flex-shrink: 0; margin-top: 2px; }
.card-title { font-size: 14px; font-weight: 600; color: var(--foreground); margin-bottom: 2px; }
.card-desc { font-size: 12px; color: var(--muted-foreground); line-height: 1.4; }
.card-actions { display: flex; gap: 8px; margin-left: auto; flex-shrink: 0; }

.empty-inline { padding: 28px; text-align: center; font-size: 13px; color: var(--muted-foreground); }
.empty-state {
  display: flex; flex-direction: column; align-items: center; gap: 10px; padding: 60px 20px; text-align: center;
}
.empty-icon { width: 42px; height: 42px; color: var(--muted-foreground); }
.empty-title { font-size: 15px; font-weight: 600; color: var(--foreground); }
.empty-desc { font-size: 13px; color: var(--muted-foreground); }

/* ===== Buttons ===== */
.btn {
  display: inline-flex; align-items: center; justify-content: center; gap: 7px;
  padding: 8px 18px; border: 1px solid transparent; border-radius: var(--radius);
  font-size: 13.5px; font-weight: 500; cursor: pointer; text-decoration: none;
  transition: all 0.2s ease; white-space: nowrap;
}
.btn:disabled { opacity: 0.4; cursor: not-allowed; pointer-events: none; }
.btn:active:not(:disabled) { transform: scale(0.96); }

.btn-primary { background: var(--primary); color: #fff; border-color: var(--primary); }
.btn-primary:hover:not(:disabled) { background: #2563eb; }

.btn-outline { background: transparent; color: var(--foreground); border-color: var(--border); }
.btn-outline:hover:not(:disabled) { border-color: var(--border-hover); background: var(--muted); }

.btn-ghost { background: transparent; color: var(--muted-foreground); border-color: transparent; }
.btn-ghost:hover:not(:disabled) { color: var(--foreground); background: var(--muted); }

.btn-danger { background: var(--destructive); color: #fff; border-color: var(--destructive); }
.btn-danger:hover:not(:disabled) { background: #dc2626; }

.btn-danger-outline { background: transparent; color: var(--destructive); border-color: rgba(239,68,68,0.3); }
.btn-danger-outline:hover:not(:disabled) { background: rgba(239,68,68,0.1); border-color: var(--destructive); }

.btn-sm { padding: 6px 14px; font-size: 12.5px; }
.btn-xs { padding: 4px 10px; font-size: 11.5px; }
.btn-block { width: 100%; }
.btn-icon { width: 15px; height: 15px; flex-shrink: 0; }

/* ===== Form inputs ===== */
.field-group { margin-bottom: 14px; }
.field-label { display: block; font-size: 12.5px; color: var(--muted-foreground); font-weight: 500; margin-bottom: 5px; }
.required-star { color: var(--destructive); }
.field-input {
  width: 100%; padding: 9px 12px; background: transparent; border: 1px solid var(--border);
  border-radius: var(--radius); color: var(--foreground); font-size: 13.5px; outline: none;
  transition: border-color 0.2s;
}
.field-input:focus { border-color: var(--primary); box-shadow: 0 0 0 3px rgba(59,130,246,0.15); }
.field-input::placeholder { color: #52525b; }
.field-input.error { border-color: var(--destructive); }
.field-error { font-size: 11.5px; color: var(--destructive); margin-top: 4px; display: block; }

textarea.field-input { resize: vertical; min-height: 50px; font-family: inherit; }

/* ===== Badges ===== */
.badge {
  display: inline-flex; align-items: center; padding: 2px 8px; border-radius: 20px;
  font-size: 10.5px; font-weight: 500; border: 1px solid var(--border);
}

/* ===== Auth card ===== */
.auth-card {
  max-width: 380px; margin: 60px auto 0;
  background: rgba(17,17,19,0.5); border: 1px solid var(--border); border-radius: var(--radius);
  padding: 36px 28px; text-align: center;
}
.auth-header { margin-bottom: 22px; }
.lock-icon { width: 38px; height: 38px; color: var(--primary); margin: 0 auto 12px; }
.auth-title { font-size: 17px; font-weight: 600; color: var(--foreground); margin-bottom: 5px; }
.auth-desc { font-size: 13px; color: var(--muted-foreground); }
.auth-input {
  width: 100%; padding: 10px 14px; background: transparent; border: 1px solid var(--border);
  border-radius: var(--radius); color: var(--foreground); font-size: 14px; outline: none;
  margin-bottom: 12px; transition: border-color 0.2s; text-align: center; letter-spacing: 4px;
}
.auth-input:focus { border-color: var(--primary); box-shadow: 0 0 0 3px rgba(59,130,246,0.15); }
.auth-error { font-size: 12.5px; color: var(--destructive); margin-bottom: 12px; }

/* ===== Table ===== */
.data-table { border: 1px solid var(--border); border-radius: var(--radius); overflow: hidden; }
.table-header { background: var(--muted); border-bottom: 1px solid var(--border); position: sticky; top: 0; z-index: 1; }
.table-scroll { max-height: 500px; overflow-y: auto; }
.table-row { display: flex; align-items: center; gap: 8px; padding: 9px 12px; font-size: 13px; }
.table-row:not(.table-header):hover { background: var(--muted); }

/* ===== Drop zone ===== */
.drop-zone {
  border: 2px dashed var(--border); border-radius: var(--radius); padding: 28px 20px;
  text-align: center; transition: all 0.2s; cursor: pointer; margin-bottom: 14px;
  background: rgba(17,17,19,0.2);
}
.drop-zone:hover { border-color: var(--primary); background: rgba(59,130,246,0.04); }
.drop-zone.dragging { border-color: var(--primary); background: rgba(59,130,246,0.08); }
.drop-zone.filled { border-style: solid; border-color: #22c55e; background: rgba(34,197,94,0.04); }
.drop-text { font-size: 13px; color: var(--foreground); margin-bottom: 4px; }
.drop-hint { font-size: 12px; color: var(--muted-foreground); margin-bottom: 14px; }

.upload-icon { width: 36px; height: 36px; color: var(--muted-foreground); margin: 0 auto 8px; }
.upload-options { display: flex; gap: 8px; justify-content: center; }
.upload-options label { position: relative; }
.upload-options input[type="file"] {
  position: absolute; inset: 0; opacity: 0; cursor: pointer;
}

/* ===== Modal ===== */
.modal-overlay {
  position: fixed; inset: 0; z-index: 100; display: flex; align-items: center; justify-content: center;
  background: rgba(0,0,0,0.7); backdrop-filter: blur(6px);
}
.modal-panel {
  background: rgba(17,17,19,0.6); border: 1px solid var(--border); border-radius: var(--radius);
  padding: 32px; text-align: center; max-width: 420px; width: 90%;
}
.modal-icon { width: 44px; height: 44px; margin: 0 auto 16px; }
.modal-title { font-size: 17px; font-weight: 600; color: var(--foreground); margin-bottom: 6px; }
.modal-sub { font-size: 13px; color: var(--muted-foreground); margin-bottom: 18px; }
.modal-track { height: 4px; background: var(--muted); border-radius: 2px; overflow: hidden; margin-bottom: 8px; }
.modal-fill { height: 100%; background: var(--primary); border-radius: 2px; transition: width 0.3s ease; }
.modal-pct { font-size: 13px; color: var(--muted-foreground); margin-bottom: 16px; }
.modal-err { font-size: 12.5px; color: var(--destructive); margin-bottom: 14px; }
.modal-actions { display: flex; gap: 8px; justify-content: center; }

.modal-enter-active, .modal-leave-active { transition: opacity 0.2s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
</style>
