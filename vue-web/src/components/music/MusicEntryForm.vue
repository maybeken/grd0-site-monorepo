<template>
  <div class="flex flex-col gap-4 rounded-xl border border-accent/20 bg-background p-4">
    <div class="flex items-center justify-between">
      <h3 class="text-lg font-bold">
        {{ is_edit ? $t('music.editor.songTitleEdit') : $t('music.editor.songTitleNew') }}
      </h3>
      <div class="flex gap-2">
        <button
          v-if="is_edit && deletable"
          class="rounded border border-red-500/40 px-3 py-1 text-sm text-red-400 hover:bg-red-500/10 disabled:opacity-50"
          :disabled="saving"
          @click="requestDelete"
        >
          {{ $t('common.delete') }}
        </button>
        <button
          class="rounded border border-accent/40 px-3 py-1 text-sm hover:bg-accent/10"
          @click="$emit('cancel')"
        >
          {{ $t('common.cancel') }}
        </button>
        <button
          class="rounded bg-accent px-3 py-1 text-sm font-semibold text-background hover:brightness-110 disabled:opacity-50"
          :disabled="saving || !can_save"
          @click="save"
        >
          {{ saving ? $t('common.saving') : $t('common.save') }}
        </button>
      </div>
    </div>

    <div class="flex flex-wrap items-center gap-2 text-sm">
      <span class="text-foreground/60">{{ $t('music.editor.videoField') }}:</span>
      <span class="font-mono">{{ form.v || $t('music.editor.videoUnset') }}</span>
      <button
        v-if="!is_edit"
        class="rounded border border-accent/40 px-2 py-1 text-xs hover:bg-accent/10"
        @click="picker = true"
      >
        {{ form.v ? $t('music.editor.changeLink') : $t('music.editor.resolveLink') }}
      </button>
    </div>

    <MusicPicker
      v-if="picker"
      :link_only="true"
      @select="onSelect"
      @close="picker = false"
    ></MusicPicker>

    <div
      v-if="duplicate"
      class="flex flex-wrap items-center gap-2 rounded border border-accent/30 bg-shade px-3 py-2 text-sm"
    >
      <span>{{ $t('music.editor.duplicateSong') }}</span>
      <button class="underline hover:text-accent" @click="editExisting">
        {{ $t('music.editor.editExisting') }}
      </button>
    </div>

    <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
      <div class="flex flex-col gap-1">
        <label class="text-sm font-semibold">{{ $t('music.editor.songTitleField') }}</label>
        <input
          v-model="form.title"
          class="rounded border border-accent/40 bg-background px-2 py-1 text-sm"
          :placeholder="$t('music.editor.songTitlePlaceholder')"
        />
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-sm font-semibold">{{ $t('music.editor.songArtistField') }}</label>
        <input
          v-model="form.artist"
          class="rounded border border-accent/40 bg-background px-2 py-1 text-sm"
          :placeholder="$t('music.editor.songArtistPlaceholder')"
        />
      </div>
    </div>

    <div class="flex flex-col gap-1">
      <label class="text-sm font-semibold">{{ $t('music.editor.notesField') }}</label>
      <textarea
        v-model="form.description"
        rows="3"
        class="resize-y rounded border border-accent/40 bg-background px-2 py-1 text-sm"
        :placeholder="$t('music.editor.notesPlaceholder')"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import MusicPicker from '@/components/editor/article/MusicPicker.vue'
import { upsertMusicRaw } from '@/services/music'
import { useToast } from '@/composables/useToast'
import { useUnsavedGuard } from '@/composables/useUnsavedGuard'

import type { Music } from '@/interfaces/Music'

interface Props {
  entry?: Music | null
  catalogue?: Music[]
  membership?: number
  deletable?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  entry: null,
  catalogue: () => [],
  membership: 0,
  deletable: false
})

const emit = defineEmits<{
  saved: []
  cancel: []
  delete: []
  edit: [v: string]
}>()

const { t } = useI18n()
const toast = useToast()

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

useUnsavedGuard(dirty)

const can_save = computed(() => !!form.v && !!form.title.trim())

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

function on_save_shortcut(event: KeyboardEvent): void {
  if (!(event.ctrlKey || event.metaKey) || event.key.toLowerCase() !== 's') return

  event.preventDefault()
  void save()
}

onMounted(() => window.addEventListener('keydown', on_save_shortcut))
onBeforeUnmount(() => window.removeEventListener('keydown', on_save_shortcut))
</script>
