<template>
  <EditorPage>
    <EditorPageHeader
      :title="equipment?.name || $t('editor.equipment')"
      :back-to="EQUIPMENT_ROUTE"
      :back-label="$t('editor.ui.backTo', { section: $t('editor.equipment') })"
    >
      <template #actions>
        <template v-if="equipment">
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

    <EditorError v-else-if="!equipment" :message="$t('editor.notFound')">
      <template #actions>
        <EditorButton variant="secondary" :to="EQUIPMENT_ROUTE">
          {{ $t('editor.ui.backTo', { section: $t('editor.equipment') }) }}
        </EditorButton>
      </template>
    </EditorError>

    <EquipmentForm
      v-else
      ref="equipment_form"
      :equipment="equipment"
      :standalone="false"
      :form-id="FORM_ID"
      @saved="onSaved"
      @cancel="back"
      @update:saving="saving = $event"
      @update:can-save="can_save = $event"
    ></EquipmentForm>
  </EditorPage>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

import EquipmentForm from '@/components/coffee/EquipmentForm.vue'
import EditorButton from '@/components/editor/ui/EditorButton.vue'
import EditorError from '@/components/editor/ui/EditorError.vue'
import EditorLoading from '@/components/editor/ui/EditorLoading.vue'
import EditorPage from '@/components/editor/ui/EditorPage.vue'
import EditorPageHeader from '@/components/editor/ui/EditorPageHeader.vue'
import { useToast } from '@/composables/useToast'
import { getCoffeeEquipmentItemRaw } from '@/services/coffee'

import type { BrewEquipment } from '@/interfaces/Coffee'

const EQUIPMENT_ROUTE = '/editor/coffee/equipment'
const FORM_ID = 'equipment-form'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const toast = useToast()

const equipment = ref<BrewEquipment | null>(null)
const equipment_form = ref<InstanceType<typeof EquipmentForm> | null>(null)
const loading = ref(true)
const saving = ref(false)
const can_save = ref(false)

const id = computed(() => {
  const value = route.params.id

  return typeof value === 'string' ? value : ''
})

async function load(): Promise<void> {
  loading.value = true
  equipment.value = null

  try {
    equipment.value = await getCoffeeEquipmentItemRaw(id.value)
  } catch {
    equipment.value = null
  } finally {
    loading.value = false
  }
}

function onSaved(): void {
  toast.success(t('common.done'))
}

function back(): void {
  router.push(EQUIPMENT_ROUTE)
}

function cancel(): void {
  equipment_form.value?.requestCancel()
}

watch(id, load, { immediate: true })
</script>
