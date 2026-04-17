import http from './http'
import type { ApiResponse, Notice, PageData } from '../types'

export interface NoticeQuery {
  page: number
  size: number
  keyword?: string
}

export interface NoticePayload {
  title: string
  noticeType: string
  content: string
  pinned: number
  status: number
}

export async function fetchNotices(query: NoticeQuery): Promise<PageData<Notice>> {
  const { data } = await http.get<ApiResponse<PageData<Notice>>>('/api/notices', { params: query })
  return data.data
}

export async function createNotice(payload: NoticePayload): Promise<Notice> {
  const { data } = await http.post<ApiResponse<Notice>>('/api/notices', payload)
  return data.data
}

export async function updateNotice(id: number, payload: Partial<NoticePayload>): Promise<Notice> {
  const { data } = await http.put<ApiResponse<Notice>>(`/api/notices/${id}`, payload)
  return data.data
}

export async function deleteNotice(id: number): Promise<void> {
  await http.delete(`/api/notices/${id}`)
}
