import { useMutation, useQueryClient } from '@tanstack/react-query'
import { notificationsApi } from '@/entities/notification/api'
import { notificationKeys } from '@/entities/notification/queries'

export function useMarkRead() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => notificationsApi.markRead(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: notificationKeys.all }),
  })
}

export function useMarkAllRead() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: notificationsApi.markAllRead,
    onSuccess: () => qc.invalidateQueries({ queryKey: notificationKeys.all }),
  })
}

export function useDeleteNotification() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => notificationsApi.deleteOne(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: notificationKeys.all }),
  })
}

export function useDeleteAllNotifications() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: notificationsApi.deleteAll,
    onSuccess: () => qc.invalidateQueries({ queryKey: notificationKeys.all }),
  })
}
