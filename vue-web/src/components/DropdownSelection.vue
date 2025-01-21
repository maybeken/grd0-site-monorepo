<template>
  <div id="category-btn" class="relative h-12">
    <div class="absolute top-0 h-12 z-10 w-full md:w-96 max-w-full">
      <div
        class="flex px-4 py-2 w-full rounded-lg font-black bg-background border-[1px] border-solid border-foreground cursor-default z-20"
        @click="expanded = !expanded;"
      >
        <span class="grow">{{ stylize ? stylize(selected ?? "") : selected }}</span>
        <Icon icon="mynaui:chevron-down" height="auto" />
      </div>
      <div
        class="px-4 py-2 w-full bg-background bg-opacity-50 backdrop-blur-sm border-[1px] border-solid border-foreground cursor-pointer motion-preset-slide-down motion-duration-500"
        v-for="(item, key, idx) of { all: {}, ...$props.options }" :key="key" :value="item"
        :class="displayButtonStylize(`${key}`, idx, ($props.options ? Object.keys($props.options).length : 0))"
        @click="clickButton(`${key}`)"
      >
        {{ stylize ? stylize(`${key}`) : key }}
      </div>
    </div>
  </div>

</template>

<script setup lang="ts">
import { ref, defineEmits } from 'vue';
import { useRouter } from 'vue-router';

import type { GalleryCategory } from '@/interfaces/Gallery';

interface Props {
  options?: GalleryCategory;
  selected?: string;
  stylize?: (text: string) => string;
}

const $props = defineProps<Props>();
const $emit = defineEmits<{
  select: [value?: string]
}>()

const $router = useRouter();
const selected = ref($props.selected);
const expanded = ref(false);

function displayButtonStylize(val: string, idx: number, max: number = -1) {
  let style: string[] = [];

  if (expanded.value) {
    if (idx === 0) {
      style = [...style, 'rounded-t-lg'];
    } else if (idx === max) {
      style = [...style, 'rounded-b-lg'];
    }

    if (val === selected.value) {
      style = [...style, 'font-black'];
    }

    return style;
  }

  return ['hidden'];
}

function clickButton(val: string) {
  $emit('select', val === 'all' ? val : `/gallery/${val}`);
  selected.value = val;
  expanded.value = !expanded.value;
}
</script>