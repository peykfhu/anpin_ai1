/**
 * Referral/Commission System API endpoints
 * Handles referral code management, commission tracking, and invite records
 */

import { apiClient } from './client'
import type { BasePaginationResponse } from '@/types'

// ==================== Types ====================

export interface CommissionTier {
  min_monthly_recharge: number
  max_monthly_recharge: number // 0 means unlimited
  rate: number
}

export interface ReferralInfo {
  enabled: boolean
  referral_code: string
  commission_rate: number
  total_commission: number
  total_invitees: number
  total_recharged: number
  balance: number
  commission_tiers: CommissionTier[]
}

export interface InviteRecord {
  id: number
  invitee_email: string
  invitee_username: string
  registered_at: string
  total_recharged: number
  total_commission: number
  status: string
}

export interface CommissionRecord {
  id: number
  invitee_email: string
  invitee_username: string
  order_id: string
  recharge_amount: number
  commission_rate: number
  commission_amount: number
  status: string
  created_at: string
}

export interface WithdrawRequest {
  amount: number
}

export interface WithdrawResponse {
  message: string
  balance: number
}

// ==================== API Functions ====================

/**
 * Get current user's referral info (code, stats)
 */
export async function getReferralInfo(): Promise<ReferralInfo> {
  const { data } = await apiClient.get<ReferralInfo>('/user/referral/info')
  return data
}

/**
 * Generate or regenerate referral code
 */
export async function generateReferralCode(): Promise<{ referral_code: string }> {
  const { data } = await apiClient.post<{ referral_code: string }>('/user/referral/generate-code')
  return data
}

/**
 * Get invite records (paginated)
 */
export async function getInviteRecords(
  page: number = 1,
  pageSize: number = 10,
): Promise<BasePaginationResponse<InviteRecord>> {
  const { data } = await apiClient.get<BasePaginationResponse<InviteRecord>>(
    '/user/referral/invites',
    { params: { page, page_size: pageSize } },
  )
  return data
}

/**
 * Get commission records (paginated)
 */
export async function getCommissionRecords(
  page: number = 1,
  pageSize: number = 10,
): Promise<BasePaginationResponse<CommissionRecord>> {
  const { data } = await apiClient.get<BasePaginationResponse<CommissionRecord>>(
    '/user/referral/commissions',
    { params: { page, page_size: pageSize } },
  )
  return data
}

/**
 * Withdraw commission to balance
 */
export async function withdrawCommission(req: WithdrawRequest): Promise<WithdrawResponse> {
  const { data } = await apiClient.post<WithdrawResponse>('/user/referral/withdraw', req)
  return data
}

export const referralAPI = {
  getReferralInfo,
  generateReferralCode,
  getInviteRecords,
  getCommissionRecords,
  withdrawCommission,
}

export default referralAPI
