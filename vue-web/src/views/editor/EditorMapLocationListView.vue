<template>
  <EditorPage>
    <EditorPageHeader
      :title="$t('editor.mapLocations')"
      :description="$t('map.editor.locationsHint')"
    >
      <template #actions>
        <EditorButton v-if="!show_form" variant="primary" @click="openForm">
          <Icon icon="mynaui:plus" height="1rem"></Icon>
          {{ $t('map.editor.newLocation') }}
        </EditorButton>
      </template>
    </EditorPageHeader>

    <EditorSection v-if="show_form" :title="$t('map.editor.newLocation')">
      <form class="flex flex-col gap-4" @submit.prevent="create">
        <EditorField
          :label="$t('map.editor.titleField')"
          html-for="new-map-location-title"
          required
          :help="
            generated_slug
              ? $t('map.editor.slugField') + ': ' + generated_slug
              : $t('map.editor.slugImmutable')
          "
        >
          <EditorInput
            id="new-map-location-title"
            ref="title_input"
            v-model="new_title"
            :placeholder="$t('map.editor.titlePlaceholder')"
            :disabled="creating"
            required
          ></EditorInput>
        </EditorField>

        <p v-if="collision" class="text-sm font-medium text-red-400" role="alert">
          {{ $t('map.editor.slugCollision', { slug: generated_slug }) }}
        </p>

        <div class="flex justify-end gap-2">
          <EditorButton variant="secondary" :disabled="creating" @click="closeForm">
            {{ $t('common.cancel') }}
          </EditorButton>
          <EditorButton
            type="submit"
            variant="primary"
            :disabled="creating || !generated_slug"
            :loading="creating"
          >
            {{ $t('map.editor.create') }}
          </EditorButton>
        </div>
      </form>
    </EditorSection>

    <EditorError v-else-if="load_failed" :message="t('map.editor.failedLoad')">
      <template #actions>
        <EditorButton variant="secondary" @click="fetchLocations">
          {{ $t('common.retry') }}
        </EditorButton>
      </template>
    </EditorError>

    <EditorLoading v-else-if="loading"></EditorLoading>

    <div v-else class="flex flex-col gap-2">
      <EditorEntityRow
        v-for="location in locations"
        :key="location.slug"
        :title="location.title"
        :to="`/editor/travel/${location.slug}`"
        @delete="remove(location)"
      >
        <span
          v-if="location.subtitle"
          class="truncate text-sm font-medium text-foreground/60"
        >
          {{ location.subtitle }}
        </span>
        <span class="flex items-center gap-2">
          <span class="truncate font-mono text-sm font-medium text-foreground/60">
            {{ location.slug }}
          </span>
          <span
            class="rounded-full px-2 py-0.5 text-sm font-medium uppercase tracking-wide"
            :class="
              location.published
                ? 'bg-emerald-600/20 text-emerald-400'
                : 'bg-orange-500/20 text-orange-400'
            "
          >
            {{ location.published ? $t('map.editor.published') : $t('map.editor.draft') }}
          </span>
        </span>
      </EditorEntityRow>

      <EditorEmptyState
        v-if="!locations.length"
        :message="$t('map.editor.noLocations')"
      ></EditorEmptyState>
    </div>
  </EditorPage>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import EditorEntityRow from '@/components/editor/EditorEntityRow.vue'
import EditorButton from '@/components/editor/ui/EditorButton.vue'
import EditorEmptyState from '@/components/editor/ui/EditorEmptyState.vue'
import EditorError from '@/components/editor/ui/EditorError.vue'
import EditorField from '@/components/editor/ui/EditorField.vue'
import EditorInput from '@/components/editor/ui/EditorInput.vue'
import EditorLoading from '@/components/editor/ui/EditorLoading.vue'
import EditorPage from '@/components/editor/ui/EditorPage.vue'
import EditorPageHeader from '@/components/editor/ui/EditorPageHeader.vue'
import EditorSection from '@/components/editor/ui/EditorSection.vue'
import Icon from '@/components/shared/Icon.vue'
import {
  deleteMapLocation,
  getAllMapLocationsRaw,
  toKebabCase,
  upsertMapLocation
} from '@/services/travelersMap'
import { useToast } from '@/composables/useToast'

import type { MapLocation } from '@/interfaces/TravelersMap'

const { t } = useI18n()
const toast = useToast()
const router = useRouter()

const locations = ref<MapLocation[]>([])
const loading = ref(true)
const load_failed = ref(false)
const show_form = ref(false)
const creating = ref(false)
const new_title = ref('')
const title_input = ref<HTMLInputElement | null>(null)

const generated_slug = computed(() => toKebabCase(new_title.value))
const collision = computed(
  () =>
    !!generated_slug.value &&
    locations.value.some((location) => location.slug === generated_slug.value)
)

async function fetchLocations(): Promise<void> {
  loading.value = true
  load_failed.value = false

  try {
    const response = await getAllMapLocationsRaw()

    locations.value = Array.isArray(response) ? response : []
  } catch {
    locations.value = []
    load_failed.value = true
  } finally {
    loading.value = false
  }
}

async function openForm(): Promise<void> {
  new_title.value = ''
  show_form.value = true

  await nextTick()
  title_input.value?.focus()
}

function closeForm(): void {
  show_form.value = false
  new_title.value = ''
}

async function create(): Promise<void> {
  const slug = generated_slug.value
  const title = new_title.value.trim()

  if (!slug || !title || creating.value) return

  if (locations.value.some((location) => location.slug === slug)) {
    toast.error(t('map.editor.slugCollision', { slug }))
    return
  }

  creating.value = true

  try {
    await upsertMapLocation({
      slug,
      title,
      icon: '',
      color: 'default',
      pos: [0, 0],
      published: false
    })

    show_form.value = false
    new_title.value = ''
    toast.success(t('common.done'))
    await router.push(`/editor/travel/${slug}`)
  } catch {
    toast.error(t('map.editor.failedCreate'))
  } finally {
    creating.value = false
  }
}

async function remove(location: MapLocation): Promise<void> {
  if (!location.slug || !confirm(t('map.editor.confirmDelete', { title: location.title }))) {
    return
  }

  try {
    await deleteMapLocation(location.slug)
    locations.value = locations.value.filter((item) => item.slug !== location.slug)
    toast.success(t('common.deleted'))
  } catch {
    toast.error(t('map.editor.failedDelete'))
  }
}

onMounted(fetchLocations)
</script>
