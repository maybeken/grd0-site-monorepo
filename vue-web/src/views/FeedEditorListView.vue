<template>
  <div class="flex flex-col gap-6 px-4 md:px-16">
    <div class="flex justify-between items-center">
      <p class="text-2xl font-bold">{{ $t('feed.list.title') }}<CursorBlink /></p>
      <router-link
        to="/feed/editor/new"
        class="px-4 py-2 border border-dotted rounded-full bg-accent/20 hover:bg-accent/30"
      >
        {{ $t('feed.list.newEntry') }}
      </router-link>
    </div>

    <div v-if="loading" class="flex flex-col gap-4">
      <div v-for="i in 3" :key="i" class="bg-shade rounded-2xl p-4">
        <Skeleton h="md" w="2/3" :loading="true"></Skeleton>
        <Skeleton h="sm" w="1/3" :loading="true"></Skeleton>
      </div>
    </div>

    <div v-else-if="error" class="text-center py-12">
      <p class="text-red-400">{{ $t('feed.loadError') }}</p>
      <button
        class="mt-4 px-4 py-2 border border-dotted rounded-full hover:bg-accent/20"
        @click="loadEntries"
      >
        {{ $t('feed.retry') }}
      </button>
    </div>

    <div v-else-if="entries.length === 0" class="text-center py-12">
      <p class="text-secondary">{{ $t('feed.noEntries') }}</p>
    </div>

    <div v-else class="flex flex-col gap-4">
      <div
        v-for="entry in entries"
        :key="entry.id"
        class="bg-shade rounded-2xl p-4 flex justify-between items-center"
      >
        <div class="flex flex-col gap-1">
          <div class="flex items-center gap-2">
            <h3 class="font-bold text-lg">{{ entry.title }}</h3>
            <span
              class="text-xs px-2 py-1 rounded-full"
              :class="entry.published_at ? 'bg-emerald-900/30 text-emerald-400' : 'bg-yellow-900/30 text-yellow-400'"
            >
              {{ entry.published_at ? $t('feed.list.published') : $t('feed.list.draft') }}
            </span>
          </div>
          <p v-if="entry.subtitle" class="text-secondary text-sm">{{ entry.subtitle }}</p>
          <p class="text-secondary text-xs">
            {{ dayjs(entry.created_at).calendar() }}
          </p>
        </div>
        <div class="flex gap-2">
          <router-link
            :to="`/feed/editor/${entry.slug}`"
            class="px-3 py-1 border border-dotted rounded-full hover:bg-accent/20 text-sm"
          >
            {{ $t('feed.list.editEntry') }}
          </router-link>
          <button
            class="px-3 py-1 border border-dotted rounded-full hover:bg-red-500/20 text-red-400 text-sm"
            @click="deleteEntry(entry.slug)"
          >
            {{ $t('feed.list.deleteEntry') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import dayjs from 'dayjs'
import calendar from 'dayjs/plugin/calendar'

dayjs.extend(calendar)

import { adminInstance } from '@/services/api'
import type { FeedEntry } from '@/interfaces/Feed'

const { t } = useI18n()

const entries = ref<FeedEntry[]>([])
const loading = ref(true)
const error = ref(false)

async function loadEntries() {
  loading.value = true
  error.value = false

  try {
    const response = await adminInstance.Get<FeedEntry[]>('/feed/all').send()
    entries.value = response as unknown as FeedEntry[]
  } catch {
    error.value = true
  } finally {
    loading.value = false
  }
}

async function deleteEntry(slug: string) {
  if (!confirm(t('feed.list.confirmDelete'))) return

  try {
    await adminInstance.Delete(`/feed/${slug}`).send()
    entries.value = entries.value.filter((e) => e.slug !== slug)
  } catch {
    alert(t('feed.list.deleteFailed'))
  }
}

onMounted(() => {
  loadEntries()
})
</script>
