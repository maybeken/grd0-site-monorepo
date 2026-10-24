<template>
  <EditorPage>
    <EditorPageHeader :title="$t('editor.equipment')">
      <template #actions>
        <EditorButton v-if="!show_form" variant="primary" @click="show_form = true">
          <Icon icon="mynaui:plus" height="1rem"></Icon>
          {{ $t('coffee.editor.newEquipment') }}
        </EditorButton>
      </template>
    </EditorPageHeader>

    <EquipmentForm
      v-if="show_form"
      :equipment="null"
      @saved="onSaved"
      @cancel="show_form = false"
    ></EquipmentForm>

    <EditorError v-else-if="load_failed" :message="t('editor.loadFailed')">
      <template #actions>
        <EditorButton variant="secondary" @click="fetchEquipment">
          {{ $t('common.retry') }}
        </EditorButton>
      </template>
    </EditorError>

    <EditorLoading v-else-if="loading"></EditorLoading>

    <div v-else class="flex flex-col gap-2">
      <EditorEntityRow
        v-for="item in equipment"
        :key="item.id"
        :title="item.name"
        :to="`/editor/coffee/equipment/${item.id}`"
        @delete="remove(item)"
      >
        <span v-if="item.type" class="truncate text-sm font-medium text-foreground/60">
          {{ item.type }}
        </span>
      </EditorEntityRow>

      <EditorEmptyState
        v-if="!equipment.length"
        :message="$t('coffee.editor.noEquipment')"
      ></EditorEmptyState>
    </div>
  </EditorPage>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import EquipmentForm from '@/components/coffee/EquipmentForm.vue'
import EditorEntityRow from '@/components/editor/EditorEntityRow.vue'
import EditorButton from '@/components/editor/ui/EditorButton.vue'
import EditorEmptyState from '@/components/editor/ui/EditorEmptyState.vue'
import EditorError from '@/components/editor/ui/EditorError.vue'
import EditorLoading from '@/components/editor/ui/EditorLoading.vue'
import EditorPage from '@/components/editor/ui/EditorPage.vue'
import EditorPageHeader from '@/components/editor/ui/EditorPageHeader.vue'
import Icon from '@/components/shared/Icon.vue'
import { deleteCoffeeEquipmentRaw, getCoffeeEquipmentRaw } from '@/services/coffee'
import { useToast } from '@/composables/useToast'

import type { BrewEquipment } from '@/interfaces/Coffee'

const { t } = useI18n()
const toast = useToast()

const equipment = ref<BrewEquipment[]>([])
const show_form = ref(false)
const loading = ref(true)
const load_failed = ref(false)

async function fetchEquipment(): Promise<void> {
  loading.value = true
  load_failed.value = false

  try {
    equipment.value = await getCoffeeEquipmentRaw()
  } catch {
    equipment.value = []
    load_failed.value = true
  } finally {
    loading.value = false
  }
}

async function onSaved(): Promise<void> {
  show_form.value = false
  toast.success(t('common.done'))
  await fetchEquipment()
}

async function remove(item: BrewEquipment): Promise<void> {
  if (!item.id || !confirm(t('common.confirmDeleteItem', { name: item.name }))) return

  try {
    await deleteCoffeeEquipmentRaw(item.id)
    toast.success(t('common.deleted'))
    await fetchEquipment()
  } catch {
    toast.error(t('coffee.editor.failedDeleteEquipment'))
  }
}

onMounted(fetchEquipment)
</script>
