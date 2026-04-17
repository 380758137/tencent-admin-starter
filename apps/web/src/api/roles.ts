import http from './http'
import type { ApiResponse, PageData, Role } from '../types'

export interface RoleQuery {
  page: number
  size: number
  keyword?: string
}

export interface RolePayload {
  name: string
  roleKey: string
  permissions: string
  dataScope: 'all' | 'dept' | 'self'
  status: number
  remark: string
}

export async function fetchRoles(query: RoleQuery): Promise<PageData<Role>> {
  const { data } = await http.get<ApiResponse<PageData<Role>>>('/api/roles', { params: query })
  return data.data
}

export async function createRole(payload: RolePayload): Promise<Role> {
  const { data } = await http.post<ApiResponse<Role>>('/api/roles', payload)
  return data.data
}

export async function updateRole(id: number, payload: Partial<RolePayload>): Promise<Role> {
  const { data } = await http.put<ApiResponse<Role>>(`/api/roles/${id}`, payload)
  return data.data
}

export async function deleteRole(id: number): Promise<void> {
  await http.delete(`/api/roles/${id}`)
}
