export interface CollectionContribution {
  id: string
  collection_id: string
  user_id: string
  paid: boolean
  paid_at?: string
  status: 'not_paid' | 'pending' | 'paid'
  receipt_url?: string
}

export interface Collection {
  id: string
  title: string
  per_person_amount: number
  created_by: string
  created_at: string
  contributions: CollectionContribution[]
}
