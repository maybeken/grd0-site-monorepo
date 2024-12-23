<template>
  <select
    class="px-4 py-2 rounded bg-background border-[1px] border-solid border-foreground"
    v-model="selected"
    @change="$emit('select', selected)"
  >
    <option v-for="item of $props.options" :key="item" :value="item">
      {{ stylize ? stylize(item) : item }}
    </option>
  </select>
</template>

<script setup lang="ts">
import { ref, watchEffect } from 'vue';

interface Props {
  options?: string[];
  selected?: string;
  stylize?: (text: string) => string;
}

const $props = defineProps<Props>();
const selected = ref($props.selected);

watchEffect(() => {
  if (!selected.value && $props.selected) selected.value = $props.selected;
})
</script>