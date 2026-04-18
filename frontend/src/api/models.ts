import { apiClient } from './client'

export interface PublicModelPricing {
  model: string
  platform: string
  tags: string[]
  input_price: number | null
  output_price: number | null
  cache_write_price: number | null
  cache_read_price: number | null
  billing_mode: string
  rate_multiplier: number
}

export interface ModelPricingResponse {
  models: PublicModelPricing[]
  updated_at: string
}

export const modelsAPI = {
  getPublicPricing() {
    return apiClient.get<ModelPricingResponse>('/settings/public/model-pricing')
  }
}
