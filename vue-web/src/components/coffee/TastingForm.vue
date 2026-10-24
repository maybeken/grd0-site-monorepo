<template>
  <EditorSection
    :title="
      standalone
        ? isEdit
          ? $t('coffee.tastingForm.titleEdit')
          : $t('coffee.tastingForm.titleNew')
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
        :disabled="saving || !dirty"
        :loading="saving"
      >
        {{ $t('common.save') }}
      </EditorButton>
    </template>

    <form :id="formId" class="flex flex-col gap-5" @submit.prevent="save">
      <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
        <EditorField :label="$t('coffee.tastingForm.bean')" html-for="tasting-bean">
          <template #label-action>
            <EditorButton
              v-if="!showQuickAddBean"
              size="sm"
              variant="secondary"
              @click="showQuickAddBean = true"
            >
              {{ $t('coffee.tastingForm.addNew') }}
            </EditorButton>
          </template>
          <EditorSelect id="tasting-bean" v-model="form.bean.id">
            <option value="">{{ $t('coffee.tastingForm.selectBean') }}</option>
            <option v-for="bean in beans" :key="bean.id" :value="bean.id">
              {{ bean.roaster }} {{ bean.name }} ({{ bean.roast_date?.split('T')[0]?.replace(/-/g, '/') }})
            </option>
          </EditorSelect>
        </EditorField>

        <EditorField :label="$t('coffee.tastingForm.equipment')" html-for="tasting-equipment">
          <template #label-action>
            <EditorButton
              v-if="!showQuickAddEquipment"
              size="sm"
              variant="secondary"
              @click="showQuickAddEquipment = true"
            >
              {{ $t('coffee.tastingForm.addNew') }}
            </EditorButton>
          </template>
          <EditorSelect id="tasting-equipment" v-model="form.equipment.id">
            <option value="">{{ $t('coffee.tastingForm.selectEquipment') }}</option>
            <option
              v-for="item in (equipmentList || []).filter((entry) => entry.type !== 'Grinder')"
              :key="item.id"
              :value="item.id"
            >
              {{ item.name }}
            </option>
          </EditorSelect>
        </EditorField>
      </div>

      <div
        v-if="showQuickAddBean"
        class="flex flex-col gap-3 rounded-lg border border-accent/20 bg-accent/5 p-3"
      >
        <EditorField :label="$t('coffee.tastingForm.beanName')" html-for="quick-bean-name" required>
          <EditorInput id="quick-bean-name" v-model="quickBean.name" required></EditorInput>
        </EditorField>
        <div class="grid grid-cols-1 gap-3 md:grid-cols-2">
          <EditorField :label="$t('coffee.tastingForm.roaster')" html-for="quick-bean-roaster">
            <EditorInput id="quick-bean-roaster" v-model="quickBean.roaster"></EditorInput>
          </EditorField>
          <EditorField :label="$t('coffee.tastingForm.origin')" html-for="quick-bean-origin">
            <EditorInput id="quick-bean-origin" v-model="quickBean.origin"></EditorInput>
          </EditorField>
        </div>
        <div class="flex justify-end gap-2">
          <EditorButton size="sm" variant="secondary" @click="cancelQuickBean">
            {{ $t('common.cancel') }}
          </EditorButton>
          <EditorButton
            size="sm"
            variant="primary"
            :disabled="quickBeanSaving || !quickBean.name.trim()"
            :loading="quickBeanSaving"
            @click="saveQuickBean"
          >
            {{ $t('common.add') }}
          </EditorButton>
        </div>
      </div>

      <div
        v-if="showQuickAddEquipment"
        class="flex flex-col gap-3 rounded-lg border border-accent/20 bg-accent/5 p-3"
      >
        <EditorField
          :label="$t('coffee.tastingForm.equipmentName')"
          html-for="quick-equipment-name"
          required
        >
          <EditorInput id="quick-equipment-name" v-model="quickEquipment.name" required></EditorInput>
        </EditorField>
        <EditorField :label="$t('coffee.equipmentForm.type')" html-for="quick-equipment-type">
          <EditorSelect id="quick-equipment-type" v-model="quick_type_value">
            <option value="">{{ $t('common.selectType') }}</option>
            <option v-for="type in equipmentTypes" :key="type.value" :value="type.value">
              {{ $t(type.label) }}
            </option>
          </EditorSelect>
        </EditorField>
        <div class="flex justify-end gap-2">
          <EditorButton size="sm" variant="secondary" @click="cancelQuickEquipment">
            {{ $t('common.cancel') }}
          </EditorButton>
          <EditorButton
            size="sm"
            variant="primary"
            :disabled="quickEquipmentSaving || !quickEquipment.name.trim()"
            :loading="quickEquipmentSaving"
            @click="saveQuickEquipment"
          >
            {{ $t('common.add') }}
          </EditorButton>
        </div>
      </div>

      <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
        <EditorField :label="$t('coffee.tastingForm.date')" html-for="tasting-date">
          <EditorInput id="tasting-date" v-model="form.tasted_at" type="date"></EditorInput>
        </EditorField>

        <EditorField :label="$t('coffee.tastingForm.rating')" html-for="tasting-rating">
          <EditorInput
            id="tasting-rating"
            v-model="form.rating"
            type="number"
            :number="true"
            :disabled="idkFields.rating"
          ></EditorInput>
        </EditorField>
      </div>

      <fieldset class="rounded border border-accent/20 p-3">
        <legend class="px-1 text-lg font-medium">{{ $t('coffee.tastingForm.brewRecipe') }}</legend>

        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3">
          <EditorField :label="$t('coffee.tastingForm.grinder')" html-for="tasting-grinder">
            <template #label-action>
              <IdkToggle v-model="idkFields.grinder"></IdkToggle>
            </template>
            <EditorSelect
              id="tasting-grinder"
              v-model="form.grinder!.id"
              :disabled="idkFields.grinder"
            >
              <option value="">{{ $t('coffee.tastingForm.selectEquipment') }}</option>
              <option
                v-for="item in (equipmentList || []).filter((entry) => entry.type === 'Grinder')"
                :key="item.id"
                :value="item.id"
              >
                {{ item.name }}
              </option>
            </EditorSelect>
          </EditorField>

          <EditorField :label="$t('coffee.tastingForm.grindSize')" html-for="tasting-grind-size">
            <template #label-action>
              <IdkToggle v-model="idkFields.grind_size"></IdkToggle>
            </template>
            <EditorSelect
              id="tasting-grind-size"
              v-model="form.grind_size"
              :disabled="idkFields.grind_size"
            >
              <option value="">{{ $t('coffee.tastingForm.selectGrind') }}</option>
              <option v-for="g in GRIND_SIZES" :key="g.value" :value="g.value">
                {{ $t(g.label) }}
              </option>
            </EditorSelect>
          </EditorField>

          <EditorField :label="$t('coffee.tastingForm.grindSetting')" html-for="tasting-grind-setting">
            <template #label-action>
              <IdkToggle v-model="idkFields.grind_setting"></IdkToggle>
            </template>
            <EditorInput
              id="tasting-grind-setting"
              v-model="form.grind_setting"
              type="number"
              :number="true"
              :disabled="idkFields.grind_setting"
            ></EditorInput>
          </EditorField>

          <EditorField :label="$t('coffee.tastingForm.coffeeDose')" html-for="tasting-coffee-dose">
            <template #label-action>
              <IdkToggle v-model="idkFields.coffee_dose"></IdkToggle>
            </template>
            <EditorInput
              id="tasting-coffee-dose"
              v-model="form.coffee_dose"
              type="number"
              :number="true"
              :disabled="idkFields.coffee_dose"
            ></EditorInput>
          </EditorField>

          <EditorField :label="$t('coffee.tastingForm.waterIn')" html-for="tasting-water-in">
            <template #label-action>
              <IdkToggle v-model="idkFields.water_in"></IdkToggle>
            </template>
            <EditorInput
              id="tasting-water-in"
              v-model="form.water_in"
              type="number"
              :number="true"
              :disabled="idkFields.water_in"
            ></EditorInput>
          </EditorField>

          <EditorField :label="$t('coffee.tastingForm.coffeeOut')" html-for="tasting-coffee-out">
            <template #label-action>
              <IdkToggle v-model="idkFields.coffee_out"></IdkToggle>
            </template>
            <EditorInput
              id="tasting-coffee-out"
              v-model="form.coffee_out"
              type="number"
              :number="true"
              :disabled="idkFields.coffee_out"
            ></EditorInput>
          </EditorField>

          <EditorField :label="$t('coffee.tastingForm.ratio')">
            <p class="flex min-h-9 items-center text-lg font-medium">
              {{ computedRatio != null ? '1:' + computedRatio.toFixed(1) : $t('coffee.tastingForm.idk') }}
            </p>
          </EditorField>

          <EditorField :label="$t('coffee.tastingForm.brewTime')" html-for="tasting-brew-time">
            <template #label-action>
              <IdkToggle v-model="idkFields.brew_time"></IdkToggle>
            </template>
            <BrewTimeInput
              id="tasting-brew-time"
              v-model="form.brew_time"
              :disabled="idkFields.brew_time"
            ></BrewTimeInput>
          </EditorField>

          <EditorField :label="$t('coffee.tastingForm.waterTemp')" html-for="tasting-water-temp">
            <template #label-action>
              <IdkToggle v-model="idkFields.water_temperature"></IdkToggle>
            </template>
            <EditorInput
              id="tasting-water-temp"
              v-model="form.water_temperature"
              type="number"
              :number="true"
              :disabled="idkFields.water_temperature"
            ></EditorInput>
          </EditorField>
        </div>
      </fieldset>

      <fieldset class="rounded border border-accent/20 p-3">
        <legend class="px-1 text-lg font-medium">{{ $t('coffee.tastingForm.tasteProfile') }}</legend>
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-4">
          <EditorField
            v-for="dim in tasteDimensions"
            :key="dim.key"
            :label="$t(dim.label)"
            :html-for="`tasting-${dim.key}`"
          >
            <EditorInput
              :id="`tasting-${dim.key}`"
              :model-value="tasteValue(dim.key)"
              type="number"
              :number="true"
              :disabled="!!idkFields[dim.key as string]"
              @update:model-value="setTasteValue(dim.key, $event)"
            ></EditorInput>
          </EditorField>
        </div>
      </fieldset>

      <EditorField :label="$t('coffee.tastingForm.notes')" html-for="tasting-notes">
        <EditorTextarea id="tasting-notes" v-model="form.overall_notes"></EditorTextarea>
      </EditorField>

      <label class="flex items-center gap-2 text-lg font-medium">
        <input type="checkbox" v-model="form.pinned" class="h-4 w-4 accent-accent" />
        <span>{{ $t('coffee.tastingForm.pinTasting') }}</span>
      </label>
    </form>
  </EditorSection>
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
import BrewTimeInput from './BrewTimeInput.vue'
import IdkToggle from './IdkToggle.vue'
import { useConfirmDiscard } from '@/composables/useConfirmDiscard'
import { useSaveShortcut } from '@/composables/useSaveShortcut'
import { useToast } from '@/composables/useToast'
import { useUnsavedGuard } from '@/composables/useUnsavedGuard'
import { EQUIPMENT_TYPES, GRIND_SIZES, TASTE_DIMENSIONS } from '@/helpers/coffee'
import {
  upsertCoffeeBeanRaw,
  upsertCoffeeEquipmentRaw,
  upsertCoffeeTastingRaw
} from '@/services/coffee'

