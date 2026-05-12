<template>
  <div class="calc-result">
    <div class="result-header">
      <div class="result-label-row">
        <svg viewBox="0 0 24 24" class="result-head-icon" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
        <span class="result-label">人天计算结果</span>
      </div>
      <span class="result-badge" :class="severity">{{ result.person_days }} 人天</span>
    </div>
    <div class="result-stats">
      <div class="result-stat">
        <svg viewBox="0 0 24 24" class="stat-icon work" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="5"/><path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/></svg>
        <span class="stat-num work">{{ result.work_hours }}h</span>
        <span class="stat-tag">工作时间</span>
        <span class="stat-mul">×1.0</span>
      </div>
      <div class="result-stat">
        <svg viewBox="0 0 24 24" class="stat-icon overtime" fill="none" stroke="currentColor" stroke-width="2"><path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/></svg>
        <span class="stat-num overtime">{{ result.overtime_hours }}h</span>
        <span class="stat-tag">加班时间</span>
        <span class="stat-mul">×1.5</span>
      </div>
      <div class="result-stat">
        <svg viewBox="0 0 24 24" class="stat-icon holiday" fill="none" stroke="currentColor" stroke-width="2"><rect width="18" height="18" x="3" y="4" rx="2"/><path d="M16 2v4M8 2v4M3 10h18"/></svg>
        <span class="stat-num holiday">{{ result.holiday_hours }}h</span>
        <span class="stat-tag">节假日</span>
        <span class="stat-mul">×2.0</span>
      </div>
    </div>
    <details class="detail-section" v-if="result.detail && result.detail.length">
      <summary class="detail-toggle">
        <svg viewBox="0 0 24 24" class="detail-chevron" fill="none" stroke="currentColor" stroke-width="2"><polyline points="6 9 12 15 18 9"/></svg>
        查看计算明细 ({{ result.detail.length }} 条)
      </summary>
      <div class="detail-lines">
        <div class="detail-line" v-for="(d, i) in result.detail" :key="i"
          :class="lineClass(d)">
          {{ d }}
        </div>
      </div>
    </details>
  </div>
</template>

<script>
export default {
  name: 'CalcResult',
  props: {
    result: { type: Object, required: true },
    severity: { type: String, default: '' }
  },
  methods: {
    lineClass(d) {
      if (d.includes('节假日')) return 'holiday-line'
      if (d.includes('非工作时间')) return 'overtime-line'
      if (d.includes('工作时间')) return 'work-line'
      return ''
    }
  }
}
</script>

<style scoped>
.calc-result {
  border-left: 3px solid var(--primary);
  padding-left: 12px;
}

.result-header {
  display: flex; align-items: center; justify-content: space-between;
  margin-bottom: 12px; padding-bottom: 10px;
  border-bottom: 1px solid rgba(39,39,42,0.4);
}
.result-label-row { display: flex; align-items: center; gap: 8px; }
.result-head-icon { width: 17px; height: 17px; color: var(--primary); flex-shrink: 0; }
.result-label { font-size: 13px; color: var(--muted-foreground); font-weight: 500; }

.result-badge {
  font-size: 22px; font-weight: 700; padding: 5px 16px;
  border-radius: 8px; letter-spacing: -0.3px;
}
.result-badge.high { background: rgba(34,197,94,0.12); color: #22c55e; }
.result-badge.mid { background: rgba(245,158,11,0.12); color: #f59e0b; }
.result-badge.low { background: rgba(59,130,246,0.12); color: var(--primary); }

.result-stats { display: flex; gap: 10px; }
.result-stat {
  flex: 1; display: flex; flex-direction: column; align-items: center;
  gap: 5px; padding: 12px 8px; border-radius: var(--radius);
  background: rgba(24,24,27,0.35); border: 1px solid rgba(39,39,42,0.3);
  transition: background .15s;
}
.result-stat:hover { background: rgba(24,24,27,0.5); }
.stat-icon { width: 20px; height: 20px; flex-shrink: 0; margin-bottom: 2px; }
.stat-icon.work { color: #22c55e; }
.stat-icon.overtime { color: #f59e0b; }
.stat-icon.holiday { color: var(--destructive); }

.stat-num { font-size: 19px; font-weight: 700; line-height: 1; }
.stat-num.work { color: #22c55e; }
.stat-num.overtime { color: #f59e0b; }
.stat-num.holiday { color: var(--destructive); }
.stat-tag { font-size: 10.5px; color: var(--muted-foreground); font-weight: 500; }
.stat-mul { font-size: 10px; color: var(--muted-foreground); opacity: 0.6; }

.detail-section { margin-top: 12px; }
.detail-toggle {
  display: flex; align-items: center; gap: 5px;
  font-size: 12px; color: var(--muted-foreground); cursor: pointer; user-select: none;
  padding: 6px 10px; border-radius: 6px;
  transition: all .15s;
}
.detail-toggle:hover { color: var(--foreground); background: rgba(59,130,246,0.06); }
.detail-chevron { width: 14px; height: 14px; flex-shrink: 0; transition: transform .2s; }
details[open] .detail-chevron { transform: rotate(180deg); }

.detail-lines { margin-top: 8px; display: flex; flex-direction: column; gap: 2px; max-height: 200px; overflow-y: auto; padding-left: 4px; }
.detail-line {
  font-size: 11.5px; padding: 5px 8px; border-radius: 4px;
  font-family: 'JetBrains Mono','Consolas',monospace;
  background: rgba(24,24,27,0.3); color: var(--muted-foreground);
  border-left: 2px solid transparent;
}
.detail-line.holiday-line { color: var(--destructive); background: rgba(239,68,68,0.06); border-left-color: var(--destructive); }
.detail-line.overtime-line { color: #f59e0b; background: rgba(245,158,11,0.06); border-left-color: #f59e0b; }
.detail-line.work-line { color: #22c55e; background: rgba(34,197,94,0.06); border-left-color: #22c55e; }
</style>
