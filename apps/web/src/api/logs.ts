import http from './http'
import type { ApiResponse, LoginLog, OperationLog, PageData } from '../types'

export interface OperationLogQuery {
  page: number
  size: number
  username?: string
  method?: string
  path?: string
  status?: number
}

export interface LoginLogQuery {
  page: number
  size: number
  username?: string
  status?: number
}

export async function fetchOperationLogs(query: OperationLogQuery): Promise<PageData<OperationLog>> {
  const { data } = await http.get<ApiResponse<PageData<OperationLog>>>('/api/operation-logs', { params: query })
  return data.data
}

export async function fetchLoginLogs(query: LoginLogQuery): Promise<PageData<LoginLog>> {
  const { data } = await http.get<ApiResponse<PageData<LoginLog>>>('/api/login-logs', { params: query })
  return data.data
}
