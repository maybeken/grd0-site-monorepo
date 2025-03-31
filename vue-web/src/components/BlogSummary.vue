<template>
  <BlogCard>
    <div class="relative py-4 px-4 md:px-16">
      <div class="flex gap-6">
        <div class="flex flex-col">
          <div class="h-24 w-24">
            <ProfileIcon class="mx-auto" :email="$props?.author?.email ?? ''"></ProfileIcon>
          </div>
          <Skeleton h="xs" w="sm" :loading="loading">
            <p class="text-secondary text-center">{{ $props?.author?.display_name }}</p>
          </Skeleton>
        </div>

        <div class="flex flex-col gap-2 hover:underline my-auto w-full">
          <router-link :to="`/blog/${$props?.uri}`">
            <Skeleton h="md" w="1/3" :loading="loading">
              <h1 class="font-black">{{ $props?.title }}</h1>
            </Skeleton>
            <Skeleton h="sm" w="2/3" :loading="loading">
              <h2>{{ $props.subtitle }}</h2>
            </Skeleton>
          </router-link>
        </div>
      </div>

      <hr class="border-accent my-2" />

      <div class="content px-2 py-2 text-justify max-h-[66vh] overflow-hidden">
        <Skeleton h="md" w="full" :loading="loading">
          <MarkdownDisplay :md="$props?.content" :lazy_loading="false"></MarkdownDisplay>
        </Skeleton>
      </div>

      <div v-if="!loading"
        class="absolute bottom-0 left-0 h-28 w-full z-10 rounded-xl bg-linear-to-t from-background from-33%">
        <router-link :to="`/blog/${$props?.uri}`">
          <div class="flex flex-col-reverse h-28 w-full hover:cursor-pointer hover:underline">
            <div class="mx-auto my-8">
              <p class="text-2xl">Continue to Read</p>
            </div>
          </div>
        </router-link>
      </div>

      <hr class="border-accent my-2" />

      <div class="flex justify-end">

        <Skeleton h="sm" w="xl" :loading="loading">
          <p class="text-accent italic">Posted At: {{ dayjs($props?.created_at).calendar() }} | Last Edited: {{
            dayjs($props?.updated_at).calendar() }}</p>
        </Skeleton>
      </div>
    </div>
  </BlogCard>
</template>

<script setup lang="ts">
import dayjs from 'dayjs';
import calendar from 'dayjs/plugin/calendar';

dayjs.extend(calendar);

import type { BlogPost } from '@/interfaces/Blog';

interface Props extends BlogPost {
  loading: boolean,
};

const $props = defineProps<Props>();
</script>

<style lang="postcss" scoped>
.content:first-letter {
  @apply text-2xl font-semibold;
}
</style>