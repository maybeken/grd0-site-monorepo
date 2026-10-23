<template>
  <Map.OlMap
    class="relative h-96 w-full overflow-hidden rounded-xl"
    @singleclick="onSingleClick"
  >
    <Map.OlView :center="center" :zoom="zoom" projection="EPSG:3857" />

    <Layers.OlTileLayer>
      <Sources.OlSourceOSM />
    </Layers.OlTileLayer>

    <Layers.OlVectorLayer :style="marker_style">
      <Sources.OlSourceVector>
        <Map.OlFeature v-if="marker_position">
          <Geometries.OlGeomPoint :coordinates="marker_position" />
        </Map.OlFeature>
      </Sources.OlSourceVector>
    </Layers.OlVectorLayer>

    <Map.OlOverlay v-if="marker_position" :position="marker_position" :auto-pan="false">
      <div
        class="max-w-72 select-none rounded-r-2xl rounded-b-2xl px-4 py-2"
        :class="[
          parseTailwindColor('bg', color),
          parseTailwindColor('text', text_color || 'default')
        ]"
      >
        <div class="flex gap-1 pb-1">
          <Icon
            v-if="icon"
            :class="title ? ['mr-1'] : []"
            :icon="icon"
            height="auto"
          />
          <span class="text-base">{{ title || '' }}</span>
        </div>
        <p v-if="subtitle" class="text-xs font-light italic text-justify">
          {{ subtitle }}
        </p>
      </div>
    </Map.OlOverlay>
  </Map.OlMap>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Map, Layers, Sources, Geometries } from 'vue3-openlayers'
import Style from 'ol/style/Style'
import CircleStyle from 'ol/style/Circle'
import Fill from 'ol/style/Fill'
import Stroke from 'ol/style/Stroke'

import { epsg3857toEpsg4326, epsg4326toEpsg3857, isValidLatitude, isValidLongitude } from '@/helpers/map'
import { parseTailwindColor } from '@/helpers/mapPalette'

interface SingleClickEvent {
  coordinate: number[]
}

interface Props {
  latitude?: number | null
  longitude?: number | null
  title?: string
  subtitle?: string
  icon?: string
  color?: string
  text_color?: string
}

const props = withDefaults(defineProps<Props>(), {
  latitude: null,
  longitude: null,
  title: '',
  subtitle: '',
  icon: '',
  color: 'default',
  text_color: ''
})

const emit = defineEmits<{
  select: [coordinates: number[]]
}>()

const marker_style = new Style({
  image: new CircleStyle({
    radius: 8,
    fill: new Fill({ color: '#f97316' }),
    stroke: new Stroke({ color: '#ffffff', width: 2 })
  })
})

const marker_position = computed<number[] | null>(() => {
  if (
    props.latitude === null ||
    props.longitude === null ||
    !isValidLatitude(props.latitude) ||
    !isValidLongitude(props.longitude)
  ) {
    return null
  }

  return epsg4326toEpsg3857([props.latitude, props.longitude])
})

const center = marker_position.value ?? [0, 0]
const zoom = marker_position.value ? 4 : 2

function onSingleClick(event: SingleClickEvent): void {
  emit('select', epsg3857toEpsg4326(event.coordinate))
}
</script>
