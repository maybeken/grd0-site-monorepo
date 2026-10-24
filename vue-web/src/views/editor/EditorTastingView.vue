<template>
  <EditorPage>
    <EditorPageHeader
      :title="tasting?.bean?.name || $t('editor.tastings')"
      :back-to="TASTINGS_ROUTE"
      :back-label="$t('editor.ui.backTo', { section: $t('editor.tastings') })"
    >
      <template #actions>
        <template v-if="tasting">
          <EditorButton variant="secondary" @click="cancel">
            {{ $t('common.cancel') }}
          </EditorButton>
          <EditorButton
            type="submit"
            :form="FORM_ID"
            variant="primary"
            :disabled="saving || !dirty"
            :loading="saving"
          >
            {{ $t('common.save') }}
          </EditorButton>
        </template>
      </template>
    </EditorPageHeader>

    <EditorLoading v-if="loading"></EditorLoading>

    <EditorError v-else-if="!tasting" :message="$t('editor.notFound')">
      <template #actions>
        <EditorButton variant="secondary" :to="TASTINGS_ROUTE">
          {{ $t('editor.ui.backTo', { section: $t('editor.tastings') }) }}
        </EditorButton>
      </template>
    </EditorError>

    <TastingForm
      v-else
      ref="tasting_form"
      :tasting="tasting"
      :beans="beans"
      :equipment-list="equipment"
      :standalone="false"
      :form-id="FORM_ID"
      @saved="onSaved"
      @cancel="back"
      @bean-added="fetchBeans"
      @equipment-added="fetchEquipment"
      @update:saving="saving = $event"
      @update:dirty="dirty = $event"
    ></TastingForm>
  </EditorPage>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

import TastingForm from '@/components/coffee/TastingForm.vue'
import EditorButton from '@/components/editor/ui/EditorButton.vue'
import EditorError from '@/components/editor/ui/EditorError.vue'
import EditorLoading from '@/components/editor/ui/EditorLoading.vue'
import EditorPage from '@/components/editor/ui/EditorPage.vue'
import EditorPageHeader from '@/components/editor/ui/EditorPageHeader.vue'
import { useToast } from '@/composables/useToast'
import {
  getCoffeeBeansRaw,
  getCoffeeEquipmentRaw,
  getCoffeeTastingRaw
} from '@/services/coffee'

import type { BrewEquipment, CoffeeBean, TastingNote } from '@/interfaces/Coffee'

const TASTINGS_ROUTE = '/editor/coffee/tastings'
const FORM_ID = 'tasting-form'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const toast = useToast()

const tasting = ref<TastingNote | null>(null)
const tasting_form = ref<InstanceType<typeof TastingForm> | null>(null)
const beans = ref<CoffeeBean[]>([])
const equipment = ref<BrewEquipment[]>([])
const loading = ref(true)
const saving = ref(false)
const dirty = ref(false)

const id = computed(() => {
  const value = route.params.id

  return typeof value === 'string' ? value : ''
})

async function fetchBeans(): Promise<void> {
  beans.value = await getCoffeeBeansRaw()
}

async function fetchEquipment(): Promise<void> {
  equipment.value = await getCoffeeEquipmentRaw()
}

async function load(): Promise<void> {
  loading.value = true
  tasting.value = null

  try {
    const [tasting_value] = await Promise.all([
      getCoffeeTastingRaw(id.value),
      fetchBeans(),
      fetchEquipment()
    ])

    tasting.value = tasting_value
  } catch {
    tasting.value = null
  } finally {
    loading.value = false
  }
}

function onSaved(): void {
  toast.success(t('common.done'))
}

function back(): void {
  router.push(TASTINGS_ROUTE)
}

function cancel(): void {
  tasting_form.value?.requestCancel()
}

watch(id, load, { immediate: true })
</script>
