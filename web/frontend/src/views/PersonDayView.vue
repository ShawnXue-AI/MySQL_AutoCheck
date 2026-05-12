<template>
  <div v-if="!authenticated" class="auth-card">
    <div class="auth-header">
      <svg viewBox="0 0 24 24" class="lock-icon" v-once><path d="M18 8h-1V6c0-2.76-2.24-5-5-5S7 3.24 7 6v2H6c-1.1 0-2 .9-2 2v10c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V10c0-1.1-.9-2-2-2zm-6 9c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2zm3.1-9H8.9V6c0-1.71 1.39-3.1 3.1-3.1s3.1 1.39 3.1 3.1v2z" fill="currentColor"/></svg>
      <h2 class="auth-title">管理员认证</h2>
      <p class="auth-desc">请输入管理员密码以访问人天记录</p>
    </div>
    <input type="password" class="auth-input" v-model="password" placeholder="管理员密码" @keyup.enter="doAuth" />
    <p class="auth-error" v-if="authError">{{ authError }}</p>
    <button class="btn btn-primary btn-block" @click="doAuth" :disabled="authing">
      <svg viewBox="0 0 24 24" class="btn-icon" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><rect width="18" height="11" x="3" y="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
      {{ authing ? '验证中...' : '验证' }}
    </button>
  </div>

  <div v-else class="pdv-main">
    <div class="pdv-top">
      <section class="card form-card">
        <div class="card-header">
          <svg viewBox="0 0 24 24" class="card-icon" v-once fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5"/><path d="M18.37 2.63a2.12 2.12 0 013 3L13 14l-4 1 1-4 8.37-8.37z"/></svg>
          <div>
            <h2 class="card-title">{{ editingId ? '编辑记录' : '新增人天记录' }}</h2>
            <p class="card-desc">// 输入客户名称和时间范围，自动计算人天</p>
          </div>
          <div class="card-actions">
            <button v-if="editingId" class="btn btn-outline btn-sm" style="color: var(--accent-red); border-color: rgba(248,81,73,0.3);" @click="resetForm">
              <svg viewBox="0 0 24 24" class="btn-icon" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 6L6 18M6 6l12 12"/></svg>
              取消编辑
            </button>
          </div>
        </div>
        <div class="form-body">
          <div class="pdv-form-grid">
            <div class="field-row">
              <label class="field-label">客户名称</label>
              <input type="text" class="field-input" v-model="form.customerName" placeholder="请输入客户名称" />
            </div>
            <div class="field-row field-row-span">
              <label class="field-label">工作内容</label>
              <textarea class="field-input field-textarea" v-model="form.workContent" placeholder="请输入本次工作内容描述（可选）" rows="2"></textarea>
            </div>
            <div class="field-row">
              <label class="field-label">开始时间</label>
              <input type="datetime-local" class="field-input field-dt" v-model="form.startTime" @change="autoCalc" />
            </div>
            <div class="field-row">
              <label class="field-label">结束时间</label>
              <input type="datetime-local" class="field-input field-dt" v-model="form.endTime" @change="autoCalc" />
            </div>
          </div>
          <div class="form-actions">
            <button class="btn btn-outline btn-sm" @click="doCalc" :disabled="!canCalc">
              <svg viewBox="0 0 24 24" class="btn-icon" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M12 6v6l4 2"/></svg>
              计算人天
            </button>
            <button class="btn btn-primary btn-sm" @click="doSave" :disabled="!canSave">
              <svg viewBox="0 0 24 24" class="btn-icon" fill="none" stroke="currentColor" stroke-width="2"><path d="M19 21H5a2 2 0 01-2-2V5a2 2 0 012-2h11l5 5v11a2 2 0 01-2 2z"/><polyline points="17 21 17 13 7 13 7 21"/><polyline points="7 3 7 8 15 8"/></svg>
              {{ editingId ? '更新' : '保存记录' }}
            </button>
          </div>
        </div>
      </section>

      <section class="card">
        <div class="card-header">
          <svg viewBox="0 0 24 24" class="card-icon" style="color: var(--accent-amber)" v-once fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="18" x="3" y="4" rx="2"/><path d="M16 2v4M8 2v4M3 10h18"/></svg>
          <div>
            <h2 class="card-title">节假日管理</h2>
            <p class="card-desc">// 周末自动计入，节假日系数 x2.0</p>
          </div>
        </div>
        <button class="btn btn-outline btn-sm btn-block-wide" @click="fetchHolidaysNow" :disabled="fetchingHolidays">
          <svg viewBox="0 0 24 24" class="btn-icon" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" v-if="!fetchingHolidays"><path d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.66 0 3-4.03 3-9s-1.34-9-3-9m0 18c-1.66 0-3-4.03-3-9s1.34-9 3-9m-9 9a9 9 0 019-9"/></svg>
          <span class="pdv-spin" v-else></span>
          {{ fetchingHolidays ? '同步中...' : '从国务院日历同步节假日' }}
        </button>
        <p class="pdv-sync-msg" v-if="syncMsg">{{ syncMsg }}</p>
        <div class="pdv-h-add">
          <input type="date" class="field-input" v-model="holidayDate" />
          <input type="text" class="field-input" v-model="holidayName" placeholder="名称（可选）" />
          <button class="btn btn-outline btn-sm" @click="addHoliday" :disabled="!holidayDate">
            <svg viewBox="0 0 24 24" class="btn-icon" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 5v14M5 12h14"/></svg>
            添加
          </button>
        </div>
        <div class="pdv-h-list" v-if="holidays.length">
          <div class="pdv-h-row" v-for="h in holidays" :key="h.date">
            <span class="pdv-h-date">{{ h.date }}</span>
            <span class="pdv-h-name">{{ h.name || '-' }}</span>
            <button class="btn btn-ghost btn-xs" style="color: var(--accent-red)" @click="delHoliday(h.date)">
              <svg viewBox="0 0 24 24" class="btn-icon" style="width:13px;height:13px" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 6h18M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/></svg>
              删除
            </button>
          </div>
        </div>
        <div class="empty-inline" v-else style="padding:14px">暂无节假日，请点击同步</div>
      </section>

      <section v-if="calcResult" class="card calc-card">
        <CalcResult :result="calcResult" :severity="resultSeverity" />
      </section>
    </div>

    <section class="card">
      <div class="card-header">
        <svg viewBox="0 0 24 24" class="card-icon" v-once fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>
        <div>
          <h2 class="card-title">人天记录列表</h2>
          <p class="card-desc">$ find records  →  {{ records.length }} 条，合计 {{ totalPersonDays }} 人天</p>
        </div>
        <div class="card-actions">
          <button class="btn btn-outline btn-sm" @click="fetchRecords">
            <svg viewBox="0 0 24 24" class="btn-icon" fill="none" stroke="currentColor" stroke-width="2"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/></svg>
            刷新
          </button>
        </div>
      </div>
      <div v-if="records.length === 0" class="empty-inline">暂无记录</div>
      <div v-else class="pdv-wrap">
        <table class="pdv-tbl">
          <thead><tr>
            <th class="col-name">客户名称</th>
            <th class="col-work">工作内容</th>
            <th class="col-time">开始时间</th>
            <th class="col-time">结束时间</th>
            <th class="col-num">人天数</th>
            <th class="col-num">工作h</th>
            <th class="col-num">加班h</th>
            <th class="col-num">假日h</th>
            <th class="col-act">操作</th>
          </tr></thead>
          <tbody>
            <tr v-for="r in records" :key="r.id">
              <td class="col-name" :title="r.customer_name">{{ r.customer_name }}</td>
              <td class="col-work" :title="r.work_content">{{ r.work_content || '-' }}</td>
              <td class="col-time">{{ fmt(r.start_time) }}</td>
              <td class="col-time">{{ fmt(r.end_time) }}</td>
              <td class="col-num days">{{ r.person_days }}</td>
              <td class="col-num">{{ r.work_hours }}</td>
              <td class="col-num ot">{{ r.overtime_hours }}</td>
              <td class="col-num hl">{{ r.holiday_hours }}</td>
              <td class="col-act">
                <button class="btn btn-ghost btn-xs" @click="editRecord(r)">
                  <svg viewBox="0 0 24 24" class="btn-icon" style="width:13px;height:13px" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
                  编辑
                </button>
                <button class="btn btn-ghost btn-xs" style="color: var(--accent-red)" @click="delRecord(r.id)">
                  <svg viewBox="0 0 24 24" class="btn-icon" style="width:13px;height:13px" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 6h18M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/></svg>
                  删除
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </div>
</template>

