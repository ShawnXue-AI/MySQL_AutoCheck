<template>
  <div class="admin-view">
    <!-- Auth Screen -->
    <div v-if="!authenticated" class="auth-card">
      <div class="auth-header">
        <svg viewBox="0 0 48 48" class="lock-icon"><rect x="14" y="20" width="20" height="20" rx="3" fill="none" stroke="currentColor" stroke-width="2.5"/><path d="M18 20V14a6 6 0 1112 0v6" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/><circle cx="24" cy="30" r="2" fill="currentColor"/></svg>
        <h2 class="auth-title">Admin Auth</h2>
        <p class="auth-desc">// 请输入管理员密码以继续</p>
      </div>
      <div class="input-prefix">
        <span class="input-prefix-icon">$</span>
        <input type="password" class="field-input has-prefix" v-model="password" placeholder="password" @keyup.enter="doAuth" />
      </div>
      <p class="auth-error" v-if="authError">{{ authError }}</p>
      <button class="btn btn-primary btn-block" @click="doAuth" :disabled="authing">
        {{ authing ? 'Authenticating...' : 'Authenticate' }}
      </button>
    </div>

    <!-- Admin Content -->
    <div v-else class="admin-content">
      <!-- Running Tasks -->
      <section class="card">
        <div class="card-header">
          <svg viewBox="0 0 24 24" class="card-icon" style="color: var(--accent-cyan);"><path d="M13 2.03v2.02c4.39.54 7.5 4.53 6.96 8.92-.46 3.64-3.32 6.53-6.96 6.96v2.02c5.49-.66 9.5-5.6 8.84-11.09-.55-4.58-4.09-8.19-8.84-8.83zM4.26 17.86c1.94 2.86 5.06 4.7 8.74 4.76v-2.02c-2.41-.12-4.53-1.44-5.74-3.28l-3 .54zm.03-11.72l3 .54c1.21-1.84 3.33-3.16 5.74-3.28V1.38c-3.68.06-6.8 1.9-8.74 4.76z" fill="currentColor"/></svg>
          <div>
            <h2 class="card-title">运行中任务</h2>
            <p class="card-desc">$ ps aux | grep analysis  →  {{ runningTasks.length }} 进程</p>
          </div>
          <button class="btn btn-ghost btn-sm" @click="pollRunning">刷新</button>
        </div>
        <div v-if="runningTasks.length === 0" class="empty-inline">
          <span class="terminal-prompt">[INFO]</span> 暂无运行中的任务
        </div>
        <div v-else class="running-list">
          <div class="running-row" v-for="t in runningTasks" :key="t.id">
            <span class="status-dot" :class="t.status === 'running' ? 'active' : 'idle'"></span>
            <span class="r-name">{{ t.customer_name || '未命名' }}</span>
            <code class="r-file">{{ t.file_name }}</code>
            <span class="badge" :class="t.status === 'running' ? 'badge-info' : 'badge-muted'">{{ statusLabel(t.status) }}</span>
            <button class="btn btn-xs" style="border: 1px solid var(--accent-red); color: var(--accent-red); background: transparent;" @click="cancelRunning(t.id)">kill</button>
          </div>
        </div>
      </section>

      <!-- File Management -->
      <section class="card">
        <div class="card-header">
          <svg viewBox="0 0 24 24" class="card-icon"><path d="M20 6h-8l-2-2H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2zm-6 10H6v-2h8v2zm4-4H6v-2h12v2z" fill="currentColor"/></svg>
          <div>
            <h2 class="card-title">文件管理</h2>
            <p class="card-desc">$ ls web_data/  →  {{ files.length }} 项</p>
          </div>
          <div class="card-actions">
            <button class="btn btn-ghost btn-sm" @click="pollFiles">刷新</button>
            <button class="btn btn-danger btn-sm" :disabled="selectedPaths.length === 0" @click="handleCleanup">
              清理 ({{ selectedPaths.length }})
            </button>
          </div>
        </div>
        <div v-if="files.length === 0" class="empty-inline">
          <span class="terminal-prompt">[INFO]</span> 暂无文件
        </div>
        <div v-else class="file-table">
          <div class="file-row file-header">
            <label class="file-check">
              <input type="checkbox" :checked="selectAll" @change="toggleSelectAll" />
            </label>
            <span class="f-name h">name</span>
            <span class="f-customer h">customer</span>
            <span class="f-size h">size</span>
            <span class="f-date h">modified</span>
          </div>
          <div class="file-scroll">
            <div class="file-row" v-for="(f, i) in files" :key="i">
              <label class="file-check">
                <input type="checkbox" :value="f.path" v-model="selectedPaths" />
              </label>
              <span class="f-name" :title="f.name">{{ f.name }}</span>
              <span class="f-customer">{{ f.customer }}</span>
              <span class="f-size">{{ formatSize(f.size) }}</span>
              <span class="f-date">{{ f.modified }}</span>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import { adminAuth, getRunningTasks, cancelTask, getAdminFiles, adminCleanup } from '../api/index.js'
