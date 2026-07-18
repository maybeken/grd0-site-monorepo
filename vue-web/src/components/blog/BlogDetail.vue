<template>
  <div class="py-4 w-full">
    <h1 class="font-black">
      <Skeleton h="md" w="1/3" :loading="loading">{{ $props?.content?.title }}</Skeleton>
    </h1>
    <h2 class="py-2 text-secondary">
      <Skeleton h="sm" w="2/3" :loading="loading">{{ $props?.content?.subtitle }}</Skeleton>
    </h2>
  </div>

  <div class="flex gap-4">
    <div class="h-24 w-24 md:h-16 md:w-16">
      <ProfileIcon :email="$props?.content?.author?.email ?? ''"></ProfileIcon>
    </div>
    <div class="my-auto text-sm md:text-md">
      <Skeleton h="sm" w="lg" :loading="loading">
        <p class="text-secondary text-lg">{{ $props?.content?.author?.display_name }}</p>
      </Skeleton>
      <div class="md:flex gap-2 text-secondary italic">
        <Skeleton h="xs" w="md" :loading="loading">
          <p>{{ $t('blog.postedAt') }} {{ dayjs($props?.content?.created_at).calendar() }}</p>
        </Skeleton>
        <p class="hidden md:block px-px">|</p>
        <Skeleton h="xs" w="md" :loading="loading">
          <p>{{ $t('blog.lastEdited') }} {{ dayjs($props?.content?.updated_at).calendar() }}</p>
        </Skeleton>
      </div>
    </div>
  </div>

  <hr class="border-accent my-2" />

  <div class="content px-2 py-2 text-justify">
    <Skeleton h="lg" w="full" :loading="loading">
      <MarkdownDisplay :md="$props?.content?.content"></MarkdownDisplay>
    </Skeleton>
  </div>
</template>

<script setup lang="ts">
import dayjs from 'dayjs'
import calendar from 'dayjs/plugin/calendar'

dayjs.extend(calendar)

import type { BlogPost } from '@/interfaces/Blog'

const $props = defineProps<{
  content?: BlogPost
  loading: boolean
}>()
</script>
