<template>
  <AppLayout>
    <div class="space-y-6">
      <!-- Loading -->
      <div v-if="loading" class="flex items-center justify-center py-20">
        <div class="h-8 w-8 animate-spin rounded-full border-4 border-primary-500 border-t-transparent" />
      </div>

      <!-- Error -->
      <div v-else-if="error" class="card p-12 text-center">
        <div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-red-100 dark:bg-red-900/30">
          <svg class="h-8 w-8 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z" />
          </svg>
        </div>
        <p class="text-gray-500 dark:text-gray-400">{{ error }}</p>
        <button class="btn btn-primary mt-4" @click="fetchData">{{ t('common.retry') }}</button>
      </div>

      <template v-else>
        <!-- Stats Summary -->
        <div class="grid grid-cols-2 gap-3 sm:grid-cols-4">
          <div class="flex items-center gap-3 rounded-xl border border-gray-200 bg-white p-4 dark:border-dark-700 dark:bg-dark-800">
            <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-primary-500/10">
              <svg class="h-5 w-5 text-primary-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M21 7.5l-9-5.25L3 7.5m18 0l-9 5.25m9-5.25v9l-9 5.25M3 7.5l9 5.25M3 7.5v9l9 5.25m0-9v9" />
              </svg>
            </div>
            <div>
              <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ allModels.length }}</p>
              <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('modelPlaza.totalModels') }}</p>
            </div>
          </div>
          <div
            v-for="p in platformList"
            :key="p"
            class="flex items-center gap-3 rounded-xl border border-gray-200 bg-white p-4 dark:border-dark-700 dark:bg-dark-800"
          >
            <div :class="['flex h-10 w-10 items-center justify-center rounded-lg', platformBgClass(p)]">
              <ModelIcon :model="platformSampleModel(p)" size="20px" />
            </div>
            <div>
              <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ platformCount(p) }}</p>
              <p class="text-xs text-gray-500 dark:text-dark-400">{{ platformLabel(p) }}</p>
            </div>
          </div>
        </div>

        <!-- Filter Bar -->
        <div class="flex flex-wrap items-center gap-3">
          <!-- Platform Tabs -->
          <div class="flex flex-1 flex-wrap gap-1 rounded-xl bg-gray-100 p-1 dark:bg-dark-800">
            <button
              v-for="tab in platformTabs"
              :key="tab"
              class="rounded-lg px-3 py-2 text-sm font-medium transition-all"
              :class="
                activePlatform === tab
                  ? 'bg-white text-gray-900 shadow dark:bg-dark-700 dark:text-white'
                  : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300'
              "
              @click="activePlatform = tab"
            >
              {{ tab === 'all' ? t('modelPlaza.all') : platformLabel(tab) }}
              <span class="ml-1 text-xs opacity-60">{{ tab === 'all' ? allModels.length : platformCount(tab) }}</span>
            </button>
          </div>

          <!-- Search -->
          <div class="relative w-full sm:w-64">
            <svg class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
            </svg>
            <input
              v-model="searchQuery"
              type="text"
              :placeholder="t('modelPlaza.searchPlaceholder')"
              class="w-full rounded-xl border border-gray-200 bg-white py-2 pl-9 pr-3 text-sm outline-none transition-colors focus:border-primary-500 focus:ring-1 focus:ring-primary-500 dark:border-dark-700 dark:bg-dark-800 dark:text-white dark:focus:border-primary-400"
            />
          </div>
        </div>

        <!-- Empty filtered state -->
        <div v-if="filteredModels.length === 0 && allModels.length > 0" class="card py-16 text-center">
          <p class="text-gray-500 dark:text-gray-400">{{ t('modelPlaza.noResults') }}</p>
        </div>

        <!-- Empty state -->
        <div v-else-if="allModels.length === 0" class="card py-16 text-center">
          <div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-gray-100 dark:bg-dark-700">
            <svg class="h-8 w-8 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21 7.5l-9-5.25L3 7.5m18 0l-9 5.25m9-5.25v9l-9 5.25M3 7.5l9 5.25M3 7.5v9l9 5.25m0-9v9" />
            </svg>
          </div>
          <p class="text-gray-500 dark:text-gray-400">{{ t('modelPlaza.noModels') }}</p>
        </div>

        <!-- Model Cards Grid -->
        <div v-else class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <div
            v-for="model in filteredModels"
            :key="model.model"
            class="group overflow-hidden rounded-2xl border border-gray-200 bg-white transition-all hover:shadow-lg hover:-translate-y-0.5 dark:border-dark-700 dark:bg-dark-800"
          >
            <!-- Model Header -->
            <div class="p-4 pb-3">
              <div class="mb-2 flex items-center gap-2.5">
                <ModelIcon :model="model.model" size="22px" />
                <h3 class="min-w-0 flex-1 truncate text-sm font-bold text-gray-900 dark:text-white" :title="model.model">
                  {{ model.model }}
                </h3>
              </div>
              <div class="flex flex-wrap items-center gap-1.5">
                <span
                  :class="[
                    'rounded-md border px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wide',
                    getPlatformTagClass(model.platform),
                  ]"
                >
                  {{ model.platform }}
                </span>
                <span
                  v-for="tag in model.tags"
                  :key="tag"
                  class="rounded-md bg-blue-50 px-2 py-0.5 text-[10px] font-medium text-blue-600 dark:bg-blue-900/30 dark:text-blue-400"
                >
                  {{ tag }}
                </span>
              </div>
            </div>

            <!-- Pricing Grid -->
            <div class="border-t border-gray-100 bg-gray-50/50 p-4 dark:border-dark-700 dark:bg-dark-800/50">
              <div class="grid grid-cols-2 gap-x-4 gap-y-2.5">
                <div v-if="model.input_price != null">
                  <p class="text-[10px] font-medium uppercase tracking-wider text-gray-400 dark:text-dark-500">
                    {{ t('modelPlaza.inputPrice') }}
                  </p>
                  <p class="text-sm font-bold text-gray-900 dark:text-white">${{ formatPrice(model.input_price) }}<span class="text-xs font-normal text-gray-400">/M</span></p>
                </div>
                <div v-if="model.output_price != null">
                  <p class="text-[10px] font-medium uppercase tracking-wider text-gray-400 dark:text-dark-500">
                    {{ t('modelPlaza.outputPrice') }}
                  </p>
                  <p class="text-sm font-bold text-gray-900 dark:text-white">${{ formatPrice(model.output_price) }}<span class="text-xs font-normal text-gray-400">/M</span></p>
                </div>
                <div v-if="model.cache_write_price != null">
                  <p class="text-[10px] font-medium uppercase tracking-wider text-gray-400 dark:text-dark-500">
                    {{ t('modelPlaza.cacheWritePrice') }}
                  </p>
                  <p class="text-sm font-bold text-gray-900 dark:text-white">${{ formatPrice(model.cache_write_price) }}<span class="text-xs font-normal text-gray-400">/M</span></p>
                </div>
                <div v-if="model.cache_read_price != null">
                  <p class="text-[10px] font-medium uppercase tracking-wider text-gray-400 dark:text-dark-500">
                    {{ t('modelPlaza.cacheReadPrice') }}
                  </p>
                  <p class="text-sm font-bold text-gray-900 dark:text-white">${{ formatPrice(model.cache_read_price) }}<span class="text-xs font-normal text-gray-400">/M</span></p>
                </div>
              </div>

              <!-- Rate Multiplier -->
              <div v-if="model.rate_multiplier !== 1" class="mt-3 flex justify-end">
                <span class="rounded-full bg-amber-100 px-2 py-0.5 text-[10px] font-semibold text-amber-700 dark:bg-amber-900/30 dark:text-amber-400">
                  ×{{ Number(model.rate_multiplier.toPrecision(10)) }}
                </span>
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
import { modelsAPI, type PublicModelPricing } from '@/api/models'
import AppLayout from '@/components/layout/AppLayout.vue'
import ModelIcon from '@/components/common/ModelIcon.vue'

