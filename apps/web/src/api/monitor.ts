import http from './http'
import type { ApiResponse, MonitorOverview } from '../types'

export async function fetchMonitorOverview(): Promise<MonitorOverview> {
  const { data } = await http.get<ApiResponse<MonitorOverview>>('/api/monitor/overview')
  return data.data
}
