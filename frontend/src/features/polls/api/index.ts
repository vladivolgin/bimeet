import { apiFetch } from '@/shared/api/client'

export const pollsApi = {
  list: (eventId: string) => apiFetch<unknown[]>(`/api/events/${eventId}/polls`),
  create: (eventId: string, data: { question: string; options: string[] }) =>
    apiFetch<unknown>(`/api/events/${eventId}/polls`, { method: 'POST', body: JSON.stringify(data) }),
  vote: (eventId: string, pollId: string, optionId: string) =>
    apiFetch<unknown>(`/api/events/${eventId}/polls/${pollId}/vote`, {
      method: 'POST',
      body: JSON.stringify({ option_id: optionId }),
    }),
}