import type { BrewEquipment, CoffeeBean, TastingNote } from '@/interfaces/Coffee'

interface Props {
  tasting?: TastingNote | null
  beans: CoffeeBean[]
  equipmentList: BrewEquipment[]
  standalone?: boolean
  formId?: string
}

const props = withDefaults(defineProps<Props>(), {
  tasting: null,
  standalone: true,
  formId: 'tasting-form'
})

const emit = defineEmits<{
  saved: []
  cancel: []
  beanAdded: [bean: CoffeeBean]
  equipmentAdded: [equipment: BrewEquipment]
  'update:saving': [saving: boolean]
  'update:can-save': [can_save: boolean]
  'update:dirty': [dirty: boolean]
}>()

const { t } = useI18n()
const toast = useToast()
const confirmDiscard = useConfirmDiscard()

const isEdit = computed(() => !!props.tasting?.id)
const saving = ref(false)

const showQuickAddBean = ref(false)
const quickBeanSaving = ref(false)
const quickBean = reactive({
  name: '',
  roaster: null as string | null,
  origin: null as string | null
})

const showQuickAddEquipment = ref(false)
const quickEquipmentSaving = ref(false)
const quickEquipment = reactive({
  name: '',
  type: null as string | null
})

const equipmentTypes = EQUIPMENT_TYPES
const tasteDimensions = TASTE_DIMENSIONS.map((dimension) => ({
  key: `taste_${dimension.value}` as keyof TastingNote,
  label: dimension.label
}))

