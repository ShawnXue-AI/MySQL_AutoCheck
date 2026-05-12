<template>
  <div class="home-view">
    <div class="two-col">
      <div class="col-left">
        <!-- Upload Card -->
        <section class="card upload-card">
          <div class="card-header">
            <svg viewBox="0 0 24 24" class="card-icon"><path d="M9 16h6v-6H9v6zm-4 4V4h10l4 4v12H5z" fill="currentColor"/></svg>
            <div>
              <h2 class="card-title">上传巡检日志</h2>
              <p class="card-desc">$ upload --format .log .zip .tar.gz</p>
            </div>
          </div>

          <div class="field-group">
            <label class="field-label">客户名称 <span class="required">*</span></label>
            <div class="input-prefix">
              <span class="input-prefix-icon">@</span>
              <input type="text" class="field-input has-prefix" :class="{ error: showCustomerError }" v-model="customerName" placeholder="customer_name" required @input="showCustomerError = false" />
            </div>
            <span class="field-error" v-if="showCustomerError">错误: 客户名称为必填字段</span>
          </div>

          <div class="drop-zone" :class="{ dragging: isDragging, filled: uploadedFiles.length > 0 }"
            @dragenter.prevent="isDragging = true" @dragover.prevent="isDragging = true"
            @dragleave.prevent="isDragging = false" @drop.prevent="handleDrop">
            <div v-if="uploadedFiles.length === 0" class="drop-placeholder">
              <div class="drop-terminal">
                <span class="terminal-prompt">$</span>
                <span class="terminal-cursor">drop files here</span>
              </div>
              <p class="drop-hint">// 或使用下方按钮选择文件/文件夹</p>
              <div class="upload-options">
                <label class="btn btn-outline">
                  <svg viewBox="0 0 24 24" class="btn-icon"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8l-6-6zm-1 7V3.5L18.5 9H13z" fill="currentColor"/></svg>
                  选择文件
                  <input type="file" multiple @change="handleFileSelect" />
                </label>
                <label class="btn btn-outline">
                  <svg viewBox="0 0 24 24" class="btn-icon"><path d="M20 6h-8l-2-2H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2zm-6 10H6v-2h8v2zm4-4H6v-2h12v2z" fill="currentColor"/></svg>
                  选择文件夹
                  <input type="file" multiple webkitdirectory @change="handleFileSelect" />
                </label>
              </div>
            </div>
            <div v-else class="drop-result">
              <svg viewBox="0 0 24 24" class="check-icon"><circle cx="12" cy="12" r="10" fill="var(--accent-green-subtle)" stroke="var(--accent-green)" stroke-width="2"/><path d="M8 12l3 3 5-5" stroke="var(--accent-green)" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/></svg>
              <span class="drop-result-text">{{ uploadedFiles.length }} 个文件已就绪</span>
            </div>
          </div>

          <div v-if="uploadedFiles.length > 0" class="file-list">
            <div class="file-list-header">
              <span class="flh-name">filename</span>
              <span class="flh-size">size</span>
            </div>
            <div class="file-item" v-for="(f, i) in uploadedFiles.slice(0, 8)" :key="i">
              <span class="file-index">{{ String(i + 1).padStart(2, '0') }}</span>
              <svg viewBox="0 0 24 24" class="file-icon"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8l-6-6z" fill="currentColor" opacity="0.8"/><path d="M14 2v6h6" fill="currentColor" opacity="0.4"/></svg>
              <span class="file-name">{{ f.name || f }}</span>
              <span class="file-size" v-if="f.size">{{ formatSize(f.size) }}</span>
            </div>
            <div v-if="uploadedFiles.length > 8" class="file-more">
              <span class="terminal-prompt">...</span> 还有 {{ uploadedFiles.length - 8 }} 个文件
            </div>
          </div>

          <div class="action-row">
            <button class="btn btn-primary" :disabled="busy || uploadedFiles.length === 0" @click="submitUpload">
              <svg viewBox="0 0 24 24" class="btn-icon"><path d="M9 16h6v-6h4l-7-7-7 7h4v6zm-4 2h14v2H5v-2z" fill="currentColor"/></svg>
              {{ busy ? '上传中...' : '上传到服务器' }}
            </button>
            <button v-if="uploadedFiles.length > 0" class="btn btn-ghost" @click="resetUpload">
              <svg viewBox="0 0 24 24" class="btn-icon"><path d="M17.65 6.35A7.958 7.958 0 0012 4a8 8 0 100 16 7.96 7.96 0 005.66-2.34A8 8 0 0012 4zM10 16v-4H7l5-5 5 5h-3v4h-4z" fill="currentColor"/></svg>
              重新选择
            </button>
          </div>
        </section>

        <!-- Analysis Options Card -->
        <section class="card">
          <div class="card-header">
            <svg viewBox="0 0 24 24" class="card-icon"><path d="M19.43 12.98c.04-.32.07-.64.07-.98s-.03-.66-.07-.98l2.11-1.65c.19-.15.24-.42.12-.64l-2-3.46c-.12-.22-.39-.3-.61-.22l-2.49 1c-.52-.4-1.08-.73-1.69-.98l-.38-2.65C14.46 2.18 14.25 2 14 2h-4c-.25 0-.46.18-.49.42l-.38 2.65c-.61.25-1.17.59-1.69.98l-2.49-1c-.23-.09-.49 0-.61.22l-2 3.46c-.13.22-.07.49.12.64l2.11 1.65c-.04.32-.07.65-.07.98s.03.66.07.98l-2.11 1.65c-.19.15-.24.42-.12.64l2 3.46c.12.22.39.3.61.22l2.49-1c.52.4 1.08.73 1.69.98l.38 2.65c.03.24.24.42.49.42h4c.25 0 .46-.18.49-.42l.38-2.65c.61-.25 1.17-.59 1.69-.98l2.49 1c.23.09.49 0 .61-.22l2-3.46c.12-.22.07-.49-.12-.64l-2.11-1.65zM12 15.5c-1.93 0-3.5-1.57-3.5-3.5s1.57-3.5 3.5-3.5 3.5 1.57 3.5 3.5-1.57 3.5-3.5 3.5z" fill="currentColor"/></svg>
            <div>
              <h2 class="card-title">分析选项</h2>
              <p class="card-desc">// 选择需要生成的报告类型</p>
            </div>
          </div>
          <div class="toggle-group">
            <label class="toggle-row">
              <div class="toggle-info">
                <span class="toggle-label">Word 报告</span>
                <span class="toggle-desc">生成 .docx 格式巡检详细报告（含自动目录）</span>
              </div>
              <div class="toggle-switch" :class="{ on: options.generate_word }" @click="options.generate_word = !options.generate_word">
                <div class="toggle-knob"></div>
              </div>
            </label>
            <label class="toggle-row">
              <div class="toggle-info">
                <span class="toggle-label">Excel 统计表</span>
                <span class="toggle-desc">生成 .xlsx 格式巡检统计报表</span>
              </div>
              <div class="toggle-switch" :class="{ on: options.generate_excel }" @click="options.generate_excel = !options.generate_excel">
                <div class="toggle-knob"></div>
              </div>
            </label>
          </div>
          <div class="field-group" style="margin-top: 8px;">
            <label class="field-label">报告字体</label>
            <select class="field-select" v-model="reportFont">
              <option value="微软雅黑">微软雅黑</option>
              <option value="宋体">宋体</option>
              <option value="仿宋">仿宋</option>
              <option value="黑体">黑体</option>
              <option value="楷体">楷体</option>
              <option value="Arial">Arial</option>
              <option value="Times New Roman">Times New Roman</option>
            </select>
          </div>
          <button class="btn btn-accent btn-block" :disabled="!uploadComplete || busy" @click="startAnalysis">
            <svg viewBox="0 0 24 24" class="btn-icon"><path d="M8 5v14l11-7z" fill="currentColor"/></svg>
            {{ busy ? '执行中...' : '开始分析' }}
          </button>
        </section>
      </div>

      <div class="col-right">
        <!-- Stats Card -->
        <section class="card stat-card">
          <div class="card-header">
            <svg viewBox="0 0 24 24" class="card-icon"><path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-7 3c1.93 0 3.5 1.57 3.5 3.5S13.93 13 12 13s-3.5-1.57-3.5-3.5S10.07 6 12 6zm7 13H5v-.23c0-.62.28-1.2.76-1.58C7.47 15.82 9.64 15 12 15s4.53.82 6.24 2.19c.48.38.76.97.76 1.58V19z" fill="currentColor"/></svg>
            <div>
              <h2 class="card-title">实时状态</h2>
              <p class="card-desc">$ stats --watch</p>
            </div>
          </div>
          <div class="stat-grid">
            <div class="stat-box">
              <span class="stat-num blue">{{ stats.active_users }}</span>
              <span class="stat-num-label">在线人数</span>
            </div>
            <div class="stat-box">
              <span class="stat-num green">{{ stats.total_completed }}</span>
              <span class="stat-num-label">累计完成</span>
            </div>
            <div class="stat-box">
              <span class="stat-num amber">{{ stats.running_tasks }}</span>
              <span class="stat-num-label">运行中任务</span>
            </div>
          </div>
        </section>

        <!-- Running Tasks Card -->
        <section class="card" v-if="runningTasks.length > 0">
          <div class="card-header">
            <svg viewBox="0 0 24 24" class="card-icon" style="color: var(--accent-cyan);"><path d="M13 2.03v2.02c4.39.54 7.5 4.53 6.96 8.92-.46 3.64-3.32 6.53-6.96 6.96v2.02c5.49-.66 9.5-5.6 8.84-11.09-.55-4.58-4.09-8.19-8.84-8.83zM4.26 17.86c1.94 2.86 5.06 4.7 8.74 4.76v-2.02c-2.41-.12-4.53-1.44-5.74-3.28l-3 .54zm.03-11.72l3 .54c1.21-1.84 3.33-3.16 5.74-3.28V1.38c-3.68.06-6.8 1.9-8.74 4.76z" fill="currentColor"/></svg>
            <div>
              <h2 class="card-title">运行中任务</h2>
              <p class="card-desc">$ ps aux | grep analysis  →  {{ runningTasks.length }} 个进程</p>
            </div>
          </div>
          <div class="running-list">
            <div class="running-item" v-for="t in runningTasks" :key="t.id">
              <span class="status-dot active"></span>
              <span class="running-name">{{ t.customer_name || '未命名' }}</span>
              <code class="running-file">{{ t.file_name }}</code>
              <button class="btn btn-xs" style="border: 1px solid var(--accent-red); color: var(--accent-red); background: transparent;" @click="cancelRunningTask(t.id)" title="终止进程">
                kill
              </button>
            </div>
          </div>
        </section>
      </div>
    </div>

    <ProcessModal
      :visible="modalVisible"
      :phase="modalPhase"
      :progress="modalProgress"
      :title-text="modalTitle"
      :sub-text="modalSub"
      :error-msg="modalError"
      :download-url="modalDownloadUrl"
      :cancelling="modalCancelling"
      @cancel="handleModalCancel"
      @close="handleModalClose"
    />
  </div>