const { t } = useI18n()

const loading = ref(true)
const error = ref('')
const allModels = ref<PublicModelPricing[]>([])
const searchQuery = ref('')
const activePlatform = ref('all')

function getPlatformTagClass(platform: string): string {
  switch (platform.toLowerCase()) {
    case 'anthropic':
      return 'bg-orange-500/10 text-orange-600 border-orange-500/30 dark:text-orange-400'
    case 'openai':
      return 'bg-green-500/10 text-green-600 border-green-500/30 dark:text-green-400'
    case 'google':
    case 'gemini':
      return 'bg-blue-500/10 text-blue-600 border-blue-500/30 dark:text-blue-400'
    case 'deepseek':
      return 'bg-indigo-500/10 text-indigo-600 border-indigo-500/30 dark:text-indigo-400'
    case 'meta':
      return 'bg-sky-500/10 text-sky-600 border-sky-500/30 dark:text-sky-400'
    case 'mistral':
      return 'bg-yellow-500/10 text-yellow-600 border-yellow-500/30 dark:text-yellow-400'
    default:
      return 'bg-slate-500/10 text-slate-600 border-slate-500/30 dark:text-slate-400'
  }
}

function platformBgClass(p: string): string {
  switch (p.toLowerCase()) {
    case 'anthropic': return 'bg-orange-500/10'
    case 'openai': return 'bg-green-500/10'
    case 'google': case 'gemini': return 'bg-blue-500/10'
    case 'deepseek': return 'bg-indigo-500/10'
    default: return 'bg-gray-500/10'
  }
}

