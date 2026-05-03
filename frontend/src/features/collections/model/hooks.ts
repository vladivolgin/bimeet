import { useMutation, useQueryClient } from '@tanstack/react-query'
import { collectionsApi } from '@/features/collections/api'
import { eventKeys } from '@/entities/event/queries'

export function useAddCollection(eventId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: { title: string; per_person_amount: number }) => collectionsApi.create(eventId, data),
    onSuccess: () => qc.invalidateQueries({ queryKey: eventKeys.detail(eventId) }),
  })
}

export function useRemoveCollection(eventId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (collectionId: string) => collectionsApi.delete(eventId, collectionId),
    onSuccess: () => qc.invalidateQueries({ queryKey: eventKeys.detail(eventId) }),
  })
}

export function useSubmitContribution(eventId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ collectionId, file }: { collectionId: string; file: File }) =>
      collectionsApi.contribute(eventId, collectionId, file),
    onSuccess: () => qc.invalidateQueries({ queryKey: eventKeys.detail(eventId) }),
  })
}

export function useConfirmContribution(eventId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ collectionId, contribId }: { collectionId: string; contribId: string }) =>
      collectionsApi.confirmContribution(eventId, collectionId, contribId),
    onSuccess: () => qc.invalidateQueries({ queryKey: eventKeys.detail(eventId) }),
  })
}

export function useRejectContribution(eventId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ collectionId, contribId }: { collectionId: string; contribId: string }) =>
      collectionsApi.rejectContribution(eventId, collectionId, contribId),
    onSuccess: () => qc.invalidateQueries({ queryKey: eventKeys.detail(eventId) }),
  })
}

export function useMarkPaid(eventId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ collectionId, userId }: { collectionId: string; userId: string }) =>
      collectionsApi.markPaid(eventId, collectionId, userId),
    onSuccess: () => qc.invalidateQueries({ queryKey: eventKeys.detail(eventId) }),
  })
}
