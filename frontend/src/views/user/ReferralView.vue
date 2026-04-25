<template>
  <AppLayout>
    <div class="space-y-6">
      <!-- Loading -->
      <div v-if="loading" class="flex items-center justify-center py-12">
        <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
      </div>

      <template v-else-if="referralInfo && !referralInfo.enabled">
        <!-- Referral Disabled Banner -->
        <div class="flex flex-col items-center justify-center py-16">
          <div class="flex h-16 w-16 items-center justify-center rounded-2xl bg-gray-100 dark:bg-dark-700">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" class="h-8 w-8 text-gray-400 dark:text-dark-500">
              <path stroke-linecap="round" stroke-linejoin="round" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
            </svg>
          </div>
          <h3 class="mt-4 text-lg font-semibold text-gray-900 dark:text-white">{{ t('referral.referralDisabled') }}</h3>
          <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">{{ t('referral.referralDisabledDesc') }}</p>
        </div>
      </template>

      <template v-else>
        <!-- Header Banner with gradient -->
        <div class="referral-banner">
          <div class="referral-banner-bg"></div>
          <div class="relative z-10 flex flex-col items-start gap-4 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <h2 class="text-xl font-bold text-white sm:text-2xl">{{ t('referral.title') }}</h2>
              <p class="mt-1 text-sm text-emerald-100/80">{{ t('referral.description') }}</p>
            </div>
            <div class="flex items-center gap-2">
              <button
                v-if="referralInfo?.referral_code"
                class="btn-glass"
                @click="copyReferralLink"
              >
                <Icon name="link" size="sm" class="mr-1.5" />
                {{ t('referral.copyLink') }}
              </button>
            </div>
          </div>
        </div>

        <!-- Referral Code Card -->
        <div class="card p-5">
          <div class="flex items-center gap-3 mb-4">
            <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-emerald-500/10 text-emerald-500">
              <Icon name="key" size="md" />
            </div>
            <div>
              <h3 class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('referral.myReferralCode') }}</h3>
            </div>
          </div>

          <div v-if="referralInfo?.referral_code" class="flex flex-col gap-3 sm:flex-row sm:items-center">
            <div class="referral-code-box">
              <span class="font-mono text-lg font-bold tracking-wider text-emerald-600 dark:text-emerald-400">
                {{ referralInfo.referral_code }}
              </span>
            </div>
            <div class="flex gap-2">
              <button class="btn btn-primary btn-sm" @click="copyReferralCode">
                <Icon name="copy" size="sm" class="mr-1" />
                {{ t('referral.copyCode') }}
              </button>
              <button class="btn btn-secondary btn-sm" @click="regenerateCode" :disabled="regenerating">
                <Icon name="refresh" size="sm" class="mr-1" :class="{ 'animate-spin': regenerating }" />
                {{ t('referral.regenerateCode') }}
              </button>
            </div>
          </div>
          <div v-else class="flex items-center gap-3">
            <button class="btn btn-primary" @click="generateCode" :disabled="generating">
              <Icon name="plus" size="sm" class="mr-1.5" />
              {{ t('referral.generateCode') }}
            </button>
          </div>
        </div>

        <!-- Stats Grid -->
        <div class="grid grid-cols-2 gap-4 lg:grid-cols-4">
          <div class="referral-stat-card">
            <div class="referral-stat-icon bg-emerald-500/10 text-emerald-500">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" class="h-5 w-5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v12m-3-2.818l.879.659c1.171.879 3.07.879 4.242 0 1.172-.879 1.172-2.303 0-3.182C13.536 12.219 12.768 12 12 12c-.725 0-1.45-.22-2.003-.659-1.106-.879-1.106-2.303 0-3.182s2.9-.879 4.006 0l.415.33M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div>
              <p class="referral-stat-label">{{ t('referral.totalCommission') }}</p>
              <p class="referral-stat-value text-emerald-600 dark:text-emerald-400">
                ${{ formatMoney(referralInfo?.total_commission ?? 0) }}
              </p>
            </div>
          </div>

          <div class="referral-stat-card">
            <div class="referral-stat-icon bg-blue-500/10 text-blue-500">
              <Icon name="users" size="sm" />
            </div>
            <div>
              <p class="referral-stat-label">{{ t('referral.totalInvitees') }}</p>
              <p class="referral-stat-value text-blue-600 dark:text-blue-400">
                {{ referralInfo?.total_invitees ?? 0 }}
              </p>
            </div>
          </div>

          <div class="referral-stat-card">
            <div class="referral-stat-icon bg-violet-500/10 text-violet-500">
              <Icon name="trendingUp" size="sm" />
            </div>
            <div>
              <p class="referral-stat-label">{{ t('referral.totalRecharged') }}</p>
              <p class="referral-stat-value text-violet-600 dark:text-violet-400">
                ${{ formatMoney(referralInfo?.total_recharged ?? 0) }}
              </p>
            </div>
          </div>

          <div class="referral-stat-card">
            <div class="referral-stat-icon bg-amber-500/10 text-amber-500">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" class="h-5 w-5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 18.75a60.07 60.07 0 0115.797 2.101c.727.198 1.453-.342 1.453-1.096V18.75M3.75 4.5v.75A.75.75 0 013 6h-.75m0 0v-.375c0-.621.504-1.125 1.125-1.125H20.25M2.25 6v9m18-10.5v.75c0 .414.336.75.75.75h.75m-1.5-1.5h.375c.621 0 1.125.504 1.125 1.125v9.75c0 .621-.504 1.125-1.125 1.125h-.375m1.5-1.5H21a.75.75 0 00-.75.75v.75m0 0H3.75m0 0h-.375a1.125 1.125 0 01-1.125-1.125V15m1.5 1.5v-.75A.75.75 0 003 15h-.75M15 10.5a3 3 0 11-6 0 3 3 0 016 0zm3 0h.008v.008H18V10.5zm-12 0h.008v.008H6V10.5z" />
              </svg>
            </div>
            <div>
              <p class="referral-stat-label">{{ t('referral.commissionRate') }}</p>
              <p class="referral-stat-value text-amber-600 dark:text-amber-400">
                {{ commissionRateRange }}
              </p>
            </div>
          </div>
        </div>

        <!-- Commission Tiers Card -->
        <div class="card p-5">
          <div class="flex items-center gap-3 mb-5">
            <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-violet-500/10 text-violet-500">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" class="h-5 w-5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z" />
              </svg>
            </div>
            <div>
              <h3 class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('referral.commissionTiers') }}</h3>
              <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('referral.commissionTiersDesc') }}</p>
            </div>
          </div>

          <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
            <div
              v-for="(tier, idx) in tierCards"
              :key="idx"
              :class="[
                'referral-tier-card',
                currentTierIndex === idx ? 'referral-tier-card-active referral-tier-card-active-' + tier.color : ''
              ]"
            >
              <div class="referral-tier-header">
                <span :class="['referral-tier-badge', 'referral-tier-badge-' + tier.color]">
                  {{ tier.label }}
                </span>
                <span v-if="currentTierIndex === idx" class="referral-tier-current">
                  ✦ {{ t('referral.currentTier') }}
                </span>
              </div>
              <div :class="['referral-tier-rate', 'referral-tier-rate-' + tier.color]">
                {{ tier.rateDisplay }}<span class="text-2xl">%</span>
              </div>
              <p class="referral-tier-condition">{{ tier.condition }}</p>
              <div class="referral-tier-bar">
                <div :class="['referral-tier-bar-fill', 'referral-tier-bar-' + tier.color]" :style="{ width: tierProgress(idx) + '%' }"></div>
              </div>
            </div>
          </div>

          <div class="mt-4 flex items-center gap-2 rounded-lg bg-gray-50 px-4 py-3 dark:bg-dark-800/50">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" class="h-4 w-4 shrink-0 text-gray-400">
              <path stroke-linecap="round" stroke-linejoin="round" d="M11.25 11.25l.041-.02a.75.75 0 011.063.852l-.708 2.836a.75.75 0 001.063.853l.041-.021M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9-3.75h.008v.008H12V8.25z" />
            </svg>
            <span class="text-xs text-gray-500 dark:text-dark-400">
              {{ t('referral.monthlyRechargeNote') }}
            </span>
          </div>
        </div>

        <!-- Withdraw Card -->
        <div v-if="(referralInfo?.total_commission ?? 0) > 0" class="card p-5">
          <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
            <div class="flex items-center gap-3">
              <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-emerald-500/10 text-emerald-500">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" class="h-5 w-5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 18.75a60.07 60.07 0 0115.797 2.101c.727.198 1.453-.342 1.453-1.096V18.75M3.75 4.5v.75A.75.75 0 013 6h-.75m0 0v-.375c0-.621.504-1.125 1.125-1.125H20.25M2.25 6v9m18-10.5v.75c0 .414.336.75.75.75h.75m-1.5-1.5h.375c.621 0 1.125.504 1.125 1.125v9.75c0 .621-.504 1.125-1.125 1.125h-.375m1.5-1.5H21a.75.75 0 00-.75.75v.75m0 0H3.75m0 0h-.375a1.125 1.125 0 01-1.125-1.125V15m1.5 1.5v-.75A.75.75 0 003 15h-.75M15 10.5a3 3 0 11-6 0 3 3 0 016 0zm3 0h.008v.008H18V10.5zm-12 0h.008v.008H6V10.5z" />
                </svg>
              </div>
              <div>
                <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('referral.availableBalance') }}</p>
                <p class="text-xl font-bold text-emerald-600 dark:text-emerald-400">
                  ${{ formatMoney(referralInfo?.total_commission ?? 0) }}
                </p>
              </div>
            </div>
            <button
              class="btn btn-primary"
              @click="handleWithdraw"
              :disabled="withdrawing || (referralInfo?.total_commission ?? 0) <= 0"
            >
              <Icon name="download" size="sm" class="mr-1.5" />
              {{ withdrawing ? '...' : t('referral.withdrawToBalance') }}
            </button>
          </div>
        </div>

        <!-- How It Works -->
        <div class="card p-5">
          <h3 class="mb-4 text-base font-semibold text-gray-900 dark:text-white">
            {{ t('referral.howItWorks') }}
          </h3>
          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <div v-for="(step, idx) in steps" :key="idx" class="referral-step-card">
              <div class="referral-step-number">{{ idx + 1 }}</div>
              <h4 class="mt-3 text-sm font-semibold text-gray-900 dark:text-white">{{ step.title }}</h4>
              <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">{{ step.desc }}</p>
            </div>
          </div>
        </div>

        <!-- Tabs: Invite Records / Commission Details -->
        <div class="card overflow-hidden">
          <div class="border-b border-gray-200 dark:border-dark-700">
            <nav class="flex -mb-px">
              <button
                v-for="tab in tabs"
                :key="tab.key"
                :class="[
                  'referral-tab',
                  activeTab === tab.key ? 'referral-tab-active' : 'referral-tab-inactive'
                ]"
                @click="activeTab = tab.key"
              >
                {{ tab.label }}
              </button>
            </nav>
          </div>

          <!-- Invite Records Table -->
          <div v-if="activeTab === 'invites'" class="overflow-x-auto">
            <div v-if="inviteRecords.length === 0" class="flex flex-col items-center justify-center py-12">
              <Icon name="users" size="xl" class="mb-3 text-gray-300 dark:text-dark-600" />
              <p class="text-sm font-medium text-gray-500 dark:text-dark-400">{{ t('referral.noInvites') }}</p>
              <p class="mt-1 text-xs text-gray-400 dark:text-dark-500">{{ t('referral.noInvitesDesc') }}</p>
            </div>
            <table v-else class="w-full text-left text-sm">
              <thead>
                <tr class="border-b border-gray-100 dark:border-dark-700">
                  <th class="px-5 py-3 text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('referral.inviteeEmail') }}</th>
                  <th class="px-5 py-3 text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('referral.registeredAt') }}</th>
                  <th class="px-5 py-3 text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('referral.rechargeAmount') }}</th>
                  <th class="px-5 py-3 text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('referral.commissionAmount') }}</th>
                  <th class="px-5 py-3 text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('referral.commissionStatus') }}</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100 dark:divide-dark-700">
                <tr v-for="record in inviteRecords" :key="record.id" class="hover:bg-gray-50 dark:hover:bg-dark-750">
                  <td class="px-5 py-3">
                    <div class="flex items-center gap-2">
                      <div class="flex h-7 w-7 items-center justify-center rounded-full bg-emerald-500/10 text-xs font-medium text-emerald-600 dark:text-emerald-400">
                        {{ (record.invitee_email || record.invitee_username || '?')[0].toUpperCase() }}
                      </div>
                      <span class="text-gray-900 dark:text-white">{{ record.invitee_email || record.invitee_username }}</span>
                    </div>
                  </td>
                  <td class="px-5 py-3 text-gray-500 dark:text-dark-400">{{ formatDate(record.registered_at) }}</td>
                  <td class="px-5 py-3 font-medium text-gray-900 dark:text-white">${{ formatMoney(record.total_recharged) }}</td>
                  <td class="px-5 py-3 font-medium text-emerald-600 dark:text-emerald-400">${{ formatMoney(record.total_commission) }}</td>
                  <td class="px-5 py-3">
                    <span :class="statusClass(record.status)">
                      {{ statusLabel(record.status) }}
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
            <div v-if="invitePagination.pages > 1" class="border-t border-gray-100 px-5 py-3 dark:border-dark-700">
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500 dark:text-dark-400">
                  {{ invitePagination.total }} {{ t('referral.totalInvitees').toLowerCase() }}
                </span>
                <div class="flex gap-1">
                  <button
                    v-for="p in invitePagination.pages"
                    :key="p"
                    :class="['referral-page-btn', invitePage === p ? 'referral-page-btn-active' : '']"
                    @click="invitePage = p; loadInvites()"
                  >{{ p }}</button>
                </div>
              </div>
            </div>
          </div>

          <!-- Commission Records Table -->
          <div v-if="activeTab === 'commissions'" class="overflow-x-auto">
            <div v-if="commissionRecords.length === 0" class="flex flex-col items-center justify-center py-12">
              <Icon name="chart" size="xl" class="mb-3 text-gray-300 dark:text-dark-600" />
              <p class="text-sm font-medium text-gray-500 dark:text-dark-400">{{ t('referral.noCommissions') }}</p>
              <p class="mt-1 text-xs text-gray-400 dark:text-dark-500">{{ t('referral.noCommissionsDesc') }}</p>
            </div>
            <table v-else class="w-full text-left text-sm">
              <thead>
                <tr class="border-b border-gray-100 dark:border-dark-700">
                  <th class="px-5 py-3 text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('referral.inviteeEmail') }}</th>
                  <th class="px-5 py-3 text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('referral.orderId') }}</th>
                  <th class="px-5 py-3 text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('referral.rechargeAmount') }}</th>
                  <th class="px-5 py-3 text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('referral.commissionRate') }}</th>
                  <th class="px-5 py-3 text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('referral.commissionAmount') }}</th>
                  <th class="px-5 py-3 text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('common.time', '时间') }}</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100 dark:divide-dark-700">
                <tr v-for="record in commissionRecords" :key="record.id" class="hover:bg-gray-50 dark:hover:bg-dark-750">
                  <td class="px-5 py-3 text-gray-900 dark:text-white">{{ record.invitee_email || record.invitee_username }}</td>
                  <td class="px-5 py-3">
                    <span class="font-mono text-xs text-gray-500 dark:text-dark-400">{{ record.order_id }}</span>
                  </td>
                  <td class="px-5 py-3 font-medium text-gray-900 dark:text-white">${{ formatMoney(record.recharge_amount) }}</td>
                  <td class="px-5 py-3 text-gray-500 dark:text-dark-400">{{ (record.commission_rate * 100).toFixed(0) }}%</td>
                  <td class="px-5 py-3 font-medium text-emerald-600 dark:text-emerald-400">+${{ formatMoney(record.commission_amount) }}</td>
                  <td class="px-5 py-3 text-gray-500 dark:text-dark-400">{{ formatDate(record.created_at) }}</td>
                </tr>
              </tbody>
            </table>
            <div v-if="commissionPagination.pages > 1" class="border-t border-gray-100 px-5 py-3 dark:border-dark-700">
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500 dark:text-dark-400">
                  {{ commissionPagination.total }} {{ t('referral.commissionRecords').toLowerCase() }}
                </span>
                <div class="flex gap-1">
                  <button
                    v-for="p in commissionPagination.pages"
                    :key="p"
                    :class="['referral-page-btn', commissionPage === p ? 'referral-page-btn-active' : '']"
                    @click="commissionPage = p; loadCommissions()"
                  >{{ p }}</button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import {
  getReferralInfo,
  generateReferralCode,
  getInviteRecords,
  getCommissionRecords,
  withdrawCommission,
  type ReferralInfo,
  type InviteRecord,
  type CommissionRecord,
  type CommissionTier,
} from '@/api/referral'