<script>
import { ref, reactive, computed } from 'vue'
import { useAuth } from '../composables/useAuth.js'
import { createPersonDay, getPersonDays, updatePersonDay, deletePersonDay, calculatePersonDay, getHolidays, addHoliday, deleteHoliday, fetchHolidaysAPI } from '../api/index.js'
import CalcResult from '../components/CalcResult.vue'

const fmtDateTime = (t) => {
  if (!t) return '-'
  const d = new Date(t)
  if (isNaN(d.getTime())) return t
  const p = n => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}`
}
const toLocalDt = (t) => {
  if (!t) return ''
  const d = new Date(t)
  if (isNaN(d.getTime())) return ''
  const p = n => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())}T${p(d.getHours())}:${p(d.getMinutes())}`
}

export default {
  name: 'PersonDayView',
  components: { CalcResult },
  setup() {
    const { authenticated, password, authing, authError, doAuth } = useAuth()
    const editingId = ref(null)
    const records = ref([])
    const holidays = ref([])
    const calcResult = ref(null)
    const calcBusy = ref(false)
    const fetchingHolidays = ref(false)
    const syncMsg = ref('')
    const holidayDate = ref('')
    const holidayName = ref('')
    const form = reactive({ customerName: '', workContent: '', startTime: '', endTime: '' })

    const canCalc = computed(() => form.customerName.trim() && form.startTime && form.endTime)
    const canSave = computed(() => canCalc.value && calcResult.value)
    const totalPersonDays = computed(() => records.value.reduce((s, r) => s + (r.person_days || 0), 0).toFixed(2))
    const resultSeverity = computed(() => {
      if (!calcResult.value) return ''
      const d = calcResult.value.person_days
      return d >= 3 ? 'high' : d >= 1 ? 'mid' : 'low'
    })

    async function fetchRecords() {
      try { const { data } = await getPersonDays(); records.value = data.records || [] } catch {}
    }
    async function fetchHolidaysList() {
      try { const { data } = await getHolidays(); holidays.value = data.holidays || [] } catch {}
    }

    async function doCalc() {
      if (!canCalc.value) return
      calcBusy.value = true
      try {
        const { data } = await calculatePersonDay(form.startTime.replace('T', ' '), form.endTime.replace('T', ' '))
        calcResult.value = data
      } catch { calcResult.value = null }
      finally { calcBusy.value = false }
    }
    function autoCalc() {
      if (!form.startTime || !form.endTime || !form.customerName.trim()) { calcResult.value = null; return }
      doCalc()
    }

    async function doSave() {
      if (!canSave.value) return
      const st = form.startTime.replace('T', ' '), et = form.endTime.replace('T', ' ')
      try {
        if (editingId.value) await updatePersonDay(editingId.value, form.customerName.trim(), st, et, form.workContent.trim())
        else await createPersonDay(form.customerName.trim(), st, et, form.workContent.trim())
        resetForm()
        fetchRecords()
      } catch (e) { alert(e.response?.data?.error || '保存失败') }
    }

    function editRecord(r) {
      editingId.value = r.id
      form.customerName = r.customer_name
      form.workContent = r.work_content || ''
      form.startTime = toLocalDt(r.start_time)
      form.endTime = toLocalDt(r.end_time)
      doCalc()
    }
    function resetForm() {
      editingId.value = null
      Object.assign(form, { customerName: '', workContent: '', startTime: '', endTime: '' })
      calcResult.value = null
    }

    async function delRecord(id) {
      if (!confirm('确认删除该记录？')) return
      try { await deletePersonDay(id); fetchRecords() } catch (e) { alert(e.response?.data?.error || '删除失败') }
    }

    async function addHoliday() {
      if (!holidayDate.value) return
      try { await addHoliday(holidayDate.value, holidayName.value.trim()); holidayDate.value = ''; holidayName.value = ''; fetchHolidaysList() }
      catch (e) { alert(e.response?.data?.error || '添加失败') }
    }
    async function delHoliday(date) {
      try { await deleteHoliday(date); fetchHolidaysList() } catch (e) { alert(e.response?.data?.error || '删除失败') }
    }
    async function fetchHolidaysNow() {
      fetchingHolidays.value = true; syncMsg.value = ''
      try {
        const { data } = await fetchHolidaysAPI()
        holidays.value = data.holidays || []
        let msg = data.message || '同步完成'
        if (data.debug?.length) msg += '\n[调试] ' + data.debug.join(' | ')
        if (data.errors?.length) msg += '\n[错误] ' + data.errors.join('; ')
        syncMsg.value = msg
        setTimeout(() => { syncMsg.value = '' }, 12000)
      } catch (e) { syncMsg.value = '同步失败: ' + (e.message || e) }
      finally { fetchingHolidays.value = false }
    }

    fetchRecords()
    fetchHolidaysList()

    return {
      authenticated, password, authing, authError, doAuth,
      editingId, records, holidays, calcResult, calcBusy,
      fetchingHolidays, syncMsg, holidayDate, holidayName, form,
      canCalc, canSave, totalPersonDays, resultSeverity,
      fetchRecords, doCalc, autoCalc, doSave,
      editRecord, resetForm, delRecord,
      addHoliday, delHoliday, fetchHolidaysNow,
      fmt: fmtDateTime
    }
  }
}
</script>

