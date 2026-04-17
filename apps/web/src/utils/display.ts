export function roleLabel(roleKey?: string): string {
  if (!roleKey) return '-'
  if (roleKey === 'admin') return '超级管理员'
  if (roleKey === 'operator') return '操作员'
  return roleKey
}

export function dataScopeLabel(scope?: string): string {
  if (scope === 'all') return '全部数据'
  if (scope === 'dept') return '本部门及子部门'
  return '仅本人'
}

export function menuTypeLabel(menuType?: string): string {
  if (menuType === 'button') return '按钮'
  if (menuType === 'menu') return '菜单'
  return menuType || '-'
}

export function noticeTypeLabel(noticeType?: string): string {
  if (noticeType === 'announcement') return '公告'
  if (noticeType === 'notice') return '通知'
  return noticeType || '-'
}

export function jobResultLabel(result?: string): string {
  if (!result) return '-'
  if (result === 'success') return '成功'
  if (result === 'failed') return '失败'
  return result
}

export function triggerTypeLabel(triggerType?: string): string {
  if (triggerType === 'manual') return '手动触发'
  if (triggerType === 'schedule') return '定时触发'
  return triggerType || '-'
}

export function dbHealthLabel(health?: string): string {
  if (health === 'ok') return '正常'
  if (!health) return '未知'
  return health
}