</template>

<script>
import { uploadFiles, uploadArchive, startAnalysis, getStats, getTaskStatus, cancelTask, getRunningTasks, buildZipDownloadUrl, resetSession } from '../api/index.js'
import { toastSuccess, toastError } from '../components/Toast.vue'
import ProcessModal from '../components/ProcessModal.vue'

export default {
  name: 'HomeView',
  components: { ProcessModal },
  data() {
    return {
      isDragging: false,
      customerName: '',
      showCustomerError: false,
      selectedFiles: [],
      uploadedFiles: [],
      uploadComplete: false,
      busy: false,
      modalVisible: false,
      modalPhase: 'upload',
      modalProgress: 0,
      modalTitle: '',
      modalSub: '',
      modalError: '',
      modalDownloadUrl: '',
      modalCancelling: false,
      reportFont: '微软雅黑',
      options: { generate_word: true, generate_excel: false },
      currentTaskId: null,
      taskPollTimer: null,
      runningTasks: [],
      stats: { active_users: 0, total_completed: 0, running_tasks: 0 },
      statsTimer: null,
      runningTimer: null
    }
  },
  mounted() {
    this.pollStats(); this.pollRunning()
    this.statsTimer = setInterval(() => this.pollStats(), 5000)
    this.runningTimer = setInterval(() => this.pollRunning(), 5000)
  },
  beforeUnmount() {
    if (this.taskPollTimer) clearInterval(this.taskPollTimer)
    if (this.statsTimer) clearInterval(this.statsTimer)
    if (this.runningTimer) clearInterval(this.runningTimer)
  },
  methods: {
    async pollStats() {
      try { const { data } = await getStats(); this.stats = data; this.$emit('stats-update', data) } catch {}
    },
    async pollRunning() {
      try { const { data } = await getRunningTasks(); this.runningTasks = data.tasks || [] } catch {}
    },
    handleDrop(e) {
      this.isDragging = false
      const items = e.dataTransfer.items
      if (!items) { this.handleFiles(e.dataTransfer.files); return }
      const files = []
      const processEntry = (entry, path) => {
        return new Promise((resolve) => {
          if (entry.isFile) { entry.file(f => { f._fullPath = path; files.push(f); resolve() }) }
          else if (entry.isDirectory) { const reader = entry.createReader(); reader.readEntries(entries => { Promise.all(entries.map(e => processEntry(e, path + '/' + e.name))).then(resolve) }) }
          else { resolve() }
        })
      }
      Promise.all([...items].filter(i => i.webkitGetAsEntry).map(i => processEntry(i.webkitGetAsEntry(), i.webkitGetAsEntry().name))).then(() => { if (files.length) this.handleFiles(files) })
    },
    handleFileSelect(e) { this.handleFiles(e.target.files) },
    handleFiles(fileList) {
      const maxLogSize = 10 * 1024 * 1024
      const maxArchiveSize = 30 * 1024 * 1024
      const files = Array.from(fileList).filter(f => {
        if (!f.name.endsWith('.log') && !f.name.endsWith('.zip') && !f.name.endsWith('.tar.gz') && !f.name.endsWith('.tgz') && !(f.size > 0 && f._fullPath)) return false
        const isArchive = f.name.endsWith('.zip') || f.name.endsWith('.tar.gz') || f.name.endsWith('.tgz')
        if (isArchive && f.size > maxArchiveSize) { toastError(`压缩包 ${f.name} 超过30MB限制`); return false }
        if (!isArchive && f.size > maxLogSize) { toastError(`日志文件 ${f.name} 超过10MB限制`); return false }
        return true
      })
      if (!files.length) return
      resetSession()
      this.selectedFiles = files; this.uploadedFiles = files; this.uploadComplete = false
    },
    resetUpload() { this.uploadedFiles = []; this.selectedFiles = []; this.uploadComplete = false },
    formatSize(bytes) { if (bytes < 1024) return bytes + ' B'; if (bytes < 1048576) return (bytes/1024).toFixed(1) + ' KB'; return (bytes/1048576).toFixed(1) + ' MB' },
    async submitUpload() {
      if (!this.customerName.trim()) { this.showCustomerError = true; toastError('请输入客户名称'); return }
      this.showCustomerError = false; this.busy = true
      this.modalVisible = true; this.modalPhase = 'upload'; this.modalProgress = 0; this.modalTitle = '正在上传文件'; this.modalSub = ''; this.modalError = ''; this.modalCancelling = false; this.currentTaskId = null
      try {
        this.modalProgress = 15
        const af = this.selectedFiles.find(f => f.name.endsWith('.zip') || f.name.endsWith('.tar.gz') || f.name.endsWith('.tgz'))
        this.modalProgress = 40
        await (af ? uploadArchive(af) : uploadFiles(this.selectedFiles))
        this.modalProgress = 100
        this.modalTitle = '上传完成'
        this.modalSub = '即将返回主页面...'
        this.uploadComplete = true
        this.busy = false
        toastSuccess('上传成功')
        setTimeout(() => {
          this.modalVisible = false
          this.modalProgress = 0
        }, 1200)
      } catch (e) {
        this.modalPhase = 'error'; this.modalTitle = '上传失败'; this.modalError = e.response?.data?.error || e.message
        this.uploadComplete = false
        this.busy = false
      }
    },
    async startAnalysis() {
      this.busy = true; this.modalVisible = true
      this.modalPhase = 'running'; this.modalProgress = 10; this.modalTitle = '正在提交分析'; this.modalSub = ''; this.modalError = ''; this.modalCancelling = false; this.modalDownloadUrl = ''
      const name = this.customerName.trim() || '未命名客户'
      try {
        const { data } = await startAnalysis(name, this.options, this.reportFont)
        this.currentTaskId = data.task_id
        this.modalTitle = '正在生成报告'; this.modalSub = name; this.modalProgress = 20
        this.taskPollTimer = setInterval(() => this.pollTask(data.task_id), 2000)
      } catch (e) {
        this.modalPhase = 'error'; this.modalTitle = '启动失败'; this.modalError = e.response?.data?.error || e.message
        this.busy = false
      }
    },
    async pollTask(taskId) {
      try {
        const { data } = await getTaskStatus(taskId)
        this.stats = data.stats; this.$emit('stats-update', data.stats)
        const t = data.task
        const pct = { pending: 20, running: 60, completed: 100, failed: 100, cancelled: 100 }[t.status] || 20
        this.modalProgress = pct
        if (t.progress) this.modalSub = t.progress
        if (t.status === 'completed') {
          clearInterval(this.taskPollTimer)
          this.modalPhase = 'done'; this.modalTitle = '报告生成完成'; this.modalSub = t.zip_file || ''; this.modalProgress = 100
          this.modalDownloadUrl = buildZipDownloadUrl(t.id)
          this.busy = false
          toastSuccess('报告生成完成')
        } else if (t.status === 'failed') {
          clearInterval(this.taskPollTimer)
          this.modalPhase = 'error'; this.modalTitle = '分析失败'; this.modalError = t.error || '未知错误'
          this.busy = false
        } else if (t.status === 'cancelled') {
          clearInterval(this.taskPollTimer)
          this.handleModalClose()
        }
      } catch {}
    },
    async handleModalCancel() {
      if (this.modalPhase === 'upload') {
        this.modalCancelling = true
        await new Promise(r => setTimeout(r, 800))
        this.modalVisible = false
        this.modalCancelling = false
        return
      }
      if (this.modalPhase === 'running') {
        this.modalCancelling = true
        if (this.currentTaskId) {
          try { await cancelTask(this.currentTaskId) } catch {}
          await new Promise(r => setTimeout(r, 1500))
        }
        this.modalVisible = false
        toastSuccess('任务已终止，文件已清理')
        this.resetAll()
      }
    },
    handleModalClose() {
      this.modalVisible = false
      clearInterval(this.taskPollTimer)
      this.busy = false
      this.modalProgress = 0
      this.modalCancelling = false
      this.currentTaskId = null
    },
    resetAll() {
      clearInterval(this.taskPollTimer)
      this.busy = false; this.uploadComplete = false; this.uploadedFiles = []; this.selectedFiles = []
      this.modalVisible = false; this.modalProgress = 0; this.modalDownloadUrl = ''; this.modalCancelling = false
      this.currentTaskId = null
    },
    async cancelRunningTask(taskId) {
      try { await cancelTask(taskId); toastSuccess('任务已终止') } catch (e) { toastError('取消失败') }
    }
  }
}
</script>

