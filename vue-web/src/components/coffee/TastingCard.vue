<template>
  <div class="p-2 h-fit bg-background rounded-xl overflow-hidden flex flex-col border border-accent/20">
    <div class="flex items-top md:items-center p-2 gap-1 font-bold text-sm">
      <Icon v-if="tasting.pinned" icon="mynaui:pin" height="14" />
      <div class="hidden md:block">
        <span>{{ formatDate(tasting.tasted_at) }}</span>
        <span class="select-none">|</span>
        <span class="opacity-70">{{ equipment?.name }}</span>
      </div>
      <div class="md:hidden">
        <p>{{ formatDate(tasting.tasted_at, true) }}</p>
        <p class="opacity-70">{{ equipment?.name }}</p>
      </div>

      <span v-if="tasting.rating != null" class="ml-auto"> {{ tasting.rating }}/10 </span>
    </div>

    <div class="flex items-center">
      <div class="hidden md:block w-full h-32 overflow-hidden md:flex-shrink-0">
        <img
          :src="equipmentImage"
          :alt="equipment?.type || $t('coffee.tastingCard.equipment')"
          class="w-full h-full object-contain grayscale brightness-400 dotify"
        />
      </div>

      <div class="md:hidden w-2/2 my-0 w-1/2 px-2">
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

    <div class="flex flex-1 min-h-0">
      <div class="flex flex-col gap-1 p-2 w-full md:w-1/2 overflow-hidden text-sm md:text-xs">
        <div class="font-semibold truncate">{{ tasting.bean?.name || $t('coffee.tastingCard.idk') }}</div>
        <div>
          {{ tasting.bean?.origin || $t('coffee.tastingCard.idk') }}
          <span v-if="tasting.bean?.roaster"> · {{ tasting.bean.roaster }}</span>
        </div>
        <div class="opacity-70">
          <span :class="{ 'opacity-50': tasting.bean.process == null }">{{
            tasting.bean.process || $t('coffee.tastingCard.idk')
          }}</span>
        </div>

        <div
          class="grid grid-cols-2 gap-x-2 place-content-between my-1 text-xs md:block md:space-y-0.5"
        >
          <div>
            {{ $t('coffee.tastingCard.variety') }}
            <span :class="{ 'opacity-50': tasting.bean.variety == null }">{{
              tasting.bean.variety || $t('coffee.tastingCard.idk')
            }}</span>
          </div>
          <div>
            {{ $t('coffee.tastingCard.altitude') }}
            <span :class="{ 'opacity-50': tasting.bean.altitude == null }">{{
              tasting.bean.altitude || $t('coffee.tastingCard.idk')
            }}</span>
          </div>
          <div>
            {{ $t('coffee.tastingCard.grinder') }}
            <span :class="{ 'opacity-50': tasting.grinder?.name == null }">{{
              tasting.grinder?.name || $t('coffee.tastingCard.idk')
            }}</span>
          </div>
          <div>
            {{ $t('coffee.tastingCard.grindSize') }}
            <span :class="{ 'opacity-50': tasting.grind_size == null }">{{
              grindSizeLabel(tasting.grind_size)
            }}</span>
          </div>
          <div>
            {{ $t('coffee.tastingCard.grindSetting') }}
            <span :class="{ 'opacity-50': tasting.grind_setting == null }">{{
              tasting.grind_setting ?? $t('coffee.tastingCard.idk')
            }}</span>
          </div>
          <div>
            {{ $t('coffee.tastingCard.dose') }}
            <span :class="{ 'opacity-50': tasting.coffee_dose == null }">{{
              tasting.coffee_dose != null ? tasting.coffee_dose + ' g' : $t('coffee.tastingCard.idk')
            }}</span>
          </div>
          <div>
            {{ $t('coffee.tastingCard.water') }}
            <span :class="{ 'opacity-50': tasting.water_in == null }">{{
              tasting.water_in != null ? tasting.water_in + ' ml' : $t('coffee.tastingCard.idk')
            }}</span>
          </div>
          <div>
            {{ $t('coffee.tastingCard.yield') }}
            <span :class="{ 'opacity-50': tasting.coffee_out == null }">{{
              tasting.coffee_out != null ? tasting.coffee_out + ' ml' : $t('coffee.tastingCard.idk')
            }}</span>
          </div>
          <div>
            {{ $t('coffee.tastingCard.ratio') }}
            <span :class="{ 'opacity-50': tasting.ratio == null }">{{
              tasting.ratio != null ? '1:' + tasting.ratio.toFixed(1) : $t('coffee.tastingCard.idk')
            }}</span>
          </div>
          <div>
            {{ $t('coffee.tastingCard.brewTime') }}
            <span :class="{ 'opacity-50': tasting.brew_time == null }">{{
              formatBrewTime(tasting.brew_time)
            }}</span>
          </div>
          <div>
            {{ $t('coffee.tastingCard.waterTemp') }}
            <span :class="{ 'opacity-50': tasting.water_temperature == null }">{{
              tasting.water_temperature != null
                ? tasting.water_temperature + '°C'
                : $t('coffee.tastingCard.idk')
            }}</span>
          </div>
        </div>
      </div>

      <div class="hidden md:block my-0 w-1/2 px-2">
        <TasteRadarChart
          :taste-fruity="tasting.taste_fruity"
          :taste-sour="tasting.taste_sour"
          :taste-fermented="tasting.taste_fermented"
          :taste-sweetness="tasting.taste_sweetness"
          :taste-nutty="tasting.taste_nutty"
          :taste-spice="tasting.taste_spice"
          :taste-floral="tasting.taste_floral"
          :taste-green="tasting.taste_green"
          :taste-tobacco="tasting.taste_tobacco"
          :taste-bitter="tasting.taste_bitter"
        />
      </div>
    </div>

    <div
      v-if="tasting.overall_notes"
      class="px-2 pb-2 text-md md:text-xs opacity-80 overflow-hidden max-h-16 flex-shrink-0"
    >
      {{ tasting.overall_notes }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import dayjs from 'dayjs'
import type { TastingNote } from '@/interfaces/Coffee'
import { EQUIPMENT_IMAGE_MAP, DEFAULT_EQUIPMENT_IMAGE, GRIND_SIZES } from '@/helpers/coffee'
import TasteRadarChart from './TasteRadarChart.vue'

const { t, locale } = useI18n()

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

function grindSizeLabel(val: string | null | undefined): string {
  if (val == null) return t('coffee.tastingCard.idk')
  const found = GRIND_SIZES.find((g) => g.value === val)
  return found ? t(found.label) : val
}

function formatDate(dateStr: string, mini: boolean = false): string {
  if (!dateStr) return t('coffee.tastingCard.idk')
  const d = dayjs(dateStr)

  if (mini) {
    return d.locale(locale.value).format('D/M/YYYY')
  }

  return d.locale(locale.value).format('MMM D, YYYY')
}

function formatBrewTime(seconds: number | null | undefined): string {
  if (seconds == null) return t('coffee.tastingCard.idk')
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins}:${String(secs).padStart(2, '0')}`
}
</script>

<style lang="postcss" scoped>
.dotify {
  -webkit-mask-image: url('/assets/dot-matrix.svg');
  mask-image: url('/assets/dot-matrix.svg');
  mask-size: 3rem;
}
</style>
