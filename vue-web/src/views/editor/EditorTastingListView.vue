<template>
  <EditorPage>
    <EditorPageHeader :title="$t('editor.tastings')">
      <template #actions>
        <EditorButton v-if="!show_form" variant="primary" @click="show_form = true">
          <Icon icon="mynaui:plus" height="1rem"></Icon>
          {{ $t('coffee.editor.newTasting') }}
        </EditorButton>
      </template>
    </EditorPageHeader>

    <TastingForm
      v-if="show_form"
      :tasting="null"
      :beans="beans"
      :equipment-list="equipment"
      @saved="onSaved"
      @cancel="show_form = false"
      @bean-added="fetchBeans"
      @equipment-added="fetchEquipment"
    ></TastingForm>

    <EditorError v-else-if="load_failed" :message="t('editor.loadFailed')">
      <template #actions>
        <EditorButton variant="secondary" @click="fetchAll">
          {{ $t('common.retry') }}
        </EditorButton>
      </template>
    </EditorError>

    <EditorLoading v-else-if="loading"></EditorLoading>

    <div v-else class="flex flex-col gap-2">
      <EditorEntityRow
        v-for="tasting in tastings"
        :key="tasting.id"
        :title="tasting.bean?.name || $t('coffee.editor.unknownBean')"
        :to="`/editor/coffee/tastings/${tasting.id}`"
        @delete="remove(tasting)"
      >
        <span v-if="metaLabel(tasting)" class="truncate text-sm font-medium text-foreground/60">
          {{ metaLabel(tasting) }}
        </span>
      </EditorEntityRow>

      <EditorEmptyState
        v-if="!tastings.length"
        :message="$t('coffee.editor.noTastingsYet')"
      ></EditorEmptyState>
    </div>
  </EditorPage>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import TastingForm from '@/components/coffee/TastingForm.vue'
import EditorEntityRow from '@/components/editor/EditorEntityRow.vue'
import EditorButton from '@/components/editor/ui/EditorButton.vue'
import EditorEmptyState from '@/components/editor/ui/EditorEmptyState.vue'
import EditorError from '@/components/editor/ui/EditorError.vue'
import EditorLoading from '@/components/editor/ui/EditorLoading.vue'
import EditorPage from '@/components/editor/ui/EditorPage.vue'
import EditorPageHeader from '@/components/editor/ui/EditorPageHeader.vue'
import Icon from '@/components/shared/Icon.vue'
import {
  deleteCoffeeTastingRaw,
  getCoffeeBeansRaw,
  getCoffeeEquipmentRaw,
  getCoffeeTastingsRaw
} from '@/services/coffee'
import { useToast } from '@/composables/useToast'

import type { BrewEquipment, CoffeeBean, TastingNote } from '@/interfaces/Coffee'

const { t } = useI18n()
const toast = useToast()

const beans = ref<CoffeeBean[]>([])
const equipment = ref<BrewEquipment[]>([])
const tastings = ref<TastingNote[]>([])
const show_form = ref(false)
const loading = ref(true)
const load_failed = ref(false)

function metaLabel(tasting: TastingNote): string {
  const parts = [tasting.equipment?.name, tasting.tasted_at?.split('T')[0]].filter(Boolean)
  let label = parts.join(t('coffee.editor.dotSeparator'))

  if (tasting.rating != null) {
    label += ` ${t('coffee.editor.rating', { rating: tasting.rating })}`
  }

  return label
}

async function fetchBeans(): Promise<void> {
  beans.value = await getCoffeeBeansRaw()
}

async function fetchEquipment(): Promise<void> {
  equipment.value = await getCoffeeEquipmentRaw()
}

async function fetchTastings(): Promise<void> {
  tastings.value = await getCoffeeTastingsRaw()
}

async function fetchAll(): Promise<void> {
  loading.value = true
  load_failed.value = false

  try {
    await Promise.all([fetchBeans(), fetchEquipment(), fetchTastings()])
  } catch {
    tastings.value = []
    load_failed.value = true
  } finally {
    loading.value = false
  }
}

async function onSaved(): Promise<void> {
  show_form.value = false
  toast.success(t('common.done'))
  await fetchAll()
}

async function remove(tasting: TastingNote): Promise<void> {
  const label =
    tasting.bean?.name || tasting.tasted_at?.split('T')[0] || t('coffee.editor.unknownBean')

  if (!tasting.id || !confirm(t('coffee.editor.confirmDeleteTasting', { label }))) return

  try {
    await deleteCoffeeTastingRaw(tasting.id)
    toast.success(t('common.deleted'))
    await fetchTastings()
  } catch {
    toast.error(t('coffee.editor.failedDeleteTasting'))
  }
}

onMounted(fetchAll)
</script>
