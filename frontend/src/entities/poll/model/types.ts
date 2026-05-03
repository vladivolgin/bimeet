export interface PollOption {
  id: string
  label: string
  votes: string[] // voter user IDs
}

export interface Poll {
  id: string
  question: string
  options: PollOption[]
  created_by: string
  created_at: string
}
