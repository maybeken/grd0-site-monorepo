<template>
  <div class="p-4 lg:max-w-[66%] mx-auto">
    <div class="py-4">
      <h1 class="font-black">
        <Skeleton :h="12" :w="'1/3'" :loading="loading">{{ $props?.content?.title }}</Skeleton>
      </h1>
      <h2 class="py-2 text-secondary">
        <Skeleton :h="8" :w="'2/3'" :loading="loading">{{ $props?.content?.subtitle }}</Skeleton>
      </h2>
    </div>

    <div class="flex gap-4">
      <ProfileIcon class="w-24 h-24 md:h-16 md:w-16"
        :src="$props?.content?.author?.email ? `${ASSET_URL}/profile/${$props?.content?.author?.email}.jpg` : ''">
      </ProfileIcon>
      <div class="my-auto text-sm md:text-md">
        <Skeleton :w="32" :h="8" :loading="loading">
          <p class="text-secondary">{{ $props?.content?.author?.display_name }}</p>
        </Skeleton>
        <Skeleton :w="96" :h="8" :loading="loading">
          <p class="text-secondary italic">
            Posted At: {{ dayjs($props?.content?.created_at).calendar() }} |
            Last Edited: {{
              dayjs($props?.content?.updated_at).calendar() }}
          </p>
        </Skeleton>
      </div>
    </div>

    <hr class="border-accent my-2" />

    <div class="content px-2 py-2 text-justify">
      <Skeleton class="h-10 w-full" :loading="loading">
        <MarkdownDisplay :md="$props?.content?.content"></MarkdownDisplay>
      </Skeleton>
    </div>
  </div>
</template>

<script setup lang="ts">
import dayjs from 'dayjs';
import calendar from 'dayjs/plugin/calendar';

dayjs.extend(calendar);

const ASSET_URL = import.meta.env.VITE_ASSETS_URL;

import type { BlogPost } from '@/interfaces/Blog';

const $props = defineProps<{
  content?: BlogPost,
  loading: boolean,
}>();
</script>