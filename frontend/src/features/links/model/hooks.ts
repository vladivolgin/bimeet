import { useMutation, useQueryClient } from '@tanstack/react-query'
import { linksApi } from '@/features/links/api'
import { eventKeys } from '@/entities/event/queries'

export function useAddLink(eventId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: { title: string; url: string }) => linksApi.create(eventId, data),
    onSuccess: () => qc.invalidateQueries({ queryKey: eventKeys.detail(eventId) }),
  })
}

export function useRemoveLink(eventId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (linkId: string) => linksApi.delete(eventId, linkId),
    onSuccess: () => qc.invalidateQueries({ queryKey: eventKeys.detail(eventId) }),
  })
}
