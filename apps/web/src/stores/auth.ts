import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { fetchMyMenus, login, me } from '../api/auth'
import type { AuthUser, LoginPayload, Menu } from '../types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const user = ref<AuthUser | null>(null)
  const menus = ref<Menu[]>([])

  const isLoggedIn = computed(() => !!token.value)
  const roleKeys = computed(() =>
    (user.value?.role || '')
      .split(',')
      .map((item) => item.trim())
      .filter(Boolean)
  )
  const isAdmin = computed(() => roleKeys.value.includes('admin'))
  const permissions = computed(() => user.value?.permissions || [])

  async function signIn(payload: LoginPayload) {
    const result = await login(payload)
    token.value = result.token
    user.value = result.user
    localStorage.setItem('token', result.token)
    try {
      menus.value = await fetchMyMenus()
    } catch {
      menus.value = []
    }
  }

  async function loadMe() {
    if (!token.value) return
    user.value = await me()
  }

  async function loadMenus() {
    if (!token.value) return
    menus.value = await fetchMyMenus()
  }

  function logout() {
    token.value = ''
    user.value = null
    menus.value = []
    localStorage.removeItem('token')
  }

  function hasPerm(perm: string) {
    if (isAdmin.value) return true
    return permissions.value.includes('*') || permissions.value.includes(perm)
  }

  return {
    token,
    user,
    menus,
    isLoggedIn,
    isAdmin,
    permissions,
    hasPerm,
    signIn,
    loadMe,
    loadMenus,
    logout
  }
})
