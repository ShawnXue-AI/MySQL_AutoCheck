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
        <span class="stat-mul">x1.0</span>
      </div>
      <div class="result-stat">
        <svg viewBox="0 0 24 24" class="stat-icon overtime" fill="none" stroke="currentColor" stroke-width="2"><path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/></svg>
        <span class="stat-num overtime">{{ result.overtime_hours }}h</span>
        <span class="stat-tag">加班时间</span>
        <span class="stat-mul">x1.5</span>
      </div>
      <div class="result-stat">
        <svg viewBox="0 0 24 24" class="stat-icon holiday" fill="none" stroke="currentColor" stroke-width="2"><rect width="18" height="18" x="3" y="4" rx="2"/><path d="M16 2v4M8 2v4M3 10h18"/></svg>
        <span class="stat-num holiday">{{ result.holiday_hours }}h</span>
        <span class="stat-tag">节假日</span>
        <span class="stat-mul">x2.0</span>
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
  border-left: 3px solid var(--accent-blue);
  padding-left: 14px;
}

.result-header {
  display: flex; align-items: center; justify-content: space-between;
  margin-bottom: 14px; padding-bottom: 10px;
  border-bottom: 1px solid var(--border-default);
}
.result-label-row { display: flex; align-items: center; gap: 8px; }
.result-head-icon { width: 18px; height: 18px; color: var(--accent-blue); flex-shrink: 0; }
.result-label { font-size: 13px; color: var(--text-tertiary); font-weight: 500; }

.result-badge {
  font-family: var(--font-code);
  font-size: 22px; font-weight: 700; padding: 5px 16px;
  border-radius: var(--radius-sm); letter-spacing: -0.3px;
}
.result-badge.high { background: var(--accent-green-subtle); color: var(--accent-green); }
.result-badge.mid { background: var(--accent-amber-subtle); color: var(--accent-amber); }
.result-badge.low { background: var(--accent-blue-subtle); color: var(--accent-blue); }

.result-stats { display: flex; gap: 10px; }
.result-stat {
  flex: 1; display: flex; flex-direction: column; align-items: center;
  gap: 5px; padding: 14px 8px; border-radius: var(--radius-md);
  background: var(--bg-tertiary); border: 1px solid var(--border-default);
  transition: background var(--transition-fast);
}
.result-stat:hover { background: var(--bg-hover); }

.stat-icon { width: 20px; height: 20px; flex-shrink: 0; margin-bottom: 2px; }
.stat-icon.work { color: var(--accent-green); }
.stat-icon.overtime { color: var(--accent-amber); }
.stat-icon.holiday { color: var(--accent-red); }

.stat-num {
  font-family: var(--font-code);
  font-size: 19px; font-weight: 700; line-height: 1;
  font-variant-numeric: tabular-nums;
}
.stat-num.work { color: var(--accent-green); }
.stat-num.overtime { color: var(--accent-amber); }
.stat-num.holiday { color: var(--accent-red); }
.stat-tag { font-size: 10.5px; color: var(--text-tertiary); font-weight: 500; }
.stat-mul { font-family: var(--font-code); font-size: 10px; color: var(--text-muted); }

.detail-section { margin-top: 14px; }
.detail-toggle {
  display: flex; align-items: center; gap: 5px;
  font-size: 12px; color: var(--text-tertiary); cursor: pointer; user-select: none;
  padding: 6px 10px; border-radius: var(--radius-sm);
  transition: all var(--transition-fast);
}
.detail-toggle:hover { color: var(--text-secondary); background: var(--accent-blue-subtle); }
.detail-chevron { width: 14px; height: 14px; flex-shrink: 0; transition: transform 0.2s; }
details[open] .detail-chevron { transform: rotate(180deg); }

.detail-lines {
  margin-top: 8px; display: flex; flex-direction: column; gap: 2px;
  max-height: 200px; overflow-y: auto; padding-left: 4px;
}
.detail-line {
  font-family: var(--font-code);
  font-size: 11.5px; padding: 5px 8px; border-radius: 4px;
  background: var(--bg-tertiary); color: var(--text-tertiary);
  border-left: 2px solid transparent;
}
.detail-line.holiday-line { color: var(--accent-red); background: var(--accent-red-subtle); border-left-color: var(--accent-red); }
.detail-line.overtime-line { color: var(--accent-amber); background: var(--accent-amber-subtle); border-left-color: var(--accent-amber); }
.detail-line.work-line { color: var(--accent-green); background: var(--accent-green-subtle); border-left-color: var(--accent-green); }
</style>
