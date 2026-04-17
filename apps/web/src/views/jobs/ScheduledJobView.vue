<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import {
  createScheduledJob,
  deleteScheduledJob,
  fetchScheduledJobLogs,
  fetchScheduledJobs,
  runScheduledJob,
  updateScheduledJob,
  type ScheduledJobPayload
} from '../../api/scheduled-jobs'
import type { ScheduledJob, ScheduledJobLog } from '../../types'
import { jobResultLabel, triggerTypeLabel } from '../../utils/display'

const loading = ref(false)
const list = ref<ScheduledJob[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(10)
const keyword = ref('')

const logLoading = ref(false)
const logs = ref<ScheduledJobLog[]>([])
const logTotal = ref(0)
const logPage = ref(1)
const logSize = ref(8)
const selectedJobId = ref<number | undefined>()

const visible = ref(false)
const editingId = ref<number | null>(null)
const form = reactive<ScheduledJobPayload>({
  name: '',
  cronExpr: '0 */5 * * * *',
  command: 'echo hello',
  status: 1,
  remark: ''
})

const jobColumns = [
  { colKey: 'id', title: 'ID', width: 70 },
  { colKey: 'name', title: '任务名称' },
  { colKey: 'cronExpr', title: 'Cron 表达式' },
  { colKey: 'command', title: '命令' },
  { colKey: 'status', title: '状态', width: 90 },
  { colKey: 'lastResult', title: '最近结果', width: 100 },
  { colKey: 'lastRunAt', title: '最近执行时间', width: 180 },
  { colKey: 'actions', title: '操作', width: 220 }
]

const logColumns = [
  { colKey: 'id', title: 'ID', width: 70 },
  { colKey: 'jobId', title: '任务ID', width: 90 },
  { colKey: 'triggerType', title: '触发方式', width: 100 },
  { colKey: 'status', title: '结果', width: 100 },
  { colKey: 'message', title: '消息' },
  { colKey: 'runAt', title: '执行时间', width: 180 }
]

function resetForm() {
  form.name = ''
  form.cronExpr = '0 */5 * * * *'
  form.command = 'echo hello'
  form.status = 1
  form.remark = ''
}

async function loadJobs() {
  loading.value = true
  try {
    const res = await fetchScheduledJobs({ page: page.value, size: size.value, keyword: keyword.value || undefined })
    list.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

async function loadLogs() {
  logLoading.value = true
  try {
    const res = await fetchScheduledJobLogs({
      page: logPage.value,
      size: logSize.value,
      jobId: selectedJobId.value
    })
    logs.value = res.list
    logTotal.value = res.total
  } finally {
    logLoading.value = false
  }
}

function openCreate() {
  editingId.value = null
  resetForm()
  visible.value = true
}

function openEdit(row: ScheduledJob) {
  editingId.value = row.id
  form.name = row.name
  form.cronExpr = row.cronExpr
  form.command = row.command
  form.status = row.status
  form.remark = row.remark
  visible.value = true
}

async function submit() {
  try {
    if (editingId.value) {
      await updateScheduledJob(editingId.value, form)
      MessagePlugin.success('更新成功')
    } else {
      await createScheduledJob(form)
      MessagePlugin.success('创建成功')
    }
    visible.value = false
    await loadJobs()
  } catch {
    MessagePlugin.error('保存失败')
  }
}

async function remove(row: ScheduledJob) {
  try {
    await deleteScheduledJob(row.id)
    MessagePlugin.success('删除成功')
    await loadJobs()
    await loadLogs()
  } catch {
    MessagePlugin.error('删除失败')
  }
}

async function runNow(row: ScheduledJob) {
  try {
    await runScheduledJob(row.id)
    MessagePlugin.success('触发成功')
    selectedJobId.value = row.id
    await Promise.all([loadJobs(), loadLogs()])
  } catch {
    MessagePlugin.error('触发失败')
  }
}

function filterLogsByJob(row: ScheduledJob) {
  selectedJobId.value = row.id
  logPage.value = 1
  loadLogs()
}

function clearLogFilter() {
  selectedJobId.value = undefined
  logPage.value = 1
  loadLogs()
}

onMounted(async () => {
  await Promise.all([loadJobs(), loadLogs()])
})
</script>

<template>
  <div class="soft-panel">
    <h2 class="page-title">定时任务</h2>
    <div class="page-subtitle">维护任务配置，支持手动触发并追踪执行日志。</div>
    <div style="height: 14px" />

    <div class="toolbar">
      <t-input v-model="keyword" placeholder="搜索任务名/命令" style="width: 280px" />
      <t-button @click="loadJobs">查询</t-button>
      <div class="grow" />
      <t-button theme="primary" @click="openCreate">新建任务</t-button>
    </div>

    <div class="data-table-wrap">
      <t-table row-key="id" :data="list" :columns="jobColumns" :loading="loading" bordered>
        <template #status="{ row }">
          <t-tag class="table-status" :theme="row.status === 1 ? 'success' : 'warning'">
            {{ row.status === 1 ? '启用' : '停用' }}
          </t-tag>
        </template>
        <template #lastResult="{ row }">
          <t-tag :theme="row.lastResult === 'success' ? 'success' : row.lastResult ? 'danger' : 'default'" variant="outline">
            {{ jobResultLabel(row.lastResult) }}
          </t-tag>
        </template>
        <template #actions="{ row }">
          <t-space>
            <t-link theme="primary" hover="color" @click="openEdit(row)">编辑</t-link>
            <t-link theme="success" hover="color" @click="runNow(row)">执行</t-link>
            <t-link theme="warning" hover="color" @click="filterLogsByJob(row)">日志</t-link>
            <t-link theme="danger" hover="color" @click="remove(row)">删除</t-link>
          </t-space>
        </template>
      </t-table>
    </div>

    <div style="margin-top: 12px; display: flex; justify-content: flex-end">
      <t-pagination v-model:current="page" v-model:page-size="size" :total="total" @change="loadJobs" />
    </div>

    <div style="height: 18px" />
    <div class="toolbar">
      <div style="font-size: 15px; font-weight: 600">执行日志</div>
      <t-tag v-if="selectedJobId" theme="primary" variant="outline">当前筛选任务ID: {{ selectedJobId }}</t-tag>
      <div class="grow" />
      <t-button variant="outline" @click="clearLogFilter">清空筛选</t-button>
    </div>
    <div class="data-table-wrap">
      <t-table row-key="id" :data="logs" :columns="logColumns" :loading="logLoading" bordered>
        <template #triggerType="{ row }">
          {{ triggerTypeLabel(row.triggerType) }}
        </template>
        <template #status="{ row }">
          <t-tag :theme="row.status === 'success' ? 'success' : 'danger'" variant="outline">{{ jobResultLabel(row.status) }}</t-tag>
        </template>
      </t-table>
    </div>
    <div style="margin-top: 12px; display: flex; justify-content: flex-end">
      <t-pagination v-model:current="logPage" v-model:page-size="logSize" :total="logTotal" @change="loadLogs" />
    </div>

    <t-dialog v-model:visible="visible" :header="editingId ? '编辑任务' : '新建任务'" width="680px" @confirm="submit">
      <t-form layout="vertical">
        <t-form-item label="任务名称">
          <t-input v-model="form.name" />
        </t-form-item>
        <t-form-item label="Cron 表达式">
          <t-input v-model="form.cronExpr" />
        </t-form-item>
        <t-form-item label="命令">
          <t-input v-model="form.command" />
        </t-form-item>
        <t-form-item label="状态">
          <t-select v-model="form.status">
            <t-option :value="1" label="启用" />
            <t-option :value="0" label="停用" />
          </t-select>
        </t-form-item>
        <t-form-item label="备注">
          <t-input v-model="form.remark" />
        </t-form-item>
      </t-form>
    </t-dialog>
  </div>
</template>
