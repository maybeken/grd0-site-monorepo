<template>
  <div class="flex w-full flex-col gap-6">
    <EditorSection
      :title="
        standalone ? (isEdit ? $t('grandTour.editor.form.titleEdit') : $t('grandTour.editor.form.titleNew')) : ''
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

      <form :id="formId" class="flex flex-col gap-6" novalidate @submit.prevent="save">
        <EditorSection :title="$t('grandTour.editor.form.identity')">
          <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
            <EditorField
              :label="$t('grandTour.editor.form.make')"
              html-for="vehicle-make"
              required
              :error="errors.make"
            >
              <EditorInput id="vehicle-make" v-model="form.make" required></EditorInput>
            </EditorField>

            <EditorField
              :label="$t('grandTour.editor.form.model')"
              html-for="vehicle-model"
              required
              :error="errors.model"
            >
              <EditorInput id="vehicle-model" v-model="form.model" required></EditorInput>
            </EditorField>

            <EditorField
              :label="$t('grandTour.editor.form.year')"
              html-for="vehicle-year"
              required
              :error="errors.year"
            >
              <EditorInput
                id="vehicle-year"
                v-model="form.year"
                type="number"
                number
                required
              ></EditorInput>
            </EditorField>

            <EditorField :label="$t('grandTour.editor.form.bodyType')" html-for="vehicle-body-type">
              <EditorSelect id="vehicle-body-type" v-model="form.body_type">
                <option value="">{{ $t('grandTour.editor.form.none') }}</option>
                <option v-for="option in BODY_TYPES" :key="option.value" :value="option.value">
                  {{ $t(option.label) }}
                </option>
              </EditorSelect>
            </EditorField>
          </div>
        </EditorSection>

        <EditorSection :title="$t('grandTour.editor.form.specs')">
          <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
            <EditorField :label="$t('grandTour.editor.form.engine')" html-for="vehicle-engine">
              <EditorInput id="vehicle-engine" v-model="form.engine"></EditorInput>
            </EditorField>

            <EditorField :label="$t('grandTour.editor.form.power')" html-for="vehicle-power">
              <EditorInput id="vehicle-power" v-model="form.power"></EditorInput>
            </EditorField>

            <EditorField
              :label="$t('grandTour.editor.form.transmission')"
              html-for="vehicle-transmission"
            >
              <EditorInput id="vehicle-transmission" v-model="form.transmission"></EditorInput>
            </EditorField>

            <EditorField
              :label="$t('grandTour.editor.form.drivetrain')"
              html-for="vehicle-drivetrain"
            >
              <EditorInput id="vehicle-drivetrain" v-model="form.drivetrain"></EditorInput>
            </EditorField>

            <EditorField :label="$t('grandTour.editor.form.fuel')" html-for="vehicle-fuel">
              <EditorInput id="vehicle-fuel" v-model="form.fuel"></EditorInput>
            </EditorField>

            <EditorField :label="$t('grandTour.editor.form.rating')" html-for="vehicle-rating">
              <EditorInput
                id="vehicle-rating"
                v-model="form.rating"
                type="number"
                number
              ></EditorInput>
            </EditorField>
          </div>
        </EditorSection>

        <EditorSection :title="$t('grandTour.editor.form.location')">
          <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
            <EditorField
              :label="$t('grandTour.editor.form.country')"
              html-for="vehicle-country"
              :help="$t('grandTour.editor.form.countryHint')"
            >
              <CountryField
                id="vehicle-country"
                :code="form.country_code"
                :country="form.country"
                @select="setCountry"
              ></CountryField>
            </EditorField>

            <EditorField :label="$t('grandTour.editor.form.city')" html-for="vehicle-city">
              <EditorInput id="vehicle-city" v-model="form.city"></EditorInput>
            </EditorField>

            <EditorField :label="$t('grandTour.editor.form.drivenAt')" html-for="vehicle-driven-at">
              <EditorInput id="vehicle-driven-at" v-model="form.driven_at" type="date"></EditorInput>
            </EditorField>

            <EditorField
              :label="$t('grandTour.editor.form.sortOrder')"
              html-for="vehicle-sort-order"
              :help="$t('grandTour.editor.form.sortOrderHint')"
            >
              <EditorInput
                id="vehicle-sort-order"
                v-model="form.sort_order"
                type="number"
                number
              ></EditorInput>
            </EditorField>
          </div>

          <EditorField
            class="mt-4"
            :label="$t('grandTour.editor.form.notes')"
            html-for="vehicle-notes"
          >
            <EditorTextarea id="vehicle-notes" v-model="form.notes" :rows="4"></EditorTextarea>
          </EditorField>
        </EditorSection>

        <EditorSection :title="$t('grandTour.editor.cover.title')">
          <div class="flex flex-col gap-4">
            <div class="flex flex-wrap items-center gap-4">
              <div class="h-24 w-40 shrink-0 text-foreground md:h-28 md:w-48">
                <div v-if="cover_preview" class="flex h-full w-full items-center justify-center">
                  <img
                    v-if="cover_preview.kind === 'cutout'"
                    :src="cover_preview.src"
                    :alt="`${form.make} ${form.model}`"
                    class="h-full w-full object-contain"
                  />
                  <SilhouetteCover
                    v-else
                    :src="cover_preview.src"
                    :alt="`${form.make} ${form.model}`"
                  ></SilhouetteCover>
                </div>
                <div
                  v-else
                  class="flex h-full w-full items-center justify-center rounded border border-dashed border-accent/40 text-lg font-bold uppercase tracking-widest text-foreground/40"
                >
                  {{ $t('grandTour.unknownCover') }}
                </div>
              </div>
              <div class="flex min-w-0 flex-col gap-1">
                <p class="text-lg font-medium">{{ $t('grandTour.editor.cover.current') }}</p>
                <p class="text-sm font-medium text-foreground/60">{{ cover_status }}</p>
              </div>
            </div>

            <VehicleCoverTool @generated="onGenerated" @discarded="onGeneratedDiscarded"></VehicleCoverTool>
          </div>
        </EditorSection>
      </form>
    </EditorSection>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch, watchEffect } from 'vue'
