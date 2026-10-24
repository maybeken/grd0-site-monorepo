<template>
  <EditorSection
    :title="
      standalone
        ? is_edit
          ? $t('music.editor.titleEdit')
          : $t('music.editor.titleNew')
        : ''
    "
  >
    <template v-if="standalone" #actions>
      <EditorButton variant="secondary" @click="requestCancel">
        {{ $t('common.cancel') }}
      </EditorButton>
      <EditorButton
        type="submit"
        :form="formId"
        variant="primary"
        :disabled="saving || !can_save"
        :loading="saving"
      >
        {{ $t('common.save') }}
      </EditorButton>
    </template>

    <form :id="formId" class="flex flex-col gap-4" @submit.prevent="save">
      <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
        <EditorField :label="$t('music.editor.titleField')" html-for="playlist-title" required>
          <EditorInput
            id="playlist-title"
            v-model="form.title"
            :placeholder="$t('music.editor.titlePlaceholder')"
            required
          ></EditorInput>
        </EditorField>

        <EditorField :label="$t('music.editor.sortingField')" html-for="playlist-sorting">
          <EditorInput
            id="playlist-sorting"
            v-model="sorting_value"
            type="number"
            :number="true"
          ></EditorInput>
        </EditorField>
      </div>

      <EditorField :label="$t('music.editor.descriptionField')" html-for="playlist-description">
        <EditorTextarea
          id="playlist-description"
          v-model="form.description"
          :rows="2"
          :placeholder="$t('music.editor.descriptionPlaceholder')"
        ></EditorTextarea>
      </EditorField>

      <div class="flex flex-col gap-2 rounded-lg border border-accent/20 p-3">
        <div class="flex flex-wrap items-center justify-between gap-2">
          <p class="text-lg font-medium">
            {{
              $t(
                'music.editor.songCount',
                { count: selected_ids.length },
                selected_ids.length
              )
            }}
          </p>
          <EditorButton size="sm" variant="secondary" @click="picker = true">
            {{ $t('music.editor.addSong') }}
          </EditorButton>
        </div>

        <MusicPicker
          v-if="picker"
          :link_only="true"
          @select="onAddSong"
          @close="picker = false"
        ></MusicPicker>

        <EditorEmptyState
          v-if="!catalogue_songs || catalogue_songs.length === 0"
          :message="$t('music.editor.noSongsAvailable')"
        ></EditorEmptyState>

        <div v-else class="flex max-h-72 flex-col overflow-y-auto rounded border border-accent/20">
          <label
            v-for="song in catalogue_songs"
            :key="song.id || song.v"
            class="flex cursor-pointer items-center gap-3 border-b border-accent/10 px-3 py-2 hover:bg-shade"
          >
            <input
              type="checkbox"
              class="h-4 w-4 accent-accent"
              :checked="song.id ? selected_ids.includes(song.id) : false"
              @change="toggle(song)"
            />
            <img
              class="h-10 w-10 rounded object-cover"
              :src="`https://i.ytimg.com/vi_webp/${song.v}/default.webp`"
              loading="lazy"
            />
            <span class="flex min-w-0 flex-1 flex-col">
              <span class="truncate text-lg font-medium">{{ song.title }}</span>
              <span class="truncate text-sm font-medium text-foreground/60">
                {{ song.artist || $t('music.unknown') }}
              </span>
            </span>
            <EditorButton
              size="sm"
              variant="destructive"
              :disabled="!song.id || !selected_ids.includes(song.id)"
              @click.stop.prevent="removeFromPlaylist(song)"
            >
              {{ $t('music.editor.removeFromPlaylist') }}
            </EditorButton>
          </label>
        </div>
      </div>
    </form>
  </EditorSection>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch, watchEffect } from 'vue'
import { useI18n } from 'vue-i18n'

import MusicPicker from '@/components/editor/article/MusicPicker.vue'
import EditorButton from '@/components/editor/ui/EditorButton.vue'
import EditorEmptyState from '@/components/editor/ui/EditorEmptyState.vue'
import EditorField from '@/components/editor/ui/EditorField.vue'
import EditorInput from '@/components/editor/ui/EditorInput.vue'
import EditorSection from '@/components/editor/ui/EditorSection.vue'
import EditorTextarea from '@/components/editor/ui/EditorTextarea.vue'
import { useConfirmDiscard } from '@/composables/useConfirmDiscard'
import { useSaveShortcut } from '@/composables/useSaveShortcut'
import { useToast } from '@/composables/useToast'
import { useUnsavedGuard } from '@/composables/useUnsavedGuard'
import { getMusicRaw, upsertMusicRaw, upsertPlaylistRaw } from '@/services/music'

