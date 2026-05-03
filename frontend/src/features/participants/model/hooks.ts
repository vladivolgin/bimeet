import { useMutation, useQueryClient } from '@tanstack/react-query'
import { eventMutationsApi } from '@/features/event-manage/api'
import { eventKeys } from '@/entities/event/queries'

export function useInviteParticipant(eventId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (email: string) => eventMutationsApi.inviteParticipant(eventId, email),
    onSuccess: () => qc.invalidateQueries({ queryKey: eventKeys.detail(eventId) }),
  })
}

export function useUpdateParticipantStatus(eventId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ userId, status }: { userId: string; status: 'confirmed' | 'declined' }) =>
      eventMutationsApi.updateParticipantStatus(eventId, userId, status),
    onSuccess: () => qc.invalidateQueries({ queryKey: eventKeys.detail(eventId) }),
  })
}
