<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { fetchMonitorOverview } from '../../api/monitor'
import type { MonitorOverview } from '../../types'
import { dbHealthLabel } from '../../utils/display'

const loading = ref(false)
const overview = ref<MonitorOverview | null>(null)

const runtimeCards = computed(() => {
  if (!overview.value) return []
  return [
    { label: '协程数', value: String(overview.value.runtime.goroutines), foot: '当前并发协程数' },
    { label: '堆内存已分配', value: `${(overview.value.runtime.heapAlloc / 1024 / 1024).toFixed(2)} MB`, foot: '应用当前占用' },
    { label: '堆内存系统占用', value: `${(overview.value.runtime.heapSys / 1024 / 1024).toFixed(2)} MB`, foot: '系统层面分配' },
    { label: 'Go 运行时版本', value: overview.value.runtime.goVersion, foot: '后端运行时信息' }
  ]
})

const countEntries = computed(() => {
  if (!overview.value) return []
  return Object.entries(overview.value.counts).map(([key, value]) => ({ key, value }))
})

async function loadOverview() {
  loading.value = true
  try {
    overview.value = await fetchMonitorOverview()
  } catch {
    MessagePlugin.error('监控信息拉取失败')
  } finally {
    loading.value = false
  }
}

onMounted(loadOverview)
</script>

<template>
  <div class="soft-panel">
    <h2 class="page-title">系统监控</h2>
    <div class="page-subtitle">运行时指标、数据库状态与核心数据量概览。</div>
    <div style="height: 14px" />

    <div class="toolbar">
      <t-tag :theme="overview?.database.health === 'ok' ? 'success' : 'danger'" variant="outline">
        数据库状态：{{ dbHealthLabel(overview?.database.health) }}
      </t-tag>
      <t-tag theme="primary" variant="outline">更新时间：{{ overview?.runtime.now || '-' }}</t-tag>
      <div class="grow" />
      <t-button :loading="loading" @click="loadOverview">刷新</t-button>
    </div>

    <div class="kpi-grid">
      <div v-for="item in runtimeCards" :key="item.label" class="kpi-card">
        <div class="kpi-label">{{ item.label }}</div>
        <div class="kpi-value" style="font-size: 20px">{{ item.value }}</div>
        <div class="kpi-foot">{{ item.foot }}</div>
      </div>
    </div>

    <div style="height: 14px" />
    <div class="data-table-wrap">
      <t-table row-key="key" :data="countEntries" :loading="loading" bordered :columns="[
        { colKey: 'key', title: '统计项' },
        { colKey: 'value', title: '数量' }
      ]" />
    </div>
  </div>
</template>
