<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { listDepartment, type DepartmentDTO } from './api.generated'

const data = ref<DepartmentDTO[]>([])
const loading = ref(false)
const page = ref(1)
const size = ref(10)
const total = ref(0)

const columns = [
  { colKey: 'id', title: 'ID' },
  { colKey: 'name', title: '部门名称' },
  { colKey: 'code', title: '部门编码' },
  { colKey: 'manager', title: '负责人' },
  { colKey: 'status', title: '状态' }
]

async function loadData() {
  loading.value = true
  try {
    const res = await listDepartment(page.value, size.value)
    data.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>

<template>
  <div class="page-card">
    <t-typography-title :level="4" style="margin: 0 0 12px">部门管理（生成示例）</t-typography-title>
    <div class="data-table-wrap">
      <t-table row-key="id" :data="data" :columns="columns" :loading="loading" bordered />
    </div>
    <div style="margin-top: 12px; display: flex; justify-content: flex-end">
      <t-pagination v-model:current="page" v-model:page-size="size" :total="total" @change="loadData" />
    </div>
  </div>
</template>
