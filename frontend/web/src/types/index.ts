export interface User {
  id: number
  username: string
  nickname: string
  avatar: string
  email: string
  phone: string
  role: 'user' | 'admin'
  status: number
  created_at: string
}

export interface Image {
  id: number
  user_id: number
  title: string
  description: string
  url: string
  thumbnail: string
  width: number
  height: number
  size: number
  format: string
  status: number
  created_at: string
}

export interface Category {
  id: number
  name: string
  parent_id: number
  sort: number
  status: number
}
