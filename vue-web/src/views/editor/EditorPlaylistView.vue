<template>
  <EditorPage>
    <EditorPageHeader
      :title="playlist?.title || $t('editor.playlists')"
      :back-to="PLAYLISTS_ROUTE"
      :back-label="$t('editor.ui.backTo', { section: $t('editor.playlists') })"
    >
      <template #actions>
        <template v-if="playlist">
          <EditorButton variant="secondary" @click="cancel">
            {{ $t('common.cancel') }}
          </EditorButton>
          <EditorButton
            type="submit"
            :form="FORM_ID"
            variant="primary"
            :disabled="saving || !can_save"
            :loading="saving"
          >
            {{ $t('common.save') }}
          </EditorButton>
        </template>
      </template>
    </EditorPageHeader>

    <EditorLoading v-if="loading"></EditorLoading>

    <EditorError v-else-if="!playlist" :message="$t('editor.notFound')">
      <template #actions>
        <EditorButton variant="secondary" :to="PLAYLISTS_ROUTE">
          {{ $t('editor.ui.backTo', { section: $t('editor.playlists') }) }}
        </EditorButton>
      </template>
    </EditorError>

    <PlaylistForm
      v-else
      ref="playlist_form"
      :playlist="playlist"
      :standalone="false"
      :form-id="FORM_ID"
      @saved="onSaved"
      @cancel="back"
      @update:saving="saving = $event"
      @update:can-save="can_save = $event"
    ></PlaylistForm>
  </EditorPage>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

import PlaylistForm from '@/components/music/PlaylistForm.vue'
import EditorButton from '@/components/editor/ui/EditorButton.vue'
import EditorError from '@/components/editor/ui/EditorError.vue'
import EditorLoading from '@/components/editor/ui/EditorLoading.vue'
import EditorPage from '@/components/editor/ui/EditorPage.vue'
import EditorPageHeader from '@/components/editor/ui/EditorPageHeader.vue'
import { useToast } from '@/composables/useToast'
import { getPlaylistRaw } from '@/services/music'

import type { Playlist } from '@/interfaces/Music'

const PLAYLISTS_ROUTE = '/editor/music'
const FORM_ID = 'playlist-form'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const toast = useToast()

const playlist = ref<Playlist | null>(null)
const playlist_form = ref<InstanceType<typeof PlaylistForm> | null>(null)
const loading = ref(true)
const saving = ref(false)
const can_save = ref(false)

const id = computed(() => {
  const value = route.params.id

  return typeof value === 'string' ? value : ''
})

async function load(): Promise<void> {
  loading.value = true
  playlist.value = null

  try {
    playlist.value = await getPlaylistRaw(id.value)
  } catch {
    playlist.value = null
  } finally {
    loading.value = false
  }
}

function onSaved(): void {
  toast.success(t('common.done'))
}

function back(): void {
  router.push(PLAYLISTS_ROUTE)
}

function cancel(): void {
  playlist_form.value?.requestCancel()
}

watch(id, load, { immediate: true })
</script>
