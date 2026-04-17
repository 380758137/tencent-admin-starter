import http from './http'
import type { ApiResponse, DictionaryItem, PageData } from '../types'

export interface DictionaryQuery {
  page: number
  size: number
  keyword?: string
  dictType?: string
}

export interface DictionaryPayload {
  dictType: string
  dictLabel: string
  dictValue: string
  sort: number
  status: number
  remark: string
}

export async function fetchDictionaryItems(query: DictionaryQuery): Promise<PageData<DictionaryItem>> {
  const { data } = await http.get<ApiResponse<PageData<DictionaryItem>>>('/api/dictionary-items', { params: query })
  return data.data
}

export async function createDictionaryItem(payload: DictionaryPayload): Promise<DictionaryItem> {
  const { data } = await http.post<ApiResponse<DictionaryItem>>('/api/dictionary-items', payload)
  return data.data
}

export async function updateDictionaryItem(id: number, payload: Partial<DictionaryPayload>): Promise<DictionaryItem> {
  const { data } = await http.put<ApiResponse<DictionaryItem>>(`/api/dictionary-items/${id}`, payload)
  return data.data
}

export async function deleteDictionaryItem(id: number): Promise<void> {
  await http.delete(`/api/dictionary-items/${id}`)
}
