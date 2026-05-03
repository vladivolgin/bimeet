import { useMutation, useQueryClient } from '@tanstack/react-query'
import { eventMutationsApi } from '@/features/event-manage/api'
import { eventKeys, publicEventKeys } from '@/entities/event/queries'
import type { UpdateEventData } from '@/entities/event/model/types'

export function useCreateEvent() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: eventMutationsApi.create,
    onSuccess: () => qc.invalidateQueries({ queryKey: eventKeys.all }),
  })
}

export function useUpdateEvent(id: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: UpdateEventData) => eventMutationsApi.update(id, data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: eventKeys.all })
      qc.invalidateQueries({ queryKey: eventKeys.detail(id) })
    },
  })
}

export function useDeleteEvent() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: eventMutationsApi.delete,
    onSuccess: () => qc.invalidateQueries({ queryKey: eventKeys.all }),
  })
}

export function useCompleteEvent(id: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: () => eventMutationsApi.complete(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: eventKeys.all })
      qc.invalidateQueries({ queryKey: eventKeys.detail(id) })
    },
  })
}

export function useJoinPublicEvent() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => eventMutationsApi.joinPublicEvent(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: publicEventKeys.all })
      qc.invalidateQueries({ queryKey: eventKeys.all })
    },
  })
}

export function useConfirmAttendance() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ eventId, userId }: { eventId: string; userId: string }) =>
      eventMutationsApi.updateParticipantStatus(eventId, userId, 'confirmed'),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: eventKeys.all })
    },
  })
}
