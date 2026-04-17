import http from './http'
import type { ApiResponse, PageData, User } from '../types'

export interface UserQuery {
  page: number
  size: number
  keyword?: string
}

export interface UserPayload {
  username?: string
  password?: string
  displayName: string
  role: string
  deptId: number
  status: number
}

export async function fetchUsers(query: UserQuery): Promise<PageData<User>> {
  const { data } = await http.get<ApiResponse<PageData<User>>>('/api/users', { params: query })
  return data.data
}

export async function createUser(payload: UserPayload): Promise<User> {
  const { data } = await http.post<ApiResponse<User>>('/api/users', payload)
  return data.data
}

export async function updateUser(id: number, payload: UserPayload): Promise<User> {
  const { data } = await http.put<ApiResponse<User>>(`/api/users/${id}`, payload)
  return data.data
}

export async function deleteUser(id: number): Promise<void> {
  await http.delete(`/api/users/${id}`)
}

export async function exportUsers(): Promise<Blob> {
  const { data } = await http.get('/api/users/export', { responseType: 'blob' })
  return data
}

export async function importUsers(file: File): Promise<{ imported: number; skipped: number; failed: number }> {
  const formData = new FormData()
  formData.append('file', file)
  const { data } = await http.post<ApiResponse<{ imported: number; skipped: number; failed: number }>>('/api/users/import', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
  return data.data
}