import { toastSuccess, toastError } from '../components/Toast.vue'

export default {
  name: 'AdminView',
  data() {
    return {
      authenticated: false,
      password: '',
      authing: false,
      authError: '',
      runningTasks: [],
      files: [],
      selectedPaths: [],
      pollTimer: null,
      fileTimer: null
    }
  },
  computed: {
    selectAll: {
      get() { return this.files.length > 0 && this.selectedPaths.length === this.files.length },
      set(v) { this.selectedPaths = v ? this.files.map(f => f.path) : [] }
    }
  },
  mounted() {
    this.pollRunning()
    this.pollTimer = setInterval(() => this.pollRunning(), 5000)
  },
  beforeUnmount() {
    if (this.pollTimer) clearInterval(this.pollTimer)
    if (this.fileTimer) clearInterval(this.fileTimer)
  },
  methods: {
    async doAuth() {
      this.authing = true; this.authError = ''
      try {
        const { data } = await adminAuth(this.password)
        if (data.success) {
          this.authenticated = true
          this.pollRunning()
          this.pollFiles()
          this.pollTimer = setInterval(() => this.pollRunning(), 5000)
          this.fileTimer = setInterval(() => this.pollFiles(), 15000)
        } else {
          this.authError = data.error || '密码错误'
        }
      } catch (e) {
        this.authError = '请求失败，请检查后端服务'
      } finally { this.authing = false }
    },
    async pollRunning() {
      try { const { data } = await getRunningTasks(); this.runningTasks = data.tasks || [] } catch {}
    },
    async pollFiles() {
      try { const { data } = await getAdminFiles(); this.files = data.files || [] } catch {}
    },
    async cancelRunning(id) {
      try { await cancelTask(id); toastSuccess('已终止'); this.pollRunning() } catch (e) { toastError('终止失败') }
    },
    toggleSelectAll() {
      if (this.selectedPaths.length === this.files.length) {
        this.selectedPaths = []
      } else {
        this.selectedPaths = this.files.map(f => f.path)
      }
    },
    async handleCleanup() {
      if (!this.selectedPaths.length) return
      try {
        const { data } = await adminCleanup(this.selectedPaths)
        if (data.deleted && data.deleted.length) toastSuccess(`已清理 ${data.deleted.length} 项`)
        if (data.errors && data.errors.length) toastError(data.errors.join('\n'))
        this.selectedPaths = []
        this.pollFiles()
      } catch (e) { toastError('清理失败') }
    },
    statusLabel(s) { return { pending: '等待中', running: '运行中', completed: '已完成', failed: '失败', cancelled: '已取消' }[s] || s },
    formatSize(b) { if (b < 1024) return b + ' B'; if (b < 1048576) return (b/1024).toFixed(1) + ' KB'; return (b/1048576).toFixed(1) + ' MB' }
  }
}
</script>

