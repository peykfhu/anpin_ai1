<template>
  <div
    :class="[
      'group relative flex flex-col overflow-hidden rounded-2xl border transition-all',
      'hover:shadow-lg hover:-translate-y-0.5',
      selected
        ? 'border-primary-500 ring-2 ring-primary-500/20 dark:border-primary-400 dark:ring-primary-400/20'
        : 'border-gray-200 dark:border-dark-700',
      'bg-white dark:bg-dark-800',
    ]"
  >
    <div class="h-1.5 bg-gradient-to-r from-primary-400 to-primary-500" />

    <div class="flex flex-1 flex-col p-4">
      <div class="mb-2 flex items-start justify-between gap-2">
        <div class="min-w-0 flex-1">
          <div class="flex items-center gap-2">
            <h3 class="truncate text-base font-bold text-gray-900 dark:text-white">
              {{ title }}
            </h3>
            <span class="shrink-0 rounded-full bg-primary-500/10 px-2 py-0.5 text-[11px] font-medium text-primary-600 dark:text-primary-400">
              {{ t('payment.recharge.badge') }}
            </span>
          </div>
        </div>
      </div>

      <p class="mb-3 text-sm font-medium text-green-600 dark:text-green-400">
        +${{ creditedAmount }} {{ t('payment.recharge.balance') }}
      </p>

      <div class="flex-1" />

      <div class="mb-3 text-right">
        <div class="flex items-baseline justify-end gap-1">
          <span class="text-xs text-gray-400 dark:text-dark-500">¥</span>
          <span class="text-2xl font-extrabold tracking-tight text-gray-900 dark:text-white">
            {{ displayPrice }}
          </span>
        </div>
        <p v-if="showMultiplier" class="mt-0.5 text-[11px] text-gray-400 dark:text-dark-500">
          {{ t('payment.rechargeRatePreview', { usd: multiplier.toFixed(2) }) }}
        </p>
      </div>

      <button
        type="button"
        class="w-full rounded-xl bg-primary-500 py-2.5 text-sm font-semibold text-white transition-all hover:bg-primary-600 active:scale-[0.98] dark:bg-primary-600 dark:hover:bg-primary-500"
        @click="emit('select', amount)"
      >
        {{ t('payment.recharge.buyNow') }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps<{
  amount: number
  multiplier: number
  feeRate: number
  selected?: boolean
}>()

const emit = defineEmits<{ select: [amount: number] }>()
const { t } = useI18n()

const title = computed(() =>
  `${props.amount}${t('payment.recharge.usdSuffix')}`
)

const creditedAmount = computed(() =>
  (Math.round(props.amount * props.multiplier * 100) / 100).toFixed(2)
)

const feeAmount = computed(() =>
  props.feeRate > 0
    ? Math.ceil(((props.amount * props.feeRate) / 100) * 100) / 100
    : 0
)

const displayPrice = computed(() =>
  (props.amount + feeAmount.value).toFixed(2)
)

const showMultiplier = computed(() => props.multiplier !== 1)
</script>
