import { createAlova } from 'alova'
import adapterFetch from 'alova/fetch'
import VueHook from 'alova/vue'

const API_URL = import.meta.env.VITE_API_URL

const REFRESH_BUFFER_MS = 2 * 60 * 1000

const authInstance = createAlova({
  requestAdapter: adapterFetch(),
  statesHook: VueHook,
  responded: (response) => response.json(),
  baseURL: `${API_URL}`,
  timeout: 10000
})

export function isTokenNearExpiry(): boolean {
  const expires = sessionStorage.getItem('jwt_expires')
  if (!expires) return true
  return Date.now() >= Number(expires) - REFRESH_BUFFER_MS
}

export function clearAuth(): void {
  sessionStorage.removeItem('jwt_token')
  sessionStorage.removeItem('refresh_token')
  sessionStorage.removeItem('jwt_expires')
}

let refreshPromise: Promise<boolean> | null = null

export async function refreshAuthToken(): Promise<boolean> {
  if (refreshPromise) return refreshPromise

  refreshPromise = _doRefresh()
  try {
    return await refreshPromise
  } finally {
    refreshPromise = null
  }
}

interface RefreshResponse {
  token: string
  refresh: string
  expires_at: number
}

async function _doRefresh(): Promise<boolean> {
  const refreshToken = sessionStorage.getItem('refresh_token')
  if (!refreshToken) {
    clearAuth()
    return false
  }

  try {
    const data = (await authInstance.Post('/auth/refresh', { refresh: refreshToken })) as RefreshResponse

    sessionStorage.setItem('jwt_token', data.token)
    sessionStorage.setItem('refresh_token', data.refresh)
    sessionStorage.setItem('jwt_expires', String(data.expires_at))

    return true
  } catch {
    clearAuth()
    return false
  }
}
