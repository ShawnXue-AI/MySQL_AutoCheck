<template>
  <div class="person-day-view">
    <div v-if="!authenticated" class="auth-card">
      <div class="auth-header">
        <svg viewBox="0 0 24 24" class="lock-icon"><path d="M18 8h-1V6c0-2.76-2.24-5-5-5S7 3.24 7 6v2H6c-1.1 0-2 .9-2 2v10c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V10c0-1.1-.9-2-2-2zm-6 9c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2zm3.1-9H8.9V6c0-1.71 1.39-3.1 3.1-3.1s3.1 1.39 3.1 3.1v2z" fill="currentColor"/></svg>
        <h2 class="auth-title">管理员认证</h2>
        <p class="auth-desc">请输入管理员密码以访问人天记录</p>
      </div>
      <input type="password" class="auth-input" v-model="password" placeholder="管理员密码" @keyup.enter="doAuth" />
      <p class="auth-error" v-if="authError">{{ authError }}</p>
      <button class="btn btn-primary btn-block" @click="doAuth" :disabled="authing">
        {{ authing ? '验证中...' : '验证' }}
      </button>
    </div>

    <div v-else class="person-day-content">
      <div class="panel-grid">
        <section class="card form-card">
          <div class="card-header">
            <svg viewBox="0 0 24 24" class="card-icon"><path d="M19 3h-4.18C14.4 1.84 13.3 1 12 1c-1.3 0-2.4.84-2.82 2H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-7 0c.55 0 1 .45 1 1s-.45 1-1 1-1-.45-1-1 .45-1 1-1zm2 14H7v-2h7v2zm3-4H7v-2h10v2zm0-4H7V7h10v2z" fill="currentColor"/></svg>
            <div>
              <h2 class="card-title">{{ editingId ? '编辑记录' : '新增人天记录' }}</h2>
              <p class="card-desc">输入客户名称和时间范围，自动计算人天</p>
            </div>
            <div class="card-actions" v-if="editingId">
              <button class="btn btn-ghost btn-sm" @click="resetForm">取消编辑</button>
            </div>
          </div>
          <div class="form-body">
            <div class="form-row">
              <label class="form-label">客户名称</label>
              <input type="text" class="form-input" v-model="form.customerName" placeholder="请输入客户名称" />
            </div>
            <div class="form-row">
              <label class="form-label">工作内容</label>
              <textarea class="form-input form-textarea" v-model="form.workContent" placeholder="请输入本次工作内容描述（可选）" rows="2"></textarea>
            </div>
            <div class="form-row form-row-2col">
              <div class="form-col">
                <label class="form-label">开始时间</label>
                <input type="datetime-local" class="form-input" v-model="form.startTime" @change="autoCalculate" />
              </div>
              <div class="form-col">
                <label class="form-label">结束时间</label>
                <input type="datetime-local" class="form-input" v-model="form.endTime" @change="autoCalculate" />
              </div>
            </div>

            <div class="form-row form-actions">
              <button class="btn btn-outline btn-sm" @click="handleCalculate" :disabled="!canCalculate">
                计算人天
              </button>
              <button class="btn btn-primary btn-sm" @click="handleSave" :disabled="!canSave">
                {{ editingId ? '更新' : '保存记录' }}
              </button>
            </div>

            <div v-if="calcResult" class="calc-result">
              <div class="result-header">
                <span class="result-label">人天计算结果</span>
                <span class="result-badge" :class="resultSeverity">{{ calcResult.person_days }} 人天</span>
              </div>
              <div class="result-stats">
                <div class="result-stat">
                  <span class="stat-num work">{{ calcResult.work_hours }}h</span>
                  <span class="stat-tag">工作时间 ×1.0</span>
                </div>
                <div class="result-stat">
                  <span class="stat-num overtime">{{ calcResult.overtime_hours }}h</span>
                  <span class="stat-tag">加班时间 ×1.5</span>
                </div>
                <div class="result-stat">
                  <span class="stat-num holiday">{{ calcResult.holiday_hours }}h</span>
                  <span class="stat-tag">节假日 ×2.0</span>
                </div>
              </div>
              <details class="detail-section" v-if="calcResult.detail && calcResult.detail.length">
                <summary class="detail-toggle">查看计算明细 ({{ calcResult.detail.length }} 条)</summary>
                <div class="detail-lines">
                  <div class="detail-line" v-for="(d, i) in calcResult.detail" :key="i"
                    :class="{ 'holiday-line': d.includes('节假日'), 'overtime-line': d.includes('非工作时间'), 'work-line': d.includes('工作时间') }">
                    {{ d }}
                  </div>
                </div>
              </details>
            </div>
          </div>
        </section>

        <section class="card holidays-card">
          <div class="card-header">
            <svg viewBox="0 0 24 24" class="card-icon holiday-icon"><path d="M19 3h-1V1h-2v2H8V1H6v2H5c-1.11 0-1.99.9-1.99 2L3 19c0 1.1.89 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm0 16H5V8h14v11zM7 10h5v5H7z" fill="currentColor"/></svg>
            <div>
              <h2 class="card-title">节假日管理</h2>
              <p class="card-desc">周末自动计入，节假日系数 ×2.0</p>
            </div>
          </div>
          <div class="holiday-sync">
            <button class="btn btn-outline btn-sm btn-block-wide" @click="handleFetchHolidays" :disabled="fetchingHolidays">
              <span v-if="fetchingHolidays" class="spinner"></span>
              {{ fetchingHolidays ? '同步中...' : '从国务院日历同步节假日' }}
            </button>
            <p class="sync-hint" v-if="syncMsg">{{ syncMsg }}</p>
          </div>
          <div class="holiday-add">
            <input type="date" class="form-input" v-model="holidayForm.date" placeholder="选择日期" />
            <input type="text" class="form-input" v-model="holidayForm.name" placeholder="节假日名称（可选）" />
            <button class="btn btn-outline btn-sm" @click="handleAddHoliday" :disabled="!holidayForm.date">添加</button>
          </div>
          <div class="holiday-list" v-if="holidays.length > 0">
            <div class="holiday-row" v-for="h in holidays" :key="h.date">
              <span class="h-date">{{ h.date }}</span>
              <span class="h-name">{{ h.name || '-' }}</span>
              <button class="btn btn-ghost btn-xs danger" @click="handleDeleteHoliday(h.date)">删除</button>
            </div>
          </div>
          <div class="empty-inline" v-else>暂无额外节假日配置，请点击同步</div>
        </section>
      </div>

      <section class="card records-card">
        <div class="card-header">
          <svg viewBox="0 0 24 24" class="card-icon"><path d="M11 17h2v-6h-2v6zm1-15C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zM11 9h2V7h-2v2z" fill="currentColor"/></svg>
          <div>
            <h2 class="card-title">人天记录列表</h2>
            <p class="card-desc">共 {{ records.length }} 条记录，合计 {{ totalPersonDays }} 人天</p>
          </div>
          <button class="btn btn-ghost btn-sm" @click="fetchRecords">刷新</button>
        </div>
        <div v-if="records.length === 0" class="empty-inline">暂无记录</div>
        <div v-else class="records-table">
          <div class="r-row r-header">
            <span class="r-customer h">客户名称</span>
            <span class="r-work h">工作内容</span>
            <span class="r-time h">开始时间</span>
            <span class="r-time h">结束时间</span>
            <span class="r-num h">人天数</span>
            <span class="r-num h">工作h</span>
            <span class="r-num h">加班h</span>
            <span class="r-num h">节假日h</span>
            <span class="r-actions h">操作</span>
          </div>
          <div class="r-scroll">
            <div class="r-row" v-for="r in records" :key="r.id">
              <span class="r-customer" :title="r.customer_name">{{ r.customer_name }}</span>
              <span class="r-work" :title="r.work_content">{{ r.work_content || '-' }}</span>
              <span class="r-time">{{ formatTime(r.start_time) }}</span>
              <span class="r-time">{{ formatTime(r.end_time) }}</span>
              <span class="r-num days">{{ r.person_days }}</span>
              <span class="r-num">{{ r.work_hours }}</span>
              <span class="r-num overtime">{{ r.overtime_hours }}</span>
              <span class="r-num holiday">{{ r.holiday_hours }}</span>
              <span class="r-actions">
                <button class="btn btn-ghost btn-xs" @click="editRecord(r)">编辑</button>
                <button class="btn btn-ghost btn-xs danger" @click="handleDelete(r.id)">删除</button>
              </span>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import { adminAuth, createPersonDay, getPersonDays, updatePersonDay, deletePersonDay, calculatePersonDay, getHolidays, addHoliday, deleteHoliday, fetchHolidaysAPI } from '../api/index.js'

