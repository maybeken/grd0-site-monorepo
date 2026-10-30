<template>
  <EditorPage>
    <EditorPageHeader :title="$t('editor.vehicles')">
      <template #actions>
        <EditorButton :to="NEW_ROUTE" variant="primary">
          <Icon icon="mynaui:plus" height="1rem"></Icon>
          {{ $t('grandTour.editor.newVehicle') }}
        </EditorButton>
      </template>
    </EditorPageHeader>

    <EditorError v-if="load_failed" :message="t('editor.loadFailed')">
      <template #actions>
        <EditorButton variant="secondary" @click="fetchVehicles(true)">
          {{ $t('common.retry') }}
        </EditorButton>
      </template>
    </EditorError>

    <EditorLoading v-else-if="loading"></EditorLoading>

    <div v-else class="flex flex-col gap-2">
      <EditorEntityRow
        v-for="vehicle in vehicles"
        :key="vehicle.id"
        :title="`${vehicle.make} ${vehicle.model}`"
        :to="`${LISTING_ROUTE}/${vehicle.id}`"
        @delete="remove(vehicle)"
      >
        <span class="truncate text-sm font-medium text-foreground/60">
          {{ summary(vehicle) }}
        </span>
      </EditorEntityRow>

      <EditorEmptyState
        v-if="!vehicles.length"
        :message="$t('grandTour.editor.noVehicles')"
        :description="$t('grandTour.editor.noVehiclesHint')"
      ></EditorEmptyState>
    </div>
  </EditorPage>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import EditorEntityRow from '@/components/editor/EditorEntityRow.vue'
import EditorButton from '@/components/editor/ui/EditorButton.vue'
import EditorEmptyState from '@/components/editor/ui/EditorEmptyState.vue'
import EditorError from '@/components/editor/ui/EditorError.vue'
import EditorLoading from '@/components/editor/ui/EditorLoading.vue'
import EditorPage from '@/components/editor/ui/EditorPage.vue'
import EditorPageHeader from '@/components/editor/ui/EditorPageHeader.vue'
import Icon from '@/components/shared/Icon.vue'
import { useToast } from '@/composables/useToast'
import { deleteVehicleRaw, getVehicleCollection } from '@/services/vehicle'

import type { Vehicle } from '@/interfaces/Vehicle'

const LISTING_ROUTE = '/editor/vehicles'
const NEW_ROUTE = '/editor/vehicles/new'

const { t } = useI18n()
const toast = useToast()

const vehicles = ref<Vehicle[]>([])
const loading = ref(true)
const load_failed = ref(false)

function summary(vehicle: Vehicle): string {
  const cover = vehicle.silhouette_key
    ? t('grandTour.editor.coverAssigned')
    : t('grandTour.editor.coverMissing')
  const location = [vehicle.country, vehicle.city].filter(Boolean).join(', ')

  return [vehicle.year, location, cover].filter(Boolean).join(' · ')
}

async function fetchVehicles(force = false): Promise<void> {
  loading.value = true
  load_failed.value = false

  try {
    vehicles.value = await getVehicleCollection(force)
  } catch {
    vehicles.value = []
    load_failed.value = true
  } finally {
    loading.value = false
  }
}

async function remove(vehicle: Vehicle): Promise<void> {
  const name = `${vehicle.make} ${vehicle.model}`

  if (!vehicle.id || !confirm(t('common.confirmDeleteItem', { name }))) return

  try {
    await deleteVehicleRaw(vehicle.id)
    toast.success(t('common.deleted'))
    await fetchVehicles()
  } catch {
    toast.error(t('grandTour.editor.failedDelete'))
  }
}

onMounted(() => fetchVehicles())
</script>
