export interface User {
  id: string
  name: string
  last_name?: string
  email: string
  birth_date?: string
  city?: string
  avatar_url?: string
  created_at?: string
}

export interface ProfileStats {
  organized: number
  participated: number
  completed: number
  upcoming: number
}

export type ParticipantStatus = 'invited' | 'confirmed' | 'declined'

export interface Participant {
  id: string
  user: User
  status: ParticipantStatus
}