<style scoped>
.home-view { max-width: 1160px; margin: 0 auto; }
.two-col { display: grid; grid-template-columns: 1fr 350px; gap: 24px; align-items: start; }

/* Upload-specific card */
.upload-card { border-color: var(--border-hover); }

/* Input with prefix */
.input-prefix {
  position: relative;
}
.input-prefix-icon {
  position: absolute; left: 12px; top: 50%; transform: translateY(-50%);
  font-family: var(--font-code);
  font-size: 14px; color: var(--text-muted); font-weight: 500;
  pointer-events: none;
}
.field-input.has-prefix { padding-left: 32px; }

/* Drop Zone — terminal styled */
.drop-zone {
  border: 2px dashed var(--border-hover);
  border-radius: var(--radius-md);
  padding: 40px 20px;
  text-align: center;
  cursor: pointer;
  transition: all var(--transition-normal);
  margin-bottom: 10px;
  background: var(--bg-primary);
}
.drop-zone:hover, .drop-zone.dragging {
  border-color: var(--accent-blue);
  background: var(--accent-blue-subtle);
}
.drop-zone.filled {
  border-color: var(--accent-green);
  border-style: solid;
  background: var(--accent-green-subtle);
  padding: 18px 20px;
}

.drop-terminal {
  display: flex; align-items: center; gap: 8px;
  justify-content: center; margin-bottom: 10px;
}
.terminal-prompt {
  font-family: var(--font-code);
  color: var(--accent-green);
  font-weight: 500;
  font-size: 15px;
}
.terminal-cursor {
  font-family: var(--font-code);
  color: var(--text-primary);
  font-size: 15px;
}
.terminal-cursor::after {
  content: '_';
  color: var(--accent-blue);
  animation: blink-cursor 1s step-end infinite;
}
@keyframes blink-cursor {
  0%, 100% { opacity: 1; } 50% { opacity: 0; }
}

