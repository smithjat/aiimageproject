import { request, type ApiResponse } from './request'

export interface AIModel {
  id: number
  model_id: string
  model_name: string
  provider: string
  description: string
  capabilities: string[]
  supported_parameters: string[]
  is_active: boolean
  priority: number
  metadata: Record<string, unknown>
  parameter_constraints: Record<string, unknown>
  created_at: string
  updated_at: string
}

export interface Text2ImageRequest {
  prompt: string
  model_id?: string
  negative_prompt?: string
  width?: number
  height?: number
  cfg_scale?: number
  steps?: number
  seed?: number
  style?: string
  quality?: string
}

export interface ImageEditRequest {
  image: string
  prompt: string
  model_id?: string
  negative_prompt?: string
  mask?: string
  strength?: number
  cfg_scale?: number
  steps?: number
  seed?: number
}

export interface GenerateResponse {
  success: boolean
  message: string
  status: string
  task_id: string
  images: string[]
  quota: number
}

export const modelApi = {
  getAll(): Promise<ApiResponse<AIModel[]>> {
    return request.get('/models')
  },

  getByCapability(capability: 'generation' | 'editing'): Promise<ApiResponse<AIModel[]>> {
    return request.get(`/models/capability/${capability}`)
  },

  getByProvider(provider: string): Promise<ApiResponse<AIModel[]>> {
    return request.get(`/models/provider/${provider}`)
  },

  getById(modelId: string): Promise<ApiResponse<AIModel>> {
    return request.get(`/models/${modelId}`)
  },

  getCapabilities(): Promise<ApiResponse<string[]>> {
    return request.get('/models/capabilities')
  },
}

export const imageApi = {
  text2Image(data: Text2ImageRequest): Promise<ApiResponse<GenerateResponse>> {
    return request.post('/text2image', data)
  },

  imageEdit(data: ImageEditRequest): Promise<ApiResponse<GenerateResponse>> {
    return request.post('/image/edit', data)
  },
}
