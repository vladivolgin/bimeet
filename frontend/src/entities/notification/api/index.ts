import { apiFetch } from '@/shared/api/client'
import type { Notification } from '@/entities/notification/model/types'

export const notificationsApi = {
  list: () => apiFetch<Notification[]>('/api/notifications'),
  markRead: (id: string) => apiFetch<void>(`/api/notifications/${id}/read`, { method: 'PATCH' }),
  markAllRead: () => apiFetch<void>('/api/notifications/read-all', { method: 'POST' }),
  deleteOne: (id: string) => apiFetch<void>(`/api/notifications/${id}`, { method: 'DELETE' }),
  deleteAll: () => apiFetch<void>('/api/notifications', { method: 'DELETE' }),
}
