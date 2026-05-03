import { apiFetch } from '@/shared/api/client'
import type { EventLink } from '@/entities/event/model/types'

export const linksApi = {
  list: (eventId: string) => apiFetch<EventLink[]>(`/api/events/${eventId}/links`),
  create: (eventId: string, data: { title: string; url: string }) =>
    apiFetch<EventLink>(`/api/events/${eventId}/links`, {
      method: 'POST',
      body: JSON.stringify(data),
    }),
  delete: (eventId: string, linkId: string) =>
    apiFetch<void>(`/api/events/${eventId}/links/${linkId}`, { method: 'DELETE' }),
}
