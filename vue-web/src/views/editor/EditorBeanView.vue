<template>
  <div class="flex flex-col gap-4">
    <router-link to="/editor/coffee/beans" class="text-sm text-secondary hover:underline">
      {{ $t('editor.backToList') }}
    </router-link>

    <p v-if="loading" class="text-secondary">{{ $t('common.loading') }}</p>

    <div v-else-if="bean" class="max-w-2xl">
      <BeanForm :bean="bean" @saved="onSaved" @cancel="back"></BeanForm>
    </div>

    <div v-else class="flex flex-col items-start gap-3">
      <p class="text-red-400">{{ $t('editor.notFound') }}</p>
      <router-link to="/editor/coffee/beans" class="rounded-full border border-dotted px-3 py-1">
        {{ $t('editor.backToList') }}
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

import BeanForm from '@/components/coffee/BeanForm.vue'
import { getCoffeeBeanRaw } from '@/services/coffee'
import { useToast } from '@/composables/useToast'

import type { CoffeeBean } from '@/interfaces/Coffee'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const toast = useToast()

const bean = ref<CoffeeBean | null>(null)
const loading = ref(true)

const id = computed(() => {
  const value = route.params.id

  return typeof value === 'string' ? value : ''
})

async function load(): Promise<void> {
  loading.value = true
  bean.value = null

  try {
    bean.value = await getCoffeeBeanRaw(id.value)
  } catch {
    bean.value = null
  } finally {
    loading.value = false
  }
}

function onSaved(): void {
  toast.success(t('common.done'))
}

function back(): void {
  router.push('/editor/coffee/beans')
}

watch(id, load, { immediate: true })
</script>
