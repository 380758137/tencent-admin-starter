<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { createNotice, deleteNotice, fetchNotices, updateNotice, type NoticePayload } from '../../api/notices'
import { useCrudPerms } from '../../composables/use-perms'
import type { Notice } from '../../types'
import { noticeTypeLabel } from '../../utils/display'

const loading = ref(false)
const list = ref<Notice[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(10)
const keyword = ref('')

const visible = ref(false)
const editingId = ref<number | null>(null)
const form = reactive<NoticePayload>({
  title: '',
  noticeType: 'notice',
  content: '',
  pinned: 0,
  status: 1
})

const columns = [
  { colKey: 'id', title: 'ID', width: 70 },
  { colKey: 'title', title: '标题' },
  { colKey: 'noticeType', title: '类型', width: 110 },
  { colKey: 'pinned', title: '置顶', width: 90 },
  { colKey: 'status', title: '状态', width: 90 },
  { colKey: 'createdAt', title: '创建时间', width: 180 },
  { colKey: 'actions', title: '操作', width: 160 }
]
const { canCreate, canUpdate, canDelete } = useCrudPerms('notice')

function resetForm() {
  form.title = ''
  form.noticeType = 'notice'
  form.content = ''
  form.pinned = 0
  form.status = 1
}

async function loadNotices() {
  loading.value = true
  try {
    const res = await fetchNotices({ page: page.value, size: size.value, keyword: keyword.value || undefined })
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

function openEdit(row: Notice) {
  editingId.value = row.id
  form.title = row.title
  form.noticeType = row.noticeType
  form.content = row.content
  form.pinned = row.pinned
  form.status = row.status
  visible.value = true
}

async function submit() {
  try {
    if (editingId.value) {
      await updateNotice(editingId.value, form)
      MessagePlugin.success('更新成功')
    } else {
      await createNotice(form)
      MessagePlugin.success('创建成功')
    }
    visible.value = false
    await loadNotices()
  } catch {
    MessagePlugin.error('保存失败')
  }
}

async function remove(row: Notice) {
  try {
    await deleteNotice(row.id)
    MessagePlugin.success('删除成功')
    await loadNotices()
  } catch {
    MessagePlugin.error('删除失败')
  }
}

onMounted(loadNotices)
</script>

<template>
  <div class="soft-panel">
    <h2 class="page-title">通知公告</h2>
    <div class="page-subtitle">维护系统公告、通知类型与置顶状态。</div>
    <div style="height: 14px" />
    <div class="toolbar">
      <t-input v-model="keyword" placeholder="搜索标题/内容" style="width: 260px" />
      <t-button @click="loadNotices">查询</t-button>
      <div class="grow" />
      <t-button v-if="canCreate" theme="primary" @click="openCreate">新建公告</t-button>
    </div>

    <div class="data-table-wrap">
      <t-table row-key="id" :data="list" :columns="columns" :loading="loading" bordered>
        <template #noticeType="{ row }">
          <t-tag variant="outline" theme="primary">{{ noticeTypeLabel(row.noticeType) }}</t-tag>
        </template>
        <template #pinned="{ row }">
          <t-tag :theme="row.pinned === 1 ? 'warning' : 'default'" variant="outline">{{ row.pinned === 1 ? '是' : '否' }}</t-tag>
        </template>
        <template #status="{ row }">
          <t-tag class="table-status" :theme="row.status === 1 ? 'success' : 'warning'">
            {{ row.status === 1 ? '发布' : '下线' }}
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
      <t-pagination v-model:current="page" v-model:page-size="size" :total="total" @change="loadNotices" />
    </div>

    <t-dialog v-model:visible="visible" :header="editingId ? '编辑公告' : '新建公告'" width="680px" @confirm="submit">
      <t-form layout="vertical">
        <t-form-item label="标题">
          <t-input v-model="form.title" />
        </t-form-item>
        <t-form-item label="类型">
          <t-select v-model="form.noticeType">
            <t-option value="notice" label="通知" />
            <t-option value="announcement" label="公告" />
          </t-select>
        </t-form-item>
        <t-form-item label="内容">
          <t-textarea v-model="form.content" :autosize="{ minRows: 4, maxRows: 8 }" />
        </t-form-item>
        <t-form-item label="置顶">
          <t-select v-model="form.pinned">
            <t-option :value="1" label="是" />
            <t-option :value="0" label="否" />
          </t-select>
        </t-form-item>
        <t-form-item label="状态">
          <t-select v-model="form.status">
            <t-option :value="1" label="发布" />
            <t-option :value="0" label="下线" />
          </t-select>
        </t-form-item>
      </t-form>
    </t-dialog>
  </div>
</template>
