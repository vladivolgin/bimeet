import { useQuery } from '@tanstack/react-query'
import { eventsReadApi } from '@/entities/event/api'
import type { Event } from '@/entities/event/model/types'

export const eventKeys = {
  all: ['events'] as const,
  detail: (id: string) => ['events', id] as const,
}

export const publicEventKeys = {
  all: ['publicEvents'] as const,
}

export function useEvents() {
  const { data = [] } = useQuery({
    queryKey: eventKeys.all,
    queryFn: eventsReadApi.list,
  })
  return data
}

export function useEvent(id: string) {
  const { data } = useQuery({
    queryKey: eventKeys.detail(id),
    queryFn: () => eventsReadApi.getById(id),
    enabled: !!id,
  })
  return data as Event | undefined
}

export function usePublicEvents() {
  const { data = [] } = useQuery({
    queryKey: publicEventKeys.all,
    queryFn: eventsReadApi.listPublic,
  })
  return data
}
