<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import { useAuthStore } from '../stores/auth'
import { generatedModuleMenus } from '../generated/module-menus.generated'
import { roleLabel } from '../utils/display'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const collapsed = ref(false)

const currentPath = computed(() => route.path)
const asideWidth = computed(() => (collapsed.value ? 64 : 232))
const asideWidthPx = computed(() => `${asideWidth.value}px`)

const menuIconMap: Record<string, string> = {
  '/': 'dashboard',
  '/users': 'user',
  '/departments': 'root-list',
  '/roles': 'secured',
  '/menus': 'view-list',
  '/dictionary': 'book',
  '/system-params': 'setting-1',
  '/logs': 'chart',
  '/monitor': 'desktop',
  '/positions': 'usergroup',
  '/notices': 'notification',
  '/online-users': 'browse',
  '/scheduled-jobs': 'time'
}

type SidebarMenu = {
  path: string
  label: string
  icon: string
  parentId: number
  sort: number
}

const visibleMenus = computed<SidebarMenu[]>(() => {
  const merged: SidebarMenu[] = auth.menus
    .filter((item) => item.status === 1 && item.menuType === 'menu' && !!item.path)
    .map((item) => {
      const normalizedPath = item.path.startsWith('/') ? item.path : `/${item.path}`
      return {
        path: normalizedPath,
        label: item.name,
        icon: menuIconMap[normalizedPath] || 'layers',
        parentId: item.parentId ?? 0,
        sort: item.sort ?? 0
      }
    })

  const seen = new Set(merged.map((item) => item.path))
  for (const item of generatedModuleMenus) {
    if (!auth.isAdmin && !auth.hasPerm(item.perm)) continue
    if (seen.has(item.path)) continue
    merged.push({
      path: item.path,
      label: item.label,
      icon: 'layers',
      parentId: 0,
      sort: 900
    })
    seen.add(item.path)
  }

  merged.sort((a, b) => (a.parentId - b.parentId) || (a.sort - b.sort) || a.path.localeCompare(b.path))
  if (merged.length === 0) {
    merged.push({ path: '/', label: '工作台', icon: 'dashboard', parentId: 0, sort: 0 })
  }
  return merged
})

function navigate(path: string) {
  router.push(path)
}

function toggleSidebar() {
  collapsed.value = !collapsed.value
}

function doLogout() {
  auth.logout()
  router.push('/login')
}

onMounted(async () => {
  if (!auth.isLoggedIn) return
  if (auth.menus.length > 0) return
  try {
    await auth.loadMenus()
  } catch {
    MessagePlugin.error('菜单加载失败')
  }
})
</script>

<template>
  <t-layout style="height: 100vh; background: var(--td-bg-color-page)">
    <t-aside
      :width="asideWidth"
      class="side-aside"
      theme="dark"
      :style="{
        background: 'var(--td-gray-color-13)',
        transition:
          'width 0.24s cubic-bezier(0.38, 0, 0.24, 1), min-width 0.24s cubic-bezier(0.38, 0, 0.24, 1), max-width 0.24s cubic-bezier(0.38, 0, 0.24, 1), flex-basis 0.24s cubic-bezier(0.38, 0, 0.24, 1)',
        width: asideWidthPx,
        minWidth: asideWidthPx,
        maxWidth: asideWidthPx,
        flexGrow: 0,
        flexShrink: 0,
        flexBasis: asideWidthPx
      }"
    >
      <div style="height: 56px; padding: 0 20px; display: flex; align-items: center; border-bottom: 1px solid rgba(255, 255, 255, 0.08)">
        <div style="width: 24px; height: 24px; border-radius: 4px; background: var(--td-brand-color); color: #fff; font-size: 12px; display: flex; align-items: center; justify-content: center; font-weight: 600">
          T
        </div>
        <transition name="brand-fade">
          <div v-show="!collapsed" style="margin-left: 10px; overflow: hidden">
            <div style="font-size: 15px; font-weight: 600; color: #fff; line-height: 1.2; white-space: nowrap">后台管理系统</div>
            <div style="margin-top: 2px; color: rgba(255, 255, 255, 0.55); font-size: 11px; white-space: nowrap">TDesign Starter 风格</div>
          </div>
        </transition>
      </div>
      <t-menu :value="currentPath" :collapsed="collapsed" class="side-menu" style="width: 100%" theme="dark" @change="(value: string | number) => navigate(String(value))">
        <t-menu-item v-for="item in visibleMenus" :key="item.path" :value="item.path">
          <template #icon><t-icon :name="item.icon" /></template>
          {{ item.label }}
        </t-menu-item>
      </t-menu>
    </t-aside>
    <t-layout>
      <t-header style="height: 56px; border-bottom: 1px solid #e5e6eb; background: #fff; padding: 0 16px; display: flex; justify-content: space-between; align-items: center">
        <div style="display: flex; align-items: center; gap: 10px">
          <t-button variant="text" shape="square" @click="toggleSidebar">☰</t-button>
          <div>
            <div style="font-size: 18px; font-weight: 600; line-height: 1.1">智能后台工作台</div>
            <div style="margin-top: 2px; color: #86909c; font-size: 12px">Spec 驱动模块扩展 · 全栈 CRUD</div>
          </div>
        </div>
        <div style="display: flex; gap: 10px; align-items: center">
          <t-tag theme="primary" variant="outline">{{ roleLabel(auth.user?.role) }}</t-tag>
          <span>{{ auth.user?.displayName }}</span>
          <t-button variant="outline" size="small" @click="doLogout">退出登录</t-button>
        </div>
      </t-header>
      <t-content class="admin-content">
        <router-view />
      </t-content>
    </t-layout>
  </t-layout>
</template>