const quick_type_value = computed({
  get: () => quickEquipment.type ?? '',
  set: (value: string) => {
    quickEquipment.type = value || null
  }
})

function createEmptyForm(): TastingNote {
  return {
    bean: { name: '' },
    equipment: { name: '' },
    tasted_at: new Date().toISOString().split('T')[0]!,
    pinned: false,
    grinder: { name: '' },
    grind_size: null,
    grind_setting: null,
    coffee_dose: null,
    water_in: null,
    coffee_out: null,
    ratio: null,
    brew_time: null,
    water_temperature: null,
    overall_notes: null,
    rating: null,
    ...Object.fromEntries(TASTE_DIMENSIONS.map((dimension) => [`taste_${dimension.value}`, null]))
  }
}

const form = reactive<TastingNote>(createEmptyForm())
const baseline = ref('')

function snapshot(): string {
  return JSON.stringify(form)
}

const idkFields = reactive<Record<string, boolean>>({
  grinder: false,
  grind_size: false,
  grind_setting: false,
  coffee_dose: false,
  water_in: false,
  coffee_out: false,
  brew_time: false,
  water_temperature: false,
  rating: false
})

const computedRatio = computed(() => {
  if (form.water_in != null && form.coffee_dose != null && form.coffee_dose > 0) {
    return form.water_in / form.coffee_dose
  }

  return null
})

