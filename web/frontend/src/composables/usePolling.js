import { ref, onBeforeUnmount } from 'vue'

export function usePolling(fn, intervalMs = 5000) {
  const timer = ref(null)

  function start() {
    stop()
    fn()
    timer.value = setInterval(() => fn(), intervalMs)
  }

  function stop() {
    if (timer.value) {
      clearInterval(timer.value)
      timer.value = null
    }
  }

  function restart(newInterval) {
    if (newInterval) intervalMs = newInterval
    stop()
    timer.value = setInterval(() => fn(), intervalMs)
  }

  onBeforeUnmount(stop)

  return { start, stop, restart }
}
