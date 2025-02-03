<template>
  <Iconify :class="$props.dotify !== 'false' ? 'dotify' : ''" v-bind="$props"></Iconify>
</template>

<script setup lang="ts">
import { Icon as Iconify } from "@iconify/vue";

import type { IconProps } from "@iconify/vue";

interface Props extends IconProps {
  dotify?: string,
}

const $props = defineProps<Props>();
</script>

<style lang="postcss" scoped>
.dotify {
  -webkit-mask-image: url("/assets/dot-matrix.svg");
  mask-image: url("/assets/dot-matrix.svg");
}
</style>