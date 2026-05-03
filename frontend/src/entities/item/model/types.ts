import type { User } from '@/entities/user/model/types'

export interface Item {
  id: string
  name: string
  assigned_to?: User | null
}
