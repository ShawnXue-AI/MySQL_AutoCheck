<template>
  <transition name="modal">
    <div v-if="visible" class="modal-overlay" @click.self="() => {}">
      <div class="modal-panel">
        <div class="modal-body">
          <!-- Phase icons -->
          <svg v-if="phase === 'upload'" viewBox="0 0 48 48" class="modal-icon up"><path d="M24 16v12m0 0l-4-4m4 4l4-4M14 28v4a4 4 0 004 4h12a4 4 0 004-4v-4" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round"/><path d="M24 4v12m0 0l-4-4m4 4l4-4" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" opacity="0.5"/></svg>
          <svg v-else-if="phase === 'running'" viewBox="0 0 48 48" class="modal-icon spin"><path d="M24 4v4m0 32v4m20-20h-4M8 24H4m33.9-9.9l-2.8 2.8M12.9 33.9l-2.8 2.8m28.3 0l-2.8-2.8M12.9 14.1l-2.8-2.8" stroke="currentColor" stroke-width="3" fill="none" stroke-linecap="round"/></svg>
          <svg v-else-if="phase === 'done'" viewBox="0 0 48 48" class="modal-icon done"><circle cx="24" cy="24" r="18" fill="none" stroke="var(--accent-green)" stroke-width="3"/><path d="M16 24l6 6 10-10" stroke="var(--accent-green)" stroke-width="3" fill="none" stroke-linecap="round" stroke-linejoin="round"/></svg>
          <svg v-else-if="phase === 'error'" viewBox="0 0 48 48" class="modal-icon err"><circle cx="24" cy="24" r="18" fill="none" stroke="var(--accent-red)" stroke-width="3"/><path d="M18 18l12 12M30 18l-12 12" stroke="var(--accent-red)" stroke-width="3" fill="none" stroke-linecap="round"/></svg>

          <h3 class="modal-title">{{ titleText }}</h3>
          <p class="modal-sub" v-if="subText">{{ subText }}</p>

          <!-- Progress bar -->
          <div class="modal-track" v-if="phase !== 'done' && phase !== 'error'">
            <div class="modal-fill" :style="{ width: progress + '%' }"></div>
          </div>
          <p class="modal-pct" v-if="phase === 'upload'">{{ progress }}%</p>

          <!-- Terminal-style status line -->
          <div class="modal-status-line" v-if="phase === 'running'">
            <span class="terminal-prompt">$</span>
            <span>building report</span>
            <span class="status-dots">
              <span class="sdot">.</span><span class="sdot">.</span><span class="sdot">.</span>
            </span>
          </div>

          <p class="modal-err" v-if="errorMsg">{{ errorMsg }}</p>
        </div>
        <div class="modal-actions">
          <button v-if="phase === 'upload'" class="btn btn-outline btn-sm" style="border-color: rgba(248,81,73,0.25); color: var(--accent-red);" @click="emitCancel" :disabled="cancelling">
            {{ cancelling ? '取消中...' : '取消上传' }}
          </button>
          <button v-if="phase === 'running'" class="btn btn-outline btn-sm" style="border-color: rgba(248,81,73,0.25); color: var(--accent-red);" @click="emitCancel" :disabled="cancelling">
            {{ cancelling ? '终止中...' : '终止任务' }}
          </button>
          <a v-if="phase === 'done'" :href="downloadUrl" class="btn btn-primary">
            <svg viewBox="0 0 24 24" class="btn-icon"><path d="M19 9h-4V3H9v6H5l7 7 7-7zM5 18v2h14v-2H5z" fill="currentColor"/></svg>
            下载报告 (ZIP)
          </a>
          <button v-if="phase === 'done' || phase === 'error'" class="btn btn-ghost" @click="emitClose">关闭</button>
        </div>
      </div>
    </div>
  </transition>
</template>

<script>
export default {
  name: 'ProcessModal',
  props: {
    visible: { type: Boolean, default: false },
    phase: { type: String, default: 'upload' },
    progress: { type: Number, default: 0 },
    titleText: { type: String, default: '' },
    subText: { type: String, default: '' },
    errorMsg: { type: String, default: '' },
    downloadUrl: { type: String, default: '' },
    cancelling: { type: Boolean, default: false }
  },
  emits: ['cancel', 'close'],
  methods: {
    emitCancel() { this.$emit('cancel') },
    emitClose() { this.$emit('close') }
  }
}
</script>

<style scoped>
.modal-overlay {
  position: fixed; inset: 0; z-index: 9000;
  display: flex; align-items: center; justify-content: center;
  background: rgba(0, 0, 0, 0.6); backdrop-filter: blur(8px);
}

.modal-panel {
  width: 420px; max-width: 92vw;
  background: var(--bg-secondary); border: 1px solid var(--border-hover);
  border-radius: var(--radius-xl); overflow: hidden;
  box-shadow: var(--shadow-lg);
}

.modal-body {
  padding: 40px 28px 24px;
  display: flex; flex-direction: column; align-items: center; text-align: center; gap: 12px;
}

.modal-icon { width: 52px; height: 52px; }
.modal-icon.up { color: var(--accent-blue); }
.modal-icon.spin { color: var(--accent-blue); animation: spin-icon 1.8s linear infinite; }
.modal-icon.err { color: var(--accent-red); }

@keyframes spin-icon { to { transform: rotate(360deg); } }

.modal-title { font-family: var(--font-display); font-size: 17px; font-weight: 600; color: var(--text-primary); }
.modal-sub { font-family: var(--font-code); font-size: 12.5px; color: var(--text-tertiary); }

.modal-track {
  width: 100%; height: 4px; background: var(--bg-input);
  border-radius: 2px; overflow: hidden;
}
.modal-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--accent-blue), var(--accent-cyan));
  border-radius: 2px; transition: width 0.4s ease;
}
.modal-pct {
  font-family: var(--font-code); font-size: 24px; font-weight: 600;
  color: var(--text-primary); letter-spacing: 1px;
  font-variant-numeric: tabular-nums;
}

.modal-status-line {
  display: flex; align-items: center; gap: 6px;
  font-family: var(--font-code); font-size: 12.5px; color: var(--text-tertiary);
}
.terminal-prompt { color: var(--accent-green); font-weight: 500; }
.status-dots { display: flex; gap: 2px; }
.sdot { animation: dot-blink 1.4s ease-in-out infinite; }
.sdot:nth-child(2) { animation-delay: 0.2s; }
.sdot:nth-child(3) { animation-delay: 0.4s; }
@keyframes dot-blink {
  0%, 100% { opacity: 0.2; }
  50% { opacity: 1; }
}

.modal-err {
  font-family: var(--font-code); font-size: 12.5px; color: var(--accent-red);
  background: var(--accent-red-subtle); padding: 8px 14px;
  border-radius: var(--radius-sm); width: 100%;
}

.modal-actions {
  display: flex; justify-content: center; gap: 10px;
  padding: 16px 28px 24px; border-top: 1px solid var(--border-default);
}

/* Transitions */
.modal-enter-active { transition: all 0.25s ease-out; }
.modal-leave-active { transition: all 0.2s ease-in; }
.modal-enter-from { opacity: 0; }
.modal-enter-from .modal-panel { transform: scale(0.92); opacity: 0; }
.modal-leave-to { opacity: 0; }
.modal-leave-to .modal-panel { transform: scale(0.92); opacity: 0; }
</style>
