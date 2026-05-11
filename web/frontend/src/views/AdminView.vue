<template>
  <div class="admin-view">
    <div v-if="!authenticated" class="auth-card">
      <div class="auth-header">
        <svg viewBox="0 0 24 24" class="lock-icon"><path d="M18 8h-1V6c0-2.76-2.24-5-5-5S7 3.24 7 6v2H6c-1.1 0-2 .9-2 2v10c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V10c0-1.1-.9-2-2-2zm-6 9c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2zm3.1-9H8.9V6c0-1.71 1.39-3.1 3.1-3.1s3.1 1.39 3.1 3.1v2z" fill="currentColor"/></svg>
        <h2 class="auth-title">管理员认证</h2>
        <p class="auth-desc">请输入管理员密码以访问后端管理</p>
      </div>
      <input type="password" class="auth-input" v-model="password" placeholder="管理员密码" @keyup.enter="doAuth" />
      <p class="auth-error" v-if="authError">{{ authError }}</p>
      <button class="btn btn-primary btn-block" @click="doAuth" :disabled="authing">
        {{ authing ? '验证中...' : '验证' }}
      </button>
    </div>

    <div v-else class="admin-content">
      <section class="card">
        <div class="card-header">
          <svg viewBox="0 0 24 24" class="card-icon active-icon"><path d="M13 2.03v2.02c4.39.54 7.5 4.53 6.96 8.92-.46 3.64-3.32 6.53-6.96 6.96v2.02c5.49-.66 9.5-5.6 8.84-11.09-.55-4.58-4.09-8.19-8.84-8.83zM4.26 17.86c1.94 2.86 5.06 4.7 8.74 4.76v-2.02c-2.41-.12-4.53-1.44-5.74-3.28l-3 .54zm.03-11.72l3 .54c1.21-1.84 3.33-3.16 5.74-3.28V1.38c-3.68.06-6.8 1.9-8.74 4.76z" fill="currentColor"/></svg>
          <div>
            <h2 class="card-title">运行中任务</h2>
            <p class="card-desc">{{ runningTasks.length }} 个任务执行中</p>
          </div>
          <button class="btn btn-ghost btn-sm" @click="pollRunning">刷新</button>
        </div>
        <div v-if="runningTasks.length === 0" class="empty-inline">暂无运行中的任务</div>
        <div v-else class="running-list">
          <div class="running-row" v-for="t in runningTasks" :key="t.id">
            <span class="r-dot" :class="t.status"></span>
            <span class="r-name">{{ t.customer_name || '未命名' }}</span>
            <span class="r-file">{{ t.file_name }}</span>
            <span class="r-badge" :class="t.status">{{ statusLabel(t.status) }}</span>
            <button class="btn btn-danger btn-xs" @click="cancelRunning(t.id)">终止</button>
          </div>
        </div>
      </section>

      <section class="card">
        <div class="card-header">
          <svg viewBox="0 0 24 24" class="card-icon"><path d="M20 6h-8l-2-2H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2zm-6 10H6v-2h8v2zm4-4H6v-2h12v2z" fill="currentColor"/></svg>
          <div>
            <h2 class="card-title">文件管理</h2>
            <p class="card-desc">web_data 目录文件（{{ files.length }} 项）</p>
          </div>
          <div class="card-actions">
            <button class="btn btn-ghost btn-sm" @click="pollFiles">刷新</button>
            <button class="btn btn-danger btn-sm" :disabled="selectedPaths.length === 0" @click="handleCleanup">
              清理选中 ({{ selectedPaths.length }})
            </button>
          </div>
        </div>
        <div v-if="files.length === 0" class="empty-inline">暂无文件</div>
        <div v-else class="file-table">
          <div class="file-row file-header">
            <label class="file-check">
              <input type="checkbox" :checked="selectAll" @change="toggleSelectAll" title="全选/取消" />
            </label>
            <span class="f-name h">文件/目录名</span>
            <span class="f-customer h">客户名称</span>
            <span class="f-size h">大小</span>
            <span class="f-date h">修改日期</span>
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
        console.log('[Auth] 发送密码:', this.password)
        const { data } = await adminAuth(this.password)
        console.log('[Auth] 收到的响应:', data)
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
        console.error('[Auth] 请求异常:', e)
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
  background: #161b22; border: 1px solid #21262d; border-radius: 12px; padding: 40px 32px;
  text-align: center;
}

