<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { createUser, deleteUser, exportUsers, fetchUsers, importUsers, updateUser, type UserPayload } from '../../api/users'
import { fetchRoles } from '../../api/roles'
import { useAuthStore } from '../../stores/auth'
import type { Role, User } from '../../types'
import { dataScopeLabel, roleLabel } from '../../utils/display'

const auth = useAuthStore()
const loading = ref(false)
const users = ref<User[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(10)
const keyword = ref('')

const visible = ref(false)
const editingId = ref<number | null>(null)
const form = reactive<UserPayload>({
  username: '',
  password: '',
  displayName: '',
  role: 'operator',
  deptId: 0,
  status: 1
})
const roleCatalog = ref<Role[]>([
  {
    id: 0,
    name: '超级管理员',
    roleKey: 'admin',
    permissions: '*',
    dataScope: 'all',
    status: 1,
    remark: '',
    createdAt: '',
    updatedAt: ''
  },
  {
    id: 0,
    name: '操作员',
    roleKey: 'operator',
    permissions: 'user:list,department:list,position:list,notice:list,online-user:list',
    dataScope: 'dept',
    status: 1,
    remark: '',
    createdAt: '',
    updatedAt: ''
  }
])
const roleKeyword = ref('')
const assignRoleVisible = ref(false)
const draftRoles = ref<string[]>(['operator'])
const fileInput = ref<HTMLInputElement | null>(null)
const canCreate = computed(() => auth.hasPerm('user:create'))
const canUpdate = computed(() => auth.hasPerm('user:update'))
const canDelete = computed(() => auth.hasPerm('user:delete'))
const canImport = computed(() => auth.hasPerm('user:import'))
const canExport = computed(() => auth.hasPerm('user:export'))
const dataScopeText = computed(() => {
  if (auth.user?.dataScope === 'all') return '全部数据'
  if (auth.user?.dataScope === 'dept') return '本部门及子部门数据'
  return '仅本人数据'
})
const filteredRoles = computed(() => {
  const kw = roleKeyword.value.trim().toLowerCase()
  if (!kw) return roleCatalog.value
  return roleCatalog.value.filter((item) => item.name.toLowerCase().includes(kw) || item.roleKey.toLowerCase().includes(kw))
})
const selectedRoles = computed(() => {
  const selected = new Set(splitRoleKeys(form.role))
  return roleCatalog.value.filter((item) => selected.has(item.roleKey))
})
const draftSelectedRoles = computed(() => {
  const selected = new Set(draftRoles.value)
  return roleCatalog.value.filter((item) => selected.has(item.roleKey))
})

const columns = [
  { colKey: 'id', title: 'ID', width: 70 },
  { colKey: 'username', title: '用户名' },
  { colKey: 'displayName', title: '显示名' },
  { colKey: 'role', title: '角色' },
  { colKey: 'deptId', title: '部门ID', width: 90 },
  { colKey: 'status', title: '状态' },
  { colKey: 'actions', title: '操作', width: 180 }
]
const roleAssignColumns = [
  { colKey: 'name', title: '角色名称' },
  { colKey: 'roleKey', title: '角色标识', width: 140 },
  { colKey: 'dataScope', title: '数据范围', width: 140 },
  { colKey: 'permissions', title: '权限数', width: 100 },
  { colKey: 'actions', title: '分配', width: 120 }
]

function resetForm() {
  form.username = ''
  form.password = ''
  form.displayName = ''
  form.role = 'operator'
  form.deptId = 0
  form.status = 1
  draftRoles.value = ['operator']
  roleKeyword.value = ''
}

async function loadUsers() {
  loading.value = true
  try {
    const res = await fetchUsers({ page: page.value, size: size.value, keyword: keyword.value || undefined })
    users.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

async function loadRoleOptions() {
  try {
    const res = await fetchRoles({ page: 1, size: 200 })
    if (res.list.length === 0) return
    roleCatalog.value = res.list
  } catch {
    // Keep built-in fallback options when role list loading fails.
  }
}

function rolePermissionCount(raw: string): number {
  const value = raw.trim()
  if (!value) return 0
  if (value === '*') return 1
  return value.split(',').map((item) => item.trim()).filter(Boolean).length
}

function openCreate() {
  editingId.value = null
  resetForm()
  visible.value = true
}

function openEdit(row: User) {
  editingId.value = row.id
  form.username = row.username
  form.password = ''
  form.displayName = row.displayName
  form.role = row.role
  draftRoles.value = splitRoleKeys(row.role)
  form.deptId = row.deptId
  form.status = row.status
  visible.value = true
}

function openAssignRoleDialog() {
  draftRoles.value = splitRoleKeys(form.role)
  roleKeyword.value = ''
  assignRoleVisible.value = true
}

function chooseDraftRole(roleKey: string) {
  if (draftRoles.value.includes(roleKey)) {
    draftRoles.value = draftRoles.value.filter((item) => item !== roleKey)
    return
  }
  draftRoles.value = [...draftRoles.value, roleKey]
}

function confirmAssignRole() {
  const normalized = joinRoleKeys(draftRoles.value)
  if (!normalized) {
    MessagePlugin.warning('至少选择一个角色')
    return
  }
  form.role = normalized
  assignRoleVisible.value = false
}

function splitRoleKeys(raw: string): string[] {
  return raw
    .split(',')
    .map((item) => item.trim())
    .filter(Boolean)
}

function joinRoleKeys(keys: string[]): string {
  return Array.from(new Set(keys.map((item) => item.trim()).filter(Boolean))).join(',')
}

function splitRoles(raw: string): string[] {
  return raw.split(',').map((item) => item.trim()).filter(Boolean)
}

function hasDraftRole(roleKey: string): boolean {
  return draftRoles.value.includes(roleKey)
}

async function submit() {
  try {
    if (editingId.value) {
      await updateUser(editingId.value, form)
      MessagePlugin.success('更新成功')
    } else {
      if (!form.username) {
        MessagePlugin.warning('新建用户必须填写用户名')
        return
      }
      await createUser(form)
      MessagePlugin.success('创建成功')
    }
    visible.value = false
    await loadUsers()
  } catch {
    MessagePlugin.error('保存失败')
  }
}

async function remove(row: User) {
  try {
    await deleteUser(row.id)
    MessagePlugin.success('删除成功')
    await loadUsers()
  } catch {
    MessagePlugin.error('删除失败')
  }
}

async function doExport() {
  try {
    const blob = await exportUsers()
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'users.csv'
    a.click()
    window.URL.revokeObjectURL(url)
    MessagePlugin.success('导出成功')
  } catch {
    MessagePlugin.error('导出失败')
  }
}

function triggerImport() {
  fileInput.value?.click()
}

async function onImportFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  try {
    const result = await importUsers(file)
    MessagePlugin.success(`导入完成：成功 ${result.imported}，跳过 ${result.skipped}，失败 ${result.failed}`)
    await loadUsers()
  } catch {
    MessagePlugin.error('导入失败')
  } finally {
    input.value = ''
  }
}

onMounted(async () => {
  await Promise.all([loadUsers(), loadRoleOptions()])
})
</script>

<template>
  <div class="soft-panel">
    <h2 class="page-title">用户管理</h2>
    <div class="page-subtitle">统一管理账号、角色和启禁状态。当前数据范围：{{ dataScopeText }}</div>
    <div style="height: 14px" />
    <div class="toolbar">
      <t-input v-model="keyword" placeholder="搜索用户名/显示名" style="width: 260px" />
      <t-button @click="loadUsers">查询</t-button>
      <div class="grow" />
      <t-button v-if="canImport" variant="outline" @click="triggerImport">导入CSV</t-button>
      <t-button v-if="canExport" variant="outline" @click="doExport">导出CSV</t-button>
      <t-button v-if="canCreate" theme="primary" @click="openCreate">新建用户</t-button>
      <input ref="fileInput" type="file" accept=".csv,text/csv" style="display: none" @change="onImportFileChange" />
    </div>

    <div class="data-table-wrap">
      <t-table row-key="id" :data="users" :columns="columns" :loading="loading" bordered>
      <template #role="{ row }">
        <t-space size="4px">
          <t-tag
            v-for="item in splitRoles(row.role)"
            :key="`${row.id}-${item}`"
            variant="outline"
            theme="primary"
          >
            {{ roleLabel(item) }}
          </t-tag>
        </t-space>
      </template>
      <template #status="{ row }">
        <t-tag class="table-status" :theme="row.status === 1 ? 'success' : 'warning'">{{ row.status === 1 ? '启用' : '禁用' }}</t-tag>
      </template>
        <template #actions="{ row }">
          <t-space>
          <t-link v-if="canUpdate" theme="primary" hover="color" @click="openEdit(row)">编辑</t-link>
          <t-link v-if="canDelete" theme="danger" hover="color" @click="remove(row)">删除</t-link>
        </t-space>
      </template>
    </t-table>
    </div>

    <div style="margin-top: 12px; display: flex; justify-content: flex-end">
      <t-pagination
        v-model:current="page"
        v-model:page-size="size"
        :total="total"
        @change="loadUsers"
      />
    </div>

    <t-dialog v-model:visible="visible" :header="editingId ? '编辑用户' : '新建用户'" width="560px" @confirm="submit">
      <t-form layout="vertical">
        <t-form-item label="用户名">
          <t-input v-model="form.username" :disabled="!!editingId" />
        </t-form-item>
        <t-form-item label="密码（留空则不改）">
          <t-input v-model="form.password" type="password" />
        </t-form-item>
        <t-form-item label="显示名">
          <t-input v-model="form.displayName" />
        </t-form-item>
        <t-form-item label="角色">
          <t-space align="center">
            <t-space size="6px">
              <t-tag v-for="item in selectedRoles" :key="item.roleKey" theme="primary" variant="outline">
                {{ item.name }}（{{ item.roleKey }}）
              </t-tag>
              <t-tag v-if="selectedRoles.length === 0" theme="warning" variant="outline">未分配角色</t-tag>
            </t-space>
            <t-button variant="outline" @click="openAssignRoleDialog">分配角色</t-button>
          </t-space>
        </t-form-item>
        <t-form-item label="状态">
          <t-select v-model="form.status">
            <t-option :value="1" label="启用" />
            <t-option :value="0" label="禁用" />
          </t-select>
        </t-form-item>
        <t-form-item label="部门ID">
          <t-input-number v-model="form.deptId" />
        </t-form-item>
      </t-form>
    </t-dialog>

    <t-dialog v-model:visible="assignRoleVisible" header="分配角色" width="860px" @confirm="confirmAssignRole">
      <div class="toolbar" style="margin-bottom: 12px">
        <t-input v-model="roleKeyword" placeholder="搜索角色名称/标识" clearable style="width: 320px" />
        <div class="grow" />
        <t-tag theme="primary" variant="outline">
          待分配角色数：{{ draftSelectedRoles.length }}
        </t-tag>
      </div>

      <div class="data-table-wrap">
        <t-table row-key="roleKey" :data="filteredRoles" :columns="roleAssignColumns" bordered>
          <template #dataScope="{ row }">
            {{ dataScopeLabel(row.dataScope) }}
          </template>
          <template #permissions="{ row }">
            {{ rolePermissionCount(row.permissions) }}
          </template>
          <template #actions="{ row }">
            <t-button
              size="small"
              :theme="hasDraftRole(row.roleKey) ? 'primary' : 'default'"
              :variant="hasDraftRole(row.roleKey) ? 'base' : 'outline'"
              @click="chooseDraftRole(row.roleKey)"
            >
              {{ hasDraftRole(row.roleKey) ? '已选择' : '选择' }}
            </t-button>
          </template>
        </t-table>
      </div>
    </t-dialog>
  </div>
</template>
