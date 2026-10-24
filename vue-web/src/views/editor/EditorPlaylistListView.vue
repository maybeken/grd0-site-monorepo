<template>
  <EditorPage>
    <EditorPageHeader :title="$t('editor.playlists')" :description="$t('music.editor.listHint')">
      <template #actions>
        <EditorButton v-if="!show_form" variant="primary" @click="show_form = true">
          <Icon icon="mynaui:plus" height="1rem"></Icon>
          {{ $t('music.editor.newPlaylist') }}
        </EditorButton>
      </template>
    </EditorPageHeader>

    <PlaylistForm
      v-if="show_form"
      :playlist="null"
      @saved="onSaved"
      @cancel="show_form = false"
    ></PlaylistForm>

    <EditorError v-else-if="load_failed" :message="t('music.editor.failedLoad')">
      <template #actions>
        <EditorButton variant="secondary" @click="fetchPlaylists">
          {{ $t('common.retry') }}
        </EditorButton>
      </template>
    </EditorError>

    <EditorLoading v-else-if="loading"></EditorLoading>

    <div v-else class="flex flex-col gap-2">
      <EditorEntityRow
        v-for="playlist in playlists"
        :key="playlist.id"
        :title="playlist.title"
        :to="`/editor/music/${playlist.id}`"
        @delete="remove(playlist)"
      >
        <span
          v-if="playlist.description"
          class="truncate text-sm font-medium text-foreground/60"
        >
          {{ playlist.description }}
        </span>
        <span class="text-sm font-medium text-foreground/60">
          {{
            $t(
              'music.editor.songCount',
              { count: playlist.songs?.length ?? 0 },
              playlist.songs?.length ?? 0
            )
          }}
        </span>
      </EditorEntityRow>

      <EditorEmptyState
        v-if="!playlists.length"
        :message="$t('music.editor.noPlaylists')"
      ></EditorEmptyState>
    </div>
  </EditorPage>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import PlaylistForm from '@/components/music/PlaylistForm.vue'
import EditorEntityRow from '@/components/editor/EditorEntityRow.vue'
import EditorButton from '@/components/editor/ui/EditorButton.vue'
import EditorEmptyState from '@/components/editor/ui/EditorEmptyState.vue'
import EditorError from '@/components/editor/ui/EditorError.vue'
import EditorLoading from '@/components/editor/ui/EditorLoading.vue'
import EditorPage from '@/components/editor/ui/EditorPage.vue'
import EditorPageHeader from '@/components/editor/ui/EditorPageHeader.vue'
import Icon from '@/components/shared/Icon.vue'
import { deletePlaylistRaw, getPlaylistsRaw } from '@/services/music'
import { useToast } from '@/composables/useToast'

import type { Playlist } from '@/interfaces/Music'

const { t } = useI18n()
const toast = useToast()

const playlists = ref<Playlist[]>([])
const show_form = ref(false)
const loading = ref(true)
const load_failed = ref(false)

async function fetchPlaylists(): Promise<void> {
  loading.value = true
  load_failed.value = false

  try {
    playlists.value = await getPlaylistsRaw()
  } catch {
    playlists.value = []
    load_failed.value = true
  } finally {
    loading.value = false
  }
}

async function onSaved(): Promise<void> {
  show_form.value = false
  toast.success(t('common.done'))
  await fetchPlaylists()
}

async function remove(playlist: Playlist): Promise<void> {
  if (!playlist.id || !confirm(t('music.editor.confirmDelete', { title: playlist.title }))) return

  try {
    await deletePlaylistRaw(playlist.id)
    toast.success(t('common.deleted'))
    await fetchPlaylists()
  } catch {
    toast.error(t('music.editor.failedDelete'))
  }
}

onMounted(fetchPlaylists)
</script>
