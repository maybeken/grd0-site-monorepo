<template>
  <div class="px-2">
    <div class="flex gap-2 py-2">
      <div class="md:flex items-end gap-2 grow">
        <p class="text-xl">Places I've Been To</p>
        <p class="text-sm">(Zoom in to see more)</p>
      </div>
      <button class="rounded-xl p-2 border-2 border-dotted" @click="map_center = [20, 70]; zoom = 2;">
        <div class="flex gap-1">
          <span class="pl-2">Re-Center</span>
          <Icon icon="mynaui:center-focus" height="auto" />
        </div>
      </button>
    </div>

    <Map
      :map_center="map_center"
      :zoom="zoom"
    >
    </Map>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
() => import("@/components/Map.vue");

const map_center = ref([20, 70]);
const zoom = ref(2);
</script>