.drop-hint {
  font-family: var(--font-code);
  font-size: 11.5px;
  color: var(--text-muted);
  margin-bottom: 16px;
}

.upload-options { display: flex; gap: 10px; justify-content: center; }

.check-icon { width: 28px; height: 28px; flex-shrink: 0; }
.drop-result { display: flex; align-items: center; gap: 12px; justify-content: center; }
.drop-result-text { font-family: var(--font-code); font-size: 14px; color: var(--accent-green); }

/* File List */
.file-list { display: flex; flex-direction: column; gap: 2px; margin-top: 2px; }

.file-list-header {
  display: flex; align-items: center; gap: 8px;
  padding: 5px 10px; font-family: var(--font-code);
  font-size: 10.5px; color: var(--text-muted);
  text-transform: uppercase; letter-spacing: 0.5px;
  border-bottom: 1px solid var(--border-default);
  margin-bottom: 4px;
}
.flh-name { flex: 1; }
.flh-size { width: 70px; text-align: right; flex-shrink: 0; }

.file-item {
  display: flex; align-items: center; gap: 8px;
  padding: 7px 10px; background: var(--bg-input);
  border-radius: var(--radius-sm); font-size: 13px;
  border: 1px solid transparent;
  transition: all var(--transition-fast);
}
.file-item:hover { border-color: var(--border-hover); }

