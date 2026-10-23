<template>
  <div class="flex flex-col gap-4">
    <router-link to="/editor/coffee/tastings" class="text-sm text-secondary hover:underline">
      {{ $t('editor.backToList') }}
    </router-link>

    <p v-if="loading" class="text-secondary">{{ $t('common.loading') }}</p>

    <div v-else-if="tasting" class="max-w-4xl">
      <TastingForm
        :tasting="tasting"
        :beans="beans"
        :equipment-list="equipment"
        @saved="onSaved"
        @cancel="back"
        @bean-added="fetchBeans"
        @equipment-added="fetchEquipment"
      ></TastingForm>
    </div>

    <div v-else class="flex flex-col items-start gap-3">
      <p class="text-red-400">{{ $t('editor.notFound') }}</p>
      <router-link to="/editor/coffee/tastings" class="rounded-full border border-dotted px-3 py-1">
        {{ $t('editor.backToList') }}
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

import TastingForm from '@/components/coffee/TastingForm.vue'
import {
  getCoffeeBeansRaw,
  getCoffeeEquipmentRaw,
  getCoffeeTastingRaw
} from '@/services/coffee'
import { useToast } from '@/composables/useToast'

import type { BrewEquipment, CoffeeBean, TastingNote } from '@/interfaces/Coffee'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const toast = useToast()

const tasting = ref<TastingNote | null>(null)
const beans = ref<CoffeeBean[]>([])
const equipment = ref<BrewEquipment[]>([])
const loading = ref(true)

const id = computed(() => {
  const value = route.params.id

  return typeof value === 'string' ? value : ''
})

async function fetchBeans(): Promise<void> {
  beans.value = await getCoffeeBeansRaw()
}

async function fetchEquipment(): Promise<void> {
  equipment.value = await getCoffeeEquipmentRaw()
}

async function load(): Promise<void> {
  loading.value = true
  tasting.value = null

  try {
    const [tasting_value] = await Promise.all([
      getCoffeeTastingRaw(id.value),
      fetchBeans(),
      fetchEquipment()
    ])

    tasting.value = tasting_value
  } catch {
    tasting.value = null
  } finally {
    loading.value = false
  }
}

function onSaved(): void {
  toast.success(t('common.done'))
}

function back(): void {
  router.push('/editor/coffee/tastings')
}

watch(id, load, { immediate: true })
</script>
