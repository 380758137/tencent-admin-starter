import http from './http'
import type { ApiResponse, AuthUser, LoginPayload, LoginResult } from '../types'

export async function login(payload: LoginPayload): Promise<LoginResult> {
  const { data } = await http.post<ApiResponse<LoginResult>>('/api/auth/login', payload)
  return data.data
}

export async function me(): Promise<AuthUser> {
  const { data } = await http.get<ApiResponse<AuthUser>>('/api/auth/me')
  return data.data
}

