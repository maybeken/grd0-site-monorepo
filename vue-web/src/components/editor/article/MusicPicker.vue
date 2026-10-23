<template>
  <div class="flex flex-col gap-3 rounded-xl border border-dotted p-3">
    <div class="flex items-center gap-2">
      <span class="font-bold">{{ $t('article.editor.musicTitle') }}</span>
      <button class="ml-auto rounded-full border border-dotted px-3 py-1" @click="$emit('close')">
        {{ $t('common.cancel') }}
      </button>
    </div>

    <form class="flex flex-col gap-2" @submit.prevent="useLink">
      <div class="flex flex-col gap-2 sm:flex-row">
        <input
          v-model="link"
          class="grow rounded-full border border-dotted bg-background px-3 py-1"
          type="text"
          inputmode="url"
          :aria-label="$t('article.editor.youtubeLinkPlaceholder')"
          :placeholder="$t('article.editor.youtubeLinkPlaceholder')"
        />
        <button
          class="rounded-full border border-dotted px-3 py-1 disabled:opacity-50"
          type="submit"
          :disabled="resolving || !link.trim()"
        >
          {{ resolving ? $t('article.editor.resolvingLink') : $t('article.editor.useLink') }}
        </button>
      </div>
      <p v-if="error" class="text-sm text-secondary">{{ error }}</p>
    </form>

    <div v-if="preview" class="flex flex-col gap-2 rounded-xl bg-shade p-3">
      <MusicPlayer
        variant="compact"
        :v="preview.v"
        :title="preview.title"
        :artist="preview.artist"
      ></MusicPlayer>
      <div class="flex justify-end gap-2">
        <button class="rounded-full border border-dotted px-3 py-1" @click="preview = null">
          {{ $t('common.cancel') }}
        </button>
        <button class="rounded-full border border-dotted px-3 py-1" @click="confirmPreview">
          {{ $t('article.editor.useThisSong') }}
        </button>
      </div>
    </div>

    <p
      v-if="!$props.link_only && !loading && (!songs || songs.length === 0)"
      class="text-sm text-secondary"
    >
      {{ $t('article.editor.noSongs') }}
    </p>

    <div v-if="!$props.link_only" class="flex flex-col max-h-[50vh] overflow-y-auto">
      <button
        v-for="song in songs"
        :key="song.v"
        class="flex items-center gap-3 border-b py-2 text-left hover:bg-shade"
        @click="select(song)"
      >
        <img
          class="h-12 w-12 rounded-lg object-cover"
          :src="`https://i.ytimg.com/vi_webp/${song.v}/default.webp`"
          loading="lazy"
        />
        <span class="flex flex-col">
          <span class="font-bold">{{ song.title }}</span>
          <span class="text-sm text-secondary">{{ song.artist || $t('music.unknown') }}</span>
        </span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

import MusicPlayer from '@/components/music/MusicPlayer.vue'
import { getMusic } from '@/services/music'

import { t } from '@/i18n'

import type { Music } from '@/interfaces/Music'

const $emit = defineEmits<{
  select: [song: { v: string; title: string; artist?: string }]
  close: []
}>()

const $props = withDefaults(
  defineProps<{
    link_only?: boolean
  }>(),
  {
    link_only: false
  }
)

const response = getMusic()
const songs = response.data
const loading = response.loading

const link = ref('')
const resolving = ref(false)
const error = ref('')
const preview = ref<{ v: string; title: string; artist?: string } | null>(null)

watch(link, () => {
  preview.value = null
})

const youtube_id = /^[A-Za-z0-9_-]{11}$/

function parseYouTubeId(input: string): string | null {
  const value = input.trim()

  if (youtube_id.test(value)) return value

  let url: URL
  try {
    url = new URL(/^https?:\/\//i.test(value) ? value : `https://${value}`)
  } catch {
    return null
  }

  const host = url.hostname.replace(/^(www|m|music)\./, '')

  if (host === 'youtu.be') {
    const id = url.pathname.split('/').filter(Boolean)[0]
    return id && youtube_id.test(id) ? id : null
  }

  if (host !== 'youtube.com' && host !== 'youtube-nocookie.com') return null

  const query_id = url.searchParams.get('v')
  if (query_id && youtube_id.test(query_id)) return query_id

  const [prefix, id] = url.pathname.split('/').filter(Boolean)
  const supported = ['embed', 'shorts', 'live']
  if (supported.includes(prefix ?? '') && id && youtube_id.test(id)) return id

  return null
}

async function fetchSong(v: string): Promise<{ title: string; artist?: string }> {
  const query = new URLSearchParams({
    url: `https://www.youtube.com/watch?v=${v}`,
    format: 'json'
  })
  const response = await fetch(`https://www.youtube.com/oembed?${query}`)

  if (!response.ok) throw new Error('oembed failed')

  const data = (await response.json()) as { title?: string; author_name?: string }

  return {
    title: data.title?.trim() || v,
    artist: data.author_name?.trim()
  }
}

async function useLink(): Promise<void> {
  if (resolving.value) return

  const v = parseYouTubeId(link.value)
  if (!v) {
    error.value = t('article.editor.invalidYoutubeLink')
    return
  }

  resolving.value = true
  error.value = ''

  try {
    const song = await fetchSong(v)
    preview.value = { v, ...song }
  } catch {
    error.value = t('article.editor.youtubeLinkFailed')
  } finally {
    resolving.value = false
  }
}

function confirmPreview(): void {
  if (!preview.value) return

  $emit('select', preview.value)
  preview.value = null
}

function select(song: Music): void {
  preview.value = null
  $emit('select', { v: song.v, title: song.title, artist: song.artist })
}
</script>