.file-index {
  font-family: var(--font-code);
  font-size: 10.5px; color: var(--text-muted);
  width: 20px; flex-shrink: 0; text-align: right;
}
.file-icon { width: 14px; height: 14px; color: var(--accent-blue); flex-shrink: 0; }
.file-name { flex: 1; color: var(--text-secondary); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-family: var(--font-code); font-size: 12px; }
.file-size { color: var(--text-muted); flex-shrink: 0; font-family: var(--font-code); font-size: 11px; width: 70px; text-align: right; }
.file-more { font-family: var(--font-code); font-size: 11.5px; color: var(--text-muted); padding: 4px 10px; display: flex; align-items: center; gap: 6px; }
.file-more .terminal-prompt { font-size: 12px; }

.action-row { display: flex; gap: 8px; margin-top: 12px; }

/* Toggles */
.toggle-group { margin-bottom: 8px; }
.toggle-row {
  display: flex; align-items: center; justify-content: space-between;
  padding: 14px 0; cursor: pointer;
  border-bottom: 1px solid var(--border-default);
}
.toggle-row:last-child { border-bottom: none; }
.toggle-info { flex: 1; }
.toggle-label { display: block; font-size: 14px; font-weight: 500; color: var(--text-secondary); margin-bottom: 2px; }
.toggle-desc { font-size: 12px; color: var(--text-muted); font-family: var(--font-code); }

