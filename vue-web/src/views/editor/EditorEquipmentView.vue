<template>
  <div class="flex flex-col gap-4">
    <router-link to="/editor/coffee/equipment" class="text-sm text-secondary hover:underline">
      {{ $t('editor.backToList') }}
    </router-link>

    <p v-if="loading" class="text-secondary">{{ $t('common.loading') }}</p>

    <div v-else-if="equipment" class="max-w-2xl">
      <EquipmentForm :equipment="equipment" @saved="onSaved" @cancel="back"></EquipmentForm>
    </div>

    <div v-else class="flex flex-col items-start gap-3">
      <p class="text-red-400">{{ $t('editor.notFound') }}</p>
      <router-link
        to="/editor/coffee/equipment"
        class="rounded-full border border-dotted px-3 py-1"
      >
        {{ $t('editor.backToList') }}
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

import EquipmentForm from '@/components/coffee/EquipmentForm.vue'
import { getCoffeeEquipmentItemRaw } from '@/services/coffee'
import { useToast } from '@/composables/useToast'

import type { BrewEquipment } from '@/interfaces/Coffee'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const toast = useToast()

const equipment = ref<BrewEquipment | null>(null)
const loading = ref(true)

const id = computed(() => {
  const value = route.params.id

  return typeof value === 'string' ? value : ''
})

async function load(): Promise<void> {
  loading.value = true
  equipment.value = null

  try {
    equipment.value = await getCoffeeEquipmentItemRaw(id.value)
  } catch {
    equipment.value = null
  } finally {
    loading.value = false
  }
}

function onSaved(): void {
  toast.success(t('common.done'))
}

function back(): void {
  router.push('/editor/coffee/equipment')
}

watch(id, load, { immediate: true })
</script>
