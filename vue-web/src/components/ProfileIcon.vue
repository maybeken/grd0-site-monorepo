<template>
  <div
    class="border-solid border-accent border-4 rounded-full overflow-clip h-full w-full bg-secondary"
  >
    <img
      v-if="email"
      class="motion-preset-fade"
      :src="`${ASSET_URL}/cdn-cgi/image/width=${cdn_config.resolution},quality=${cdn_config.quality}/profile/${$props?.email}.jpg`"
    />
  </div>
</template>

<script setup lang="ts">
const ASSET_URL = import.meta.env.VITE_ASSETS_URL;

const cdn_config = {
  resolution: 256,
  quality: 75,
};

interface Props {
  email: string;
};

const {
  email,
} = defineProps<Props>();
</script>