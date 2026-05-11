<template>
  <div class="home-view">
    <div class="two-col">
      <div class="col-left">
        <section class="card">
          <div class="card-header">
            <svg viewBox="0 0 24 24" class="card-icon"><path d="M9 16h6v-6H9v6zm-4 4V4h10l4 4v12H5z" fill="currentColor"/></svg>
            <div>
              <h2 class="card-title">上传巡检日志</h2>
              <p class="card-desc">支持 .log / .zip / .tar.gz 格式，可拖拽文件夹</p>
            </div>
          </div>

          <div class="field-group">
            <label class="field-label">客户名称 <span class="required-star">*</span></label>
            <input type="text" class="field-input" :class="{ error: showCustomerError }" v-model="customerName" placeholder="必填，用于命名报告文件" required @input="showCustomerError = false" />
            <span class="field-error" v-if="showCustomerError">请输入客户名称</span>
          </div>

          <div class="drop-zone" :class="{ dragging: isDragging, filled: uploadedFiles.length > 0 }"
            @dragenter.prevent="isDragging = true" @dragover.prevent="isDragging = true"
            @dragleave.prevent="isDragging = false" @drop.prevent="handleDrop">
            <div v-if="uploadedFiles.length === 0" class="drop-placeholder">
              <svg viewBox="0 0 48 48" class="upload-icon"><path d="M24 16v12m0 0l-4-4m4 4l4-4M14 28v4a4 4 0 004 4h12a4 4 0 004-4v-4" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round"/><path d="M24 4v12m0 0l-4-4m4 4l4-4" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" opacity="0.5"/></svg>
              <p class="drop-text">拖拽文件或文件夹到此处</p>
              <p class="drop-hint">或点击下方按钮选择</p>
              <div class="upload-options">
                <label class="btn btn-outline">
                  <svg viewBox="0 0 24 24" class="btn-icon"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8l-6-6zm-1 7V3.5L18.5 9H13z" fill="currentColor"/></svg>
                  上传文件
                  <input type="file" multiple @change="handleFileSelect" />
                </label>
                <label class="btn btn-outline">
                  <svg viewBox="0 0 24 24" class="btn-icon"><path d="M20 6h-8l-2-2H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2zm-6 10H6v-2h8v2zm4-4H6v-2h12v2z" fill="currentColor"/></svg>
                  上传文件夹
                  <input type="file" multiple webkitdirectory @change="handleFileSelect" />
                </label>
              </div>
            </div>
            <div v-else class="drop-result">
              <svg viewBox="0 0 24 24" class="check-svg"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z" fill="#3fb950"/></svg>
              <span class="drop-text">已选择 {{ uploadedFiles.length }} 个文件</span>
            </div>
          </div>

          <div v-if="uploadedFiles.length > 0" class="file-list">
            <div class="file-item" v-for="(f, i) in uploadedFiles.slice(0, 8)" :key="i">
              <svg viewBox="0 0 24 24" class="file-icon"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8l-6-6z" fill="currentColor" opacity="0.8"/><path d="M14 2v6h6" fill="currentColor" opacity="0.4"/></svg>
              <span class="file-name">{{ f.name || f }}</span>
              <span class="file-size" v-if="f.size">{{ formatSize(f.size) }}</span>
            </div>
            <div v-if="uploadedFiles.length > 8" class="file-more">... 还有 {{ uploadedFiles.length - 8 }} 个文件</div>
          </div>

          <div class="action-row">
            <button class="btn btn-primary" :disabled="busy || uploadedFiles.length === 0" @click="submitUpload">
              <svg viewBox="0 0 24 24" class="btn-icon"><path d="M9 16h6v-6h4l-7-7-7 7h4v6zm-4 2h14v2H5v-2z" fill="currentColor"/></svg>
              {{ busy ? '处理中...' : '上传到服务器' }}
            </button>
            <button v-if="uploadedFiles.length > 0" class="btn btn-ghost" @click="resetUpload">重新选择</button>
          </div>
        </section>

        <section class="card">
          <div class="card-header">
            <svg viewBox="0 0 24 24" class="card-icon"><path d="M19.43 12.98c.04-.32.07-.64.07-.98s-.03-.66-.07-.98l2.11-1.65c.19-.15.24-.42.12-.64l-2-3.46c-.12-.22-.39-.3-.61-.22l-2.49 1c-.52-.4-1.08-.73-1.69-.98l-.38-2.65C14.46 2.18 14.25 2 14 2h-4c-.25 0-.46.18-.49.42l-.38 2.65c-.61.25-1.17.59-1.69.98l-2.49-1c-.23-.09-.49 0-.61.22l-2 3.46c-.13.22-.07.49.12.64l2.11 1.65c-.04.32-.07.65-.07.98s.03.66.07.98l-2.11 1.65c-.19.15-.24.42-.12.64l2 3.46c.12.22.39.3.61.22l2.49-1c.52.4 1.08.73 1.69.98l.38 2.65c.03.24.24.42.49.42h4c.25 0 .46-.18.49-.42l.38-2.65c.61-.25 1.17-.59 1.69-.98l2.49 1c.23.09.49 0 .61-.22l2-3.46c.12-.22.07-.49-.12-.64l-2.11-1.65zM12 15.5c-1.93 0-3.5-1.57-3.5-3.5s1.57-3.5 3.5-3.5 3.5 1.57 3.5 3.5-1.57 3.5-3.5 3.5z" fill="currentColor"/></svg>
            <div>
              <h2 class="card-title">分析选项</h2>
              <p class="card-desc">选择需要生成的报告类型</p>
            </div>
          </div>
          <div class="toggle-group">
            <label class="toggle-row">
              <div class="toggle-info">
                <span class="toggle-label">Word 报告</span>
                <span class="toggle-desc">.docx 格式巡检详细报告（含自动目录）</span>
              </div>
              <div class="toggle" :class="{ on: options.generate_word }" @click="options.generate_word = !options.generate_word"><div class="toggle-knob"></div></div>
            </label>
            <label class="toggle-row">
              <div class="toggle-info">
                <span class="toggle-label">Excel 统计表</span>
                <span class="toggle-desc">.xlsx 格式巡检统计报表</span>
              </div>
              <div class="toggle" :class="{ on: options.generate_excel }" @click="options.generate_excel = !options.generate_excel"><div class="toggle-knob"></div></div>
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
            {{ busy ? '处理中...' : '开始分析' }}
          </button>
        </section>
      </div>

      <div class="col-right">
        <section class="card stat-card">
          <div class="card-header">
            <svg viewBox="0 0 24 24" class="card-icon"><path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-7 3c1.93 0 3.5 1.57 3.5 3.5S13.93 13 12 13s-3.5-1.57-3.5-3.5S10.07 6 12 6zm7 13H5v-.23c0-.62.28-1.2.76-1.58C7.47 15.82 9.64 15 12 15s4.53.82 6.24 2.19c.48.38.76.97.76 1.58V19z" fill="currentColor"/></svg>
            <div>
              <h2 class="card-title">实时状态</h2>
            </div>
          </div>
          <div class="stat-grid">
            <div class="stat-box">
              <span class="stat-num green">{{ stats.active_users }}</span>
              <span class="stat-num-label">在线人数</span>
            </div>
            <div class="stat-box">
              <span class="stat-num">{{ stats.total_completed }}</span>
              <span class="stat-num-label">累计完成</span>
            </div>
            <div class="stat-box">
              <span class="stat-num blue">{{ stats.running_tasks }}</span>
              <span class="stat-num-label">运行中任务</span>
            </div>
          </div>
        </section>

        <section class="card" v-if="runningTasks.length > 0">
          <div class="card-header">
            <svg viewBox="0 0 24 24" class="card-icon active-icon"><path d="M13 2.03v2.02c4.39.54 7.5 4.53 6.96 8.92-.46 3.64-3.32 6.53-6.96 6.96v2.02c5.49-.66 9.5-5.6 8.84-11.09-.55-4.58-4.09-8.19-8.84-8.83zM4.26 17.86c1.94 2.86 5.06 4.7 8.74 4.76v-2.02c-2.41-.12-4.53-1.44-5.74-3.28l-3 .54zm.03-11.72l3 .54c1.21-1.84 3.33-3.16 5.74-3.28V1.38c-3.68.06-6.8 1.9-8.74 4.76z" fill="currentColor"/></svg>
            <div>
              <h2 class="card-title">运行中任务</h2>
              <p class="card-desc">{{ runningTasks.length }} 个任务执行中</p>
            </div>
          </div>
          <div class="running-list">
            <div class="running-item" v-for="t in runningTasks" :key="t.id">
              <span class="running-dot"></span>
              <span class="running-name">{{ t.customer_name || '未命名' }}</span>
              <span class="running-file">{{ t.file_name }}</span>
              <button class="running-cancel" @click="cancelRunningTask(t.id)" title="终止">终止</button>
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

    // --- Modal orchestration ---
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