<style scoped>
.pdv-main { max-width: 1200px; margin: 0 auto; display: flex; flex-direction: column; gap: 16px; }

.pdv-top {
  display: grid; grid-template-columns: 1fr 340px; gap: 16px; align-items: stretch;
}
.calc-card { grid-column: 1 / -1; }

.form-card { display: flex; flex-direction: column; }
.form-body { flex: 1; display: flex; flex-direction: column; gap: 14px; }
.pdv-form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px 14px; }
.field-row { display: flex; flex-direction: column; gap: 5px; }
.field-row-span { grid-column: 1 / -1; }
.field-label { font-size: 12.5px; color: var(--text-tertiary); font-weight: 500; }
.field-input {
  padding: 9px 12px; background: var(--bg-input); border: 1px solid var(--border-hover);
  border-radius: var(--radius-sm); color: var(--text-primary); font-size: 13.5px; outline: none;
  transition: border-color var(--transition-normal); width: 100%; font-family: var(--font-body);
}
.field-input:focus { border-color: var(--accent-blue); box-shadow: var(--shadow-glow-blue); }
.field-input::placeholder { color: var(--text-muted); }
.field-textarea { resize: vertical; min-height: 54px; }
.field-dt {
  font-family: var(--font-code); font-size: 13px; letter-spacing: 0.3px;
}
.form-actions { display: flex; gap: 8px; justify-content: center; margin-top: auto; padding: 16px 0 4px; }

