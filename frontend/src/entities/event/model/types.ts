import type { User, Participant } from '@/entities/user/model/types'
import type { Collection } from '@/entities/collection/model/types'
import type { Poll } from '@/entities/poll/model/types'
import type { Item } from '@/entities/item/model/types'
import type { Carpool } from '@/entities/carpool/model/types'

export type EventCategory = 'ordinary' | 'business'

export interface EventLink {
  id: string
  event_id: string
  title: string
  url: string
  created_by: string
  created_at: string
}

export interface ChangeLogEntry {
  id: string
  changed_by: User
  field_name: string
  old_value: string
  new_value: string
  changed_at: string
}

export interface Event {
  id: string
  title: string
  description: string
  date_start: string
  date_end: string
  location: string
  category: EventCategory
  dress_code?: string
  status: 'active' | 'completed'
  is_public: boolean
  max_guests?: number
  confirmed_count: number
  organizer: User
  invite_token: string
  created_at: string
  updated_at: string
  my_status?: string
  participants: Participant[]
  collections: Collection[]
  polls: Poll[]
  items: Item[]
  carpools: Carpool[]
  links: EventLink[]
  change_log: ChangeLogEntry[]
}

export interface EventListItem {
  id: string
  title: string
  description: string
  date_start: string
  date_end: string
  location: string
  category: EventCategory
  dress_code?: string
  status: 'active' | 'completed'
  is_public: boolean
  max_guests?: number
  confirmed_count: number
  organizer: { id: string; name: string; email: string }
  invite_token: string
  created_at: string
  my_status: string
  participants: Array<{ id: string; user: { id: string; name: string; email: string }; status: string }>
}

export interface PublicEventListItem extends EventListItem {
  is_participant: boolean
}

export interface InviteEventInfo {
  id: string
  title: string
  description: string
  date_start: string
  date_end: string
  location: string
  category: EventCategory
  dress_code?: string
  organizer: { id: string; name: string; email: string }
  confirmed_count: number
}

export interface CreateEventData {
  title: string
  description: string
  date_start: string
  date_end: string
  location: string
  category: EventCategory
  dress_code?: string
}

export interface UpdateEventData {
  title?: string
  description?: string
  date_start?: string
  date_end?: string
  location?: string
  category?: EventCategory
  dress_code?: string
}
