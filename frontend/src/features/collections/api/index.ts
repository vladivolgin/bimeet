import { apiFetch } from '@/shared/api/client'
import type { CollectionContribution } from '@/entities/collection/model/types'

export const collectionsApi = {
  list: (eventId: string) => apiFetch<unknown[]>(`/api/events/${eventId}/collections`),
  create: (eventId: string, data: { title: string; per_person_amount: number }) =>
    apiFetch<unknown>(`/api/events/${eventId}/collections`, { method: 'POST', body: JSON.stringify(data) }),
  delete: (eventId: string, collectionId: string) =>
    apiFetch<void>(`/api/events/${eventId}/collections/${collectionId}`, { method: 'DELETE' }),
  contribute: (eventId: string, collectionId: string, file: File) => {
    const form = new FormData()
    form.append('receipt', file)
    return apiFetch<CollectionContribution>(
      `/api/events/${eventId}/collections/${collectionId}/contribute`,
      { method: 'POST', body: form },
    )
  },
  confirmContribution: (eventId: string, collectionId: string, contribId: string) =>
    apiFetch<CollectionContribution>(
      `/api/events/${eventId}/collections/${collectionId}/contributions/${contribId}/confirm`,
      { method: 'POST' },
    ),
  rejectContribution: (eventId: string, collectionId: string, contribId: string) =>
    apiFetch<CollectionContribution>(
      `/api/events/${eventId}/collections/${collectionId}/contributions/${contribId}/reject`,
      { method: 'POST' },
    ),
  markPaid: (eventId: string, collectionId: string, userId: string) =>
    apiFetch<CollectionContribution>(
      `/api/events/${eventId}/collections/${collectionId}/contributions/mark-paid`,
      { method: 'POST', body: JSON.stringify({ user_id: userId }) },
    ),
  summary: (eventId: string) => apiFetch<unknown>(`/api/events/${eventId}/collections/summary`),
}
