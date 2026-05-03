import { apiFetch } from '@/shared/api/client'
import type { Event, CreateEventData, UpdateEventData } from '@/entities/event/model/types'

export const eventMutationsApi = {
  create: (data: CreateEventData) =>
    apiFetch<Event>('/api/events', { method: 'POST', body: JSON.stringify(data) }),

  update: (id: string, data: UpdateEventData) =>
    apiFetch<Event>(`/api/events/${id}`, { method: 'PUT', body: JSON.stringify(data) }),

  delete: (id: string) =>
    apiFetch<void>(`/api/events/${id}`, { method: 'DELETE' }),

  complete: (id: string) =>
    apiFetch<void>(`/api/events/${id}/complete`, { method: 'POST' }),

  inviteParticipant: (eventId: string, email: string) =>
    apiFetch<{ status: string; message: string }>(`/api/events/${eventId}/participants`, { method: 'POST', body: JSON.stringify({ email }) }),

  updateParticipantStatus: (eventId: string, userId: string, status: 'confirmed' | 'declined') =>
    apiFetch(`/api/events/${eventId}/participants/${userId}`, { method: 'PATCH', body: JSON.stringify({ status }) }),

  joinPublicEvent: (id: string) =>
    apiFetch<{ status: string }>(`/api/events/${id}/join`, { method: 'POST' }),
}
