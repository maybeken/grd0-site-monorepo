<template>
  <Iconify class="dotify" v-bind="$props"></Iconify>
</template>

<script setup lang="ts">
import { Icon as Iconify } from "@iconify/vue";

import type { IconProps } from "@iconify/vue";

const $props = defineProps<IconProps>();
</script>

<style type="postcss" scoped>
.dotify {
  -webkit-mask-image: url("/dot-matrix.svg");
  mask-image: url("/dot-matrix.svg");
}
</style>