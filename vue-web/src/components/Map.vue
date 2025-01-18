<template>
  <ol-map class="relative w-full h-[100vw] md:h-[50rem] max-h-screen rounded-2xl overflow-hidden">
    <div v-if="loading" class="absolute top-0 left-0 bg-background opacity-50 w-full h-full z-10 motion-preset-fade motion-duration-1000">
      <div class="flex w-full h-full">
        <div class="mx-auto my-auto">
          <v-icon scale="10" animation="spin-pulse" speed="slow" name="md-hourglassbottom"></v-icon>
        </div>
      </div>
    </div>

    <ol-view ref="view" :center="epsg4326toEpsg3857($prop.map_center)" :zoom="$prop.zoom" projection="EPSG:3857"
      @change:resolution="resolution = $event.oldValue" />
    <ol-tile-layer>
      <ol-source-osm />
    </ol-tile-layer>

    <ol-context-menu-control :items="contextMenuItems" />

    <ol-overlay v-if="!loading" :position="item.pos"
      v-for="(item, idx) in locations.map((val) => { return { ...val, pos: epsg4326toEpsg3857(val.pos) } })" :key="idx"
      :autoPan="false">
      <div
        v-show="isDisplayOverlay(resolution, item.display_at, item.hide_at)"
        class="px-4 py-2 rounded-r-2xl rounded-b-2xl motion-preset-fade motion-duration-1000"
        :class="[parseTailwindColor('bg', item.color), parseTailwindColor('text', item.text_color)]">
        <v-icon v-if="item.icon" class="w-base" :class="item.title ? [`mr-1`] : []" :name="item.icon"></v-icon>
        <span class="text-base">{{ item.title || "" }}</span>
        <p v-if="item.subtitle" class="text-xs text-center italic font-light">{{ item.subtitle }}</p>
      </div>
    </ol-overlay>
  </ol-map>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { getMapLocation } from "@/services/travelersMap";

import type { Item } from "ol-contextmenu";

interface Props {
  map_center: [number, number],
  zoom: number,
};

const $prop = defineProps<Props>();

const resolution = ref(-1);
const contextMenuItems = ref<Item[]>([
  {
    text: `Coordinate`,
    callback: (val) => {
      const coordinates = epsg3857toEpsg4326(val.coordinate);
      const title = prompt("Enter the title:");

      navigator.clipboard.writeText(JSON.stringify({
        title,
        pos: coordinates,
      }));

      alert('Saved to clipboard!');
    },
  }
]);

const response = getMapLocation();
const locations = response?.data;
const loading = response?.loading;

function epsg4326toEpsg3857(coordinates: number[]) {
  // Parse coordinates as North-East
  let y = coordinates[0];
  let x = coordinates[1];
  x = (x * 20037508.34) / 180;
  y =
    Math.log(Math.tan(((90 + y) * Math.PI) / 360)) /
    (Math.PI / 180);
  y = (y * 20037508.34) / 180;
  return [x, y];
}

function epsg3857toEpsg4326(coordinates: number[]) {
  let x = coordinates[0];
  let y = coordinates[1];
  x = (x * 180) / 20037508.34;
  y = (y * 180) / 20037508.34;
  y = (Math.atan(Math.pow(Math.E, y * (Math.PI / 180))) * 360) / Math.PI - 90;
  return [y, x];
}

function parseTailwindColor(type: string, color: string = "default"): string {
  const background_styles: { [key: string]: string } = {
    default: 'bg-background',
    red: 'bg-red-700',
    green: 'bg-emerald-600',
    blue: 'bg-sky-600',
    sweden: 'bg-[#005293]',
  };

  const text_styles: { [key: string]: string } = {
    sweden: 'text-[#FFCD00]',
  };

  if (type === 'text') {
    return text_styles[color] ?? '';
  }

  return background_styles[color] ?? background_styles['default'];
}

function isDisplayOverlay(resolution: number, display_at?: number, hide_at?: number): boolean {

  const hide_flag = hide_at ? (!resolution ? true : resolution > hide_at) : true;
  const display_flag = display_at ? (!resolution ? false : display_at >= resolution) : true;

  return hide_flag && display_flag;
}
</script>