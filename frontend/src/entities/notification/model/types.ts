export interface Notification {
  id: string
  event_id?: string
  type: string
  message: string
  is_read: boolean
  created_at: string
}