import { useI18n } from 'vue-i18n'

import EditorButton from '@/components/editor/ui/EditorButton.vue'
import EditorField from '@/components/editor/ui/EditorField.vue'
import EditorInput from '@/components/editor/ui/EditorInput.vue'
import EditorSection from '@/components/editor/ui/EditorSection.vue'
import EditorSelect from '@/components/editor/ui/EditorSelect.vue'
import EditorTextarea from '@/components/editor/ui/EditorTextarea.vue'
import CountryField from '@/components/vehicle/CountryField.vue'
import SilhouetteCover from '@/components/vehicle/SilhouetteCover.vue'
import VehicleCoverTool from '@/components/vehicle/VehicleCoverTool.vue'
import { useConfirmDiscard } from '@/composables/useConfirmDiscard'
import { useSaveShortcut } from '@/composables/useSaveShortcut'
import { useToast } from '@/composables/useToast'
import { useUnsavedGuard } from '@/composables/useUnsavedGuard'
import { assetUrl, BODY_TYPES } from '@/helpers/vehicle'
import {
  requestVehicleAttachmentRaw,
  uploadVehicleCover,
  upsertVehicleRaw
} from '@/services/vehicle'

import type { GeneratedCover, Vehicle } from '@/interfaces/Vehicle'

interface FormState {
  id?: string
  make: string
  model: string
  year: number | string | null
  body_type: string | null
  engine: string | null
  power: string | null
  transmission: string | null
  drivetrain: string | null
  fuel: string | null
  country_code: string | null
  country: string | null
  city: string | null
  driven_at: string | null
  notes: string | null
  rating: number | string | null
  sort_order: number | string | null
  cutout_key: string | null
  silhouette_key: string | null
}

interface Props {
  vehicle?: Vehicle | null
  standalone?: boolean
  formId?: string
}

const props = withDefaults(defineProps<Props>(), {
  vehicle: null,
  standalone: true,
  formId: 'vehicle-form'
})

const emit = defineEmits<{
  saved: [vehicle: Vehicle]
  cancel: []
  'update:saving': [saving: boolean]
  'update:can-save': [can_save: boolean]
  'update:dirty': [dirty: boolean]
}>()

const { t } = useI18n()
const toast = useToast()
const confirmDiscard = useConfirmDiscard()

const isEdit = computed(() => !!props.vehicle?.id)
const saving = ref(false)
const pending_cover = ref<GeneratedCover | null>(null)

function createEmptyForm(): FormState {
  return {
    make: '',
    model: '',
    year: null,
    body_type: null,
    engine: null,
    power: null,
    transmission: null,
    drivetrain: null,
    fuel: null,
    country_code: null,
    country: null,
    city: null,
    driven_at: null,
    notes: null,
    rating: null,
    sort_order: null,
    cutout_key: null,
    silhouette_key: null
  }
}

function fromVehicle(vehicle: Vehicle): FormState {
  return {
    ...createEmptyForm(),
    ...vehicle,
    driven_at: vehicle.driven_at ? (vehicle.driven_at.split('T')[0] ?? null) : null
  }
}

const form = reactive<FormState>(createEmptyForm())
const errors = reactive<{ make: string; model: string; year: string }>({
  make: '',
  model: '',
  year: ''
})
const baseline = ref('')

function snapshot(): string {
  return JSON.stringify(form)
}

