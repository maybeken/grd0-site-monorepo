<template>
  <EditorSection
    :title="
      standalone
        ? isEdit
          ? $t('coffee.equipmentForm.titleEdit')
          : $t('coffee.equipmentForm.titleNew')
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
        <EditorField :label="$t('coffee.equipmentForm.name')" html-for="equipment-name" required>
          <EditorInput
            id="equipment-name"
            v-model="form.name"
            :placeholder="$t('coffee.equipmentForm.namePlaceholder')"
            required
          ></EditorInput>
        </EditorField>

        <EditorField :label="$t('coffee.equipmentForm.type')" html-for="equipment-type">
          <EditorSelect id="equipment-type" v-model="type_value">
            <option value="">{{ $t('common.selectType') }}</option>
            <option v-for="type in allTypes" :key="type.value" :value="type.value">
              {{ $t(type.label) }}
            </option>
          </EditorSelect>
        </EditorField>
      </div>

      <EditorField :label="$t('coffee.equipmentForm.description')" html-for="equipment-description">
        <EditorTextarea
          id="equipment-description"
          v-model="form.description"
          :placeholder="$t('coffee.equipmentForm.descriptionPlaceholder')"
        ></EditorTextarea>
      </EditorField>
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
import { useConfirmDiscard } from '@/composables/useConfirmDiscard'
import { useSaveShortcut } from '@/composables/useSaveShortcut'
import { useToast } from '@/composables/useToast'
import { useUnsavedGuard } from '@/composables/useUnsavedGuard'
import { EQUIPMENT_TYPES } from '@/helpers/coffee'
import { upsertCoffeeEquipmentRaw } from '@/services/coffee'

import type { BrewEquipment } from '@/interfaces/Coffee'

interface Props {
  equipment?: BrewEquipment | null
  standalone?: boolean
  formId?: string
}

const props = withDefaults(defineProps<Props>(), {
  equipment: null,
  standalone: true,
  formId: 'equipment-form'
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

const isEdit = computed(() => !!props.equipment?.id)
const saving = ref(false)
const allTypes = [...EQUIPMENT_TYPES]

function createEmptyForm(): BrewEquipment {
  return {
    name: '',
    type: null,
    description: null
  }
}

const form = reactive<BrewEquipment>(createEmptyForm())
const baseline = ref('')

const type_value = computed({
  get: () => form.type ?? '',
  set: (value: string) => {
    form.type = value || null
  }
})

function snapshot(): string {
  return JSON.stringify(form)
}

watch(
  () => props.equipment,
  (equipment) => {
    if (equipment) {
      Object.assign(form, createEmptyForm(), equipment)
    } else {
      Object.assign(form, createEmptyForm())
    }

    baseline.value = snapshot()
  },
  { immediate: true }
)

const dirty = computed(() => snapshot() !== baseline.value)
const can_save = computed(() => !!form.name.trim() && dirty.value)

useUnsavedGuard(dirty)
useSaveShortcut(save, { enabled: () => can_save.value && !saving.value })

watchEffect(() => emit('update:saving', saving.value))
watchEffect(() => emit('update:can-save', can_save.value))
watchEffect(() => emit('update:dirty', dirty.value))

function requestCancel(): void {
  if (!confirmDiscard(dirty.value)) return

  baseline.value = snapshot()
  emit('cancel')
}

async function save(): Promise<void> {
  if (!can_save.value || saving.value) return

  saving.value = true

  try {
    await upsertCoffeeEquipmentRaw(form)
    baseline.value = snapshot()
    emit('saved')
  } catch {
    toast.error(t('coffee.equipmentForm.failedSave'))
  } finally {
    saving.value = false
  }
}

defineExpose({ requestCancel })
</script>
