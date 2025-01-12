<template>
  <BlogCard>
    <div class="relative p-4">
      <div class="flex gap-6">
        <div class="flex flex-col">
          <div class="w-full">
            <ProfileIcon class="mx-auto w-16" :src="`//${ASSET_URL}/profile/${$props.author.email}.jpg`"></ProfileIcon>
          </div>
          <p class="text-secondary">Author: {{ $props.author.display_name }}</p>
        </div>

        <div class="flex flex-col gap-2 hover:underline my-auto">
          <router-link :to="`/blog/${$props.uri}`">
            <h1 class="font-black">{{ $props.title }}</h1>
            <h2>{{ $props.subtitle }}</h2>
          </router-link>
        </div>
      </div>

      <hr class="border-accent my-2" />
      
      <div class="content px-2 py-2 text-justify max-h-[66vh] overflow-hidden">
        <MarkdownDisplay :md="$props.content"></MarkdownDisplay>
      </div>

      <div class="absolute bottom-0 left-0 h-28 w-full z-10 rounded-xl bg-gradient-to-t from-background from-33%">
        <router-link :to="`/blog/${$props.uri}`">
          <div class="flex flex-col flex-col-reverse h-28 w-full hover:cursor-hand hover:underline">
            <div class="mx-auto my-8">
              <p class="text-2xl">Continue to Read</p>
            </div>
          </div>
        </router-link>
      </div>

      <hr class="border-accent my-2" />

      <div class="flex justify-end">
        <p class="text-accent italic">
          Posted At: {{ dayjs($props.created_at).calendar() }} | Last Edited: {{ dayjs($props.updated_at).calendar() }}</p>
      </div>
    </div>
  </BlogCard>
</template>

<script setup lang="ts">
import dayjs from 'dayjs';
import calendar from 'dayjs/plugin/calendar';

dayjs.extend(calendar);

import type { BlogPost } from '@/interfaces/Blog';

const $props = defineProps<BlogPost>();
const ASSET_URL = import.meta.env.VITE_ASSETS_DOMAIN;
</script>

<style type="poscss" scoped>
.content::first-letter {
  @apply text-2xl font-semibold;
}
</style>