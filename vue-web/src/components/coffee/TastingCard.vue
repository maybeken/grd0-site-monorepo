<template>
  <div class="h-fit bg-background rounded-xl overflow-hidden flex flex-col border border-accent/20">
    <div class="flex items-center p-2 gap-1 font-bold text-sm">
      <Icon v-if="tasting.pinned" icon="mynaui:pin" height="14" />
      <span>{{ formatDate(tasting.tasted_at) }}</span>
      <span class="select-none">|</span>
      <span class="opacity-70">{{ equipment?.name }}</span>
      <span v-if="tasting.rating != null" class="ml-auto">
        {{ tasting.rating }}/10
      </span>
    </div>

    <div class="w-full h-32 overflow-hidden flex-shrink-0">
      <img
        :src="equipmentImage"
        :alt="equipment?.type || 'Equipment'"
        class="w-full h-full object-contain grayscale brightness-400"
      />
    </div>

    <div class="flex flex-1 min-h-0">
      <div class="flex flex-col gap-1 p-2 w-1/2 overflow-hidden text-xs">
        <div class="font-semibold truncate">{{ tasting.bean?.name || 'IDK ¯\_(ツ)_/¯' }}</div>
        <div class="opacity-70 truncate">
          {{ tasting.bean?.origin || 'IDK ¯\_(ツ)_/¯' }}
          <span v-if="tasting.bean?.roaster"> · {{ tasting.bean.roaster }}</span>
        </div>

        <div class="mt-1 space-y-0.5">
          <div>Grind: <span :class="{ 'opacity-50': tasting.grind_size == null }">{{ tasting.grind_size || 'IDK ¯\_(ツ)_/¯' }}</span></div>
          <div>Grind Setting: <span :class="{ 'opacity-50': tasting.grind_setting == null }">{{ tasting.grind_setting ?? 'IDK ¯\_(ツ)_/¯' }}</span></div>
          <div>Dose: <span :class="{ 'opacity-50': tasting.coffee_dose == null }">{{ tasting.coffee_dose != null ? tasting.coffee_dose + 'g' : 'IDK ¯\_(ツ)_/¯' }}</span></div>
          <div>Water: <span :class="{ 'opacity-50': tasting.water_in == null }">{{ tasting.water_in != null ? tasting.water_in + 'ml' : 'IDK ¯\_(ツ)_/¯' }}</span></div>
          <div>Yield: <span :class="{ 'opacity-50': tasting.coffee_out == null }">{{ tasting.coffee_out != null ? tasting.coffee_out + 'ml' : 'IDK ¯\_(ツ)_/¯' }}</span></div>
          <div>Ratio: <span :class="{ 'opacity-50': tasting.ratio == null }">{{ tasting.ratio != null ? '1:' + tasting.ratio.toFixed(1) : 'IDK ¯\_(ツ)_/¯' }}</span></div>
          <div>Time: <span :class="{ 'opacity-50': tasting.brew_time == null }">{{ formatBrewTime(tasting.brew_time) }}</span></div>
        </div>
      </div>

      <div class="my-0 w-1/2 px-2">
        <TasteRadarChart
          :taste-fruity="tasting.taste_fruity"
          :taste-sour="tasting.taste_sour"
          :taste-sweetness="tasting.taste_sweetness"
          :taste-nutty="tasting.taste_nutty"
          :taste-spice="tasting.taste_spice"
          :taste-floral="tasting.taste_floral"
          :taste-green="tasting.taste_green"
        />
      </div>
    </div>

    <div
      v-if="tasting.overall_notes"
      class="px-2 pb-2 text-xs opacity-80 overflow-hidden max-h-16 flex-shrink-0"
    >
      {{ tasting.overall_notes }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { TastingNote } from '@/interfaces/Coffee'
import { EQUIPMENT_IMAGE_MAP, DEFAULT_EQUIPMENT_IMAGE } from '@/helpers/coffee'
import TasteRadarChart from './TasteRadarChart.vue'

interface Props {
  tasting: TastingNote
}

const props = defineProps<Props>()

const equipment = computed(() => props.tasting.equipment)

const equipmentImage = computed(() => {
  const type = equipment.value?.type
  if (type && type in EQUIPMENT_IMAGE_MAP) {
    return EQUIPMENT_IMAGE_MAP[type]
  }
  return DEFAULT_EQUIPMENT_IMAGE
})

function formatDate(dateStr: string): string {
  if (!dateStr) return 'IDK ¯\_(ツ)_/¯'
  const d = new Date(dateStr)
  return d.toLocaleDateString('en-GB', { month: 'short', day: 'numeric', year: 'numeric' })
}

function formatBrewTime(seconds: number | null | undefined): string {
  if (seconds == null) return 'IDK ¯\_(ツ)_/¯'
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins}:${String(secs).padStart(2, '0')}`
}
</script>