/* Cards */
.card { background: #161b22; border: 1px solid #21262d; border-radius: 12px; padding: 24px; margin-bottom: 20px; }
.card-header { display: flex; align-items: flex-start; gap: 12px; margin-bottom: 18px; }
.card-icon { width: 19px; height: 19px; color: #3fb950; flex-shrink: 0; margin-top: 2px; }
.card-icon.active-icon { color: #58a6ff; }
.card-title { font-size: 15px; font-weight: 600; color: #f0f6fc; margin-bottom: 2px; letter-spacing: -0.2px; }
.card-desc { font-size: 12.5px; color: #8b949e; }

/* Fields */
.field-group { margin-bottom: 14px; }
.field-label { display: block; font-size: 13px; font-weight: 500; color: #c9d1d9; margin-bottom: 6px; }
.field-input { width: 100%; padding: 9px 13px; background: #0d1117; border: 1px solid #30363d; border-radius: 8px; color: #f0f6fc; font-size: 13.5px; outline: none; transition: border-color 0.2s; }
.field-input:focus { border-color: #238636; box-shadow: 0 0 0 3px rgba(35,134,54,0.12); }
.field-input::placeholder { color: #484f58; }
.field-input.error { border-color: #da3633; }
.field-input.error:focus { box-shadow: 0 0 0 3px rgba(218,54,51,0.15); }
.required-star { color: #f85149; font-weight: 700; }
.field-error { display: block; font-size: 11.5px; color: #f85149; margin-top: 4px; }
.field-select { width: 100%; padding: 9px 13px; background: #0d1117; border: 1px solid #30363d; border-radius: 8px; color: #f0f6fc; font-size: 13.5px; outline: none; cursor: pointer; transition: border-color 0.2s; appearance: auto; }
.field-select:focus { border-color: #238636; box-shadow: 0 0 0 3px rgba(35,134,54,0.12); }
.field-select option { background: #161b22; color: #c9d1d9; }

/* Drop Zone */
.drop-zone { border: 2px dashed #30363d; border-radius: 10px; padding: 32px 20px; text-align: center; cursor: pointer; transition: all 0.2s; margin-bottom: 10px; }
.drop-zone:hover, .drop-zone.dragging { border-color: #238636; background: rgba(35,134,54,0.05); }
.drop-zone.filled { border-color: #238636; background: rgba(35,134,54,0.06); padding: 16px 20px; }
.upload-icon { width: 44px; height: 44px; margin: 0 auto 10px; color: #3fb950; }
.check-svg { width: 28px; height: 28px; flex-shrink: 0; }
.drop-text { font-size: 14px; color: #c9d1d9; margin-bottom: 4px; }
.drop-hint { font-size: 12px; color: #484f58; margin-bottom: 14px; }
.drop-result { display: flex; align-items: center; gap: 12px; justify-content: center; }

/* Buttons */
.btn { display: inline-flex; align-items: center; gap: 8px; padding: 9px 20px; border: none; border-radius: 8px; font-size: 14px; font-weight: 500; cursor: pointer; transition: all 0.15s; text-decoration: none; }
.btn:disabled { opacity: 0.4; cursor: not-allowed; }
.btn-icon { width: 16px; height: 16px; }
.btn-primary { background: #238636; color: #fff; }
.btn-primary:hover:not(:disabled) { background: #2ea043; }
.btn-accent { background: linear-gradient(135deg, #238636, #2ea043); color: #fff; }
.btn-accent:hover:not(:disabled) { background: linear-gradient(135deg, #1f7a2e, #238636); }
.btn-block { width: 100%; justify-content: center; padding: 11px; font-size: 15px; margin-top: 4px; }
.btn-outline { background: transparent; border: 1px solid #30363d; color: #c9d1d9; padding: 8px 18px; }
.btn-outline:hover { border-color: #238636; color: #3fb950; }
.btn-outline input { display: none; }
.upload-options { display: flex; gap: 10px; justify-content: center; }
.btn-ghost { background: transparent; color: #8b949e; padding: 9px 14px; }
.btn-ghost:hover { color: #c9d1d9; }
.action-row { display: flex; gap: 8px; margin-top: 12px; }

/* File List */
.file-list { display: flex; flex-direction: column; gap: 4px; }
.file-item { display: flex; align-items: center; gap: 8px; padding: 7px 10px; background: #0d1117; border-radius: 6px; font-size: 13px; }
.file-icon { width: 15px; height: 15px; color: #3fb950; flex-shrink: 0; }
.file-name { flex: 1; color: #c9d1d9; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.file-size { color: #484f58; flex-shrink: 0; }
.file-more { font-size: 12px; color: #484f58; padding: 4px 10px; }

/* Toggles */
.toggle-group { margin-bottom: 8px; }
.toggle-row { display: flex; align-items: center; justify-content: space-between; padding: 14px 0; cursor: pointer; border-bottom: 1px solid #21262d; }
.toggle-row:last-child { border-bottom: none; }
.toggle-info { flex: 1; }
.toggle-label { display: block; font-size: 14px; font-weight: 500; color: #c9d1d9; margin-bottom: 2px; }
.toggle-desc { font-size: 12px; color: #484f58; }
.toggle { width: 42px; height: 24px; background: #30363d; border-radius: 12px; position: relative; transition: background 0.2s; cursor: pointer; flex-shrink: 0; }
.toggle.on { background: #238636; }
.toggle-knob { width: 18px; height: 18px; background: #fff; border-radius: 50%; position: absolute; top: 3px; left: 3px; transition: transform 0.2s; box-shadow: 0 1px 3px rgba(0,0,0,0.15); }
.toggle.on .toggle-knob { transform: translateX(18px); }

/* Stats */
.stat-card { background: linear-gradient(135deg, #161b22 0%, #1a2129 100%); }
.stat-grid { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 8px; }
.stat-box { text-align: center; padding: 18px 8px; background: #0d1117; border-radius: 8px; border: 1px solid #21262d; }
.stat-num { display: block; font-size: 26px; font-weight: 700; color: #f0f6fc; line-height: 1.1; margin-bottom: 4px; }
.stat-num.green { color: #3fb950; }
.stat-num.blue { color: #58a6ff; }
.stat-num-label { font-size: 10.5px; color: #8b949e; text-transform: uppercase; letter-spacing: 0.3px; }

/* Running Tasks */
.running-list { display: flex; flex-direction: column; gap: 6px; }
.running-item { display: flex; align-items: center; gap: 8px; padding: 8px 10px; background: #0d1117; border-radius: 6px; font-size: 13px; }
.running-dot { width: 6px; height: 6px; background: #3fb950; border-radius: 50%; flex-shrink: 0; animation: pulse-dot 2s infinite; }
.running-name { font-weight: 500; color: #c9d1d9; }
.running-file { color: #484f58; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; flex: 1; }
.running-cancel { background: none; border: 1px solid #da3633; color: #f85149; padding: 3px 10px; border-radius: 4px; font-size: 12px; cursor: pointer; transition: all 0.15s; }
.running-cancel:hover { background: #da3633; color: #fff; }

@keyframes pulse-dot { 0%,100% { opacity: 1; } 50% { opacity: 0.3; } }

@media (max-width: 900px) { .two-col { grid-template-columns: 1fr; } }
</style>
