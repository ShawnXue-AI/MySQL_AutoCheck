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
            <button v-if="uploadedFiles.length > 0" class="btn btn-ghost" @click="resetUpload">
              <svg viewBox="0 0 24 24" class="btn-icon" fill="none" stroke="currentColor" stroke-width="2"><polyline points="23 4 23 10 17 10"/><polyline points="1 20 1 10 7 10"/><path d="M3.51 9a9 9 0 0114.85-3.36L23 10"/></svg>
              重新选择
            </button>
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
          <button class="btn btn-accent btn-block" :disabled="!uploadComplete || busy" @click="startAnalyze">
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
import { ref, reactive, onBeforeUnmount } from 'vue'
import { uploadFiles, uploadArchive, startAnalysis, getStats, getTaskStatus, cancelTask, getRunningTasks, buildZipDownloadUrl, resetSession } from '../api/index.js'
import { toastSuccess, toastError } from '../components/Toast.vue'
import { formatSize } from '../composables/utils.js'
import ProcessModal from '../components/ProcessModal.vue'

export default {
  name: 'HomeView',
  components: { ProcessModal },
  setup(props, { emit }) {
    const isDragging = ref(false)
    const customerName = ref('')
    const showCustomerError = ref(false)
    const selectedFiles = ref([])
    const uploadedFiles = ref([])
    const uploadComplete = ref(false)
    const busy = ref(false)
    const modalVisible = ref(false)
    const modalPhase = ref('upload')
    const modalProgress = ref(0)
    const modalTitle = ref('')
    const modalSub = ref('')
    const modalError = ref('')
    const modalDownloadUrl = ref('')
    const modalCancelling = ref(false)
    const reportFont = ref('微软雅黑')
    const options = reactive({ generate_word: true, generate_excel: false })
    const currentTaskId = ref(null)
    let taskPollTimer = null
    const runningTasks = ref([])
    const stats = reactive({ active_users: 0, total_completed: 0, running_tasks: 0 })
    let statsTimer = null
    let runningTimer = null

    function pollStats() {
      getStats().then(({ data }) => { Object.assign(stats, data); emit('stats-update', data) }).catch(() => {})
    }
    function pollRunning() {
      getRunningTasks().then(({ data }) => { runningTasks.value = data.tasks || [] }).catch(() => {})
    }

    pollStats(); pollRunning()
    statsTimer = setInterval(() => pollStats(), 5000)
    runningTimer = setInterval(() => pollRunning(), 5000)

    onBeforeUnmount(() => {
      if (taskPollTimer) clearInterval(taskPollTimer)
      clearInterval(statsTimer)
      clearInterval(runningTimer)
    })

    function handleDrop(e) {
      isDragging.value = false
      const items = e.dataTransfer.items
      if (!items) { handleFiles(e.dataTransfer.files); return }
      const files = []
      const processEntry = (entry, path) => {
        return new Promise((resolve) => {
          if (entry.isFile) { entry.file(f => { f._fullPath = path; files.push(f); resolve() }) }
          else if (entry.isDirectory) { const reader = entry.createReader(); reader.readEntries(entries => { Promise.all(entries.map(e => processEntry(e, path + '/' + e.name))).then(resolve) }) }
          else { resolve() }
        })
      }
      Promise.all([...items].filter(i => i.webkitGetAsEntry).map(i => processEntry(i.webkitGetAsEntry(), i.webkitGetAsEntry().name))).then(() => { if (files.length) handleFiles(files) })
    }
    function handleFileSelect(e) { handleFiles(e.target.files) }
    function handleFiles(fileList) {
      const maxLogSize = 10 * 1024 * 1024, maxArchiveSize = 30 * 1024 * 1024
      const files = Array.from(fileList).filter(f => {
        if (!f.name.endsWith('.log') && !f.name.endsWith('.zip') && !f.name.endsWith('.tar.gz') && !f.name.endsWith('.tgz') && !(f.size > 0 && f._fullPath)) return false
        const isArchive = f.name.endsWith('.zip') || f.name.endsWith('.tar.gz') || f.name.endsWith('.tgz')
        if (isArchive && f.size > maxArchiveSize) { toastError(`压缩包 ${f.name} 超过30MB限制`); return false }
        if (!isArchive && f.size > maxLogSize) { toastError(`日志文件 ${f.name} 超过10MB限制`); return false }
        return true
      })
      if (!files.length) return
      resetSession()
      selectedFiles.value = files; uploadedFiles.value = files; uploadComplete.value = false
    }
    function resetUpload() { uploadedFiles.value = []; selectedFiles.value = []; uploadComplete.value = false }

    async function submitUpload() {
      if (!customerName.value.trim()) { showCustomerError.value = true; toastError('请输入客户名称'); return }
      showCustomerError.value = false; busy.value = true
      modalVisible.value = true; modalPhase.value = 'upload'; modalProgress.value = 0
      modalTitle.value = '正在上传文件'; modalSub.value = ''; modalError.value = ''; modalCancelling.value = false; currentTaskId.value = null
      try {
        modalProgress.value = 15
        const af = selectedFiles.value.find(f => f.name.endsWith('.zip') || f.name.endsWith('.tar.gz') || f.name.endsWith('.tgz'))
        modalProgress.value = 40
        await (af ? uploadArchive(af) : uploadFiles(selectedFiles.value))
        modalProgress.value = 100; modalTitle.value = '上传完成'; modalSub.value = '即将返回主页面...'
        uploadComplete.value = true; busy.value = false
        toastSuccess('上传成功')
        setTimeout(() => { modalVisible.value = false; modalProgress.value = 0 }, 1200)
      } catch (e) {
        modalPhase.value = 'error'; modalTitle.value = '上传失败'; modalError.value = e.response?.data?.error || e.message
        uploadComplete.value = false; busy.value = false
      }
    }
    async function startAnalyze() {
      busy.value = true; modalVisible.value = true
      modalPhase.value = 'running'; modalProgress.value = 10; modalTitle.value = '正在提交分析'; modalSub.value = ''; modalError.value = ''; modalCancelling.value = false; modalDownloadUrl.value = ''
      const name = customerName.value.trim() || '未命名客户'
      try {
        const { data } = await startAnalysis(name, options, reportFont.value)
        currentTaskId.value = data.task_id
        modalTitle.value = '正在生成报告'; modalSub.value = name; modalProgress.value = 20
        taskPollTimer = setInterval(() => pollTask(data.task_id), 2000)
      } catch (e) {
        modalPhase.value = 'error'; modalTitle.value = '启动失败'; modalError.value = e.response?.data?.error || e.message
        busy.value = false
      }
    }
    async function pollTask(taskId) {
      try {
        const { data } = await getTaskStatus(taskId)
        Object.assign(stats, data.stats); emit('stats-update', data.stats)
        const t = data.task
        modalProgress.value = { pending: 20, running: 60, completed: 100, failed: 100, cancelled: 100 }[t.status] || 20
        if (t.progress) modalSub.value = t.progress
        if (t.status === 'completed') {
          clearInterval(taskPollTimer)
          modalPhase.value = 'done'; modalTitle.value = '报告生成完成'; modalSub.value = t.zip_file || ''; modalProgress.value = 100
          modalDownloadUrl.value = buildZipDownloadUrl(t.id)
          busy.value = false
          toastSuccess('报告生成完成')
        } else if (t.status === 'failed') {
          clearInterval(taskPollTimer)
          modalPhase.value = 'error'; modalTitle.value = '分析失败'; modalError.value = t.error || '未知错误'
          busy.value = false
        } else if (t.status === 'cancelled') {
          clearInterval(taskPollTimer)
          handleModalClose()
        }
      } catch {}
    }
    async function handleModalCancel() {
      if (modalPhase.value === 'upload') {
        modalCancelling.value = true
        await new Promise(r => setTimeout(r, 800))
        modalVisible.value = false
        modalCancelling.value = false
        return
      }
      if (modalPhase.value === 'running') {
        modalCancelling.value = true
        if (currentTaskId.value) {
          try { await cancelTask(currentTaskId.value) } catch {}
          await new Promise(r => setTimeout(r, 1500))
        }
        modalVisible.value = false
        toastSuccess('任务已终止，文件已清理')
        resetAll()
      }
    }
    function handleModalClose() {
      modalVisible.value = false
      clearInterval(taskPollTimer)
      busy.value = false; modalProgress.value = 0; modalCancelling.value = false; currentTaskId.value = null
    }
    function resetAll() {
      clearInterval(taskPollTimer)
      busy.value = false; uploadComplete.value = false; uploadedFiles.value = []; selectedFiles.value = []
      modalVisible.value = false; modalProgress.value = 0; modalDownloadUrl.value = ''; modalCancelling.value = false
      currentTaskId.value = null
    }
    async function cancelRunningTask(taskId) {
      try { await cancelTask(taskId); toastSuccess('任务已终止') } catch { toastError('取消失败') }
    }

    return {
      isDragging, customerName, showCustomerError, selectedFiles, uploadedFiles, uploadComplete,
      busy, modalVisible, modalPhase, modalProgress, modalTitle, modalSub, modalError,
      modalDownloadUrl, modalCancelling, reportFont, options, currentTaskId,
      runningTasks, stats,
      handleDrop, handleFileSelect, resetUpload, formatSize,
      submitUpload, startAnalyze, handleModalCancel, handleModalClose, cancelRunningTask
    }
  }
}
</script>

