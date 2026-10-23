<template>
  <Map.OlMap class="relative w-full h-[100vw] md:h-[50rem] max-h-lvh rounded-2xl overflow-hidden">
    <div
      v-if="loading"
      class="absolute top-0 left-0 bg-background opacity-50 w-full h-full z-10 motion-preset-fade motion-duration-1000"
    >
      <div class="flex w-full h-full">
        <div class="mx-auto my-auto">
          <Icon
            class="w-64 h-64 motion-preset-spin motion-duration-[5s]"
            icon="mynaui:spinner"
            dotify="false"
          />
        </div>
      </div>
    </div>

    <Map.OlView
      ref="view"
      :center="epsg4326toEpsg3857($prop.map_center)"
      :zoom="$prop.zoom"
      projection="EPSG:3857"
      @change:resolution="onResolutionChange"
    />
    <Layers.OlTileLayer>
      <Sources.OlSourceOSM />
    </Layers.OlTileLayer>

    <MapControls.OlContextMenuControl :items="contextMenuItems" />

    <Map.OlOverlay
      v-if="!loading"
      :position="item.pos"
      v-for="(item, idx) in locations.map((val) => {
        return { ...val, pos: epsg4326toEpsg3857(val.pos) }
      })"
      :key="idx"
      :autoPan="false"
    >
      <div
        v-show="isDisplayOverlay(resolution, item.display_at, item.hide_at)"
        class="px-4 py-2 rounded-r-2xl rounded-b-2xl motion-preset-fade motion-duration-1000 max-w-72 select-none"
        :class="[parseTailwindColor('bg', item.color), parseTailwindColor('text', item.text_color)]"
        @wheel.prevent=""
      >
        <div class="flex gap-1 pb-1">
          <Icon
            v-if="item.icon"
            :class="item.title ? [`mr-1`] : []"
            :icon="item.icon"
            height="auto"
          />
          <span class="text-base">{{ item.title || '' }}</span>
        </div>
        <p v-if="item.subtitle" class="text-xs italic font-light text-justify">
          {{ item.subtitle }}
        </p>
      </div>
    </Map.OlOverlay>
  </Map.OlMap>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { getMapLocation } from '@/services/travelersMap'
import { useToast } from '@/composables/useToast'
import { epsg3857toEpsg4326, epsg4326toEpsg3857 } from '@/helpers/map'
import { parseTailwindColor } from '@/helpers/mapPalette'

import { Map, Layers, Sources, MapControls } from 'vue3-openlayers'

import type { Item } from 'ol-contextmenu'

const { t } = useI18n()
const toast = useToast()

interface Props {
  map_center: number[]
  zoom: number
}

const $prop = defineProps<Props>()

const resolution = ref(-1)
const contextMenuItems = ref<Item[]>([
  {
    text: `Coordinate`,
    callback: (val) => {
      const coordinates = epsg3857toEpsg4326(val.coordinate)
      const title = prompt(t('map.enterTitle'))

      navigator.clipboard.writeText(
        JSON.stringify({
          title,
          pos: coordinates
        })
      )

      toast.success(t('common.savedToClipboard'))
    }
  }
])

const response = getMapLocation()
const locations = response?.data
const loading = response?.loading

function onResolutionChange(event: {
  target: { getResolution: () => number | undefined } | null
}): void {
  const value = event.target?.getResolution()

  resolution.value = typeof value === 'number' ? value : -1
}

function isDisplayOverlay(resolution: number, display_at?: number, hide_at?: number): boolean {
  const hide_flag = hide_at ? (!resolution ? true : resolution > hide_at) : true
  const display_flag = display_at ? (!resolution ? false : display_at >= resolution) : true

  return hide_flag && display_flag
}
</script>
