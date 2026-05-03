import { apiFetch } from '@/shared/api/client'
import type { User, ProfileStats } from '@/entities/user/model/types'

const API_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'

export interface AuthResponse {
  token: string
  user: User
}

export interface UpdateProfileData {
  name: string
  last_name: string
  birth_date: string
  city: string
}

export const authApi = {
  register: (name: string, email: string, password: string) =>
    apiFetch<AuthResponse>('/api/auth/register', {
      method: 'POST',
      body: JSON.stringify({ name, email, password }),
    }),
  login: (email: string, password: string) =>
    apiFetch<AuthResponse>('/api/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    }),
  forgotPassword: (email: string) =>
    apiFetch<{ status: string }>('/api/auth/forgot-password', {
      method: 'POST',
      body: JSON.stringify({ email }),
    }),
  resetPassword: (token: string, password: string) =>
    apiFetch<{ status: string }>('/api/auth/reset-password', {
      method: 'POST',
      body: JSON.stringify({ token, password }),
    }),
  getMe: () =>
    apiFetch<User>('/api/auth/me'),
  updateProfile: (data: UpdateProfileData) =>
    apiFetch<User>('/api/auth/me', {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
  getStats: () =>
    apiFetch<ProfileStats>('/api/auth/me/stats'),
  uploadAvatar: async (file: File): Promise<User> => {
    const token = localStorage.getItem('token')
    const form = new FormData()
    form.append('avatar', file)
    const res = await fetch(`${API_URL}/api/auth/me/avatar`, {
      method: 'POST',
      headers: token ? { Authorization: `Bearer ${token}` } : {},
      body: form,
    })
    if (!res.ok) {
      const body = await res.json().catch(() => ({ error: 'Unknown error' }))
      throw new Error(body.error ?? `HTTP ${res.status}`)
    }
    return res.json()
  },
  deleteAvatar: () =>
    apiFetch<User>('/api/auth/me/avatar', { method: 'DELETE' }),
}
