<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { createMenu, deleteMenu, fetchMenus, updateMenu, type MenuPayload } from '../../api/menus'
import type { Menu } from '../../types'
import { menuTypeLabel } from '../../utils/display'

const loading = ref(false)
const list = ref<Menu[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(10)
const keyword = ref('')

const visible = ref(false)
const editingId = ref<number | null>(null)
const form = reactive<MenuPayload>({
  parentId: 0,
  name: '',
  menuType: 'menu',
  path: '',
  component: '',
  perms: '',
  sort: 0,
  status: 1
})

const columns = [
  { colKey: 'id', title: 'ID', width: 70 },
  { colKey: 'parentId', title: '父级ID', width: 90 },
  { colKey: 'name', title: '菜单名称' },
  { colKey: 'menuType', title: '类型', width: 90 },
  { colKey: 'path', title: '路由路径' },
  { colKey: 'perms', title: '权限标识' },
  { colKey: 'sort', title: '排序', width: 90 },
  { colKey: 'status', title: '状态', width: 90 },
  { colKey: 'actions', title: '操作', width: 160 }
]

function resetForm() {
  form.parentId = 0
  form.name = ''
  form.menuType = 'menu'
  form.path = ''
  form.component = ''
  form.perms = ''
  form.sort = 0
  form.status = 1
}

async function loadMenus() {
  loading.value = true
  try {
    const res = await fetchMenus({ page: page.value, size: size.value, keyword: keyword.value || undefined })
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

function openEdit(row: Menu) {
  editingId.value = row.id
  form.parentId = row.parentId
  form.name = row.name
  form.menuType = row.menuType
  form.path = row.path
  form.component = row.component
  form.perms = row.perms
  form.sort = row.sort
  form.status = row.status
  visible.value = true
}

async function submit() {
  try {
    if (editingId.value) {
      await updateMenu(editingId.value, form)
      MessagePlugin.success('更新成功')
    } else {
      await createMenu(form)
      MessagePlugin.success('创建成功')
    }
    visible.value = false
    await loadMenus()
  } catch {
    MessagePlugin.error('保存失败')
  }
}

async function remove(row: Menu) {
  try {
    await deleteMenu(row.id)
    MessagePlugin.success('删除成功')
    await loadMenus()
  } catch {
    MessagePlugin.error('删除失败')
  }
}

onMounted(loadMenus)
</script>

<template>
  <div class="soft-panel">
    <h2 class="page-title">菜单管理</h2>
    <div class="page-subtitle">维护路由菜单、按钮权限与显示顺序。</div>
    <div style="height: 14px" />

    <div class="toolbar">
      <t-input v-model="keyword" placeholder="搜索菜单名称/路径/权限" style="width: 280px" />
      <t-button @click="loadMenus">查询</t-button>
      <div class="grow" />
      <t-button theme="primary" @click="openCreate">新建菜单</t-button>
    </div>

    <div class="data-table-wrap">
      <t-table row-key="id" :data="list" :columns="columns" :loading="loading" bordered>
        <template #menuType="{ row }">
          <t-tag :theme="row.menuType === 'button' ? 'warning' : 'primary'" variant="outline">{{ menuTypeLabel(row.menuType) }}</t-tag>
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
      <t-pagination v-model:current="page" v-model:page-size="size" :total="total" @change="loadMenus" />
    </div>

    <t-dialog v-model:visible="visible" :header="editingId ? '编辑菜单' : '新建菜单'" width="620px" @confirm="submit">
      <t-form layout="vertical">
        <t-form-item label="父级ID">
          <t-input-number v-model="form.parentId" />
        </t-form-item>
        <t-form-item label="菜单名称">
          <t-input v-model="form.name" />
        </t-form-item>
        <t-form-item label="菜单类型">
          <t-select v-model="form.menuType">
            <t-option value="menu" label="菜单" />
            <t-option value="button" label="按钮" />
          </t-select>
        </t-form-item>
        <t-form-item label="路由路径">
          <t-input v-model="form.path" />
        </t-form-item>
        <t-form-item label="组件路径">
          <t-input v-model="form.component" />
        </t-form-item>
        <t-form-item label="权限标识">
          <t-input v-model="form.perms" />
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
      </t-form>
    </t-dialog>
  </div>
</template>
