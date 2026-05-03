import type { User } from '@/entities/user/model/types'

export interface Carpool {
  id: string
  driver: User
  seats_available: number
  departure_point: string
  passengers: User[]
}