<style scoped>
.admin-view { max-width: 1100px; margin: 0 auto; }

/* Auth */
.auth-card {
  max-width: 420px; margin: 60px auto 0;
  background: var(--bg-secondary); border: 1px solid var(--border-default);
  border-radius: var(--radius-lg); padding: 40px 32px;
  text-align: center;
}
.auth-header { margin-bottom: 24px; }
.lock-icon { width: 44px; height: 44px; color: var(--accent-blue); margin: 0 auto 16px; }
.auth-title { font-family: var(--font-display); font-size: 20px; font-weight: 600; color: var(--text-primary); margin-bottom: 6px; }
.auth-desc { font-family: var(--font-code); font-size: 12px; color: var(--text-tertiary); }
.auth-error { font-size: 13px; color: var(--accent-red); margin: 10px 0; }

.input-prefix { position: relative; margin-bottom: 12px; }
.input-prefix-icon {
  position: absolute; left: 14px; top: 50%; transform: translateY(-50%);
  font-family: var(--font-code); font-size: 14px; color: var(--text-muted); font-weight: 500;
  pointer-events: none;
}
.field-input.has-prefix { padding-left: 34px; text-align: left; letter-spacing: 2px; }

/* Card overrides */
.card-actions { display: flex; gap: 8px; margin-left: auto; flex-shrink: 0; }

.empty-inline {
  padding: 20px; text-align: center; font-size: 13px; color: var(--text-muted);
  font-family: var(--font-code);
}
.terminal-prompt {
  color: var(--accent-green); font-weight: 500;
}

/* Running Tasks */
.running-list { display: flex; flex-direction: column; gap: 4px; }
.running-row {
  display: flex; align-items: center; gap: 10px; padding: 10px 14px;
  background: var(--bg-primary); border-radius: var(--radius-sm); font-size: 13px;
  border: 1px solid var(--border-default);
}
.r-name { font-weight: 500; color: var(--text-secondary); min-width: 80px; }
.r-file { font-size: 11px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; flex: 1; }

/* File Table */
.file-table {
  border: 1px solid var(--border-default); border-radius: var(--radius-md);
  overflow: hidden;
}
.file-header {
  background: var(--bg-primary); border-bottom: 1px solid var(--border-default);
  position: sticky; top: 0; z-index: 1;
}
.file-header .file-check input { margin: 0; }
.file-scroll { max-height: 420px; overflow-y: auto; }
.file-row {
  display: flex; align-items: center; gap: 10px; padding: 9px 14px;
  font-size: 13px;
}
.file-row:not(.file-header) { background: var(--bg-primary); }
.file-row:not(.file-header):hover { background: var(--bg-tertiary); }
.file-check { width: 28px; flex-shrink: 0; cursor: pointer; text-align: center; }
.file-check input { accent-color: var(--accent-blue); cursor: pointer; }
.f-name { flex: 2; color: var(--text-secondary); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-family: var(--font-code); font-size: 12px; }
.f-name.h { color: var(--text-tertiary); font-weight: 600; font-size: 10.5px; text-transform: uppercase; letter-spacing: 0.4px; }
.f-customer { flex: 1; color: var(--accent-blue); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-weight: 500; font-family: var(--font-code); font-size: 11.5px; }
.f-customer.h { color: var(--text-tertiary); font-weight: 600; font-size: 10.5px; text-transform: uppercase; letter-spacing: 0.4px; }
.f-size { width: 80px; color: var(--text-tertiary); text-align: right; flex-shrink: 0; font-family: var(--font-code); font-size: 11.5px; }
.f-size.h { font-size: 10.5px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.4px; }
.f-date { width: 150px; color: var(--text-muted); text-align: right; flex-shrink: 0; font-family: var(--font-code); font-size: 11px; }
.f-date.h { font-size: 10.5px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.4px; }
</style>
