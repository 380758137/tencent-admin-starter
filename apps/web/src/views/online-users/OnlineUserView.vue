<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { fetchOnlineUsers } from '../../api/online-users'
import type { OnlineUser } from '../../types'

const loading = ref(false)
const list = ref<OnlineUser[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(10)
const keyword = ref('')

const columns = [
  { colKey: 'username', title: '用户名' },
  { colKey: 'lastActiveAt', title: '最后活跃时间' },
  { colKey: 'requestCount', title: '近30分钟请求数' }
]

async function loadOnlineUsers() {
  loading.value = true
  try {
    const res = await fetchOnlineUsers({ page: page.value, size: size.value, keyword: keyword.value || undefined })
    list.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

onMounted(loadOnlineUsers)
</script>

<template>
  <div class="soft-panel">
    <h2 class="page-title">在线用户</h2>
    <div class="page-subtitle">基于最近 30 分钟操作日志聚合在线活跃用户。</div>
    <div style="height: 14px" />

    <div class="toolbar">
      <t-input v-model="keyword" placeholder="搜索用户名" style="width: 260px" />
      <t-button @click="loadOnlineUsers">查询</t-button>
      <div class="grow" />
      <t-button variant="outline" @click="loadOnlineUsers">刷新</t-button>
    </div>

    <div class="data-table-wrap">
      <t-table row-key="username" :data="list" :columns="columns" :loading="loading" bordered />
    </div>

    <div style="margin-top: 12px; display: flex; justify-content: flex-end">
      <t-pagination v-model:current="page" v-model:page-size="size" :total="total" @change="loadOnlineUsers" />
    </div>
  </div>
</template>
