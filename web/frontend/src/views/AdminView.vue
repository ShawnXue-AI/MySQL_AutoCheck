<template>
  <div class="admin-view">
    <div v-if="!authenticated" class="auth-card">
      <div class="auth-header">
        <svg viewBox="0 0 24 24" class="lock-icon" v-once><path d="M18 8h-1V6c0-2.76-2.24-5-5-5S7 3.24 7 6v2H6c-1.1 0-2 .9-2 2v10c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V10c0-1.1-.9-2-2-2zm-6 9c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2zm3.1-9H8.9V6c0-1.71 1.39-3.1 3.1-3.1s3.1 1.39 3.1 3.1v2z" fill="currentColor"/></svg>
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
          <svg viewBox="0 0 24 24" class="card-icon active-icon" v-once><path d="M13 2.03v2.02c4.39.54 7.5 4.53 6.96 8.92-.46 3.64-3.32 6.53-6.96 6.96v2.02c5.49-.66 9.5-5.6 8.84-11.09-.55-4.58-4.09-8.19-8.84-8.83zM4.26 17.86c1.94 2.86 5.06 4.7 8.74 4.76v-2.02c-2.41-.12-4.53-1.44-5.74-3.28l-3 .54zm.03-11.72l3 .54c1.21-1.84 3.33-3.16 5.74-3.28V1.38c-3.68.06-6.8 1.9-8.74 4.76z" fill="currentColor"/></svg>
          <div>
            <h2 class="card-title">运行中任务</h2>
            <p class="card-desc">{{ runningTasks.length }} 个任务执行中</p>
          </div>
          <button class="btn btn-outline btn-sm" @click="pollRunning">
            <svg viewBox="0 0 24 24" class="btn-icon" fill="none" stroke="currentColor" stroke-width="2"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/></svg>
            刷新
          </button>
        </div>
        <div v-if="runningTasks.length === 0" class="empty-inline">暂无运行中的任务</div>
        <div v-else class="running-list">
          <div class="running-row" v-for="t in runningTasks" :key="t.id">
            <span class="r-dot" :class="t.status"></span>
            <span class="r-name">{{ t.customer_name || '未命名' }}</span>
            <span class="r-file">{{ t.file_name }}</span>
            <span class="r-badge" :class="t.status">{{ statusLabel(t.status) }}</span>
            <button class="btn btn-danger btn-xs" @click="cancelRunning(t.id)">
              <svg viewBox="0 0 24 24" class="btn-icon" style="width:13px;height:13px" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 6L6 18M6 6l12 12"/></svg>
              终止
            </button>
          </div>
        </div>
      </section>

      <section class="card">
        <div class="card-header">
          <svg viewBox="0 0 24 24" class="card-icon" v-once><path d="M20 6h-8l-2-2H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2zm-6 10H6v-2h8v2zm4-4H6v-2h12v2z" fill="currentColor"/></svg>
          <div>
            <h2 class="card-title">文件管理</h2>
            <p class="card-desc">web_data 目录文件（{{ files.length }} 项）</p>
          </div>
          <div class="card-actions">
            <button class="btn btn-outline btn-sm" @click="pollFiles">
              <svg viewBox="0 0 24 24" class="btn-icon" fill="none" stroke="currentColor" stroke-width="2"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/></svg>
              刷新
            </button>
            <button class="btn btn-danger btn-sm" :disabled="selectedPaths.length === 0" @click="handleCleanup">
              <svg viewBox="0 0 24 24" class="btn-icon" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 6h18M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/></svg>
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
import { ref, computed } from 'vue'
import { useAuth } from '../composables/useAuth.js'
import { formatSize, statusLabel } from '../composables/utils.js'
import { getRunningTasks, cancelTask, getAdminFiles, adminCleanup } from '../api/index.js'
import { toastSuccess, toastError } from '../components/Toast.vue'

