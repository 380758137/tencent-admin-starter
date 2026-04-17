import http from './http'
import type { ApiResponse, Department, PageData } from '../types'

export interface DepartmentQuery {
  page: number
  size: number
  keyword?: string
}

export interface DepartmentPayload {
  name: string
  code: string
  parentId: number
  manager: string
  status: number
}

export async function fetchDepartments(query: DepartmentQuery): Promise<PageData<Department>> {
  const { data } = await http.get<ApiResponse<PageData<Department>>>('/api/departments', { params: query })
  return data.data
}

export async function createDepartment(payload: DepartmentPayload): Promise<Department> {
  const { data } = await http.post<ApiResponse<Department>>('/api/departments', payload)
  return data.data
}

export async function updateDepartment(id: number, payload: DepartmentPayload): Promise<Department> {
  const { data } = await http.put<ApiResponse<Department>>(`/api/departments/${id}`, payload)
  return data.data
}

export async function deleteDepartment(id: number): Promise<void> {
  await http.delete(`/api/departments/${id}`)
}
