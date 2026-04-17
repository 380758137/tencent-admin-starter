<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import {
  createSystemParam,
  deleteSystemParam,
  fetchSystemParams,
  updateSystemParam,
  type SystemParamPayload
} from '../../api/system-params'
import type { SystemParam } from '../../types'

const loading = ref(false)
const list = ref<SystemParam[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(10)
const keyword = ref('')

const visible = ref(false)
const editingId = ref<number | null>(null)
const form = reactive<SystemParamPayload>({
  paramKey: '',
  paramValue: '',
  paramName: '',
  status: 1,
  remark: ''
})

const columns = [
  { colKey: 'id', title: 'ID', width: 70 },
  { colKey: 'paramName', title: '参数名称' },
  { colKey: 'paramKey', title: '参数键' },
  { colKey: 'paramValue', title: '参数值' },
  { colKey: 'status', title: '状态', width: 90 },
  { colKey: 'remark', title: '备注' },
  { colKey: 'actions', title: '操作', width: 160 }
]

function resetForm() {
  form.paramKey = ''
  form.paramValue = ''
  form.paramName = ''
  form.status = 1
  form.remark = ''
}

async function loadParams() {
  loading.value = true
  try {
    const res = await fetchSystemParams({ page: page.value, size: size.value, keyword: keyword.value || undefined })
    list.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingId.value = null
  resetForm()
  visible.value = true
}

function openEdit(row: SystemParam) {
  editingId.value = row.id
  form.paramKey = row.paramKey
  form.paramValue = row.paramValue
  form.paramName = row.paramName
  form.status = row.status
  form.remark = row.remark
  visible.value = true
}

async function submit() {
  try {
    if (editingId.value) {
      await updateSystemParam(editingId.value, {
        paramName: form.paramName,
        paramValue: form.paramValue,
        status: form.status,
        remark: form.remark
      })
      MessagePlugin.success('更新成功')
    } else {
      await createSystemParam(form)
      MessagePlugin.success('创建成功')
    }
    visible.value = false
    await loadParams()
  } catch {
    MessagePlugin.error('保存失败')
  }
}

async function remove(row: SystemParam) {
  try {
    await deleteSystemParam(row.id)
    MessagePlugin.success('删除成功')
    await loadParams()
  } catch {
    MessagePlugin.error('删除失败')
  }
}

onMounted(loadParams)
</script>

<template>
  <div class="soft-panel">
    <h2 class="page-title">参数中心</h2>
    <div class="page-subtitle">集中管理平台动态参数与开关。</div>
    <div style="height: 14px" />

    <div class="toolbar">
      <t-input v-model="keyword" placeholder="搜索参数键/名称" style="width: 260px" />
      <t-button @click="loadParams">查询</t-button>
      <div class="grow" />
      <t-button theme="primary" @click="openCreate">新增参数</t-button>
    </div>

    <div class="data-table-wrap">
      <t-table row-key="id" :data="list" :columns="columns" :loading="loading" bordered>
        <template #status="{ row }">
          <t-tag class="table-status" :theme="row.status === 1 ? 'success' : 'warning'">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </t-tag>
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
      <t-pagination v-model:current="page" v-model:page-size="size" :total="total" @change="loadParams" />
    </div>

    <t-dialog v-model:visible="visible" :header="editingId ? '编辑参数' : '新增参数'" width="620px" @confirm="submit">
      <t-form layout="vertical">
        <t-form-item label="参数名称">
          <t-input v-model="form.paramName" />
        </t-form-item>
        <t-form-item label="参数键">
          <t-input v-model="form.paramKey" :disabled="!!editingId" />
        </t-form-item>
        <t-form-item label="参数值">
          <t-input v-model="form.paramValue" />
        </t-form-item>
        <t-form-item label="状态">
          <t-select v-model="form.status">
            <t-option :value="1" label="启用" />
            <t-option :value="0" label="禁用" />
          </t-select>
        </t-form-item>
        <t-form-item label="备注">
          <t-input v-model="form.remark" />
        </t-form-item>
      </t-form>
    </t-dialog>
  </div>
</template>
