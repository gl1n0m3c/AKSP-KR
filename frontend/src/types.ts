export interface User {
  id: string
  first_name: string
  last_name: string
  patronymic: string
  email: string
  phone: string | null
  telegram: string | null
  team_id?: string | null
  unit_id?: string | null
  position_id?: string | null
}

export interface Team {
  id: string
  name: string
  description: string
}

export interface Position {
  id: string
  name: string
  description: string
}

export interface Unit {
  id: string
  name: string
  description: string
}

export interface EventItem {
  id: string
  title: string
  description: string
  scheduled_at: string
  organizer_id: string
  team_id?: string | null
  participants: string[]
}

export interface MeResponse extends User {
  login: string
}

