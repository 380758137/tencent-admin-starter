import http from './http'
import type { ApiResponse, OnlineUser, PageData } from '../types'

export interface OnlineUserQuery {
  page: number
  size: number
  keyword?: string
}

export async function fetchOnlineUsers(query: OnlineUserQuery): Promise<PageData<OnlineUser>> {
  const { data } = await http.get<ApiResponse<PageData<OnlineUser>>>('/api/online-users', { params: query })
  return data.data
}