const { t } = useI18n()
const appStore = useAppStore()

const loading = ref(true)
const generating = ref(false)
const regenerating = ref(false)
const withdrawing = ref(false)

const referralInfo = ref<ReferralInfo | null>(null)
const activeTab = ref<'invites' | 'commissions'>('invites')

const inviteRecords = ref<InviteRecord[]>([])
const invitePage = ref(1)
const invitePagination = ref({ total: 0, pages: 1 })

const commissionRecords = ref<CommissionRecord[]>([])
const commissionPage = ref(1)
const commissionPagination = ref({ total: 0, pages: 1 })

const tabs = computed(() => [
  { key: 'invites' as const, label: t('referral.inviteRecords') },
  { key: 'commissions' as const, label: t('referral.commissionRecords') },
])

const steps = computed(() => [
  { title: t('referral.step1Title'), desc: t('referral.step1Desc') },
  { title: t('referral.step2Title'), desc: t('referral.step2Desc') },
  { title: t('referral.step3Title'), desc: t('referral.step3Desc') },
  { title: t('referral.step4Title'), desc: t('referral.step4Desc') },
])

// Commission tier thresholds: <500 => 10%, 500–3000 => 12%, 3000+ => 15%
const tierThresholds = [
  { min: 0, max: 500, rate: 10 },
  { min: 500, max: 3000, rate: 12 },
  { min: 3000, max: Infinity, rate: 15 },
]

