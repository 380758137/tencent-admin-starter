<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { createDepartment, deleteDepartment, fetchDepartments, updateDepartment, type DepartmentPayload } from '../../api/departments'
import type { Department } from '../../types'

const loading = ref(false)
const list = ref<Department[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(10)
const keyword = ref('')
const departmentOptions = ref<Department[]>([])

const visible = ref(false)
const editingId = ref<number | null>(null)
const form = reactive<DepartmentPayload>({
  name: '',
  code: '',
  parentId: 0,
  manager: '',
  status: 1
})

const columns = [
  { colKey: 'id', title: 'ID', width: 70 },
  { colKey: 'name', title: '部门名称' },
  { colKey: 'parentId', title: '父部门' },
  { colKey: 'code', title: '编码' },
  { colKey: 'manager', title: '负责人' },
  { colKey: 'status', title: '状态' },
  { colKey: 'actions', title: '操作', width: 160 }
]

const parentOptions = computed(() => departmentOptions.value.filter((item) => item.id !== editingId.value))

function resetForm() {
  form.name = ''
  form.code = ''
  form.parentId = 0
  form.manager = ''
  form.status = 1
}

async function loadDepartments() {
  loading.value = true
  try {
    const res = await fetchDepartments({ page: page.value, size: size.value, keyword: keyword.value || undefined })
    list.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

async function loadDepartmentOptions() {
  const pageSize = 200
  let current = 1
  let totalCount = 0
  const all: Department[] = []

  do {
    const res = await fetchDepartments({ page: current, size: pageSize })
    totalCount = res.total
    all.push(...res.list)
    current += 1
    if (res.list.length === 0) break
  } while (all.length < totalCount)

  departmentOptions.value = all
}

function parentName(parentId: number): string {
  if (parentId === 0) return '-'
  const parent = departmentOptions.value.find((item) => item.id === parentId)
  return parent ? parent.name : `#${parentId}`
}

function openCreate() {
  editingId.value = null
  resetForm()
  visible.value = true
}

function openEdit(row: Department) {
  editingId.value = row.id
  form.name = row.name
  form.code = row.code
  form.parentId = row.parentId
  form.manager = row.manager
  form.status = row.status
  visible.value = true
}

async function submit() {
  try {
    if (editingId.value) {
      await updateDepartment(editingId.value, form)
      MessagePlugin.success('更新成功')
    } else {
      await createDepartment(form)
      MessagePlugin.success('创建成功')
    }
    visible.value = false
    await Promise.all([loadDepartments(), loadDepartmentOptions()])
  } catch {
    MessagePlugin.error('保存失败')
  }
}

async function remove(row: Department) {
  try {
    await deleteDepartment(row.id)
    MessagePlugin.success('删除成功')
    await Promise.all([loadDepartments(), loadDepartmentOptions()])
  } catch {
    MessagePlugin.error('删除失败')
  }
}

onMounted(async () => {
  await Promise.all([loadDepartments(), loadDepartmentOptions()])
})
</script>

<template>
  <div class="soft-panel">
    <h2 class="page-title">部门管理</h2>
    <div class="page-subtitle">维护组织结构与负责人信息。</div>
    <div style="height: 14px" />
    <div class="toolbar">
      <t-input v-model="keyword" placeholder="搜索部门名称/编码" style="width: 260px" />
      <t-button @click="loadDepartments">查询</t-button>
      <div class="grow" />
      <t-button theme="primary" @click="openCreate">新建部门</t-button>
    </div>

    <div class="data-table-wrap">
      <t-table row-key="id" :data="list" :columns="columns" :loading="loading" bordered>
      <template #parentId="{ row }">
        {{ parentName(row.parentId) }}
      </template>
      <template #status="{ row }">
        <t-tag class="table-status" :theme="row.status === 1 ? 'success' : 'warning'">{{ row.status === 1 ? '启用' : '禁用' }}</t-tag>
      </template>
      <template #actions="{ row }">
        <t-space>
          <t-link theme="primary" hover="color" @click="openEdit(row)">编辑</t-link>
          <t-link theme="danger" hover="color" @click="remove(row)">删除</t-link>
        </t-space>
      </template>
    </t-table>
    </div>

    <div style="margin-top: 12px; display: flex; justify-content: flex-end">
      <t-pagination
        v-model:current="page"
        v-model:page-size="size"
        :total="total"
        @change="loadDepartments"
      />
    </div>

    <t-dialog v-model:visible="visible" :header="editingId ? '编辑部门' : '新建部门'" width="560px" @confirm="submit">
      <t-form layout="vertical">
        <t-form-item label="部门名称">
          <t-input v-model="form.name" />
        </t-form-item>
        <t-form-item label="部门编码">
          <t-input v-model="form.code" />
        </t-form-item>
        <t-form-item label="父部门">
          <t-select v-model="form.parentId">
            <t-option :value="0" label="无（顶级部门）" />
            <t-option v-for="item in parentOptions" :key="item.id" :value="item.id" :label="item.name" />
          </t-select>
        </t-form-item>
        <t-form-item label="负责人">
          <t-input v-model="form.manager" />
        </t-form-item>
        <t-form-item label="状态">
          <t-select v-model="form.status">
            <t-option :value="1" label="启用" />
            <t-option :value="0" label="禁用" />
          </t-select>
        </t-form-item>
      </t-form>
    </t-dialog>
  </div>
</template>
