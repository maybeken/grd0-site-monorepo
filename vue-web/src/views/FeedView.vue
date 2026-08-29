<template>
  <div class="flex flex-col gap-6">
    <div>
      <p class="text-2xl font-bold">{{ $t('feed.title') }}<CursorBlink /></p>
    </div>

    <template v-if="initialLoading">
      <BlogCard v-for="i in 3" :key="i">
        <div class="py-4 px-4 md:px-16">
          <Skeleton h="md" w="2/3" :loading="true"></Skeleton>
          <Skeleton h="sm" w="1/3" :loading="true"></Skeleton>
          <Skeleton h="lg" w="full" :loading="true"></Skeleton>
        </div>
      </BlogCard>
    </template>

    <template v-else-if="error">
      <div class="flex flex-col items-center gap-4 py-12">
        <p class="text-secondary">{{ $t('feed.loadError') }}</p>
        <button
          class="px-6 py-2 border border-dotted rounded-full hover:bg-accent/20"
          @click="retry"
        >
          {{ $t('feed.retry') }}
        </button>
      </div>
    </template>

    <template v-else-if="entries.length === 0">
      <div class="flex flex-col items-center py-12">
        <p class="text-secondary">{{ $t('feed.noEntries') }}</p>
      </div>
    </template>

    <template v-else>
      <BlogCard v-for="entry in entries" :key="entry.id">
        <div class="py-4 px-4 md:px-16">
          <div class="flex gap-6">
            <div class="flex flex-col">
              <div class="h-24 w-24">
                <ProfileIcon :email="entry.author?.email ?? ''"></ProfileIcon>
              </div>
              <p class="text-secondary text-center">{{ entry.author?.display_name }}</p>
            </div>

            <div class="flex flex-col gap-2 my-auto w-full">
              <h1 class="font-black text-xl">{{ entry.title }}</h1>
              <h2 v-if="entry.subtitle" class="text-secondary">{{ entry.subtitle }}</h2>
              <p class="text-accent italic text-sm">
                {{ $t('blog.postedAt') }} {{ dayjs(entry.created_at).calendar() }} | {{ $t('blog.lastEdited') }}
                {{ dayjs(entry.updated_at).calendar() }}
              </p>
            </div>
          </div>

          <hr class="border-accent my-2" />

          <template v-if="expandedSlugs.has(entry.slug)">
            <div 
              class="motion-safe:motion-translate-y-in-25 motion-safe:motion-opacity-in-0 motion-safe:motion-duration-300"
            >
              <div v-for="block in expandedData[entry.slug]?.blocks ?? entry.blocks" :key="block.id">
                <BlogBlockFull v-if="block.type === 'blog'" :block="block" />
                <MusicBlockFull v-else-if="block.type === 'music'" :block="block" />
              </div>
            </div>
            <div class="flex justify-center pt-2">
              <button
                class="px-4 py-1 text-sm border border-dotted rounded-full hover:bg-accent/20"
                @click="collapseEntry(entry.slug)"
              >
                {{ $t('feed.collapse') }}
              </button>
            </div>
          </template>

          <template v-else>
            <div v-if="entry.cover_image" class="mb-4 rounded-xl overflow-hidden max-h-64">
              <CDNImage :uri="entry.cover_image" :resolution="768" :quality="75" />
            </div>
            <div class="content px-2 py-2 text-justify max-h-[66vh] overflow-hidden">
              <MarkdownDisplay :md="entry.summary" :lazy_loading="false" />
            </div>
            <div class="flex justify-center pt-2">
              <button
                class="px-4 py-1 text-sm border border-dotted rounded-full hover:bg-accent/20"
                @click="expandEntry(entry)"
              >
                {{ $t('blog.continueReading') }}
              </button>
            </div>
          </template>
        </div>
      </BlogCard>

      <div v-if="loadingMore" class="flex justify-center py-4">
        <Skeleton h="sm" w="lg" :loading="true"></Skeleton>
      </div>

      <div v-if="noMore && entries.length > 0" class="flex justify-center py-4">
        <p class="text-secondary text-sm italic">{{ $t('feed.noMore') }}</p>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import dayjs from 'dayjs'
import calendar from 'dayjs/plugin/calendar'

dayjs.extend(calendar)

import { listFeedRaw, getFeedBySlugRaw } from '@/services/feed'
import type { FeedEntry } from '@/interfaces/Feed'

const PAGE_SIZE = 20

const entries = ref<FeedEntry[]>([])
const initialLoading = ref(true)
const loadingMore = ref(false)
const error = ref(false)
const noMore = ref(false)
const expandedSlugs = reactive(new Set<string>())
const expandedData = reactive<Record<string, FeedEntry>>({})

async function fetchInitial() {
  initialLoading.value = true
  error.value = false
  try {
    const data = (await listFeedRaw(PAGE_SIZE).send()) as unknown as FeedEntry[]
    entries.value = data
    noMore.value = data.length < PAGE_SIZE
  } catch {
    error.value = true
  } finally {
    initialLoading.value = false
  }
}

async function loadMore() {
  if (loadingMore.value || noMore.value || entries.value.length === 0) return

  loadingMore.value = true
  const lastEntry = entries.value[entries.value.length - 1]
  if (!lastEntry) {
    loadingMore.value = false
    return
  }
  const before = lastEntry.created_at

  try {
    const data = (await listFeedRaw(PAGE_SIZE, before).send()) as unknown as FeedEntry[]
    entries.value.push(...data)
    if (data.length < PAGE_SIZE) {
      noMore.value = true
    }
  } catch {
    // silently fail pagination
  } finally {
    loadingMore.value = false
  }
}

async function expandEntry(entry: FeedEntry) {
  expandedSlugs.add(entry.slug)
  try {
    const full = (await getFeedBySlugRaw(entry.slug).send()) as unknown as FeedEntry
    expandedData[entry.slug] = full
  } catch {
    // fall back to preview data
  }
}

function collapseEntry(slug: string) {
  expandedSlugs.delete(slug)
  delete expandedData[slug]
}

function onScroll() {
  const scrollTop = document.documentElement.scrollTop || document.body.scrollTop
  const scrollHeight = document.documentElement.scrollHeight
  const clientHeight = document.documentElement.clientHeight
  if (scrollHeight - scrollTop - clientHeight < 200) {
    loadMore()
  }
}

function retry() {
  fetchInitial()
}

onMounted(() => {
  fetchInitial()
  window.addEventListener('scroll', onScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', onScroll)
})
</script>