<style scoped>
.home-view { max-width: 1160px; margin: 0 auto; }
.two-col { display: grid; grid-template-columns: 1fr 350px; gap: 24px; align-items: start; }

.field-select { width: 100%; padding: 9px 13px; background: transparent; border: 1px solid var(--border); border-radius: var(--radius); color: var(--foreground); font-size: 13.5px; outline: none; cursor: pointer; transition: border-color 0.2s; appearance: auto; }
.field-select:focus { border-color: var(--primary); box-shadow: 0 0 0 3px rgba(59,130,246,0.15); }
.field-select option { background: var(--card); color: var(--foreground); }

/* Drop Zone */
.drop-zone { border: 2px dashed var(--border); border-radius: var(--radius); padding: 32px 20px; text-align: center; cursor: pointer; transition: all 0.2s; margin-bottom: 10px; }
.drop-zone:hover, .drop-zone.dragging { border-color: var(--primary); background: rgba(59,130,246,0.05); }
.drop-zone.filled { border-color: #22c55e; background: rgba(34,197,94,0.05); padding: 16px 20px; }
.upload-icon { width: 44px; height: 44px; margin: 0 auto 10px; color: var(--primary); }
.check-svg { width: 28px; height: 28px; flex-shrink: 0; }
.drop-text { font-size: 14px; color: var(--foreground); margin-bottom: 4px; }
.drop-hint { font-size: 12px; color: var(--muted-foreground); margin-bottom: 14px; }
.drop-result { display: flex; align-items: center; gap: 12px; justify-content: center; }

.upload-options { display: flex; gap: 10px; justify-content: center; }
.upload-options label { position: relative; }
.upload-options input[type="file"] { position: absolute; inset: 0; opacity: 0; cursor: pointer; }

.btn-accent { background: linear-gradient(135deg, #2563eb, #3b82f6); color: #fff; }
.btn-accent:hover:not(:disabled) { background: linear-gradient(135deg, #1d4ed8, #2563eb); }
.action-row { display: flex; gap: 8px; margin-top: 12px; }

/* File List */
.file-list { display: flex; flex-direction: column; gap: 4px; }
.file-item { display: flex; align-items: center; gap: 8px; padding: 7px 10px; background: rgba(24,24,27,0.25); border-radius: var(--radius); font-size: 13px; border-bottom: 1px solid rgba(39,39,42,0.3); }
.file-item:last-child { border-bottom: none; }
.file-icon { width: 15px; height: 15px; color: var(--primary); flex-shrink: 0; }
.file-name { flex: 1; color: var(--foreground); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.file-size { color: var(--muted-foreground); flex-shrink: 0; }
.file-more { font-size: 12px; color: var(--muted-foreground); padding: 4px 10px; }

/* Toggles */
.toggle-group { margin-bottom: 8px; }
.toggle-row { display: flex; align-items: center; justify-content: space-between; padding: 14px 0; cursor: pointer; border-bottom: 1px solid var(--border); }
.toggle-row:last-child { border-bottom: none; }
.toggle-info { flex: 1; }
.toggle-label { display: block; font-size: 14px; font-weight: 500; color: var(--foreground); margin-bottom: 2px; }
.toggle-desc { font-size: 12px; color: var(--muted-foreground); }
.toggle { width: 42px; height: 24px; background: var(--border); border-radius: 12px; position: relative; transition: background 0.2s; cursor: pointer; flex-shrink: 0; }
.toggle.on { background: var(--primary); }
.toggle-knob { width: 18px; height: 18px; background: #fff; border-radius: 50%; position: absolute; top: 3px; left: 3px; transition: transform 0.2s; box-shadow: 0 1px 3px rgba(0,0,0,0.15); }
.toggle.on .toggle-knob { transform: translateX(18px); }

/* Stats */
.stat-card { background: linear-gradient(135deg, var(--card) 0%, rgba(24,24,27,0.8) 100%); }
.stat-grid { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 8px; }
.stat-box { text-align: center; padding: 18px 8px; background: rgba(17,17,19,0.35); border-radius: var(--radius); border: 1px solid var(--border); }
.stat-num { display: block; font-size: 26px; font-weight: 700; color: var(--foreground); line-height: 1.1; margin-bottom: 4px; }
.stat-num.green { color: #22c55e; }
.stat-num.blue { color: var(--primary); }
.stat-num-label { font-size: 10.5px; color: var(--muted-foreground); text-transform: uppercase; letter-spacing: 0.3px; }

/* Running Tasks */
.running-list { display: flex; flex-direction: column; gap: 6px; }
.running-item { display: flex; align-items: center; gap: 8px; padding: 8px 10px; background: rgba(24,24,27,0.25); border-radius: var(--radius); font-size: 13px; border-bottom: 1px solid rgba(39,39,42,0.3); }
.running-item:last-child { border-bottom: none; }
.running-dot { width: 6px; height: 6px; background: #22c55e; border-radius: 50%; flex-shrink: 0; animation: pulse-dot 2s infinite; }
.running-name { font-weight: 500; color: var(--foreground); }
.running-file { color: var(--muted-foreground); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; flex: 1; }
.running-cancel { background: none; border: 1px solid var(--destructive); color: var(--destructive); padding: 3px 10px; border-radius: 4px; font-size: 12px; cursor: pointer; transition: all 0.15s; }
.running-cancel:hover { background: var(--destructive); color: #fff; }

@keyframes pulse-dot { 0%,100% { opacity: 1; } 50% { opacity: 0.3; } }

@media (max-width: 900px) { .two-col { grid-template-columns: 1fr; } }
</style>
