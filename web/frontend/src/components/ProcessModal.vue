<template>
  <transition name="modal">
    <div v-if="visible" class="modal-overlay" @click.self="() => {}">
      <div class="modal-panel">
        <div class="modal-body">
          <svg v-if="phase === 'upload'" viewBox="0 0 48 48" class="modal-icon up"><path d="M24 16v12m0 0l-4-4m4 4l4-4M14 28v4a4 4 0 004 4h12a4 4 0 004-4v-4" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round"/><path d="M24 4v12m0 0l-4-4m4 4l4-4" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" opacity="0.5"/></svg>
          <svg v-else-if="phase === 'running'" viewBox="0 0 48 48" class="modal-icon spin"><path d="M24 4v4m0 32v4m20-20h-4M8 24H4m33.9-9.9l-2.8 2.8M12.9 33.9l-2.8 2.8m28.3 0l-2.8-2.8M12.9 14.1l-2.8-2.8" stroke="currentColor" stroke-width="3" fill="none" stroke-linecap="round"/></svg>
          <svg v-else-if="phase === 'done'" viewBox="0 0 48 48" class="modal-icon done"><circle cx="24" cy="24" r="18" fill="none" stroke="#3fb950" stroke-width="3"/><path d="M16 24l6 6 10-10" stroke="#3fb950" stroke-width="3" fill="none" stroke-linecap="round" stroke-linejoin="round"/></svg>
          <svg v-else-if="phase === 'error'" viewBox="0 0 48 48" class="modal-icon err"><circle cx="24" cy="24" r="18" fill="none" stroke="#f85149" stroke-width="3"/><path d="M18 18l12 12M30 18l-12 12" stroke="#f85149" stroke-width="3" fill="none" stroke-linecap="round"/></svg>

          <h3 class="modal-title">{{ titleText }}</h3>
          <p class="modal-sub">{{ subText }}</p>

          <div class="modal-track" v-if="phase !== 'done' && phase !== 'error' && phase !== 'cancelled'">
            <div class="modal-fill" :style="{ width: progress + '%' }"></div>
          </div>
          <p class="modal-pct" v-if="phase === 'upload'">{{ progress }}%</p>

          <p class="modal-err" v-if="errorMsg">{{ errorMsg }}</p>
        </div>
        <div class="modal-actions">
          <button v-if="phase === 'upload'" class="btn btn-danger-outline btn-sm" @click="emitCancel" :disabled="cancelling">
            {{ cancelling ? '取消中...' : '取消上传' }}
          </button>
          <button v-if="phase === 'running'" class="btn btn-danger-outline btn-sm" @click="emitCancel" :disabled="cancelling">
            {{ cancelling ? '终止中...' : '终止任务' }}
          </button>
          <a v-if="phase === 'done'" :href="downloadUrl" class="btn btn-primary">
            <svg viewBox="0 0 24 24" class="btn-icon"><path d="M19 9h-4V3H9v6H5l7 7 7-7zM5 18v2h14v-2H5z" fill="currentColor"/></svg>
            下载报告 (ZIP)
          </a>
          <button v-if="phase === 'done' || phase === 'error' || phase === 'cancelled'" class="btn btn-ghost" @click="emitClose">关闭</button>
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
  background: rgba(0,0,0,0.55); backdrop-filter: blur(6px);
}

.modal-panel {
  width: 400px; max-width: 92vw;
  background: #161b22; border: 1px solid #30363d;
  border-radius: 16px; overflow: hidden;
  box-shadow: 0 24px 64px rgba(0,0,0,0.5);
}

.modal-body {
  padding: 36px 28px 24px;
  display: flex; flex-direction: column; align-items: center; text-align: center; gap: 10px;
}

.modal-icon { width: 52px; height: 52px; }
.modal-icon.up { color: #58a6ff; }
.modal-icon.spin { color: #58a6ff; animation: spin-icon 1.8s linear infinite; }
.modal-icon.err { color: #f85149; }

@keyframes spin-icon { to { transform: rotate(360deg); } }

.modal-title { font-size: 17px; font-weight: 600; color: #f0f6fc; }

.modal-sub { font-size: 13px; color: #8b949e; }

.modal-track {
  width: 100%; height: 5px; background: #21262d; border-radius: 3px; overflow: hidden;
}

.modal-fill {
  height: 100%; background: linear-gradient(90deg, #238636, #2ea043);
  border-radius: 3px; transition: width 0.4s;
}

.modal-pct { font-size: 22px; font-weight: 700; color: #f0f6fc; letter-spacing: 1px; }

.modal-err { font-size: 13px; color: #f85149; background: rgba(218,54,51,0.08); padding: 8px 14px; border-radius: 6px; width: 100%; }

.modal-actions {
  display: flex; justify-content: center; gap: 10px;
  padding: 16px 28px 24px; border-top: 1px solid #21262d;
}

.btn { display: inline-flex; align-items: center; gap: 8px; padding: 9px 20px; border: none; border-radius: 8px; font-size: 14px; font-weight: 500; cursor: pointer; transition: all 0.15s; text-decoration: none; }
.btn:disabled { opacity: 0.4; cursor: not-allowed; }
.btn-icon { width: 16px; height: 16px; }
.btn-primary { background: #238636; color: #fff; }
.btn-primary:hover { background: #2ea043; }
.btn-ghost { background: transparent; color: #8b949e; border: 1px solid #30363d; padding: 9px 20px; }
.btn-ghost:hover { color: #c9d1d9; border-color: #484f58; }
.btn-danger-outline { background: transparent; border: 1px solid rgba(248,81,73,0.4); color: #f85149; padding: 7px 18px; font-size: 13px; }
.btn-danger-outline:hover:not(:disabled) { background: #da3633; color: #fff; border-color: #da3633; }
.btn-sm { padding: 6px 14px; font-size: 13px; }

.modal-enter-active { transition: all 0.25s ease-out; }
.modal-leave-active { transition: all 0.2s ease-in; }
.modal-enter-from { opacity: 0; }
.modal-enter-from .modal-panel { transform: scale(0.92); opacity: 0; }
.modal-leave-to { opacity: 0; }
.modal-leave-to .modal-panel { transform: scale(0.92); opacity: 0; }
</style>