const tierCards = computed(() => [
  { label: t('referral.tierBronze'), color: 'bronze', rateDisplay: '10', condition: t('referral.tierBronzeCondition') },
  { label: t('referral.tierSilver'), color: 'silver', rateDisplay: '12', condition: t('referral.tierSilverCondition') },
  { label: t('referral.tierGold'), color: 'gold', rateDisplay: '15', condition: t('referral.tierGoldCondition') },
])

const currentTierIndex = computed(() => {
  const rate = (referralInfo.value?.commission_rate ?? 0) * 100
  if (rate >= 15) return 2
  if (rate >= 12) return 1
  return 0
})

const commissionRateRange = computed(() => {
  const tiers = referralInfo.value?.commission_tiers
  if (!tiers || tiers.length === 0) {
    return ((referralInfo.value?.commission_rate ?? 0) * 100).toFixed(0) + '%'
  }
  const rates = tiers.map(t => (t.rate * 100).toFixed(0))
  const min = rates[0]
  const max = rates[rates.length - 1]
  return min === max ? `${min}%` : `${min}%–${max}%`
})

function tierProgress(tierIdx: number): number {
  const recharged = referralInfo.value?.total_recharged ?? 0
  const tier = tierThresholds[tierIdx]
  if (recharged >= tier.max) return 100
  if (recharged <= tier.min) return 0
  return Math.min(100, ((recharged - tier.min) / (tier.max - tier.min)) * 100)
}

