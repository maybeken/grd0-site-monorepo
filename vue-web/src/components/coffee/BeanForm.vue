<template>
  <EditorSection
    :title="
      standalone
        ? isEdit
          ? $t('coffee.beanForm.titleEdit')
          : $t('coffee.beanForm.titleNew')
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
        <EditorField :label="$t('coffee.beanForm.name')" html-for="bean-name" required>
          <EditorInput
            id="bean-name"
            v-model="form.name"
            :placeholder="$t('coffee.beanForm.namePlaceholder')"
            required
          ></EditorInput>
        </EditorField>

        <EditorField :label="$t('coffee.beanForm.origin')" html-for="bean-origin">
          <EditorInput
            id="bean-origin"
            v-model="form.origin"
            :placeholder="$t('coffee.beanForm.originPlaceholder')"
          ></EditorInput>
        </EditorField>

        <EditorField :label="$t('coffee.beanForm.roaster')" html-for="bean-roaster">
          <EditorInput
            id="bean-roaster"
            v-model="form.roaster"
            :placeholder="$t('coffee.beanForm.roasterPlaceholder')"
          ></EditorInput>
        </EditorField>

        <EditorField :label="$t('coffee.beanForm.roastDate')" html-for="bean-roast-date">
          <EditorInput id="bean-roast-date" v-model="form.roast_date" type="date"></EditorInput>
        </EditorField>

        <EditorField :label="$t('coffee.beanForm.variety')" html-for="bean-variety">
          <EditorInput
            id="bean-variety"
            v-model="form.variety"
            :placeholder="$t('coffee.beanForm.varietyPlaceholder')"
          ></EditorInput>
        </EditorField>

        <EditorField :label="$t('coffee.beanForm.process')" html-for="bean-process">
          <EditorInput
            id="bean-process"
            v-model="form.process"
            :placeholder="$t('coffee.beanForm.processPlaceholder')"
          ></EditorInput>
        </EditorField>

        <EditorField :label="$t('coffee.beanForm.altitude')" html-for="bean-altitude">
          <EditorInput
            id="bean-altitude"
            v-model="form.altitude"
            :placeholder="$t('coffee.beanForm.altitudePlaceholder')"
          ></EditorInput>
        </EditorField>
      </div>

      <EditorField :label="$t('coffee.beanForm.description')" html-for="bean-description">
        <EditorTextarea
          id="bean-description"
          v-model="form.description"
          :placeholder="$t('coffee.beanForm.descriptionPlaceholder')"
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
import EditorTextarea from '@/components/editor/ui/EditorTextarea.vue'
import { useConfirmDiscard } from '@/composables/useConfirmDiscard'
import { useSaveShortcut } from '@/composables/useSaveShortcut'
import { useToast } from '@/composables/useToast'
import { useUnsavedGuard } from '@/composables/useUnsavedGuard'
import { upsertCoffeeBeanRaw } from '@/services/coffee'

import type { CoffeeBean } from '@/interfaces/Coffee'

interface Props {
  bean?: CoffeeBean | null
  standalone?: boolean
  formId?: string
}

const props = withDefaults(defineProps<Props>(), {
  bean: null,
  standalone: true,
  formId: 'bean-form'
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

const isEdit = computed(() => !!props.bean?.id)
const saving = ref(false)

function createEmptyForm(): CoffeeBean {
  return {
    name: '',
    origin: null,
    roaster: null,
    roast_date: null,
    variety: null,
    process: null,
    altitude: null,
    description: null
  }
}

const form = reactive<CoffeeBean>(createEmptyForm())
const baseline = ref('')

function snapshot(): string {
  return JSON.stringify(form)
}

watch(
  () => props.bean,
  (bean) => {
    if (bean) {
      Object.assign(form, createEmptyForm(), bean)
      if (bean.roast_date) {
        form.roast_date = bean.roast_date.split('T')[0] ?? null
      }
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
    await upsertCoffeeBeanRaw({
      ...form,
      roast_date: form.roast_date ? new Date(form.roast_date).toISOString() : null
    })
    baseline.value = snapshot()
    emit('saved')
  } catch {
    toast.error(t('coffee.beanForm.failedSave'))
  } finally {
    saving.value = false
  }
}

defineExpose({ requestCancel })
</script>
