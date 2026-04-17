import http from './http'
import type { ApiResponse, Menu, PageData } from '../types'

export interface MenuQuery {
  page: number
  size: number
  keyword?: string
}

export interface MenuPayload {
  parentId: number
  name: string
  menuType: string
  path: string
  component: string
  perms: string
  sort: number
  status: number
}

export async function fetchMenus(query: MenuQuery): Promise<PageData<Menu>> {
  const { data } = await http.get<ApiResponse<PageData<Menu>>>('/api/menus', { params: query })
  return data.data
}

export async function createMenu(payload: MenuPayload): Promise<Menu> {
  const { data } = await http.post<ApiResponse<Menu>>('/api/menus', payload)
  return data.data
}

export async function updateMenu(id: number, payload: Partial<MenuPayload>): Promise<Menu> {
  const { data } = await http.put<ApiResponse<Menu>>(`/api/menus/${id}`, payload)
  return data.data
}

export async function deleteMenu(id: number): Promise<void> {
  await http.delete(`/api/menus/${id}`)
}
