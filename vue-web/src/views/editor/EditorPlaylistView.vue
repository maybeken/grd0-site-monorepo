<template>
  <div class="flex flex-col gap-4">
    <router-link to="/editor/music" class="text-sm text-secondary hover:underline">
      {{ $t('editor.backToList') }}
    </router-link>

    <p v-if="loading" class="text-secondary">{{ $t('common.loading') }}</p>

    <div v-else-if="playlist" class="max-w-4xl">
      <PlaylistForm :playlist="playlist" @saved="onSaved" @cancel="back"></PlaylistForm>
    </div>

    <div v-else class="flex flex-col items-start gap-3">
      <p class="text-red-400">{{ $t('editor.notFound') }}</p>
      <router-link to="/editor/music" class="rounded-full border border-dotted px-3 py-1">
        {{ $t('editor.backToList') }}
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

import PlaylistForm from '@/components/music/PlaylistForm.vue'
import { getPlaylistRaw } from '@/services/music'
import { useToast } from '@/composables/useToast'

import type { Playlist } from '@/interfaces/Music'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const toast = useToast()

const playlist = ref<Playlist | null>(null)
const loading = ref(true)

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
  router.push('/editor/music')
}

watch(id, load, { immediate: true })
</script>
