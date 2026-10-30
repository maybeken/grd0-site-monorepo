<template>
  <EditorPage>
    <EditorPageHeader
      :title="heading"
      :back-to="LISTING_ROUTE"
      :back-label="$t('editor.ui.backTo', { section: $t('editor.vehicles') })"
    >
      <template #actions>
        <template v-if="show_form">
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

    <EditorError v-else-if="load_failed" :message="$t('grandTour.editor.notFound')">
      <template #actions>
        <EditorButton variant="secondary" :to="LISTING_ROUTE">
          {{ $t('editor.ui.backTo', { section: $t('editor.vehicles') }) }}
        </EditorButton>
      </template>
    </EditorError>

    <VehicleForm
      v-else
      ref="vehicle_form"
      :vehicle="vehicle"
      :standalone="false"
      :form-id="FORM_ID"
      @saved="onSaved"
      @cancel="back"
      @update:saving="saving = $event"
      @update:can-save="can_save = $event"
    ></VehicleForm>
  </EditorPage>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

import EditorButton from '@/components/editor/ui/EditorButton.vue'
import EditorError from '@/components/editor/ui/EditorError.vue'
import EditorLoading from '@/components/editor/ui/EditorLoading.vue'
import EditorPage from '@/components/editor/ui/EditorPage.vue'
import EditorPageHeader from '@/components/editor/ui/EditorPageHeader.vue'
import VehicleForm from '@/components/vehicle/VehicleForm.vue'
import { useToast } from '@/composables/useToast'
import { getVehicleRaw } from '@/services/vehicle'

import type { Vehicle } from '@/interfaces/Vehicle'

const LISTING_ROUTE = '/editor/vehicles'
const NEW_ID = 'new'
const FORM_ID = 'vehicle-form'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const toast = useToast()

const vehicle = ref<Vehicle | null>(null)
const vehicle_form = ref<InstanceType<typeof VehicleForm> | null>(null)
const loading = ref(true)
const load_failed = ref(false)
const saving = ref(false)
const can_save = ref(false)

const id = computed(() => (typeof route.params.id === 'string' ? route.params.id : ''))
const is_new = computed(() => id.value === NEW_ID)

const show_form = computed(() => !loading.value && !load_failed.value)

const heading = computed(() => {
  if (is_new.value) return t('grandTour.editor.form.titleNew')
  if (vehicle.value) return `${vehicle.value.make} ${vehicle.value.model}`
  return t('editor.vehicles')
})

async function load(): Promise<void> {
  loading.value = true
  load_failed.value = false
  vehicle.value = null

  if (is_new.value) {
    loading.value = false
    return
  }

  try {
    vehicle.value = await getVehicleRaw(id.value)
    load_failed.value = vehicle.value === null
  } catch {
    load_failed.value = true
  } finally {
    loading.value = false
  }
}

function onSaved(saved: Vehicle): void {
  toast.success(t('common.done'))

  if (is_new.value && saved.id) {
    router.replace(`${LISTING_ROUTE}/${saved.id}`)
  }
}

function back(): void {
  router.push(LISTING_ROUTE)
}

function cancel(): void {
  vehicle_form.value?.requestCancel()
}

watch(id, load, { immediate: true })
</script>