export default {
  name: 'PersonDayView',
  data() {
    return {
      authenticated: false,
      password: '',
      authing: false,
      authError: '',
      form: {
        customerName: '',
        workContent: '',
        startTime: '',
        endTime: ''
      },
      editingId: null,
      calcResult: null,
      calculating: false,
      records: [],
      holidays: [],
      holidayForm: {
        date: '',
        name: ''
      },
      fetchingHolidays: false,
      syncMsg: ''
    }
  },
  computed: {
    canCalculate() {
      return this.form.customerName.trim() && this.form.startTime && this.form.endTime
    },
    canSave() {
      return this.canCalculate && this.calcResult
    },
    totalPersonDays() {
      return this.records.reduce((sum, r) => sum + (r.person_days || 0), 0).toFixed(2)
    },
    resultSeverity() {
      if (!this.calcResult) return ''
      const d = this.calcResult.person_days
      if (d >= 3) return 'high'
      if (d >= 1) return 'mid'
      return 'low'
    }
  },
  methods: {
    async doAuth() {
      this.authing = true; this.authError = ''
      try {
        const { data } = await adminAuth(this.password)
        if (data.success) {
          this.authenticated = true
          this.fetchRecords()
          this.fetchHolidays()
        } else {
          this.authError = data.error || '密码错误'
        }
      } catch (e) {
        this.authError = '请求失败，请检查后端服务'
      } finally { this.authing = false }
    },
    async fetchRecords() {
      try {
        const { data } = await getPersonDays()
        this.records = data.records || []
      } catch {}
    },
    async fetchHolidays() {
      try {
        const { data } = await getHolidays()
        this.holidays = data.holidays || []
      } catch {}
    },
    async autoCalculate() {
      if (!this.form.startTime || !this.form.endTime || !this.form.customerName.trim()) {
        this.calcResult = null
        return
      }
      await this.handleCalculate()
    },
    async handleCalculate() {
      if (!this.canCalculate) return
      this.calculating = true
      try {
        const { data } = await calculatePersonDay(
          this.form.startTime.replace('T', ' '),
          this.form.endTime.replace('T', ' ')
        )
        this.calcResult = data
      } catch (e) {
        this.calcResult = null
      } finally {
        this.calculating = false
      }
    },
    async handleSave() {
      if (!this.canSave) return
      try {
        const startTime = this.form.startTime.replace('T', ' ')
        const endTime = this.form.endTime.replace('T', ' ')
        if (this.editingId) {
          await updatePersonDay(this.editingId, this.form.customerName.trim(), startTime, endTime, this.form.workContent.trim())
        } else {
          await createPersonDay(this.form.customerName.trim(), startTime, endTime, this.form.workContent.trim())
        }
        this.resetForm()
        this.fetchRecords()
      } catch (e) {
        alert(e.response?.data?.error || '保存失败')
      }
    },
    async handleDelete(id) {
      if (!confirm('确认删除该记录？')) return
      try {
        await deletePersonDay(id)
        this.fetchRecords()
      } catch (e) {
        alert(e.response?.data?.error || '删除失败')
      }
    },
    editRecord(r) {
      this.editingId = r.id
      this.form.customerName = r.customer_name
      this.form.workContent = r.work_content || ''
      const toLocal = (t) => {
        if (!t) return ''
        const d = new Date(t)
        if (isNaN(d.getTime())) return ''
        const pad = n => String(n).padStart(2, '0')
        return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
      }
      this.form.startTime = toLocal(r.start_time)
      this.form.endTime = toLocal(r.end_time)
      this.handleCalculate()
      this.$el.querySelector('.form-card')?.scrollIntoView({ behavior: 'smooth' })
    },
    resetForm() {
      this.editingId = null
      this.form.customerName = ''
      this.form.workContent = ''
      this.form.startTime = ''
      this.form.endTime = ''
      this.calcResult = null
    },
    async handleAddHoliday() {
      if (!this.holidayForm.date) return
      try {
        await addHoliday(this.holidayForm.date, this.holidayForm.name.trim())
        this.holidayForm.date = ''
        this.holidayForm.name = ''
        this.fetchHolidays()
      } catch (e) {
        alert(e.response?.data?.error || '添加失败')
      }
    },
    async handleDeleteHoliday(date) {
      try {
        await deleteHoliday(date)
        this.fetchHolidays()
      } catch (e) {
        alert(e.response?.data?.error || '删除失败')
      }
    },
    async handleFetchHolidays() {
      this.fetchingHolidays = true
      this.syncMsg = ''
      try {
        const { data } = await fetchHolidaysAPI()
        this.holidays = data.holidays || []
        let msg = data.message || '同步完成'
        if (data.debug && data.debug.length > 0) {
          msg += '\n[调试] ' + data.debug.join(' | ')
        }
        if (data.errors && data.errors.length > 0) {
          msg += '\n[错误] ' + data.errors.join('; ')
        }
        this.syncMsg = msg
        setTimeout(() => { this.syncMsg = '' }, 15000)
      } catch (e) {
        this.syncMsg = '同步失败，请检查网络连接: ' + (e.message || e)
      } finally {
        this.fetchingHolidays = false
      }
    },
    formatTime(t) {
      if (!t) return '-'
      const d = new Date(t)
      if (isNaN(d.getTime())) return t
      const pad = n => String(n).padStart(2, '0')
      return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
    }
  }
}
</script>

