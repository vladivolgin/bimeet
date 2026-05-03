import { useMutation, useQueryClient } from '@tanstack/react-query'
import { itemsApi } from '@/features/items/api'
import { eventKeys } from '@/entities/event/queries'

export function useAddItem(eventId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (name: string) => itemsApi.create(eventId, name),
    onSuccess: () => qc.invalidateQueries({ queryKey: eventKeys.detail(eventId) }),
  })
}

export function useAssignItem(eventId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ itemId, userId }: { itemId: string; userId: string | null }) =>
      itemsApi.assign(eventId, itemId, userId),
    onSuccess: () => qc.invalidateQueries({ queryKey: eventKeys.detail(eventId) }),
  })
}
