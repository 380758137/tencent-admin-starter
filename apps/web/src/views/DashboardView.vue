<script setup lang="ts">
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
</script>

<template>
  <div class="soft-panel">
    <h2 class="page-title">欢迎回来，{{ auth.user?.displayName }}</h2>
    <div class="page-subtitle">你正在使用可扩展后台：传统管理能力 + AI 模块化开发能力。</div>
  </div>

  <div style="height: 14px" />

  <div class="kpi-grid">
    <div class="kpi-card">
      <div class="kpi-label">当前角色</div>
      <div class="kpi-value">{{ auth.user?.role }}</div>
      <div class="kpi-foot">权限中心实时生效</div>
    </div>
    <div class="kpi-card">
      <div class="kpi-label">后台能力</div>
      <div class="kpi-value">10 Modules</div>
      <div class="kpi-foot">底座 + 岗位/公告/在线用户/定时任务</div>
    </div>
    <div class="kpi-card">
      <div class="kpi-label">扩展能力</div>
      <div class="kpi-value">SPEC</div>
      <div class="kpi-foot">支持模块生成</div>
    </div>
    <div class="kpi-card">
      <div class="kpi-label">环境状态</div>
      <div class="kpi-value">ONLINE</div>
      <div class="kpi-foot">前后端服务已联通</div>
    </div>
  </div>

  <div style="height: 14px" />

  <div class="soft-panel">
    <div style="font-size: 15px; font-weight: 600; margin-bottom: 8px">操作建议</div>
    <div style="color: #86909c; line-height: 1.8">
      先在「角色管理 / 菜单管理」完善权限框架，再在「参数中心 / 字典管理」沉淀系统配置，最后通过「日志中心 / 系统监控」完成运行观察闭环。
    </div>
    <div style="margin-top: 12px">
      <t-space>
        <t-tag theme="primary" variant="outline">JWT + RBAC</t-tag>
        <t-tag theme="success" variant="outline">Spec Driven</t-tag>
        <t-tag theme="warning" variant="outline">Dry-run First</t-tag>
      </t-space>
    </div>
  </div>
</template>
