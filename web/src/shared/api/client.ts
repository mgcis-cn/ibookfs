const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1'

export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public code?: string,
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

export async function apiCall<T>(
  endpoint: string,
  method: string = 'GET',
  data?: unknown,
  requiresAuth: boolean = true,
): Promise<T> {
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
  }

  if (requiresAuth) {
    const token = localStorage.getItem('authToken')
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    method,
    headers,
    body: data ? JSON.stringify(data) : undefined,
  })

  if (!response.ok) {
    const error = await response.json().catch(() => ({ message: 'Request failed' }))
    throw new ApiError(
      error.message || error.error || 'Request failed',
      response.status,
      error.code,
    )
  }

  return response.json()
}
