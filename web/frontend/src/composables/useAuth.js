import { ref } from 'vue'
import { adminAuth } from '../api/index.js'

export function useAuth() {
  const authenticated = ref(false)
  const password = ref('')
  const authing = ref(false)
  const authError = ref('')

  async function doAuth() {
    authing.value = true
    authError.value = ''
    try {
      const { data } = await adminAuth(password.value)
      if (data.success) {
        authenticated.value = true
        return true
      } else {
        authError.value = data.error || '密码错误'
        return false
      }
    } catch (e) {
      authError.value = '请求失败，请检查后端服务'
      return false
    } finally {
      authing.value = false
    }
  }

  return { authenticated, password, authing, authError, doAuth }
}