.btn-block-wide { width: 100%; justify-content: center; }

.pdv-sync-msg {
  font-family: var(--font-code);
  font-size: 12px; color: var(--accent-green); text-align: center;
  margin-bottom: 10px; white-space: pre-line; word-break: break-all;
}
.pdv-h-add { display: flex; gap: 8px; margin: 12px 0; align-items: center; }
.pdv-h-add .field-input { flex: 1; min-width: 0; }
.pdv-h-add .field-input[type="date"] { font-family: var(--font-code); font-size: 13px; }

.pdv-h-list {
  max-height: 220px; overflow-y: auto;
  background: var(--bg-primary); border-radius: var(--radius-md);
  border: 1px solid var(--border-default);
}
.pdv-h-row {
  display: flex; align-items: center; gap: 10px;
  padding: 7px 10px; border-bottom: 1px solid var(--border-default);
  font-size: 13px;
}
.pdv-h-row:last-child { border-bottom: none; }
.pdv-h-date {
  color: var(--accent-red); font-weight: 500;
  font-family: var(--font-code); font-size: 12.5px; flex-shrink: 0;
}
.pdv-h-name {
  color: var(--text-tertiary); flex: 1; min-width: 0;
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}

.pdv-spin {
  display: inline-block; width: 14px; height: 14px;
  border: 2px solid var(--accent-blue-subtle);
  border-top-color: var(--accent-blue);
  border-radius: 50%; animation: spin 0.6s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }

.pdv-wrap {
  max-height: 520px; overflow-y: auto;
  border: 1px solid var(--border-default); border-radius: var(--radius-md);
}

.pdv-tbl { width: 100%; border-collapse: collapse; font-size: 13px; table-layout: fixed; }
.pdv-tbl thead { position: sticky; top: 0; z-index: 1; }
.pdv-tbl th {
  padding: 10px 8px; text-align: center; font-weight: 600;
  font-family: var(--font-code); font-size: 10.5px;
  color: var(--text-tertiary); text-transform: uppercase; letter-spacing: 0.5px;
  background: var(--bg-tertiary); border-bottom: 1px solid var(--border-default);
}
.pdv-tbl td {
  padding: 9px 8px; text-align: center;
  border-bottom: 1px solid var(--border-default);
  color: var(--text-tertiary);
}
.pdv-tbl tbody tr { background: var(--bg-primary); transition: background var(--transition-fast); }
.pdv-tbl tbody tr:hover { background: var(--bg-hover); }
.pdv-tbl tbody tr:last-child td { border-bottom: none; }

.col-name { width: 12%; text-align: left !important; color: var(--text-primary) !important; font-weight: 500; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.col-work { width: 17%; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.col-time {
  width: 12%; font-family: var(--font-code); font-size: 11.5px; letter-spacing: 0.2px;
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}
.col-num { width: 7%; }
.col-act { width: 12%; }
td.col-num { font-family: var(--font-code); font-variant-numeric: tabular-nums; font-size: 12px; }
td.days { color: var(--accent-green) !important; font-weight: 600; font-size: 13.5px; }
td.ot { color: var(--accent-amber); }
td.hl { color: var(--accent-red); }

/* Auth card (reuse global .auth-card styles, add specifics) */
.empty-inline { padding: 28px; text-align: center; font-size: 13px; color: var(--text-muted); font-family: var(--font-code); }
</style>
