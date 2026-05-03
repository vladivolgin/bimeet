import { useMutation, useQueryClient } from '@tanstack/react-query'
import { carpoolsApi } from '@/features/carpools/api'
import { eventKeys } from '@/entities/event/queries'

export function useAddCarpool(eventId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: { seats_available: number; departure_point: string }) =>
      carpoolsApi.create(eventId, data),
    onSuccess: () => qc.invalidateQueries({ queryKey: eventKeys.detail(eventId) }),
  })
}

export function useJoinCarpool(eventId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (carpoolId: string) => carpoolsApi.join(eventId, carpoolId),
    onSuccess: () => qc.invalidateQueries({ queryKey: eventKeys.detail(eventId) }),
  })
}
