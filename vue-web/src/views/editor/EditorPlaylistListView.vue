<template>
  <div class="flex flex-col gap-6">
    <header class="flex flex-col gap-2">
      <h1 class="text-2xl font-bold">{{ $t('editor.playlists') }}<CursorBlink /></h1>
      <p class="text-secondary">{{ $t('music.editor.listHint') }}</p>
    </header>

    <div v-if="!show_form" class="flex justify-end">
      <button
        class="rounded bg-accent px-3 py-1 text-sm font-semibold text-background hover:brightness-110"
        @click="show_form = true"
      >
        {{ $t('music.editor.newPlaylist') }}
      </button>
    </div>

    <PlaylistForm
      v-if="show_form"
      :playlist="null"
      @saved="onSaved"
      @cancel="show_form = false"
    ></PlaylistForm>

    <div v-if="!show_form" class="flex flex-col gap-2">
      <EditorEntityRow
        v-for="playlist in playlists"
        :key="playlist.id"
        :title="playlist.title"
        :to="`/editor/music/${playlist.id}`"
        @delete="remove(playlist)"
      >
        <span v-if="playlist.description" class="truncate text-xs text-foreground/60">
          {{ playlist.description }}
        </span>
        <span class="text-xs text-foreground/40">
          {{ $t('music.editor.songCount', { count: playlist.songs?.length ?? 0 }) }}
        </span>
      </EditorEntityRow>

      <p v-if="!playlists.length" class="py-4 text-center text-sm text-foreground/40">
        {{ $t('music.editor.noPlaylists') }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import PlaylistForm from '@/components/music/PlaylistForm.vue'
import EditorEntityRow from '@/components/editor/EditorEntityRow.vue'
import { deletePlaylistRaw, getPlaylistsRaw } from '@/services/music'
import { useToast } from '@/composables/useToast'

import type { Playlist } from '@/interfaces/Music'

const { t } = useI18n()
const toast = useToast()

const playlists = ref<Playlist[]>([])
const show_form = ref(false)

async function fetchPlaylists(): Promise<void> {
  try {
    playlists.value = await getPlaylistsRaw()
  } catch {
    toast.error(t('music.editor.failedLoad'))
  }
}

async function onSaved(): Promise<void> {
  show_form.value = false
  await fetchPlaylists()
}

async function remove(playlist: Playlist): Promise<void> {
  if (!playlist.id || !confirm(t('music.editor.confirmDelete', { title: playlist.title }))) return

  try {
    await deletePlaylistRaw(playlist.id)
    await fetchPlaylists()
  } catch {
    toast.error(t('music.editor.failedDelete'))
  }
}

onMounted(fetchPlaylists)
</script>
