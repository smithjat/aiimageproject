import { request, type ApiResponse, type PageData } from './request'
import type { User, Image, Category } from '@/types'

export const authApi = {
  login(data: { username: string; password: string }): Promise<ApiResponse<{ token: string; user: User }>> {
    return request.post('/auth/login', data)
  },

  register(data: { username: string; password: string; email?: string }): Promise<ApiResponse<User>> {
    return request.post('/auth/register', data)
  },

  logout(): Promise<ApiResponse<null>> {
    return request.post('/auth/logout')
  },

  refreshToken(): Promise<ApiResponse<{ token: string }>> {
    return request.post('/auth/refresh')
  },
}

export const userApi = {
  getProfile(): Promise<ApiResponse<User>> {
    return request.get('/users/profile')
  },

  updateProfile(data: Partial<User>): Promise<ApiResponse<User>> {
    return request.put('/users/profile', data)
  },
}

export const imageApi = {
  getList(params: {
    page?: number
    page_size?: number
    category_id?: number
    keyword?: string
  }): Promise<ApiResponse<PageData<Image>>> {
    return request.get('/images', { params })
  },

  getDetail(id: number): Promise<ApiResponse<Image>> {
    return request.get(`/images/${id}`)
  },

  upload(formData: FormData): Promise<ApiResponse<Image>> {
    return request.post('/images', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },

  delete(id: number): Promise<ApiResponse<null>> {
    return request.delete(`/images/${id}`)
  },
}

export const categoryApi = {
  getList(): Promise<ApiResponse<Category[]>> {
    return request.get('/categories')
  },
}
