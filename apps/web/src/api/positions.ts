import http from './http'
import type { ApiResponse, PageData, Position } from '../types'

export interface PositionQuery {
  page: number
  size: number
  keyword?: string
}

export interface PositionPayload {
  name: string
  code: string
  sort: number
  status: number
  remark: string
}

export async function fetchPositions(query: PositionQuery): Promise<PageData<Position>> {
  const { data } = await http.get<ApiResponse<PageData<Position>>>('/api/positions', { params: query })
  return data.data
}

export async function createPosition(payload: PositionPayload): Promise<Position> {
  const { data } = await http.post<ApiResponse<Position>>('/api/positions', payload)
  return data.data
}

export async function updatePosition(id: number, payload: Partial<PositionPayload>): Promise<Position> {
  const { data } = await http.put<ApiResponse<Position>>(`/api/positions/${id}`, payload)
  return data.data
}

export async function deletePosition(id: number): Promise<void> {
  await http.delete(`/api/positions/${id}`)
}
