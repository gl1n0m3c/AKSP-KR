import { ref } from 'vue'
import api from '../api'
import type { MeResponse } from '../types'

const user = ref<MeResponse | null>(null)
const loading = ref(false)

export function useAuth() {
  const error = ref<string | null>(null)
  const showAuth = ref(false)

  const loadMe = async () => {
    loading.value = true
    error.value = null
    try {
      const res = await api.get<MeResponse>('/api/me')
      user.value = res.data
    } catch (e: any) {
      user.value = null
      showAuth.value = true
    } finally {
      loading.value = false
    }
  }

  const login = async (login: string, password: string) => {
    loading.value = true
    error.value = null
    try {
      await api.post('/api/login', { login, password })
      await loadMe()
    } catch (e: any) {
      error.value = e?.response?.data || 'login error'
      throw e
    } finally {
      loading.value = false
    }
  }

  const register = async (payload: any) => {
    loading.value = true
    error.value = null
    try {
      await api.post('/api/register', payload)
      await loadMe()
    } catch (e: any) {
      error.value = e?.response?.data || 'register error'
      throw e
    } finally {
      loading.value = false
    }
  }

  return {
    user,
    loading,
    error,
    showAuth,
    loadMe,
    login,
    register,
  }
}

