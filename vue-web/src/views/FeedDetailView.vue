<template>
  <div class="flex flex-col gap-6">
    <template v-if="loading">
      <BlogCard>
        <div class="py-4 px-4 md:px-16">
          <Skeleton h="md" w="2/3" :loading="true"></Skeleton>
          <Skeleton h="sm" w="1/3" :loading="true"></Skeleton>
          <Skeleton h="lg" w="full" :loading="true"></Skeleton>
        </div>
      </BlogCard>
    </template>

    <template v-else-if="!entry">
      <NotFound />
    </template>

    <template v-else>
      <BlogCard>
        <div class="py-4 px-4 md:px-16">
          <h1 class="font-black text-2xl">{{ entry.title }}</h1>
          <h2 v-if="entry.subtitle" class="py-2 text-secondary">{{ entry.subtitle }}</h2>

          <div class="flex gap-4">
            <div class="h-24 w-24 md:h-16 md:w-16">
              <ProfileIcon :email="entry.author?.email ?? ''"></ProfileIcon>
            </div>
            <div class="my-auto text-sm md:text-md">
              <p class="text-secondary text-lg">{{ entry.author?.display_name }}</p>
              <div class="md:flex gap-2 text-secondary italic">
                <p>{{ $t('blog.postedAt') }} {{ dayjs(entry.created_at).calendar() }}</p>
                <p class="hidden md:block px-px">|</p>
                <p>{{ $t('blog.lastEdited') }} {{ dayjs(entry.updated_at).calendar() }}</p>
              </div>
            </div>
          </div>

          <hr class="border-accent my-2" />

          <div v-for="block in entry.blocks" :key="block.id">
            <BlogBlockFull v-if="block.type === 'blog'" :block="block" />
            <MusicBlockFull v-else-if="block.type === 'music'" :block="block" />
          </div>
        </div>
      </BlogCard>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import dayjs from 'dayjs'
import calendar from 'dayjs/plugin/calendar'

dayjs.extend(calendar)

import { getFeedBySlugRaw } from '@/services/feed'
import type { FeedEntry } from '@/interfaces/Feed'

const route = useRoute()
const slug = (typeof route.params.slug === 'string' ? route.params.slug : route.params.slug?.[0]) ?? ''

const entry = ref<FeedEntry | null>(null)
const loading = ref(true)

onMounted(async () => {
  if (!slug) {
    loading.value = false
    return
  }
  try {
    const data = (await getFeedBySlugRaw(slug).send()) as unknown as FeedEntry
    entry.value = data
  } catch {
    entry.value = null
  } finally {
    loading.value = false
  }
})
</script>
