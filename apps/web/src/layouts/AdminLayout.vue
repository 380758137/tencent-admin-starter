<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { roleLabel } from '../utils/display'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const collapsed = ref(false)

const currentPath = computed(() => route.path)
const asideWidth = computed(() => (collapsed.value ? 64 : 232))
const asideWidthPx = computed(() => `${asideWidth.value}px`)

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
        <t-menu-item value="/">
          <template #icon><t-icon name="dashboard" /></template>
          工作台
        </t-menu-item>
        <t-menu-item value="/users">
          <template #icon><t-icon name="user" /></template>
          用户管理
        </t-menu-item>
        <t-menu-item value="/departments">
          <template #icon><t-icon name="root-list" /></template>
          部门管理
        </t-menu-item>
        <t-menu-item value="/roles">
          <template #icon><t-icon name="secured" /></template>
          角色管理
        </t-menu-item>
        <t-menu-item value="/menus">
          <template #icon><t-icon name="view-list" /></template>
          菜单管理
        </t-menu-item>
        <t-menu-item value="/dictionary">
          <template #icon><t-icon name="book" /></template>
          字典管理
        </t-menu-item>
        <t-menu-item value="/system-params">
          <template #icon><t-icon name="setting-1" /></template>
          参数中心
        </t-menu-item>
        <t-menu-item value="/logs">
          <template #icon><t-icon name="chart" /></template>
          日志中心
        </t-menu-item>
        <t-menu-item value="/monitor">
          <template #icon><t-icon name="desktop" /></template>
          系统监控
        </t-menu-item>
        <t-menu-item value="/positions">
          <template #icon><t-icon name="usergroup" /></template>
          岗位管理
        </t-menu-item>
        <t-menu-item value="/notices">
          <template #icon><t-icon name="notification" /></template>
          通知公告
        </t-menu-item>
        <t-menu-item value="/online-users">
          <template #icon><t-icon name="browse" /></template>
          在线用户
        </t-menu-item>
        <t-menu-item value="/scheduled-jobs">
          <template #icon><t-icon name="time" /></template>
          定时任务
        </t-menu-item>
        <t-menu-item value="/generated/department">
          <template #icon><t-icon name="layers" /></template>
          部门管理（生成示例）
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
