<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import {
  createDictionaryItem,
  deleteDictionaryItem,
  fetchDictionaryItems,
  updateDictionaryItem,
  type DictionaryPayload
} from '../../api/dictionary'
import type { DictionaryItem } from '../../types'

const loading = ref(false)
const list = ref<DictionaryItem[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(10)
const keyword = ref('')
const dictTypeQuery = ref('')

const visible = ref(false)
const editingId = ref<number | null>(null)
const form = reactive<DictionaryPayload>({
  dictType: '',
  dictLabel: '',
  dictValue: '',
  sort: 0,
  status: 1,
  remark: ''
})

const columns = [
  { colKey: 'id', title: 'ID', width: 70 },
  { colKey: 'dictType', title: '字典类型' },
  { colKey: 'dictLabel', title: '字典标签' },
  { colKey: 'dictValue', title: '字典键值' },
  { colKey: 'sort', title: '排序', width: 90 },
  { colKey: 'status', title: '状态', width: 90 },
  { colKey: 'remark', title: '备注' },
  { colKey: 'actions', title: '操作', width: 160 }
]

function resetForm() {
  form.dictType = ''
  form.dictLabel = ''
  form.dictValue = ''
  form.sort = 0
  form.status = 1
  form.remark = ''
}

async function loadItems() {
  loading.value = true
  try {
    const res = await fetchDictionaryItems({
      page: page.value,
      size: size.value,
      keyword: keyword.value || undefined,
      dictType: dictTypeQuery.value || undefined
    })
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

function openEdit(row: DictionaryItem) {
  editingId.value = row.id
  form.dictType = row.dictType
  form.dictLabel = row.dictLabel
  form.dictValue = row.dictValue
  form.sort = row.sort
  form.status = row.status
  form.remark = row.remark
  visible.value = true
}

async function submit() {
  try {
    if (editingId.value) {
      await updateDictionaryItem(editingId.value, form)
      MessagePlugin.success('更新成功')
    } else {
      await createDictionaryItem(form)
      MessagePlugin.success('创建成功')
    }
    visible.value = false
    await loadItems()
  } catch {
    MessagePlugin.error('保存失败')
  }
}

async function remove(row: DictionaryItem) {
  try {
    await deleteDictionaryItem(row.id)
    MessagePlugin.success('删除成功')
    await loadItems()
  } catch {
    MessagePlugin.error('删除失败')
  }
}

onMounted(loadItems)
</script>

<template>
  <div class="soft-panel">
    <h2 class="page-title">字典管理</h2>
    <div class="page-subtitle">维护系统枚举与可配置选项。</div>
    <div style="height: 14px" />

    <div class="toolbar">
      <t-input v-model="dictTypeQuery" placeholder="按字典类型过滤" style="width: 220px" />
      <t-input v-model="keyword" placeholder="搜索标签/键值" style="width: 220px" />
      <t-button @click="loadItems">查询</t-button>
      <div class="grow" />
      <t-button theme="primary" @click="openCreate">新增字典项</t-button>
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
      <t-pagination v-model:current="page" v-model:page-size="size" :total="total" @change="loadItems" />
    </div>

    <t-dialog v-model:visible="visible" :header="editingId ? '编辑字典项' : '新增字典项'" width="620px" @confirm="submit">
      <t-form layout="vertical">
        <t-form-item label="字典类型">
          <t-input v-model="form.dictType" />
        </t-form-item>
        <t-form-item label="字典标签">
          <t-input v-model="form.dictLabel" />
        </t-form-item>
        <t-form-item label="字典键值">
          <t-input v-model="form.dictValue" />
        </t-form-item>
        <t-form-item label="排序">
          <t-input-number v-model="form.sort" />
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
