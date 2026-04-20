<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { createPosition, deletePosition, fetchPositions, updatePosition, type PositionPayload } from '../../api/positions'
import { useCrudPerms } from '../../composables/use-perms'
import type { Position } from '../../types'

const loading = ref(false)
const list = ref<Position[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(10)
const keyword = ref('')

const visible = ref(false)
const editingId = ref<number | null>(null)
const form = reactive<PositionPayload>({
  name: '',
  code: '',
  sort: 0,
  status: 1,
  remark: ''
})

const columns = [
  { colKey: 'id', title: 'ID', width: 70 },
  { colKey: 'name', title: '岗位名称' },
  { colKey: 'code', title: '岗位编码' },
  { colKey: 'sort', title: '排序', width: 90 },
  { colKey: 'status', title: '状态', width: 90 },
  { colKey: 'remark', title: '备注' },
  { colKey: 'actions', title: '操作', width: 160 }
]
const { canCreate, canUpdate, canDelete } = useCrudPerms('position')

function resetForm() {
  form.name = ''
  form.code = ''
  form.sort = 0
  form.status = 1
  form.remark = ''
}

async function loadPositions() {
  loading.value = true
  try {
    const res = await fetchPositions({ page: page.value, size: size.value, keyword: keyword.value || undefined })
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

function openEdit(row: Position) {
  editingId.value = row.id
  form.name = row.name
  form.code = row.code
  form.sort = row.sort
  form.status = row.status
  form.remark = row.remark
  visible.value = true
}

async function submit() {
  try {
    if (editingId.value) {
      await updatePosition(editingId.value, form)
      MessagePlugin.success('更新成功')
    } else {
      await createPosition(form)
      MessagePlugin.success('创建成功')
    }
    visible.value = false
    await loadPositions()
  } catch {
    MessagePlugin.error('保存失败')
  }
}

async function remove(row: Position) {
  try {
    await deletePosition(row.id)
    MessagePlugin.success('删除成功')
    await loadPositions()
  } catch {
    MessagePlugin.error('删除失败')
  }
}

onMounted(loadPositions)
</script>

<template>
  <div class="soft-panel">
    <h2 class="page-title">岗位管理</h2>
    <div class="page-subtitle">维护岗位定义、编码与排序规则。</div>
    <div style="height: 14px" />
    <div class="toolbar">
      <t-input v-model="keyword" placeholder="搜索岗位名称/编码" style="width: 260px" />
      <t-button @click="loadPositions">查询</t-button>
      <div class="grow" />
      <t-button v-if="canCreate" theme="primary" @click="openCreate">新建岗位</t-button>
    </div>

    <div class="data-table-wrap">
      <t-table row-key="id" :data="list" :columns="columns" :loading="loading" bordered>
        <template #status="{ row }">
          <t-tag class="table-status" :theme="row.status === 1 ? 'success' : 'warning'">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </t-tag>
        </template>
        <template #actions="{ row }">
          <t-space v-if="canUpdate || canDelete">
            <t-link v-if="canUpdate" theme="primary" hover="color" @click="openEdit(row)">编辑</t-link>
            <t-link v-if="canDelete" theme="danger" hover="color" @click="remove(row)">删除</t-link>
          </t-space>
          <span v-else>-</span>
        </template>
      </t-table>
    </div>

    <div style="margin-top: 12px; display: flex; justify-content: flex-end">
      <t-pagination v-model:current="page" v-model:page-size="size" :total="total" @change="loadPositions" />
    </div>

    <t-dialog v-model:visible="visible" :header="editingId ? '编辑岗位' : '新建岗位'" width="620px" @confirm="submit">
      <t-form layout="vertical">
        <t-form-item label="岗位名称">
          <t-input v-model="form.name" />
        </t-form-item>
        <t-form-item label="岗位编码">
          <t-input v-model="form.code" />
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
