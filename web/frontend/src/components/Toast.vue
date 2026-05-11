<template>
  <transition-group name="toast" tag="div" class="toast-container">
    <div v-for="item in items" :key="item.id" class="toast-item" :class="item.type">
      <svg viewBox="0 0 24 24" class="toast-icon" v-if="item.type === 'success'"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z" fill="currentColor"/></svg>
      <svg viewBox="0 0 24 24" class="toast-icon" v-else-if="item.type === 'error'"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z" fill="currentColor"/></svg>
      <svg viewBox="0 0 24 24" class="toast-icon" v-else><path d="M1 21h22L12 2 1 21zm12-3h-2v-2h2v2zm0-4h-2v-4h2v4z" fill="currentColor"/></svg>
      <span class="toast-msg">{{ item.message }}</span>
    </div>
  </transition-group>
</template>

<script>
let _instance = null

export default {
  name: 'Toast',
  data() {
    return { items: [], nextId: 0 }
  },
  created() {
    _instance = this
  },
  methods: {
    show(message, type = 'info', duration = 3500) {
      const id = ++this.nextId
      this.items.push({ id, message, type })
      setTimeout(() => {
        this.items = this.items.filter(i => i.id !== id)
      }, duration)
    }
  }
}

export function toast(msg, type) { if (_instance) _instance.show(msg, type) }
export function toastSuccess(msg) { toast(msg, 'success') }
export function toastError(msg) { toast(msg, 'error') }
</script>

<style scoped>
.toast-container {
  position: fixed; top: 24px; right: 24px; z-index: 9999;
  display: flex; flex-direction: column; gap: 10px; min-width: 300px;
}

.toast-item {
  display: flex; align-items: center; gap: 10px;
  padding: 14px 20px; border-radius: 10px;
  font-size: 14px; font-weight: 500; box-shadow: 0 8px 24px rgba(0,0,0,0.4);
}

.toast-item.success { background: #238636; color: #fff; }
.toast-item.error { background: #da3633; color: #fff; }
.toast-item.info { background: #1f6feb; color: #fff; }

.toast-icon { width: 20px; height: 20px; flex-shrink: 0; }

.toast-enter-active { transition: all 0.3s ease-out; }
.toast-leave-active { transition: all 0.25s ease-in; }
.toast-enter-from { opacity: 0; transform: translateX(80px); }
.toast-leave-to { opacity: 0; transform: translateX(80px); }
</style>
