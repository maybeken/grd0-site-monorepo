<template>
  <div v-if="authorized" class="flex flex-col gap-6 max-w-4xl mx-auto">
    <header class="flex flex-col gap-2">
      <h1 class="text-2xl font-bold">{{ $t('music.editor.title') }}<CursorBlink /></h1>
      <p class="text-secondary">{{ $t('music.editor.listHint') }}</p>
    </header>

    <div v-if="!show_form" class="flex justify-end">
      <button
        @click="openNew"
        class="px-3 py-1 rounded bg-accent text-background text-sm font-semibold hover:brightness-110"
      >
        {{ $t('music.editor.newPlaylist') }}
      </button>
    </div>

    <PlaylistForm
      v-if="show_form"
      :playlist="editing"
      @saved="onSaved"
      @cancel="closeForm"
    ></PlaylistForm>

    <div v-if="!show_form" class="flex flex-col gap-2">
      <div
        v-for="playlist in playlists"
        :key="playlist.id"
        class="flex items-center justify-between bg-background rounded-lg p-3 border border-accent/20"
      >
        <div class="flex min-w-0 flex-col">
          <span class="font-semibold text-sm">{{ playlist.title }}</span>
          <span v-if="playlist.description" class="truncate text-xs text-foreground/60">{{
            playlist.description
          }}</span>
          <span class="text-xs text-foreground/40">{{
            $t('music.editor.songCount', { count: playlist.songs?.length ?? 0 })
          }}</span>
        </div>
        <div class="flex gap-2">
          <button
            @click="editPlaylist(playlist)"
            class="px-2 py-1 rounded text-xs border border-accent/40 hover:bg-accent/10"
          >
            {{ $t('common.edit') }}
          </button>
          <button
            @click="removePlaylist(playlist)"
            class="px-2 py-1 rounded text-xs border border-red-500/40 text-red-400 hover:bg-red-500/10"
          >
            {{ $t('common.delete') }}
          </button>
        </div>
      </div>
      <p v-if="playlists.length === 0" class="text-sm text-foreground/40 text-center py-4">
        {{ $t('music.editor.noPlaylists') }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { deletePlaylistRaw, getPlaylistsRaw } from '@/services/music'
import { isLocalEnvironment } from '@/helpers/environment'
import { useToast } from '@/composables/useToast'

import type { Playlist } from '@/interfaces/Music'

const { t } = useI18n()
const toast = useToast()

const authorized = ref(false)
const playlists = ref<Playlist[]>([])
const show_form = ref(false)
const editing = ref<Playlist | null>(null)

async function fetchPlaylists(): Promise<void> {
  try {
    playlists.value = await getPlaylistsRaw()
  } catch {
    toast.error(t('music.editor.failedLoad'))
  }
}

function openNew(): void {
  editing.value = null
  show_form.value = true
}

function editPlaylist(playlist: Playlist): void {
  editing.value = playlist
  show_form.value = true
}

function closeForm(): void {
  show_form.value = false
  editing.value = null
}

async function onSaved(): Promise<void> {
  closeForm()
  await fetchPlaylists()
}

async function removePlaylist(playlist: Playlist): Promise<void> {
  if (!playlist.id || !confirm(t('music.editor.confirmDelete', { title: playlist.title }))) return

  try {
    await deletePlaylistRaw(playlist.id)
    await fetchPlaylists()
  } catch {
    toast.error(t('music.editor.failedDelete'))
  }
}

onMounted(() => {
  if (!isLocalEnvironment() && !sessionStorage.getItem('jwt_token')) {
    window.location.href = '/redirect/auth/login'
    return
  }

  authorized.value = true
  fetchPlaylists()
})
</script>