.toggle-switch {
  width: 44px; height: 24px; background: var(--border-hover);
  border-radius: 12px; position: relative; transition: background var(--transition-normal);
  cursor: pointer; flex-shrink: 0;
}
.toggle-switch.on { background: var(--accent-blue); }
.toggle-knob {
  width: 18px; height: 18px; background: #fff;
  border-radius: 50%; position: absolute; top: 3px; left: 3px;
  transition: transform var(--transition-normal);
  box-shadow: var(--shadow-sm);
}
.toggle-switch.on .toggle-knob { transform: translateX(20px); }

/* Stats */
.stat-card { background: linear-gradient(135deg, var(--bg-secondary) 0%, var(--bg-tertiary) 100%); }
.stat-grid { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 8px; }
.stat-box {
  text-align: center; padding: 18px 8px;
  background: var(--bg-primary); border-radius: var(--radius-sm);
  border: 1px solid var(--border-default);
}
.stat-num {
  display: block;
  font-family: var(--font-code);
  font-size: 28px; font-weight: 600;
  color: var(--text-primary); line-height: 1.1; margin-bottom: 4px;
  font-variant-numeric: tabular-nums;
}
.stat-num.green { color: var(--accent-green); }
.stat-num.blue { color: var(--accent-blue); }
.stat-num.amber { color: var(--accent-amber); }
.stat-num-label {
  font-size: 10px; color: var(--text-tertiary);
  text-transform: uppercase; letter-spacing: 0.5px;
  font-weight: 500;
}

/* Running Tasks */
.running-list { display: flex; flex-direction: column; gap: 4px; }
.running-item {
  display: flex; align-items: center; gap: 8px;
  padding: 9px 12px; background: var(--bg-primary);
  border-radius: var(--radius-sm); font-size: 13px;
  border: 1px solid var(--border-default);
}
.running-name { font-weight: 500; color: var(--text-secondary); min-width: 70px; }
.running-file { font-size: 11px; }

@media (max-width: 900px) { .two-col { grid-template-columns: 1fr; } }
</style>
