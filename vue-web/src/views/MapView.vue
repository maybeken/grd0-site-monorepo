<template>
  <p class="text-xl py-2">Map of Places I've Been To</p>

  <ol-map v-if="locations" class="w-full h-[50rem] max-h-screen rounded-2xl overflow-hidden">
    <ol-view
      ref="view"
      :center="epsg4326toEpsg3857(map_center)"
      :zoom="zoom"
      projection="EPSG:3857"
      @change:resolution="resolution = $event.oldValue"
    />
    <ol-tile-layer>
      <ol-source-osm />
    </ol-tile-layer>

    <ol-context-menu-control
      :items="contextMenuItems"
    />

    <ol-overlay
      :position="item.pos"
      v-for="(item, idx) in locations.map((val) => { return {...val, pos: epsg4326toEpsg3857(val.pos)}})"
      :key="idx"
      :autoPan="true"
    >
      <div
        v-show="(!item.displayAt || item.displayAt >= resolution) && (!item.hideAt || item.hideAt < resolution)"
        class="px-4 py-2 rounded-r-2xl rounded-b-2xl motion-preset-fade"
        :class="parseTailwindColor(item.color)"
      >
        <v-icon v-if="item.icon" class="w-base" :class="item.title ? [`mr-1`] : []" :name="item.icon"></v-icon>
        <span class="text-base">{{ item.title || "" }}</span>
        <p v-if="item.subtitle" class="text-xs text-center italic font-thin">{{ item.subtitle }}</p>
      </div>
    </ol-overlay>
  </ol-map>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { getMapLocation } from "@/services/travelersMap";

import type { Item } from "ol-contextmenu";

const map_center = ref([20, 70]);
const zoom = ref(2);
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

const locations = getMapLocation();

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

    default : 'bg-background',
    red : 'bg-red-700',
    green: 'bg-emerald-600',
  };

  return styles[color] ?? styles['default'];
}
</script>