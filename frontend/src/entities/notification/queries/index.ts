import { useQuery } from '@tanstack/react-query'
import { notificationsApi } from '@/entities/notification/api'

export const notificationKeys = {
  all: ['notifications'] as const,
}

export function useNotifications() {
  return useQuery({
    queryKey: notificationKeys.all,
    queryFn: notificationsApi.list,
    refetchInterval: 30_000,
  })
}
