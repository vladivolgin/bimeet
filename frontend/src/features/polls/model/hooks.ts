import { useMutation, useQueryClient } from '@tanstack/react-query'
import { pollsApi } from '@/features/polls/api'
import { eventKeys } from '@/entities/event/queries'

export function useAddPoll(eventId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: { question: string; options: string[] }) => pollsApi.create(eventId, data),
    onSuccess: () => qc.invalidateQueries({ queryKey: eventKeys.detail(eventId) }),
  })
}

export function useVote(eventId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ pollId, optionId }: { pollId: string; optionId: string }) =>
      pollsApi.vote(eventId, pollId, optionId),
    onSuccess: () => qc.invalidateQueries({ queryKey: eventKeys.detail(eventId) }),
  })
}
