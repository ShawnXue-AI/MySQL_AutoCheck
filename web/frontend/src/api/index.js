import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 60000
})

function generateUUID() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    const r = Math.random() * 16 | 0
    const v = c === 'x' ? r : (r & 0x3 | 0x8)
    return v.toString(16)
  })
}

function getSessionId() {
  let sid = sessionStorage.getItem('session_id')
  if (!sid) {
    sid = generateUUID()
    sessionStorage.setItem('session_id', sid)
  }
  return sid
}

export function resetSession() {
  sessionStorage.setItem('session_id', generateUUID())
}

export function uploadFiles(files) {
  const form = new FormData()
  form.append('session_id', getSessionId())
  files.forEach(f => form.append('files', f))
  return api.post('/upload', form)
}

export function uploadArchive(file) {
  const form = new FormData()
  form.append('session_id', getSessionId())
  form.append('archive', file)
  return api.post('/upload/archive', form)
}

export function startAnalysis(customerName, options, fontName) {
  return api.post('/analyze', {
    session_id: getSessionId(),
    customer_name: customerName,
    options,
    font_name: fontName || '微软雅黑'
  })
}

export function getTaskStatus(taskId) {
  return api.get(`/tasks/${taskId}`, {
    params: { session_id: getSessionId() }
  })
}

export function cancelTask(taskId) {
  return api.post(`/tasks/${taskId}/cancel`)
}

export function deleteTask(taskId) {
  return api.delete(`/tasks/${taskId}`)
}

export function getStats() {
  return api.get('/stats', {
    params: { session_id: getSessionId() }
  })
}

export function getTasks() {
  return api.get('/tasks')
}

export function getRunningTasks() {
  return api.get('/tasks/running')
}

export function buildDownloadUrl(taskId, filename) {
  return `/api/download/${taskId}/${filename}`
}

export function buildZipDownloadUrl(taskId) {
  return `/api/download/zip/${taskId}`
}

export function adminAuth(password) {
  return api.post('/admin/auth', { password })
}

export function getAdminFiles() {
  return api.get('/admin/files')
}

export function adminCleanup(paths) {
  return api.post('/admin/cleanup', { paths })
}

// --- Person-Day APIs ---

export function createPersonDay(customerName, startTime, endTime, workContent) {
  return api.post('/persondays', {
    customer_name: customerName,
    start_time: startTime,
    end_time: endTime,
    work_content: workContent || ''
  })
}

export function getPersonDays() {
  return api.get('/persondays')
}

export function updatePersonDay(id, customerName, startTime, endTime, workContent) {
  return api.put(`/persondays/${id}`, {
    customer_name: customerName,
    start_time: startTime,
    end_time: endTime,
    work_content: workContent || ''
  })
}

export function deletePersonDay(id) {
  return api.delete(`/persondays/${id}`)
}

export function calculatePersonDay(startTime, endTime) {
  return api.post('/persondays/calculate', {
    start_time: startTime,
    end_time: endTime
  })
}

// --- Holiday APIs ---

export function getHolidays() {
  return api.get('/holidays')
}

export function addHoliday(date, name) {
  return api.post('/holidays', { date, name })
}

export function deleteHoliday(date) {
  return api.delete(`/holidays/${encodeURIComponent(date)}`)
}

export function fetchHolidaysAPI(year) {
  return api.post('/holidays/fetch', { year: year || new Date().getFullYear() })
}