function formatMoney(val: number): string {
  return val.toFixed(2)
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString(undefined, { year: 'numeric', month: '2-digit', day: '2-digit' })
}

function statusClass(status: string): string {
  const base = 'inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium'
  switch (status) {
    case 'active':
    case 'completed':
      return `${base} bg-emerald-50 text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-400`
    case 'pending':
      return `${base} bg-amber-50 text-amber-700 dark:bg-amber-500/10 dark:text-amber-400`
    case 'disabled':
    case 'cancelled':
      return `${base} bg-red-50 text-red-700 dark:bg-red-500/10 dark:text-red-400`
    default:
      return `${base} bg-gray-50 text-gray-700 dark:bg-dark-700 dark:text-dark-400`
  }
}

function statusLabel(status: string): string {
  const key = `referral.status${status.charAt(0).toUpperCase() + status.slice(1)}` as any
  return t(key, status)
}

async function loadReferralInfo() {
  try {
    referralInfo.value = await getReferralInfo()
  } catch (e) {
    console.error('Failed to load referral info:', e)
  }
}

async function loadInvites() {
  try {
    const res = await getInviteRecords(invitePage.value, 10)
    inviteRecords.value = res.items
    invitePagination.value = { total: res.total, pages: res.pages }
  } catch (e) {
    console.error('Failed to load invite records:', e)
  }
}

