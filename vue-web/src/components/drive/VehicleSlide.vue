<template>
  <article class="flex h-full w-full flex-col items-center justify-center gap-6 px-4 py-6 md:px-12 md:py-10">
    <div class="flex min-h-0 w-full flex-1 items-center justify-center text-foreground">
      <div
        v-if="cover"
        class="pointer-events-none flex h-full w-full max-w-4xl select-none items-center justify-center"
      >
        <img
          v-if="cover.kind === 'cutout'"
          :src="cover.src"
          :alt="`${vehicle.make} ${vehicle.model}`"
          draggable="false"
          class="h-full w-full select-none object-contain"
        />
        <SilhouetteCover
          v-else
          :src="cover.src"
          :alt="`${vehicle.make} ${vehicle.model}`"
        ></SilhouetteCover>
      </div>
      <p
        v-else
        class="text-3xl font-bold uppercase tracking-[0.4em] text-foreground/40 md:text-4xl"
      >
        {{ $t('grandTour.unknownCover') }}
      </p>
    </div>

    <div class="flex w-full max-w-3xl flex-col gap-3">
      <div class="flex flex-wrap items-end justify-between gap-x-6 gap-y-2">
        <div class="flex min-w-0 flex-col">
          <p class="text-sm font-medium uppercase tracking-[0.3em] text-foreground/60">
            {{ vehicle.make }}
          </p>
          <h2 class="text-3xl font-bold leading-tight md:text-4xl">
            {{ vehicle.model }}
            <span class="font-medium text-foreground/50">{{ vehicle.year }}</span>
          </h2>
        </div>

        <p v-if="location_line" class="flex items-center gap-2 text-lg font-medium">
          <Icon
            v-if="flag"
            :icon="flag"
            dotify="false"
            height="1.25rem"
            class="shrink-0"
          ></Icon>
          <span>{{ location_line }}</span>
        </p>
      </div>

      <dl v-if="specs.length" class="grid grid-cols-2 gap-x-6 gap-y-2 md:grid-cols-4">
        <div v-for="spec in specs" :key="spec.label" class="min-w-0">
          <dt class="text-xs font-medium uppercase tracking-widest text-foreground/60">
            {{ spec.label }}
          </dt>
          <dd class="truncate text-base font-medium" :title="spec.value">{{ spec.value }}</dd>
        </div>
      </dl>

      <p v-if="vehicle.notes" class="max-w-2xl text-base leading-relaxed text-foreground/80">
        {{ vehicle.notes }}
      </p>
    </div>
  </article>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import SilhouetteCover from '@/components/vehicle/SilhouetteCover.vue'
import Icon from '@/components/shared/Icon.vue'
import { assetUrl, flagIcon } from '@/helpers/vehicle'

import type { Vehicle } from '@/interfaces/Vehicle'

const props = defineProps<{
  vehicle: Vehicle
}>()

const { t } = useI18n()

const cover = computed<{ src: string; kind: 'cutout' | 'silhouette' } | null>(() => {
  if (props.vehicle.cutout_key) {
    return { src: assetUrl(props.vehicle.cutout_key), kind: 'cutout' }
  }

  if (props.vehicle.silhouette_key) {
    return { src: assetUrl(props.vehicle.silhouette_key), kind: 'silhouette' }
  }

  return null
})
const flag = computed(() => flagIcon(props.vehicle.country_code))

const location_line = computed(() => {
  return [props.vehicle.country, props.vehicle.city].filter(Boolean).join(' · ')
})

const specs = computed(() => {
  const vehicle = props.vehicle
  const entries: { label: string; value: string | null }[] = [
    { label: t('grandTour.details.engine'), value: vehicle.engine ?? null },
    { label: t('grandTour.details.power'), value: vehicle.power ?? null },
    { label: t('grandTour.details.transmission'), value: vehicle.transmission ?? null },
    { label: t('grandTour.details.drivetrain'), value: vehicle.drivetrain ?? null },
    { label: t('grandTour.details.fuel'), value: vehicle.fuel ?? null },
    { label: t('grandTour.details.drivenAt'), value: formatDate(vehicle.driven_at) },
    {
      label: t('grandTour.details.rating'),
      value: vehicle.rating ? '★'.repeat(Math.max(0, Math.min(5, vehicle.rating))) : null
    }
  ]

  return entries.filter(
    (entry): entry is { label: string; value: string } =>
      typeof entry.value === 'string' && entry.value.trim() !== ''
  )
})

function formatDate(value?: string | null): string | null {
  if (!value) return null

  const date = new Date(value)

  if (Number.isNaN(date.getTime())) return null

  return date.toLocaleDateString(undefined, { year: 'numeric', month: 'short' })
}
</script>