watch(
  () => props.vehicle,
  (vehicle) => {
    Object.assign(form, vehicle ? fromVehicle(vehicle) : createEmptyForm())
    errors.make = ''
    errors.model = ''
    errors.year = ''
    baseline.value = snapshot()
  },
  { immediate: true }
)

const dirty = computed(() => snapshot() !== baseline.value || pending_cover.value !== null)
const can_save = computed(() => dirty.value)

const cover_preview = computed<{ src: string; kind: 'cutout' | 'silhouette' } | null>(() => {
  if (pending_cover.value) {
    return { src: pending_cover.value.cutout_url, kind: 'cutout' }
  }

  if (form.cutout_key) {
    return { src: assetUrl(form.cutout_key), kind: 'cutout' }
  }

  if (form.silhouette_key) {
    return { src: assetUrl(form.silhouette_key), kind: 'silhouette' }
  }

  return null
})

const cover_status = computed(() => {
  if (pending_cover.value) return t('grandTour.editor.cover.pending')

  if (form.silhouette_key) return t('grandTour.editor.cover.generated')

  return t('grandTour.editor.cover.missing')
})

useUnsavedGuard(dirty)
useSaveShortcut(save, { enabled: () => can_save.value && !saving.value })

watchEffect(() => emit('update:saving', saving.value))
watchEffect(() => emit('update:can-save', can_save.value))
watchEffect(() => emit('update:dirty', dirty.value))

function onGenerated(cover: GeneratedCover): void {
  pending_cover.value = cover
}

function onGeneratedDiscarded(): void {
  pending_cover.value = null
}

function setCountry(selection: { code: string | null; country: string | null }): void {
  form.country_code = selection.code
  form.country = selection.country
}

function requestCancel(): void {
  if (!confirmDiscard(dirty.value)) return

  baseline.value = snapshot()
  emit('cancel')
}

function validate(): boolean {
  errors.make = form.make.trim() ? '' : t('grandTour.editor.form.makeRequired')
  errors.model = form.model.trim() ? '' : t('grandTour.editor.form.modelRequired')
  errors.year = form.year ? '' : t('grandTour.editor.form.yearRequired')

  return !errors.make && !errors.model && !errors.year
}

function text(value: string | null): string | null {
  const trimmed = value?.trim()

  return trimmed ? trimmed : null
}

function number_or_null(value: number | string | null): number | null {
  if (value === null || value === '') return null

  const parsed = Number(value)

  return Number.isFinite(parsed) ? parsed : null
}

async function uploadGeneratedCover(
  cover: GeneratedCover
): Promise<{ cutout_key: string; silhouette_key: string }> {
  const stamp = Date.now()
  const cutout_target = await requestVehicleAttachmentRaw(`vehicle-${stamp}-cutout.png`)
  const silhouette_target = await requestVehicleAttachmentRaw(`vehicle-${stamp}-silhouette.png`)

  await uploadVehicleCover(cutout_target, cover.cutout)
  await uploadVehicleCover(silhouette_target, cover.silhouette_light)

  return { cutout_key: cutout_target.key, silhouette_key: silhouette_target.key }
}

async function save(): Promise<void> {
  if (!can_save.value || saving.value) return
  if (!validate()) return

  saving.value = true

  try {
    let cover_keys: { cutout_key: string; silhouette_key: string } | null = null

    if (pending_cover.value) {
      try {
        cover_keys = await uploadGeneratedCover(pending_cover.value)
      } catch {
        toast.error(t('grandTour.editor.cover.uploadFailed'))
        return
      }
    }

    const payload: Vehicle = {
      id: form.id,
      make: form.make.trim(),
      model: form.model.trim(),
      year: Number(form.year),
      body_type: text(form.body_type),
      engine: text(form.engine),
      power: text(form.power),
      transmission: text(form.transmission),
      drivetrain: text(form.drivetrain),
      fuel: text(form.fuel),
      country_code: text(form.country_code)?.toUpperCase() ?? null,
      country: text(form.country),
      city: text(form.city),
      driven_at: form.driven_at ? new Date(`${form.driven_at}T00:00:00Z`).toISOString() : null,
      notes: text(form.notes),
      rating: number_or_null(form.rating),
      sort_order: number_or_null(form.sort_order),
      cutout_key: cover_keys?.cutout_key ?? form.cutout_key,
      silhouette_key: cover_keys?.silhouette_key ?? form.silhouette_key
    }

    const saved = await upsertVehicleRaw(payload)

    pending_cover.value = null
    Object.assign(form, fromVehicle(saved))
    baseline.value = snapshot()
    emit('saved', saved)
  } catch {
    toast.error(t('grandTour.editor.form.failedSave'))
  } finally {
    saving.value = false
  }
}

defineExpose({ requestCancel })
</script>