.auth-header { margin-bottom: 24px; }
.lock-icon { width: 40px; height: 40px; color: #58a6ff; margin: 0 auto 14px; }
.auth-title { font-size: 18px; font-weight: 600; color: #f0f6fc; margin-bottom: 6px; }
.auth-desc { font-size: 13px; color: #8b949e; }
.auth-input {
  width: 100%; padding: 10px 14px; background: #0d1117; border: 1px solid #30363d;
  border-radius: 8px; color: #f0f6fc; font-size: 14px; outline: none; margin-bottom: 12px;
  transition: border-color 0.2s; text-align: center; letter-spacing: 4px;
}
.auth-input:focus { border-color: #58a6ff; box-shadow: 0 0 0 3px rgba(88,166,255,0.15); }
.auth-error { font-size: 13px; color: #f85149; margin-bottom: 12px; }

/* Cards */
.card { background: #161b22; border: 1px solid #21262d; border-radius: 12px; padding: 24px; margin-bottom: 20px; }
.card-header {
  display: flex; align-items: flex-start; gap: 12px; margin-bottom: 18px;
}
.card-icon { width: 19px; height: 19px; color: #3fb950; flex-shrink: 0; margin-top: 2px; }
.card-icon.active-icon { color: #58a6ff; }
.card-title { font-size: 15px; font-weight: 600; color: #f0f6fc; margin-bottom: 2px; }
.card-desc { font-size: 12.5px; color: #8b949e; }
.card-actions { display: flex; gap: 8px; margin-left: auto; flex-shrink: 0; }

.empty-inline { padding: 20px; text-align: center; font-size: 13px; color: #484f58; }

/* Buttons */
.btn { display: inline-flex; align-items: center; gap: 8px; padding: 9px 20px; border: none; border-radius: 8px; font-size: 14px; font-weight: 500; cursor: pointer; transition: all 0.15s; }
.btn:disabled { opacity: 0.4; cursor: not-allowed; }
.btn-primary { background: #238636; color: #fff; }
.btn-primary:hover:not(:disabled) { background: #2ea043; }
.btn-block { width: 100%; justify-content: center; padding: 11px; margin-top: 4px; }
.btn-ghost { background: transparent; color: #8b949e; padding: 6px 12px; font-size: 13px; }
.btn-ghost:hover { color: #c9d1d9; }
.btn-danger { background: #da3633; color: #fff; }
.btn-danger:hover:not(:disabled) { background: #c62828; }
.btn-sm { padding: 6px 16px; font-size: 13px; }
.btn-xs { padding: 4px 12px; font-size: 12px; flex-shrink: 0; }

/* Running Tasks */
.running-list { display: flex; flex-direction: column; gap: 4px; }
.running-row {
  display: flex; align-items: center; gap: 10px; padding: 10px 12px;
  background: #0d1117; border-radius: 6px; font-size: 13px;
}
.r-dot { width: 7px; height: 7px; border-radius: 50%; flex-shrink: 0; }
.r-dot.running { background: #3fb950; animation: pulse 1.5s infinite; }
.r-dot.pending { background: #484f58; }
@keyframes pulse { 0%,100% { opacity: 1; } 50% { opacity: 0.3; } }
.r-name { font-weight: 500; color: #c9d1d9; min-width: 80px; }
.r-file { color: #8b949e; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; flex: 1; }
.r-badge { font-size: 11px; padding: 2px 8px; border-radius: 20px; flex-shrink: 0; }
.r-badge.running { background: rgba(35,134,54,0.15); color: #3fb950; }
.r-badge.pending { background: #21262d; color: #8b949e; }

/* File List */
.file-table { border: 1px solid #21262d; border-radius: 8px; overflow: hidden; }
.file-header {
  background: #0d1117; border-bottom: 1px solid #21262d;
  position: sticky; top: 0; z-index: 1;
}
.file-header .file-check input { margin: 0; }
.file-scroll { max-height: 420px; overflow-y: auto; }
.file-row {
  display: flex; align-items: center; gap: 10px; padding: 9px 12px;
  font-size: 13px;
}
.file-row:not(.file-header) { background: #0d1117; }
.file-row:not(.file-header):hover { background: #161b22; }
.file-check { width: 28px; flex-shrink: 0; cursor: pointer; text-align: center; }
.file-check input { accent-color: #238636; cursor: pointer; }
.f-name { flex: 2; color: #c9d1d9; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.f-name.h { color: #8b949e; font-weight: 600; font-size: 11.5px; text-transform: uppercase; letter-spacing: 0.3px; }
.f-customer { flex: 1; color: #3fb950; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-weight: 500; }
.f-customer.h { color: #8b949e; font-weight: 600; font-size: 11.5px; text-transform: uppercase; letter-spacing: 0.3px; }
.f-size { width: 80px; color: #8b949e; text-align: right; flex-shrink: 0; }
.f-size.h { font-size: 11.5px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.3px; }
.f-date { width: 140px; color: #484f58; text-align: right; flex-shrink: 0; }
.f-date.h { font-size: 11.5px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.3px; }
</style>
