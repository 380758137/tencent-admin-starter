import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { login, me } from '../api/auth'
import type { AuthUser, LoginPayload } from '../types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const user = ref<AuthUser | null>(null)

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const permissions = computed(() => user.value?.permissions || [])

  async function signIn(payload: LoginPayload) {
    const result = await login(payload)
    token.value = result.token
    user.value = result.user
    localStorage.setItem('token', result.token)
  }

  async function loadMe() {
    if (!token.value) return
    user.value = await me()
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
  }

  function hasPerm(perm: string) {
    if (isAdmin.value) return true
    return permissions.value.includes('*') || permissions.value.includes(perm)
  }

  return {
    token,
    user,
    isLoggedIn,
    isAdmin,
    permissions,
    hasPerm,
    signIn,
    loadMe,
    logout
  }
})
