<template>
  <EditorSection
    :title="is_edit ? $t('music.editor.songTitleEdit') : $t('music.editor.songTitleNew')"
  >
    <template #actions>
      <EditorButton
        v-if="is_edit && deletable"
        variant="destructive"
        :disabled="saving"
        @click="requestDelete"
      >
        {{ $t('common.delete') }}
      </EditorButton>
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
      <div class="flex flex-wrap items-center gap-2 text-lg font-medium">
        <span class="text-foreground/60">{{ $t('music.editor.videoField') }}:</span>
        <span class="font-mono">{{ form.v || $t('music.editor.videoUnset') }}</span>
        <EditorButton
          v-if="!is_edit"
          size="sm"
          variant="secondary"
          @click="picker = true"
        >
          {{ form.v ? $t('music.editor.changeLink') : $t('music.editor.resolveLink') }}
        </EditorButton>
      </div>

      <MusicPicker
        v-if="picker"
        :link_only="true"
        @select="onSelect"
        @close="picker = false"
      ></MusicPicker>

      <div
        v-if="duplicate"
        class="flex flex-wrap items-center gap-2 rounded border border-accent/30 bg-shade px-3 py-2 text-lg"
      >
        <span>{{ $t('music.editor.duplicateSong') }}</span>
        <button
          type="button"
          class="font-medium underline hover:text-foreground"
          @click="editExisting"
        >
          {{ $t('music.editor.editExisting') }}
        </button>
      </div>

      <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
        <EditorField :label="$t('music.editor.songTitleField')" html-for="song-title" required>
          <EditorInput
            id="song-title"
            v-model="form.title"
            :placeholder="$t('music.editor.songTitlePlaceholder')"
            required
          ></EditorInput>
        </EditorField>

        <EditorField :label="$t('music.editor.songArtistField')" html-for="song-artist">
          <EditorInput
            id="song-artist"
            v-model="form.artist"
            :placeholder="$t('music.editor.songArtistPlaceholder')"
          ></EditorInput>
        </EditorField>
      </div>

      <EditorField :label="$t('music.editor.notesField')" html-for="song-notes">
        <EditorTextarea
          id="song-notes"
          v-model="form.description"
          :placeholder="$t('music.editor.notesPlaceholder')"
        ></EditorTextarea>
      </EditorField>
    </form>
  </EditorSection>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import MusicPicker from '@/components/editor/article/MusicPicker.vue'
import EditorButton from '@/components/editor/ui/EditorButton.vue'
import EditorField from '@/components/editor/ui/EditorField.vue'
import EditorInput from '@/components/editor/ui/EditorInput.vue'
import EditorSection from '@/components/editor/ui/EditorSection.vue'
import EditorTextarea from '@/components/editor/ui/EditorTextarea.vue'
import { useConfirmDiscard } from '@/composables/useConfirmDiscard'
import { useSaveShortcut } from '@/composables/useSaveShortcut'
import { useToast } from '@/composables/useToast'
import { useUnsavedGuard } from '@/composables/useUnsavedGuard'
import { upsertMusicRaw } from '@/services/music'

import type { Music } from '@/interfaces/Music'

interface Props {
  entry?: Music | null
  catalogue?: Music[]
  membership?: number
  deletable?: boolean
  formId?: string
}

const props = withDefaults(defineProps<Props>(), {
  entry: null,
  catalogue: () => [],
  membership: 0,
  deletable: false,
  formId: 'music-entry-form'
})

const emit = defineEmits<{
  saved: []
  cancel: []
  delete: []
  edit: [v: string]
}>()

const { t } = useI18n()
const toast = useToast()
const confirmDiscard = useConfirmDiscard()

const is_edit = computed(() => !!props.entry?.v)
const saving = ref(false)
const picker = ref(false)
const duplicate = ref<Music | null>(null)

const form = reactive({
  v: '',
  title: '',
  artist: '',
  description: ''
})

const baseline = ref('')

function snapshot(): string {
  return JSON.stringify(form)
}

watch(
  () => props.entry,
  (entry) => {
    form.v = entry?.v ?? ''
    form.title = entry?.title ?? ''
    form.artist = entry?.artist ?? ''
    form.description = entry?.description ?? ''
    duplicate.value = null
    baseline.value = snapshot()
  },
  { immediate: true }
)

const dirty = computed(() => snapshot() !== baseline.value)
const can_save = computed(() => !!form.v && !!form.title.trim() && dirty.value)

useUnsavedGuard(dirty)
useSaveShortcut(save, { enabled: () => can_save.value && !saving.value })

function onSelect(song: { v: string; title: string; artist?: string }): void {
  picker.value = false

  const existing = props.catalogue.find((item) => item.v === song.v)

  if (existing) {
    duplicate.value = existing
    form.v = ''
    return
  }

  duplicate.value = null
  form.v = song.v
  form.title = song.title
  form.artist = song.artist ?? ''
}

function editExisting(): void {
  if (!duplicate.value) return

  emit('edit', duplicate.value.v)
}

function requestCancel(): void {
  if (!confirmDiscard(dirty.value)) return

  baseline.value = snapshot()
  emit('cancel')
}

function requestDelete(): void {
  if (!props.entry) return

  if (
    !confirm(
      t('music.editor.confirmDeleteSong', {
        title: props.entry.title || props.entry.v,
        count: props.membership
      })
    )
  ) {
    return
  }

  baseline.value = snapshot()
  emit('delete')
}

async function save(): Promise<void> {
  if (!can_save.value || saving.value) return

  if (!is_edit.value) {
    const existing = props.catalogue.find((item) => item.v === form.v)

    if (existing) {
      duplicate.value = existing
      return
    }
  }

  form.title = form.title.trim()
  form.artist = form.artist.trim()
  form.description = form.description.trim()

  saving.value = true

  try {
    await upsertMusicRaw({
      v: form.v,
      title: form.title,
      artist: form.artist,
      description: form.description
    })

    baseline.value = snapshot()
    toast.success(t('common.done'))
    emit('saved')
  } catch {
    toast.error(t('music.editor.failedSaveSong'))
  } finally {
    saving.value = false
  }
}

defineExpose({ requestCancel })
</script>