async function loadCommissions() {
  try {
    const res = await getCommissionRecords(commissionPage.value, 10)
    commissionRecords.value = res.items
    commissionPagination.value = { total: res.total, pages: res.pages }
  } catch (e) {
    console.error('Failed to load commission records:', e)
  }
}

async function generateCode() {
  generating.value = true
  try {
    const res = await generateReferralCode()
    if (referralInfo.value) {
      referralInfo.value.referral_code = res.referral_code
    }
    appStore.showSuccess(t('referral.codeCopied'))
  } catch (e: any) {
    appStore.showError(e.response?.data?.detail || 'Failed to generate code')
  } finally {
    generating.value = false
  }
}

async function regenerateCode() {
  regenerating.value = true
  try {
    const res = await generateReferralCode()
    if (referralInfo.value) {
      referralInfo.value.referral_code = res.referral_code
    }
    appStore.showSuccess(t('referral.codeCopied'))
  } catch (e: any) {
    appStore.showError(e.response?.data?.detail || 'Failed to regenerate code')
  } finally {
    regenerating.value = false
  }
}

async function copyReferralCode() {
  if (!referralInfo.value?.referral_code) return
  try {
    await navigator.clipboard.writeText(referralInfo.value.referral_code)
    appStore.showSuccess(t('referral.codeCopied'))
  } catch {
    appStore.showError('Failed to copy')
  }
}