function platformLabel(p: string): string {
  const labels: Record<string, string> = {
    anthropic: 'Anthropic',
    openai: 'OpenAI',
    google: 'Google',
    gemini: 'Gemini',
    deepseek: 'DeepSeek',
    meta: 'Meta',
    mistral: 'Mistral',
  }
  return labels[p.toLowerCase()] || p
}

function platformSampleModel(p: string): string {
  const samples: Record<string, string> = {
    anthropic: 'claude-sonnet',
    openai: 'gpt-4',
    google: 'gemini-pro',
    gemini: 'gemini-pro',
    deepseek: 'deepseek-chat',
    meta: 'llama-3',
    mistral: 'mistral-large',
  }
  return samples[p.toLowerCase()] || p
}

const platformList = computed(() => {
  const platforms = new Set(allModels.value.map(m => m.platform.toLowerCase()))
  return [...platforms].sort()
})

const platformTabs = computed(() => ['all', ...platformList.value])

function platformCount(p: string): number {
  return allModels.value.filter(m => m.platform.toLowerCase() === p.toLowerCase()).length
}

const filteredModels = computed(() => {
  let result = allModels.value
  if (activePlatform.value !== 'all') {
    result = result.filter(m => m.platform.toLowerCase() === activePlatform.value.toLowerCase())
  }
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.trim().toLowerCase()
    result = result.filter(m => m.model.toLowerCase().includes(q) || m.platform.toLowerCase().includes(q))
  }
  return result
})

function formatPrice(price: number | null): string {
  if (price === null || price === undefined) return '-'
  if (price >= 100) return price.toFixed(0)
  if (price >= 1) return price.toFixed(2)
  if (price >= 0.01) return price.toFixed(2)
  return price.toFixed(4)
}

async function fetchData() {
  loading.value = true
  error.value = ''
  try {
    const res = await modelsAPI.getPublicPricing()
    allModels.value = res.data.models || []
  } catch (err: unknown) {
    console.error('Failed to load model pricing:', err)
    error.value = t('modelPlaza.loadError')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>
