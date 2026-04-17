export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

export interface AuthUser {
  id: number
  username: string
  displayName: string
  role: string
  deptId: number
  permissions: string[]
  dataScope: 'all' | 'dept' | 'self'
  status: number
}

export interface LoginResult {
  token: string
  user: AuthUser
}

export interface LoginPayload {
  username: string
  password: string
}

export interface PageData<T> {
  list: T[]
  total: number
  page: number
  size: number
}

export interface User {
  id: number
  username: string
  displayName: string
  role: string
  deptId: number
  status: number
  createdAt: string
  updatedAt: string
}

export interface Department {
  id: number
  name: string
  code: string
  parentId: number
  manager: string
  status: number
  createdAt: string
  updatedAt: string
}

export interface Role {
  id: number
  name: string
  roleKey: string
  permissions: string
  dataScope: 'all' | 'dept' | 'self'
  status: number
  remark: string
  createdAt: string
  updatedAt: string
}

export interface Menu {
  id: number
  parentId: number
  name: string
  menuType: string
  path: string
  component: string
  perms: string
  sort: number
  status: number
  createdAt: string
  updatedAt: string
}

export interface DictionaryItem {
  id: number
  dictType: string
  dictLabel: string
  dictValue: string
  sort: number
  status: number
  remark: string
  createdAt: string
  updatedAt: string
}

export interface SystemParam {
  id: number
  paramKey: string
  paramValue: string
  paramName: string
  status: number
  remark: string
  createdAt: string
  updatedAt: string
}

export interface OperationLog {
  id: number
  userId: number
  username: string
  method: string
  path: string
  statusCode: number
  latencyMs: number
  ip: string
  userAgent: string
  createdAt: string
}

export interface LoginLog {
  id: number
  username: string
  status: number
  message: string
  ip: string
  userAgent: string
  createdAt: string
}

export interface MonitorOverview {
  runtime: {
    goroutines: number
    heapAlloc: number
    heapSys: number
    goVersion: string
    now: string
  }
  database: {
    health: string
    stats: Record<string, unknown>
  }
  counts: Record<string, number>
}

export interface Position {
  id: number
  name: string
  code: string
  sort: number
  status: number
  remark: string
  createdAt: string
  updatedAt: string
}

export interface Notice {
  id: number
  title: string
  noticeType: string
  content: string
  pinned: number
  status: number
  createdAt: string
  updatedAt: string
}

export interface OnlineUser {
  username: string
  lastActiveAt: string
  requestCount: number
}

export interface ScheduledJob {
  id: number
  name: string
  cronExpr: string
  command: string
  status: number
  lastRunAt?: string
  lastResult?: string
  remark: string
  createdAt: string
  updatedAt: string
}

export interface ScheduledJobLog {
  id: number
  jobId: number
  triggerType: string
  status: string
  message: string
  runAt: string
}
