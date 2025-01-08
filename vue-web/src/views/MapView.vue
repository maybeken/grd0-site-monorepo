<template>
  <ol-map class="w-full h-[50rem] max-h-screen">
    <ol-view ref="view" :center="map_center" :zoom="zoom" projection="EPSG:4326" />
    <ol-tile-layer>
      <ol-source-osm />
    </ol-tile-layer>

    <ol-overlay
      :position="item.pos"
      v-for="(item, idx) in locations"
      :key="idx"
      :autoPan="true"
    >
      <div class="px-4 py-2 rounded-r-2xl rounded-b-2xl" :class="parseTailwindColor(item.color)">
        <v-icon v-if="item.icon" class="w-base" :class="item.title ? [`mr-1`] : []" :name="item.icon"></v-icon>
        <span class="text-base">{{ item.title || "" }}</span>
        <p v-if="item.subtitle" class="text-xs text-center italic font-thin">{{ item.subtitle }}</p>
      </div>
    </ol-overlay>
  </ol-map>
</template>

<script setup lang="ts">
import { ref } from "vue";

const map_center = ref([100, 0]);
const zoom = ref(3.3);

const locations = ref([
  {
    title: "I am Here",
    subtitle: "New Zealand",
    icon: "md-locationon",
    color: "green",
    pos: [169.19502661886494, -45.04620545332485],
  },
  {
    title: "Home",
    subtitle: "Hong Kong",
    icon: "md-home-round",
    color: "red",
    pos: [114.16891354088074, 22.31008594080594],
  },
]);

function parseTailwindColor(color: string = "default"): string {
  const styles: { [key: string]: string } = {
    default : 'bg-background',
    red : 'bg-red-700',
    green: 'bg-emerald-600',
  };

  return styles[color] ?? styles['default'];
}
</script>