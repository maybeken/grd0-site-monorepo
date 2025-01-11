<template>
  <div class="flex gap-2 py-2">
    <p class="text-xl py-2 grow">Map of Places I've Been To</p>
    <button class="rounded-xl px-4 border-2 border-dotted" @click="map_center = [20, 70]; zoom = 2;">
      Re-Center
      <v-icon name="md-undo"></v-icon>
    </button>
  </div>

  <Map
    :map_center="map_center"
    :zoom="zoom"
  >
  </Map>
</template>

<script setup lang="ts">
import { ref } from "vue";

const map_center = ref([20, 70]);
const zoom = ref(2);
</script>