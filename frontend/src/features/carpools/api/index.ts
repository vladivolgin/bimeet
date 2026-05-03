import { apiFetch } from '@/shared/api/client'

export const carpoolsApi = {
  list: (eventId: string) => apiFetch<unknown[]>(`/api/events/${eventId}/carpools`),
  create: (eventId: string, data: { seats_available: number; departure_point: string }) =>
    apiFetch<unknown>(`/api/events/${eventId}/carpools`, { method: 'POST', body: JSON.stringify(data) }),
  join: (eventId: string, carpoolId: string) =>
    apiFetch<unknown>(`/api/events/${eventId}/carpools/${carpoolId}/join`, { method: 'POST' }),
}