export default {
  name: 'AdminView',
  setup() {
    const { authenticated, password, authing, authError, doAuth: _doAuth } = useAuth()
    const runningTasks = ref([])
    const files = ref([])
    const selectedPaths = ref([])
    let runningPoll = null
    let filePoll = null

    const selectAll = computed({
      get: () => files.value.length > 0 && selectedPaths.value.length === files.value.length,
      set: (v) => { selectedPaths.value = v ? files.value.map(f => f.path) : [] }
    })

    async function pollRunning() {
      try { const { data } = await getRunningTasks(); runningTasks.value = data.tasks || [] } catch {}
    }
    async function pollFiles() {
      try { const { data } = await getAdminFiles(); files.value = data.files || [] } catch {}
    }

    function doAuth() {
      _doAuth().then(success => {
        if (success) {
          pollRunning()
          pollFiles()
          runningPoll = setInterval(() => pollRunning(), 5000)
          filePoll = setInterval(() => pollFiles(), 15000)
        }
      })
    }

    async function cancelRunning(id) {
      try { await cancelTask(id); toastSuccess('已终止'); pollRunning() } catch { toastError('终止失败') }
    }

    function toggleSelectAll() {
      selectedPaths.value = selectedPaths.value.length === files.value.length
        ? [] : files.value.map(f => f.path)
    }

    async function handleCleanup() {
      if (!selectedPaths.value.length) return
      try {
        const { data } = await adminCleanup(selectedPaths.value)
        if (data.deleted?.length) toastSuccess(`已清理 ${data.deleted.length} 项`)
        if (data.errors?.length) toastError(data.errors.join('\n'))
        selectedPaths.value = []
        pollFiles()
      } catch { toastError('清理失败') }
    }

    pollRunning()
    runningPoll = setInterval(() => pollRunning(), 5000)

    return {
      authenticated, password, authing, authError, doAuth,
      runningTasks, files, selectedPaths, selectAll,
      pollRunning, pollFiles, cancelRunning, toggleSelectAll, handleCleanup,
      formatSize, statusLabel,
      _cleanup: () => { clearInterval(runningPoll); clearInterval(filePoll) }
    }
  },
  beforeUnmount() {
    this._cleanup()
  }
}
</script>

<style scoped>
.admin-view { max-width: 1100px; margin: 0 auto; }

.running-list { display: flex; flex-direction: column; gap: 4px; }
.running-row {
  display: flex; align-items: center; gap: 10px; padding: 10px 12px;
  background: rgba(24,24,27,0.3); border-radius: var(--radius); font-size: 13px;
  border-bottom: 1px solid rgba(39,39,42,0.35);
}
.running-row:last-child { border-bottom: none; }
.r-dot { width: 7px; height: 7px; border-radius: 50%; flex-shrink: 0; }
.r-dot.running { background: #22c55e; animation: pulse 1.5s infinite; }
.r-dot.pending { background: var(--border); }
@keyframes pulse { 0%,100% { opacity: 1; } 50% { opacity: 0.3; } }
.r-name { font-weight: 500; color: var(--foreground); min-width: 80px; }
.r-file { color: var(--muted-foreground); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; flex: 1; }
.r-badge { font-size: 11px; padding: 2px 8px; border-radius: 20px; flex-shrink: 0; }
.r-badge.running { background: rgba(34,197,94,0.12); color: #22c55e; }
.r-badge.pending { background: var(--muted); color: var(--muted-foreground); }

.file-table { border: 1px solid var(--border); border-radius: var(--radius); overflow: hidden; }
.file-header { background: var(--muted); border-bottom: 1px solid var(--border); position: sticky; top: 0; z-index: 1; }
.file-header .file-check input { margin: 0; }
.file-scroll { max-height: 420px; overflow-y: auto; }
.file-row { display: flex; align-items: center; gap: 10px; padding: 9px 12px; font-size: 13px; }
.file-row:not(.file-header) { background: rgba(17,17,19,0.3); border-bottom: 1px solid rgba(39,39,42,0.3); }
.file-row:not(.file-header):last-child { border-bottom: none; }
.file-row:not(.file-header):hover { background: rgba(59,130,246,0.05); }
.file-check { width: 28px; flex-shrink: 0; cursor: pointer; text-align: center; }
.file-check input { accent-color: var(--primary); cursor: pointer; }
.f-name { flex: 2; color: var(--foreground); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.f-name.h { color: var(--muted-foreground); font-weight: 600; font-size: 11.5px; text-transform: uppercase; letter-spacing: 0.3px; }
.f-customer { flex: 1; color: #22c55e; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-weight: 500; }
.f-customer.h { color: var(--muted-foreground); font-weight: 600; font-size: 11.5px; text-transform: uppercase; letter-spacing: 0.3px; }
.f-size { width: 80px; color: var(--muted-foreground); text-align: right; flex-shrink: 0; }
.f-size.h { font-size: 11.5px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.3px; }
.f-date { width: 140px; color: var(--muted-foreground); text-align: right; flex-shrink: 0; opacity: 0.7; }
.f-date.h { font-size: 11.5px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.3px; }
</style>
