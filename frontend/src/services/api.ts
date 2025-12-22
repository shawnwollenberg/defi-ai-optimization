const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

export async function apiRequest(endpoint: string, options: RequestInit = {}) {
  const token = localStorage.getItem('auth_token')
  
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...options.headers,
  }
  
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }
  
  const response = await fetch(`${API_URL}${endpoint}`, {
    ...options,
    headers,
  })
  
  if (!response.ok) {
    throw new Error(`API request failed: ${response.statusText}`)
  }
  
  return response.json()
}

