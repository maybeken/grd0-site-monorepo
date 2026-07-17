<template>
  <div class="px-2">
    <div class="flex gap-2 py-2">
      <div class="md:flex items-end mb-auto gap-2 grow font-bold">
        <p class="text-xl">Places I've Been To<CursorBlink /></p>
        <p class="text-sm">(Zoom in to see more)</p>
      </div>
      <button
        class="rounded-full p-2 border-2 border-dotted"
        @click="
          () => {
            map_center = [20, 70]
            zoom = 2
          }
        "
        title="Re-center"
      >
        <div class="flex flex-col gap-2">
          <Icon class="mx-auto" icon="mynaui:redo" height="auto" />
        </div>
      </button>
    </div>

    <TravelersMap :map_center="map_center" :zoom="zoom"> </TravelersMap>
  </div>
</template>

<script setup lang="ts">
import { ref, defineAsyncComponent } from 'vue'

const TravelersMap = defineAsyncComponent(() => import('@/components/map/TravelersMap.vue'))

const map_center = ref([20, 70])
const zoom = ref(2)
</script>
