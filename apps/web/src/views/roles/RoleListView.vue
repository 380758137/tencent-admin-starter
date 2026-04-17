<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { createRole, deleteRole, fetchRoles, updateRole, type RolePayload } from '../../api/roles'
import type { Role } from '../../types'
import { dataScopeLabel } from '../../utils/display'

const loading = ref(false)
const list = ref<Role[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(10)
const keyword = ref('')

const visible = ref(false)
const editingId = ref<number | null>(null)
const form = reactive<RolePayload>({
  name: '',
  roleKey: '',
  permissions: '',
  dataScope: 'self',
  status: 1,
  remark: ''
})
const selectedPermissions = ref<string[]>([])
const customPermission = ref('')

const presetPermissionOptions = [
  { value: 'user:list', label: '用户查看' },
  { value: 'user:create', label: '用户创建' },
  { value: 'user:update', label: '用户编辑' },
  { value: 'user:delete', label: '用户删除' },
  { value: 'user:import', label: '用户导入' },
  { value: 'user:export', label: '用户导出' },
  { value: 'department:list', label: '部门查看' },
  { value: 'department:create', label: '部门创建' },
  { value: 'department:update', label: '部门编辑' },
  { value: 'department:delete', label: '部门删除' },
  { value: 'role:list', label: '角色查看' },
  { value: 'role:create', label: '角色创建' },
  { value: 'role:update', label: '角色编辑' },
  { value: 'role:delete', label: '角色删除' },
  { value: 'menu:list', label: '菜单查看' },
  { value: 'menu:create', label: '菜单创建' },
  { value: 'menu:update', label: '菜单编辑' },
  { value: 'menu:delete', label: '菜单删除' },
  { value: 'dictionary:list', label: '字典查看' },
  { value: 'dictionary:create', label: '字典创建' },
  { value: 'dictionary:update', label: '字典编辑' },
  { value: 'dictionary:delete', label: '字典删除' },
  { value: 'param:list', label: '参数查看' },
  { value: 'param:create', label: '参数创建' },
  { value: 'param:update', label: '参数编辑' },
  { value: 'param:delete', label: '参数删除' },
  { value: 'log:operation:list', label: '操作日志查看' },
  { value: 'log:login:list', label: '登录日志查看' },
  { value: 'monitor:view', label: '系统监控查看' },
  { value: 'position:list', label: '岗位查看' },
  { value: 'position:create', label: '岗位创建' },
  { value: 'position:update', label: '岗位编辑' },
  { value: 'position:delete', label: '岗位删除' },
  { value: 'notice:list', label: '公告查看' },
  { value: 'notice:create', label: '公告创建' },
  { value: 'notice:update', label: '公告编辑' },
  { value: 'notice:delete', label: '公告删除' },
  { value: 'online-user:list', label: '在线用户查看' },
  { value: 'job:list', label: '任务查看' },
  { value: 'job:create', label: '任务创建' },
  { value: 'job:update', label: '任务编辑' },
  { value: 'job:delete', label: '任务删除' },
  { value: 'job:run', label: '任务执行' },
  { value: 'job-log:list', label: '任务日志查看' }
]

const columns = [
  { colKey: 'id', title: 'ID', width: 70 },
  { colKey: 'name', title: '角色名称' },
  { colKey: 'roleKey', title: '角色标识' },
  { colKey: 'dataScope', title: '数据范围', width: 100 },
  { colKey: 'permissions', title: '权限码' },
  { colKey: 'status', title: '状态', width: 90 },
  { colKey: 'remark', title: '备注' },
  { colKey: 'actions', title: '操作', width: 160 }
]

function resetForm() {
  form.name = ''
  form.roleKey = ''
  form.permissions = ''
  form.dataScope = 'self'
  form.status = 1
  form.remark = ''
  selectedPermissions.value = []
  customPermission.value = ''
}

function splitPermissions(raw: string): string[] {
  return raw
    .split(',')
    .map((item) => item.trim())
    .filter(Boolean)
}

function permissionPreview(raw: string): string[] {
  return splitPermissions(raw).slice(0, 3)
}

function permissionTotal(raw: string): number {
  return splitPermissions(raw).length
}

function syncPermissionsToForm() {
  form.permissions = selectedPermissions.value.join(',')
}

async function loadRoles() {
  loading.value = true
  try {
    const res = await fetchRoles({ page: page.value, size: size.value, keyword: keyword.value || undefined })
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

function openEdit(row: Role) {
  editingId.value = row.id
  form.name = row.name
  form.roleKey = row.roleKey
  form.permissions = row.permissions
  selectedPermissions.value = splitPermissions(row.permissions)
  form.dataScope = row.dataScope
  form.status = row.status
  form.remark = row.remark
  customPermission.value = ''
  visible.value = true
}

function addCustomPermission() {
  const value = customPermission.value.trim()
  if (!value) return
  if (!selectedPermissions.value.includes(value)) {
    selectedPermissions.value = [...selectedPermissions.value, value]
  }
  customPermission.value = ''
}

function removePermission(code: string) {
  selectedPermissions.value = selectedPermissions.value.filter((item) => item !== code)
}

async function submit() {
  try {
    syncPermissionsToForm()
    if (editingId.value) {
      await updateRole(editingId.value, {
        name: form.name,
        roleKey: form.roleKey,
        permissions: form.permissions,
        dataScope: form.dataScope,
        status: form.status,
        remark: form.remark
      })
      MessagePlugin.success('更新成功')
    } else {
      await createRole(form)
      MessagePlugin.success('创建成功')
    }
    visible.value = false
    await loadRoles()
  } catch {
    MessagePlugin.error('保存失败')
  }
}

async function remove(row: Role) {
  try {
    await deleteRole(row.id)
    MessagePlugin.success('删除成功')
    await loadRoles()
  } catch {
    MessagePlugin.error('删除失败')
  }
}

onMounted(loadRoles)
</script>

<template>
  <div class="soft-panel">
    <h2 class="page-title">角色管理</h2>
    <div class="page-subtitle">管理系统角色定义与启禁状态。</div>
    <div style="height: 14px" />

    <div class="toolbar">
      <t-input v-model="keyword" placeholder="搜索角色名称/标识" style="width: 260px" />
      <t-button @click="loadRoles">查询</t-button>
      <div class="grow" />
      <t-button theme="primary" @click="openCreate">新建角色</t-button>
    </div>

    <div class="data-table-wrap">
      <t-table row-key="id" :data="list" :columns="columns" :loading="loading" bordered>
        <template #dataScope="{ row }">
          <t-tag variant="outline" theme="primary">{{ dataScopeLabel(row.dataScope) }}</t-tag>
        </template>
        <template #permissions="{ row }">
          <t-space size="4px">
            <t-tag v-for="item in permissionPreview(row.permissions)" :key="item" variant="outline">{{ item }}</t-tag>
            <t-tag v-if="permissionTotal(row.permissions) > 3" theme="default" variant="light">
              +{{ permissionTotal(row.permissions) - 3 }}
            </t-tag>
            <span v-if="permissionTotal(row.permissions) === 0">-</span>
          </t-space>
        </template>
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
      <t-pagination v-model:current="page" v-model:page-size="size" :total="total" @change="loadRoles" />
    </div>

    <t-dialog v-model:visible="visible" :header="editingId ? '编辑角色' : '新建角色'" width="760px" @confirm="submit">
      <t-form layout="vertical">
        <t-form-item label="角色名称">
          <t-input v-model="form.name" />
        </t-form-item>
        <t-form-item label="角色标识">
          <t-input v-model="form.roleKey" />
        </t-form-item>
        <t-form-item label="权限配置">
          <t-select
            v-model="selectedPermissions"
            :options="presetPermissionOptions"
            multiple
            filterable
            clearable
            placeholder="请选择权限码（支持搜索与多选）"
          />
          <div style="display: flex; gap: 8px; margin-top: 10px">
            <t-input v-model="customPermission" placeholder="追加自定义权限码，如 user:reset-password" />
            <t-button variant="outline" @click="addCustomPermission">添加</t-button>
          </div>
          <t-space break-line size="4px" style="margin-top: 10px">
            <t-tag v-for="item in selectedPermissions" :key="item" closable @close="removePermission(item)">
              {{ item }}
            </t-tag>
            <span v-if="selectedPermissions.length === 0" style="color: var(--td-text-color-placeholder)">未选择权限</span>
          </t-space>
        </t-form-item>
        <t-form-item label="数据范围">
          <t-select v-model="form.dataScope">
            <t-option value="all" label="全部数据" />
            <t-option value="dept" label="本部门及子部门" />
            <t-option value="self" label="仅本人" />
          </t-select>
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