async function copyReferralLink() {
  if (!referralInfo.value?.referral_code) return
  const base = window.location.origin
  const link = `${base}/register?ref=${referralInfo.value.referral_code}`
  try {
    await navigator.clipboard.writeText(link)
    appStore.showSuccess(t('referral.linkCopied'))
  } catch {
    appStore.showError('Failed to copy')
  }
}

async function handleWithdraw() {
  const amount = referralInfo.value?.total_commission ?? 0
  if (amount <= 0) return
  withdrawing.value = true
  try {
    await withdrawCommission({ amount })
    appStore.showSuccess(t('referral.withdrawSuccess'))
    await loadReferralInfo()
  } catch (e: any) {
    appStore.showError(e.response?.data?.detail || 'Withdraw failed')
  } finally {
    withdrawing.value = false
  }
}

onMounted(async () => {
  loading.value = true
  try {
    await Promise.all([loadReferralInfo(), loadInvites(), loadCommissions()])
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.referral-banner {
  @apply relative overflow-hidden rounded-2xl p-6;
  background: linear-gradient(135deg, #059669 0%, #10b981 50%, #34d399 100%);
}
.referral-banner-bg {
  @apply pointer-events-none absolute inset-0;
  background: radial-gradient(circle at 80% 20%, rgba(255,255,255,0.15) 0%, transparent 50%),
              radial-gradient(circle at 20% 80%, rgba(255,255,255,0.08) 0%, transparent 40%);
}

.btn-glass {
  @apply inline-flex items-center rounded-lg px-3 py-1.5 text-sm font-medium text-white transition;
  background: rgba(255,255,255,0.15);
  backdrop-filter: blur(8px);
}
.btn-glass:hover {
  background: rgba(255,255,255,0.25);
}

.referral-code-box {
  @apply flex items-center rounded-xl border-2 border-dashed border-emerald-300 bg-emerald-50 px-5 py-3 dark:border-emerald-700 dark:bg-emerald-900/20;
}

/* Stat cards */
.referral-stat-card {
  @apply flex items-center gap-3 rounded-xl border border-gray-100 bg-white p-4 dark:border-dark-700 dark:bg-dark-800;
}
.referral-stat-icon {
  @apply flex h-10 w-10 shrink-0 items-center justify-center rounded-xl;
}
.referral-stat-label {
  @apply text-xs text-gray-500 dark:text-dark-400;
}
.referral-stat-value {
  @apply text-lg font-bold;
}

/* Step cards */
.referral-step-card {
  @apply relative rounded-xl border border-gray-100 bg-gray-50 p-4 dark:border-dark-700 dark:bg-dark-800/50;
}
.referral-step-number {
  @apply flex h-7 w-7 items-center justify-center rounded-full bg-emerald-500 text-xs font-bold text-white;
}

/* Tabs */
.referral-tab {
  @apply px-5 py-3 text-sm font-medium transition;
}
.referral-tab-active {
  @apply border-b-2 border-emerald-500 text-emerald-600 dark:text-emerald-400;
}
.referral-tab-inactive {
  @apply border-b-2 border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 dark:text-dark-400 dark:hover:text-dark-300;
}

/* Pagination */
.referral-page-btn {
  @apply flex h-7 w-7 items-center justify-center rounded text-xs font-medium text-gray-500 transition hover:bg-gray-100 dark:text-dark-400 dark:hover:bg-dark-700;
}
.referral-page-btn-active {
  @apply bg-emerald-500 text-white hover:bg-emerald-600 dark:hover:bg-emerald-600;
}

/* Tier cards */
.referral-tier-card {
  @apply relative rounded-xl border border-gray-100 bg-white p-5 transition-all dark:border-dark-700 dark:bg-dark-800;
}
.referral-tier-card:hover {
  @apply shadow-md;
  transform: translateY(-2px);
}
.referral-tier-card-active {
  @apply ring-2 shadow-lg;
}
.referral-tier-card-active-bronze {
  @apply ring-amber-400/50;
}
.referral-tier-card-active-silver {
  @apply ring-blue-400/50;
}
.referral-tier-card-active-gold {
  @apply ring-violet-400/50;
}
.referral-tier-header {
  @apply flex items-center justify-between mb-3;
}
.referral-tier-badge {
  @apply inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-semibold;
}
.referral-tier-badge-bronze {
  @apply bg-amber-50 text-amber-700 dark:bg-amber-500/10 dark:text-amber-400;
}
.referral-tier-badge-silver {
  @apply bg-blue-50 text-blue-700 dark:bg-blue-500/10 dark:text-blue-400;
}
.referral-tier-badge-gold {
  @apply bg-violet-50 text-violet-700 dark:bg-violet-500/10 dark:text-violet-400;
}
.referral-tier-current {
  @apply text-xs font-medium text-emerald-500;
}
.referral-tier-rate {
  @apply text-4xl font-extrabold tracking-tight mb-2;
}
.referral-tier-rate-bronze {
  @apply text-amber-600 dark:text-amber-400;
}
.referral-tier-rate-silver {
  @apply text-blue-600 dark:text-blue-400;
}
.referral-tier-rate-gold {
  @apply text-violet-600 dark:text-violet-400;
}
.referral-tier-condition {
  @apply text-xs text-gray-500 dark:text-dark-400 mb-3;
}
.referral-tier-bar {
  @apply h-1.5 w-full rounded-full bg-gray-100 dark:bg-dark-700 overflow-hidden;
}
.referral-tier-bar-fill {
  @apply h-full rounded-full transition-all duration-500;
}
.referral-tier-bar-bronze {
  @apply bg-amber-400;
}
.referral-tier-bar-silver {
  @apply bg-blue-400;
}
.referral-tier-bar-gold {
  @apply bg-violet-400;
}
</style>
