<template>
  <div class="flex flex-col gap-3 rounded-xl border border-dotted p-3">
    <div class="flex items-center gap-2">
      <span class="text-lg font-bold">{{ $t('article.editor.musicTitle') }}</span>
      <EditorButton variant="secondary" size="sm" class="ml-auto" @click="$emit('close')">
        {{ $t('common.cancel') }}
      </EditorButton>
    </div>

    <form class="flex flex-col gap-2" @submit.prevent="useLink">
      <div class="flex flex-col gap-2">
        <EditorInput
          v-model="link"
          inputmode="url"
          :aria-label="$t('article.editor.youtubeLinkPlaceholder')"
          :placeholder="$t('article.editor.youtubeLinkPlaceholder')"
        ></EditorInput>
        <div class="flex justify-end">
          <EditorButton
            type="submit"
            variant="primary"
            size="sm"
            :disabled="resolving || !link.trim()"
            :loading="resolving"
          >
            {{ resolving ? $t('article.editor.resolvingLink') : $t('article.editor.useLink') }}
          </EditorButton>
        </div>
      </div>
      <p v-if="error" class="text-sm font-medium text-red-400">{{ error }}</p>
    </form>

    <div v-if="preview" class="flex flex-col gap-2 rounded-xl bg-shade p-3">
      <MusicPlayer
        variant="compact"
        :v="preview.v"
        :title="title_override || preview.title"
        :artist="artist_override || preview.artist"
      ></MusicPlayer>

      <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
        <EditorField :label="$t('music.editor.titleOverride')" html-for="music-title-override">
          <EditorInput
            id="music-title-override"
            v-model="title_override"
            :placeholder="preview.title"
          ></EditorInput>
        </EditorField>
        <EditorField :label="$t('music.editor.artistOverride')" html-for="music-artist-override">
          <EditorInput
            id="music-artist-override"
            v-model="artist_override"
            :placeholder="preview.artist"
          ></EditorInput>
        </EditorField>
      </div>

      <div class="flex justify-end gap-2">
        <EditorButton variant="secondary" size="sm" @click="clearPreview">
          {{ $t('common.cancel') }}
        </EditorButton>
        <EditorButton variant="primary" size="sm" @click="confirmPreview">
          {{ $t('article.editor.useThisSong') }}
        </EditorButton>
      </div>
    </div>

    <p
      v-if="!$props.link_only && !loading && (!songs || songs.length === 0)"
      class="text-lg font-medium text-foreground/60"
    >
      {{ $t('article.editor.noSongs') }}
    </p>

    <div v-if="!$props.link_only" class="flex max-h-[50vh] flex-col overflow-y-auto">
      <button
        v-for="song in songs"
        :key="song.v"
        class="flex items-center gap-3 border-b py-2 text-left hover:bg-shade focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-foreground/60"
        @click="select(song)"
      >
        <img
          class="h-12 w-12 rounded-lg object-cover"
          :src="`https://i.ytimg.com/vi_webp/${song.v}/default.webp`"
          loading="lazy"
        />
        <span class="flex flex-col">
          <span class="text-lg font-medium">{{ song.title }}</span>
          <span class="text-sm font-medium text-foreground/60">
            {{ song.artist || $t('music.unknown') }}
          </span>
        </span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import EditorButton from '@/components/editor/ui/EditorButton.vue'
import EditorField from '@/components/editor/ui/EditorField.vue'
import EditorInput from '@/components/editor/ui/EditorInput.vue'
import MusicPlayer from '@/components/music/MusicPlayer.vue'
import { getMusic } from '@/services/music'

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

const { t } = useI18n()

const response = getMusic()
const songs = response.data
const loading = response.loading

const link = ref('')
const resolving = ref(false)
const error = ref('')
const preview = ref<{ v: string; title: string; artist?: string } | null>(null)
const title_override = ref('')
const artist_override = ref('')

function clearPreview(): void {
  preview.value = null
  title_override.value = ''
  artist_override.value = ''
}

watch(link, clearPreview)

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
    title_override.value = song.title
    artist_override.value = song.artist ?? ''
  } catch {
    error.value = t('article.editor.youtubeLinkFailed')
  } finally {
    resolving.value = false
  }
}

function confirmPreview(): void {
  if (!preview.value) return

  const title = title_override.value.trim() || preview.value.title
  const artist = artist_override.value.trim() || preview.value.artist

  $emit('select', { v: preview.value.v, title, artist })
  clearPreview()
}

function select(song: Music): void {
  clearPreview()
  $emit('select', { v: song.v, title: song.title, artist: song.artist })
}
</script>