<style scoped>
.person-day-view { max-width: 1200px; margin: 0 auto; }

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

.person-day-content { }

.panel-grid {
  display: grid;
  grid-template-columns: 1fr 320px;
  gap: 20px;
  margin-bottom: 20px;
}

.card { background: #161b22; border: 1px solid #21262d; border-radius: 12px; padding: 24px; }
.card-header { display: flex; align-items: flex-start; gap: 12px; margin-bottom: 18px; }
.card-icon { width: 19px; height: 19px; color: #3fb950; flex-shrink: 0; margin-top: 2px; }
.card-icon.holiday-icon { color: #d29922; }
.card-title { font-size: 15px; font-weight: 600; color: #f0f6fc; margin-bottom: 2px; }
.card-desc { font-size: 12.5px; color: #8b949e; }
.card-actions { display: flex; gap: 8px; margin-left: auto; flex-shrink: 0; }

.empty-inline { padding: 20px; text-align: center; font-size: 13px; color: #484f58; }

.form-body { display: flex; flex-direction: column; gap: 14px; }
.form-row { display: flex; flex-direction: column; gap: 5px; }
.form-row-2col { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.form-col { display: flex; flex-direction: column; gap: 5px; }
.form-label { font-size: 12.5px; color: #8b949e; font-weight: 500; }
.form-input {
  padding: 9px 12px; background: #0d1117; border: 1px solid #30363d;
  border-radius: 8px; color: #f0f6fc; font-size: 13.5px; outline: none;
  transition: border-color 0.2s;
}
.form-input:focus { border-color: #58a6ff; box-shadow: 0 0 0 3px rgba(88,166,255,0.15); }
.form-input::placeholder { color: #484f58; }
.form-textarea { resize: vertical; min-height: 50px; font-family: inherit; }

.form-actions { flex-direction: row; gap: 8px; margin-top: 4px; }

.btn { display: inline-flex; align-items: center; gap: 8px; padding: 9px 20px; border: none; border-radius: 8px; font-size: 14px; font-weight: 500; cursor: pointer; transition: all 0.15s; }
.btn:disabled { opacity: 0.4; cursor: not-allowed; }
.btn-primary { background: #238636; color: #fff; }
.btn-primary:hover:not(:disabled) { background: #2ea043; }
.btn-block { width: 100%; justify-content: center; padding: 11px; margin-top: 4px; }
.btn-block-wide { width: 100%; justify-content: center; }
.btn-outline { background: transparent; color: #58a6ff; border: 1px solid #30363d; }
.btn-outline:hover:not(:disabled) { border-color: #58a6ff; background: rgba(88,166,255,0.08); }
.btn-ghost { background: transparent; color: #8b949e; padding: 6px 12px; font-size: 13px; }
.btn-ghost:hover { color: #c9d1d9; }
.btn-sm { padding: 6px 16px; font-size: 13px; }
.btn-xs { padding: 4px 12px; font-size: 12px; flex-shrink: 0; }
.btn-xs.danger { color: #f85149; }
.btn-xs.danger:hover { color: #ff7b72; background: rgba(248,81,73,0.1); }

/* Calc Result */
.calc-result {
  margin-top: 6px; padding: 16px; background: #0d1117; border: 1px solid #21262d; border-radius: 10px;
}
.result-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px; }
.result-label { font-size: 13px; color: #8b949e; font-weight: 500; }
.result-badge { font-size: 22px; font-weight: 700; padding: 4px 14px; border-radius: 8px; }
.result-badge.high { background: rgba(35,134,54,0.15); color: #3fb950; }
.result-badge.mid { background: rgba(210,153,34,0.15); color: #d29922; }
.result-badge.low { background: rgba(88,166,255,0.15); color: #58a6ff; }

.result-stats { display: flex; gap: 12px; }
.result-stat { flex: 1; display: flex; flex-direction: column; align-items: center; gap: 4px; padding: 10px 8px; background: #161b22; border-radius: 8px; }
.stat-num { font-size: 18px; font-weight: 700; }
.stat-num.work { color: #3fb950; }
.stat-num.overtime { color: #d29922; }
.stat-num.holiday { color: #f85149; }
.stat-tag { font-size: 11px; color: #484f58; }

.detail-section { margin-top: 12px; }
.detail-toggle { font-size: 12.5px; color: #58a6ff; cursor: pointer; user-select: none; }
.detail-toggle:hover { color: #79c0ff; }
.detail-lines { margin-top: 8px; display: flex; flex-direction: column; gap: 2px; max-height: 200px; overflow-y: auto; }
.detail-line { font-size: 12px; padding: 4px 8px; border-radius: 4px; font-family: 'Consolas', 'Courier New', monospace; background: #161b22; color: #8b949e; }
.detail-line.holiday-line { color: #f85149; background: rgba(248,81,73,0.08); }
.detail-line.overtime-line { color: #d29922; background: rgba(210,153,34,0.08); }
.detail-line.work-line { color: #3fb950; background: rgba(35,134,54,0.08); }

/* Holidays */
.holiday-sync { margin-bottom: 12px; }
.sync-hint { font-size: 12px; color: #3fb950; text-align: center; margin-top: 6px; white-space: pre-line; word-break: break-all; }

.holiday-add { display: flex; gap: 8px; margin-bottom: 14px; }
.holiday-add .form-input { flex: 1; min-width: 0; }

.holiday-list { display: flex; flex-direction: column; gap: 3px; max-height: 260px; overflow-y: auto; }
.holiday-row { display: flex; align-items: center; gap: 10px; padding: 7px 10px; background: #0d1117; border-radius: 6px; font-size: 13px; }
.h-date { color: #f85149; font-weight: 500; font-family: 'Consolas', monospace; flex-shrink: 0; }
.h-name { color: #8b949e; flex: 1; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

/* Spinner */
.spinner {
  display: inline-block; width: 14px; height: 14px;
  border: 2px solid rgba(88,166,255,0.25); border-top-color: #58a6ff;
  border-radius: 50%; animation: spin 0.6s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }

/* Records Table */
.records-table { border: 1px solid #21262d; border-radius: 8px; overflow: hidden; }
.r-header { background: #0d1117; border-bottom: 1px solid #21262d; position: sticky; top: 0; z-index: 1; }
.r-scroll { max-height: 500px; overflow-y: auto; }
.r-row { display: flex; align-items: center; gap: 8px; padding: 9px 12px; font-size: 13px; }
.r-row:not(.r-header) { background: #0d1117; }
.r-row:not(.r-header):hover { background: #161b22; }

.r-customer { flex: 1.5; color: #c9d1d9; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.r-customer.h { color: #8b949e; font-weight: 600; font-size: 11.5px; text-transform: uppercase; letter-spacing: 0.3px; }
.r-work { flex: 2; color: #8b949e; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.r-work.h { color: #8b949e; font-weight: 600; font-size: 11.5px; text-transform: uppercase; letter-spacing: 0.3px; }
.r-time { flex: 1.3; color: #8b949e; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-family: 'Consolas', monospace; font-size: 12px; }
.r-time.h { color: #8b949e; font-weight: 600; font-size: 11.5px; text-transform: uppercase; letter-spacing: 0.3px; font-family: inherit; }
.r-num { width: 68px; color: #8b949e; text-align: right; flex-shrink: 0; }
.r-num.h { font-size: 11.5px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.3px; }
.r-num.days { color: #3fb950; font-weight: 600; }
.r-num.overtime { color: #d29922; }
.r-num.holiday { color: #f85149; }
.r-actions { width: 110px; display: flex; gap: 4px; justify-content: flex-end; flex-shrink: 0; }
.r-actions.h { font-size: 11.5px; color: #8b949e; font-weight: 600; text-transform: uppercase; letter-spacing: 0.3px; }
</style>