import type { Music, Playlist } from '@/interfaces/Music'

interface Props {
  playlist?: Playlist | null
  standalone?: boolean
  formId?: string
}

const props = withDefaults(defineProps<Props>(), {
  playlist: null,
  standalone: true,
  formId: 'playlist-form'
})

const emit = defineEmits<{
  saved: []
  cancel: []
  'update:saving': [saving: boolean]
  'update:can-save': [can_save: boolean]
  'update:dirty': [dirty: boolean]
}>()

const { t } = useI18n()
const toast = useToast()
const confirmDiscard = useConfirmDiscard()

const catalogue_songs = ref<Music[]>([])
const catalogue_loaded = ref(false)
const picker = ref(false)

async function fetchCatalogue(): Promise<void> {
  const data = await getMusicRaw()

  if (Array.isArray(data)) {
    catalogue_songs.value = data
    catalogue_loaded.value = true
  }
}

const catalogue_ready = fetchCatalogue().catch(() => undefined)

const is_edit = computed(() => !!props.playlist?.id)
const saving = ref(false)
const selected_ids = ref<string[]>([])

const form = reactive({
  title: '',
  description: '',
  sorting: 0
})

const baseline = ref('')

const sorting_value = computed({
  get: () => form.sorting,
  set: (value: string | number | null) => {
    form.sorting = typeof value === 'number' ? value : 0
  }
})

function snapshot(): string {
  return JSON.stringify({ form, selected_ids: selected_ids.value })
}

watch(
  () => props.playlist,
  (playlist) => {
    form.title = playlist?.title ?? ''
    form.description = playlist?.description ?? ''
    form.sorting = playlist?.sorting ?? 0
    selected_ids.value = (playlist?.songs ?? [])
      .map((song) => song.id)
      .filter((id): id is string => !!id)

    baseline.value = snapshot()
  },
  { immediate: true }
)

const dirty = computed(() => snapshot() !== baseline.value)
const can_save = computed(() => !!form.title.trim() && dirty.value)

useUnsavedGuard(dirty)
useSaveShortcut(save, { enabled: () => can_save.value && !saving.value })

watchEffect(() => emit('update:saving', saving.value))
watchEffect(() => emit('update:can-save', can_save.value))
watchEffect(() => emit('update:dirty', dirty.value))

function toggle(song: Music): void {
  if (!song.id) return

  if (selected_ids.value.includes(song.id)) {
    selected_ids.value = selected_ids.value.filter((id) => id !== song.id)
  } else {
    selected_ids.value = [...selected_ids.value, song.id]
  }
}

function removeFromPlaylist(song: Music): void {
  if (!song.id) return

  selected_ids.value = selected_ids.value.filter((id) => id !== song.id)
}

function requestCancel(): void {
  if (!confirmDiscard(dirty.value)) return

  baseline.value = snapshot()
  emit('cancel')
}

async function onAddSong(song: { v: string; title: string; artist?: string }): Promise<void> {
  if (!song.v) return

  try {
    await catalogue_ready

    const existing = catalogue_songs.value.find((item) => item.v === song.v)

    if (existing) {
      if (existing.title !== song.title || (existing.artist ?? '') !== (song.artist ?? '')) {
        toast.info(t('music.editor.correctOnSongsPage'))
      }
    } else {
      await upsertMusicRaw({ v: song.v, title: song.title, artist: song.artist })
      await fetchCatalogue()
    }

    const entry = catalogue_songs.value.find((item) => item.v === song.v)

    if (entry?.id && !selected_ids.value.includes(entry.id)) {
      selected_ids.value = [...selected_ids.value, entry.id]
    }

    picker.value = false
  } catch {
    toast.error(t('music.editor.failedAddSong'))
  }
}

async function save(): Promise<void> {
  if (!can_save.value || saving.value) return

  saving.value = true

  try {
    const songs = catalogue_loaded.value
      ? catalogue_songs.value.filter((song) => song.id && selected_ids.value.includes(song.id))
      : (props.playlist?.songs ?? [])

    await upsertPlaylistRaw({
      id: props.playlist?.id,
      title: form.title.trim(),
      description: form.description,
      sorting: form.sorting,
      songs
    })

    baseline.value = snapshot()
    emit('saved')
  } catch {
    toast.error(t('music.editor.failedSave'))
  } finally {
    saving.value = false
  }
}

defineExpose({ requestCancel })
</script>
