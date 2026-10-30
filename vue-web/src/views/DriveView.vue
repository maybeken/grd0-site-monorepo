<template>
  <div class="h-full w-full">
    <h1 class="sr-only">{{ $t('grandTour.view.title') }}</h1>

    <GrandTourCarousel v-if="show_content" :vehicles="vehicles"></GrandTourCarousel>

    <div
      v-else-if="show_empty"
      class="flex h-full flex-col items-center justify-center gap-3 px-4 text-center"
    >
      <p class="text-2xl font-bold">{{ $t('grandTour.view.empty') }}<CursorBlink /></p>
      <p class="text-base font-medium text-foreground/60">{{ $t('grandTour.view.emptyHint') }}</p>
    </div>

    <div v-else class="flex h-full w-full overflow-hidden" aria-busy="true">
      <span class="sr-only">{{ $t('grandTour.view.loading') }}</span>
      <div
        v-for="placeholder in 2"
        :key="placeholder"
        class="flex h-full w-full shrink-0 flex-col items-center justify-center gap-8 px-4 py-6 md:px-12 md:py-10"
      >
        <div
          class="h-1/2 w-full max-w-4xl rounded-3xl bg-linear-to-r from-accent via-secondary to-accent opacity-40 blur-xs"
        ></div>
        <div class="flex w-full max-w-3xl flex-col gap-3">
          <Skeleton :loading="true" h="md" w="2/3"></Skeleton>
          <Skeleton :loading="true" h="sm" w="1/3"></Skeleton>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

import GrandTourCarousel from '@/components/drive/GrandTourCarousel.vue'
import CursorBlink from '@/components/shared/CursorBlink.vue'
import Skeleton from '@/components/shared/Skeleton.vue'
import { loadFlagCollection } from '@/helpers/icons'
import { getVehicleCollection } from '@/services/vehicle'

import type { Vehicle } from '@/interfaces/Vehicle'

const vehicles = ref<Vehicle[]>([])
const loading = ref(true)
const flags_ready = ref(false)

const show_content = computed(
  () => !loading.value && flags_ready.value && vehicles.value.length > 0
)
const show_empty = computed(
  () => !loading.value && flags_ready.value && vehicles.value.length === 0
)

async function loadVehicles(): Promise<void> {
  try {
    vehicles.value = await getVehicleCollection()
  } catch {
    vehicles.value = []
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await Promise.allSettled([loadFlagCollection(), loadVehicles()])
  flags_ready.value = true
})
</script>
