<template>
  <div class="silhouette" :style="style" role="img" :aria-label="alt || undefined"></div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    src: string
    alt?: string
  }>(),
  {
    alt: ''
  }
)

const style = computed<Record<string, string>>(() => {
  const url = `url("${props.src}")`

  return {
    'mask-image': url,
    '-webkit-mask-image': url
  }
})
</script>

<style scoped>
.silhouette {
  width: 100%;
  height: 100%;
  background-color: currentColor;
  mask-repeat: no-repeat;
  -webkit-mask-repeat: no-repeat;
  mask-position: center;
  -webkit-mask-position: center;
  mask-size: contain;
  -webkit-mask-size: contain;
}
</style>
