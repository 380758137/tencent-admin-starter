import http from './http'
import type { ApiResponse, PageData, SystemParam } from '../types'

export interface SystemParamQuery {
  page: number
  size: number
  keyword?: string
}

export interface SystemParamPayload {
  paramKey: string
  paramValue: string
  paramName: string
  status: number
  remark: string
}

export async function fetchSystemParams(query: SystemParamQuery): Promise<PageData<SystemParam>> {
  const { data } = await http.get<ApiResponse<PageData<SystemParam>>>('/api/system-params', { params: query })
  return data.data
}

export async function createSystemParam(payload: SystemParamPayload): Promise<SystemParam> {
  const { data } = await http.post<ApiResponse<SystemParam>>('/api/system-params', payload)
  return data.data
}

export async function updateSystemParam(id: number, payload: Partial<SystemParamPayload>): Promise<SystemParam> {
  const { data } = await http.put<ApiResponse<SystemParam>>(`/api/system-params/${id}`, payload)
  return data.data
}

export async function deleteSystemParam(id: number): Promise<void> {
  await http.delete(`/api/system-params/${id}`)
}
