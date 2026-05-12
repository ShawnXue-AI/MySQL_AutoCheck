<template>
  <div class="tasks-view">
    <div class="tasks-toolbar">
      <div class="tasks-count">
        <span class="terminal-prompt">$</span>
        <span>find tasks — {{ sortedTasks.length }} 条记录</span>
      </div>
      <button v-if="sortedTasks.length > 0" class="btn btn-ghost btn-sm" @click="refreshTasks">
        <svg viewBox="0 0 24 24" class="btn-icon"><path d="M17.65 6.35A7.958 7.958 0 0012 4a8 8 0 100 16 7.96 7.96 0 005.66-2.34A8 8 0 0012 4zM10 16v-4H7l5-5 5 5h-3v4h-4z" fill="currentColor"/></svg>
        刷新
      </button>
    </div>

    <div v-if="sortedTasks.length === 0" class="empty-state">
      <svg viewBox="0 0 48 48" class="empty-icon"><path d="M24 4C12.95 4 4 12.95 4 24s8.95 20 20 20 20-8.95 20-20S35.05 4 24 4zm-2 30h-4V18h4v16zm8 0h-4V18h4v16z" fill="currentColor"/></svg>
      <p class="empty-title">暂无历史记录</p>
      <p class="empty-desc">完成分析后，任务将显示在此处</p>
      <router-link to="/" class="btn btn-primary">生成报告</router-link>
    </div>

    <div v-else class="task-list">
      <div v-for="task in sortedTasks" :key="task.id" class="task-row">
        <div class="task-row-left">
          <span class="status-dot" :class="statusDotClass(task.status)"></span>
          <div class="task-info">
            <div class="task-row-name">{{ task.customer_name || '未命名客户' }}</div>
            <div class="task-row-meta">
              <span class="meta-time">{{ formatTime(task.created_at) }}</span>
              <span class="meta-sep">·</span>
              <code class="meta-file">{{ task.file_name }}</code>
            </div>
          </div>
        </div>
        <div class="task-row-right">
          <span class="badge" :class="statusBadgeClass(task.status)">{{ statusLabel(task.status) }}</span>
          <a v-if="task.zip_file" :href="zipUrl(task.id)" class="btn btn-outline btn-xs" title="下载">
            <svg viewBox="0 0 24 24" width="14"><path d="M19 9h-4V3H9v6H5l7 7 7-7zM5 18v2h14v-2H5z" fill="currentColor"/></svg>
          </a>
          <button class="btn btn-outline btn-xs" style="border-color: rgba(248,81,73,0.25); color: var(--accent-red);" title="删除" @click="handleDelete(task.id)">
            <svg viewBox="0 0 24 24" width="14"><path d="M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM19 4h-3.5l-1-1h-5l-1 1H5v2h14V4z" fill="currentColor"/></svg>
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
    statusDotClass(s) {
      return { completed: 'active', running: 'active', cancelled: 'idle', pending: 'idle', failed: 'error' }[s] || 'idle'
    },
    statusBadgeClass(s) {
      return { completed: 'badge-success', running: 'badge-info', cancelled: 'badge-muted', pending: 'badge-muted', failed: 'badge-error' }[s] || 'badge-muted'
    },
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
.tasks-count {
  display: flex; align-items: center; gap: 8px;
  font-size: 13px; color: var(--text-tertiary);
}
.terminal-prompt {
  font-family: var(--font-code);
  color: var(--accent-green);
  font-weight: 500;
  font-size: 13px;
}

.task-list { display: flex; flex-direction: column; gap: 3px; }

.task-row {
  display: flex; align-items: center; justify-content: space-between;
  padding: 14px 18px; background: var(--bg-secondary);
  border: 1px solid var(--border-default); border-radius: var(--radius-md);
  transition: all var(--transition-fast);
}
.task-row:hover { border-color: var(--border-hover); background: var(--bg-tertiary); }

.task-row-left { display: flex; align-items: center; gap: 12px; flex: 1; min-width: 0; }
.task-info { display: flex; flex-direction: column; gap: 3px; min-width: 0; }
.task-row-name { font-size: 14px; font-weight: 600; color: var(--text-primary); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.task-row-meta { display: flex; align-items: center; gap: 6px; font-size: 12px; color: var(--text-tertiary); }
.meta-time { font-family: var(--font-code); font-size: 11px; }
.meta-sep { color: var(--border-hover); }
.meta-file {
  font-family: var(--font-code);
  font-size: 10.5px;
  background: none;
  color: var(--text-muted);
  padding: 0;
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}

.task-row-right { display: flex; align-items: center; gap: 8px; flex-shrink: 0; }
</style>
