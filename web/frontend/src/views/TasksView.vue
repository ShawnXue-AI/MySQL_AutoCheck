<template>
  <div class="tasks-view">
    <div class="tasks-toolbar">
      <div class="tasks-count">{{ sortedTasks.length }} 条记录</div>
      <button v-if="sortedTasks.length > 0" class="btn btn-ghost btn-sm" @click="refreshTasks">刷新</button>
    </div>

    <div v-if="sortedTasks.length === 0" class="empty-state">
      <svg viewBox="0 0 48 48" class="empty-icon"><path d="M24 4C12.95 4 4 12.95 4 24s8.95 20 20 20 20-8.95 20-20S35.05 4 24 4zm-2 30h-4V18h4v16zm8 0h-4V18h4v16z" fill="currentColor"/></svg>
      <p class="empty-title">暂无历史记录</p>
      <p class="empty-desc">完成分析后，任务将显示在此处</p>
      <router-link to="/" class="btn btn-primary">去生成报告</router-link>
    </div>

    <div v-else class="task-list">
      <div v-for="task in sortedTasks" :key="task.id" class="task-row">
        <div class="task-row-left">
          <div class="task-dot" :class="task.status"></div>
          <div class="task-info">
            <div class="task-row-name">{{ task.customer_name || '未命名客户' }}</div>
            <div class="task-row-meta">
              <span>{{ formatTime(task.created_at) }}</span>
              <span class="meta-sep">·</span>
              <span class="meta-file">{{ task.file_name }}</span>
            </div>
          </div>
        </div>
        <div class="task-row-right">
          <span class="task-badge" :class="task.status">{{ statusLabel(task.status) }}</span>
          <a v-if="task.zip_file" :href="zipUrl(task.id)" class="btn-icon-btn" title="下载">
            <svg viewBox="0 0 24 24" width="16"><path d="M19 9h-4V3H9v6H5l7 7 7-7zM5 18v2h14v-2H5z" fill="currentColor"/></svg>
          </a>
          <button class="btn-icon-btn danger" title="删除" @click="handleDelete(task.id)">
            <svg viewBox="0 0 24 24" width="16"><path d="M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM19 4h-3.5l-1-1h-5l-1 1H5v2h14V4z" fill="currentColor"/></svg>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { getTasks, deleteTask, buildZipDownloadUrl } from '../api/index.js'
import { toastSuccess, toastError } from '../components/Toast.vue'

export default {
  name: 'TasksView',
  data() {
    return { tasks: [], pollTimer: null }
  },
  computed: {
    sortedTasks() {
      return [...this.tasks].sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
    }
  },
  mounted() {
    this.fetchTasks()
    this.pollTimer = setInterval(() => this.fetchTasks(), 5000)
  },
  beforeUnmount() {
    if (this.pollTimer) clearInterval(this.pollTimer)
  },
  methods: {
    async fetchTasks() { try { const { data } = await getTasks(); this.tasks = data.tasks || [] } catch {} },
    refreshTasks() { this.fetchTasks() },
    statusLabel(s) { return { pending: '等待中', running: '运行中', completed: '已完成', failed: '失败', cancelled: '已取消' }[s] || s },
    formatTime(t) { return t ? new Date(t).toLocaleString('zh-CN', { hour12: false }) : '' },
    zipUrl(taskId) { return buildZipDownloadUrl(taskId) },
    async handleDelete(taskId) {
      try {
        await deleteTask(taskId)
        this.tasks = this.tasks.filter(t => t.id !== taskId)
        toastSuccess('已删除')
      } catch (e) { toastError('删除失败') }
    }
  }
}
</script>

<style scoped>
.tasks-view { max-width: 960px; margin: 0 auto; }

.tasks-toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.tasks-count { font-size: 13px; color: #8b949e; }

.btn { display: inline-flex; align-items: center; gap: 8px; padding: 9px 20px; border: none; border-radius: 8px; font-size: 14px; font-weight: 500; cursor: pointer; text-decoration: none; transition: all 0.15s; }
.btn-ghost { background: transparent; color: #8b949e; padding: 6px 12px; font-size: 13px; }
.btn-ghost:hover { color: #c9d1d9; background: #1c2128; }
.btn-sm { padding: 6px 12px; font-size: 13px; }
.btn-primary { background: #238636; color: #fff; }
.btn-primary:hover { background: #2ea043; }

.empty-state { text-align: center; padding: 80px 20px; }
.empty-icon { width: 48px; height: 48px; margin: 0 auto 16px; color: #30363d; }
.empty-title { font-size: 16px; font-weight: 600; color: #c9d1d9; margin-bottom: 8px; }
.empty-desc { font-size: 14px; color: #8b949e; margin-bottom: 24px; }

.task-list { display: flex; flex-direction: column; gap: 4px; }
.task-row { display: flex; align-items: center; justify-content: space-between; padding: 14px 16px; background: #161b22; border: 1px solid #21262d; border-radius: 10px; transition: border-color 0.15s; }
.task-row:hover { border-color: #30363d; }

.task-row-left { display: flex; align-items: center; gap: 12px; flex: 1; min-width: 0; }
.task-dot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }
.task-dot.completed { background: #3fb950; }
.task-dot.running { background: #3fb950; animation: pulse 1.5s infinite; }
.task-dot.cancelled { background: #484f58; }
.task-dot.pending { background: #30363d; }
.task-dot.failed { background: #da3633; }
@keyframes pulse { 0%,100% { opacity: 1; } 50% { opacity: 0.3; } }

.task-info { display: flex; flex-direction: column; gap: 2px; min-width: 0; }
.task-row-name { font-size: 14px; font-weight: 600; color: #f0f6fc; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.task-row-meta { display: flex; align-items: center; gap: 6px; font-size: 12px; color: #8b949e; }
.meta-sep { color: #30363d; }
.meta-file { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

.task-row-right { display: flex; align-items: center; gap: 8px; flex-shrink: 0; }
.task-badge { font-size: 12px; font-weight: 500; padding: 2px 10px; border-radius: 20px; }
.task-badge.completed { background: rgba(35,134,54,0.15); color: #3fb950; }
.task-badge.running { background: rgba(35,134,54,0.15); color: #3fb950; }
.task-badge.cancelled { background: #21262d; color: #484f58; }
.task-badge.pending { background: #21262d; color: #8b949e; }
.task-badge.failed { background: rgba(218,54,51,0.15); color: #f85149; }

.btn-icon-btn { display: flex; align-items: center; justify-content: center; width: 32px; height: 32px; background: #21262d; border: none; border-radius: 6px; color: #8b949e; cursor: pointer; text-decoration: none; transition: all 0.15s; }
.btn-icon-btn:hover { background: #1c2128; color: #3fb950; }
.btn-icon-btn.danger:hover { background: #da3633; color: #fff; }
</style>
