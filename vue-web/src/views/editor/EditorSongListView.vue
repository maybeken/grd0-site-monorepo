<template>
  <EditorPage>
    <EditorPageHeader :title="$t('editor.songs')" :description="$t('music.editor.songsHint')">
      <template #actions>
        <EditorButton
          v-if="!show_form && !editing_v"
          variant="primary"
          @click="show_form = true"
        >
          <Icon icon="mynaui:plus" height="1rem"></Icon>
          {{ $t('music.editor.newSong') }}
        </EditorButton>
      </template>
    </EditorPageHeader>

    <MusicEntryForm
      v-if="show_form"
      :catalogue="songs"
      @saved="onSaved"
      @cancel="show_form = false"
      @edit="onEditExisting"
    ></MusicEntryForm>

    <EditorError v-else-if="load_failed" :message="t('music.editor.failedLoadSongs')">
      <template #actions>
        <EditorButton variant="secondary" @click="load(true)">
          {{ $t('common.retry') }}
        </EditorButton>
      </template>
    </EditorError>

    <EditorLoading v-else-if="loading"></EditorLoading>

    <template v-else>
      <p v-if="membership_failed" class="text-sm font-medium text-foreground/60">
        {{ $t('music.editor.failedMembership') }}
      </p>

      <div class="flex flex-col gap-2">
        <template v-for="song in songs" :key="song.v">
          <MusicEntryForm
            v-if="editing_v === song.v"
            :entry="song"
            :membership="membership(song.v)"
            :deletable="true"
            @cancel="editing_v = ''"
            @saved="onSaved"
            @delete="deleteSong(song)"
          ></MusicEntryForm>

          <div
            v-else
            class="flex items-center gap-3 rounded-lg border border-accent/20 bg-background p-3"
          >
            <img
              class="h-12 w-12 shrink-0 rounded object-cover"
              :src="`https://i.ytimg.com/vi_webp/${song.v}/default.webp`"
              loading="lazy"
            />

            <div class="flex min-w-0 flex-1 flex-col">
              <span class="truncate text-lg font-medium" :title="song.title || song.v">
                {{ song.title || song.v }}
              </span>
              <span class="truncate text-sm font-medium text-foreground/60">
                {{ song.artist || $t('music.unknown') }}
              </span>
              <span
                v-if="song.description"
                class="line-clamp-1 text-sm font-medium text-foreground/60"
              >
                {{ song.description }}
              </span>
              <span class="text-sm font-medium text-foreground/60">
                {{
                  $t(
                    'music.editor.playlistCount',
                    { count: membership(song.v) },
                    membership(song.v)
                  )
                }}
              </span>
            </div>

            <div class="flex shrink-0 gap-2">
              <EditorButton
                size="sm"
                variant="secondary"
                :disabled="!!editing_v"
                @click="editing_v = song.v"
              >
                {{ $t('common.edit') }}
              </EditorButton>
              <EditorButton
                size="sm"
                variant="destructive"
                :disabled="!!editing_v"
                @click="confirmDelete(song)"
              >
                {{ $t('common.delete') }}
              </EditorButton>
            </div>
          </div>
        </template>

        <EditorEmptyState
          v-if="!songs.length"
          :message="$t('music.editor.noSongsAvailable')"
        ></EditorEmptyState>
      </div>
    </template>
  </EditorPage>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import MusicEntryForm from '@/components/music/MusicEntryForm.vue'
import EditorButton from '@/components/editor/ui/EditorButton.vue'
import EditorEmptyState from '@/components/editor/ui/EditorEmptyState.vue'
import EditorError from '@/components/editor/ui/EditorError.vue'
import EditorLoading from '@/components/editor/ui/EditorLoading.vue'
import EditorPage from '@/components/editor/ui/EditorPage.vue'
import EditorPageHeader from '@/components/editor/ui/EditorPageHeader.vue'
import Icon from '@/components/shared/Icon.vue'
import { deleteMusicRaw, getMusicRaw, getPlaylistsRaw } from '@/services/music'
import { useToast } from '@/composables/useToast'

import type { Music } from '@/interfaces/Music'

const { t } = useI18n()
const toast = useToast()

const songs = ref<Music[]>([])
const counts = ref<Record<string, number>>({})
const loading = ref(true)
const load_failed = ref(false)
const membership_failed = ref(false)
const show_form = ref(false)
const editing_v = ref('')

function membership(v: string): number {
  return counts.value[v] ?? 0
}

async function load(initial = false): Promise<void> {
  if (initial) loading.value = true

  load_failed.value = false
  membership_failed.value = false

  const [music_result, playlist_result] = await Promise.allSettled([
    getMusicRaw(),
    getPlaylistsRaw()
  ])

  if (music_result.status === 'fulfilled' && Array.isArray(music_result.value)) {
    songs.value = music_result.value
  } else {
    songs.value = []
    load_failed.value = true
  }

  if (playlist_result.status === 'fulfilled' && Array.isArray(playlist_result.value)) {
    const next: Record<string, number> = {}

    for (const playlist of playlist_result.value) {
      for (const song of playlist.songs ?? []) {
        if (!song.v) continue
        next[song.v] = (next[song.v] ?? 0) + 1
      }
    }

    counts.value = next
  } else {
    counts.value = {}
    membership_failed.value = true
  }

  loading.value = false
}

function onEditExisting(v: string): void {
  show_form.value = false
  editing_v.value = v
}

async function onSaved(): Promise<void> {
  show_form.value = false
  editing_v.value = ''
  toast.success(t('common.done'))
  await load()
}

async function deleteSong(song: Music): Promise<void> {
  try {
    await deleteMusicRaw(song.v)

    if (editing_v.value === song.v) editing_v.value = ''
    toast.success(t('common.deleted'))
    await load()
  } catch {
    toast.error(t('music.editor.failedDeleteSong'))
  }
}

function confirmDelete(song: Music): void {
  if (
    !confirm(
      t('music.editor.confirmDeleteSong', {
        title: song.title || song.v,
        count: membership(song.v)
      })
    )
  ) {
    return
  }

  void deleteSong(song)
}

onMounted(() => load(true))
</script>