watch(computedRatio, (value) => {
  form.ratio = value
})

function tasteValue(key: keyof TastingNote): number | null {
  const value = form[key]

  return typeof value === 'number' ? value : null
}

function setTasteValue(key: keyof TastingNote, value: string | number | null): void {
  ;(form as Record<string, unknown>)[key] = value
}

watch(
  () => props.tasting,
  (tasting) => {
    if (tasting) {
      Object.assign(form, createEmptyForm(), tasting)
      form.bean = tasting.bean ?? { name: '' }
      form.equipment = tasting.equipment ?? { name: '' }
      form.grinder = tasting.grinder ?? { name: '' }
      form.tasted_at = tasting.tasted_at || new Date().toISOString().substring(0, 10)

      for (const key of Object.keys(idkFields)) {
        if (key === 'grinder') {
          idkFields[key] = !tasting.grinder?.id
        } else {
          idkFields[key] = (tasting as unknown as Record<string, unknown>)[key] == null
        }
      }
    } else {
      Object.assign(form, createEmptyForm())

      for (const key of Object.keys(idkFields)) {
        idkFields[key] = false
      }
    }

    form.ratio = computedRatio.value
    baseline.value = snapshot()
  },
  { immediate: true }
)

const dirty = computed(() => snapshot() !== baseline.value)

useUnsavedGuard(dirty)
useSaveShortcut(save, { enabled: () => dirty.value && !saving.value })

watchEffect(() => emit('update:saving', saving.value))
watchEffect(() => emit('update:can-save', true))
watchEffect(() => emit('update:dirty', dirty.value))

for (const key of Object.keys(idkFields)) {
  watch(
    () => idkFields[key],
    (idk) => {
      if (!idk) return

      if (key === 'grinder') {
        form.grinder = { name: '' }
        return
      }

      ;(form as Record<string, unknown>)[key] = null
    }
  )
}

function requestCancel(): void {
  if (!confirmDiscard(dirty.value)) return

  baseline.value = snapshot()
  emit('cancel')
}

function cancelQuickBean(): void {
  showQuickAddBean.value = false
  quickBean.name = ''
  quickBean.roaster = null
  quickBean.origin = null
}

async function saveQuickBean(): Promise<void> {
  if (!quickBean.name.trim()) return

  quickBeanSaving.value = true

  try {
    const newBean = await upsertCoffeeBeanRaw({
      name: quickBean.name.trim(),
      roaster: quickBean.roaster?.trim() || null,
      origin: quickBean.origin?.trim() || null
    })
    emit('beanAdded', newBean)
    form.bean = newBean
    cancelQuickBean()
  } catch {
    toast.error(t('coffee.tastingForm.failedAddBean'))
  } finally {
    quickBeanSaving.value = false
  }
}

function cancelQuickEquipment(): void {
  showQuickAddEquipment.value = false
  quickEquipment.name = ''
  quickEquipment.type = null
}

async function saveQuickEquipment(): Promise<void> {
  if (!quickEquipment.name.trim()) return

  quickEquipmentSaving.value = true

  try {
    const newEquipment = await upsertCoffeeEquipmentRaw({
      name: quickEquipment.name.trim(),
      type: quickEquipment.type || null
    })
    emit('equipmentAdded', newEquipment)
    form.equipment = newEquipment
    cancelQuickEquipment()
  } catch {
    toast.error(t('coffee.tastingForm.failedAddEquipment'))
  } finally {
    quickEquipmentSaving.value = false
  }
}

async function save(): Promise<void> {
  if (saving.value) return

  saving.value = true

  try {
    await upsertCoffeeTastingRaw({
      ...form,
      tasted_at: new Date(form.tasted_at).toISOString()
    })
    baseline.value = snapshot()
    emit('saved')
  } catch {
    toast.error(t('coffee.tastingForm.failedSaveTasting'))
  } finally {
    saving.value = false
  }
}

defineExpose({ requestCancel })
</script>
