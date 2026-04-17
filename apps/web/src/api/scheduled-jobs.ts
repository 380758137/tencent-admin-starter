import http from './http'
import type { ApiResponse, PageData, ScheduledJob, ScheduledJobLog } from '../types'

export interface ScheduledJobQuery {
  page: number
  size: number
  keyword?: string
}

export interface ScheduledJobPayload {
  name: string
  cronExpr: string
  command: string
  status: number
  remark: string
}

export interface ScheduledJobLogQuery {
  page: number
  size: number
  jobId?: number
}

export async function fetchScheduledJobs(query: ScheduledJobQuery): Promise<PageData<ScheduledJob>> {
  const { data } = await http.get<ApiResponse<PageData<ScheduledJob>>>('/api/scheduled-jobs', { params: query })
  return data.data
}

export async function createScheduledJob(payload: ScheduledJobPayload): Promise<ScheduledJob> {
  const { data } = await http.post<ApiResponse<ScheduledJob>>('/api/scheduled-jobs', payload)
  return data.data
}

export async function updateScheduledJob(id: number, payload: Partial<ScheduledJobPayload>): Promise<ScheduledJob> {
  const { data } = await http.put<ApiResponse<ScheduledJob>>(`/api/scheduled-jobs/${id}`, payload)
  return data.data
}

export async function deleteScheduledJob(id: number): Promise<void> {
  await http.delete(`/api/scheduled-jobs/${id}`)
}

export async function runScheduledJob(id: number): Promise<ScheduledJob> {
  const { data } = await http.post<ApiResponse<ScheduledJob>>(`/api/scheduled-jobs/${id}/run`)
  return data.data
}

export async function fetchScheduledJobLogs(query: ScheduledJobLogQuery): Promise<PageData<ScheduledJobLog>> {
  const { data } = await http.get<ApiResponse<PageData<ScheduledJobLog>>>('/api/scheduled-job-logs', { params: query })
  return data.data
}
