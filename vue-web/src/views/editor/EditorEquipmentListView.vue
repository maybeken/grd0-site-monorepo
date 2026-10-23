<template>
  <div class="flex flex-col gap-6">
    <header class="flex flex-col gap-2">
      <h1 class="text-2xl font-bold">{{ $t('editor.equipment') }}<CursorBlink /></h1>
    </header>

    <div v-if="!show_form" class="flex justify-end">
      <button
        class="rounded bg-accent px-3 py-1 text-sm font-semibold text-background hover:brightness-110"
        @click="show_form = true"
      >
        {{ $t('coffee.editor.newEquipment') }}
      </button>
    </div>

    <EquipmentForm
      v-if="show_form"
      :equipment="null"
      @saved="onSaved"
      @cancel="show_form = false"
    ></EquipmentForm>

    <div v-if="!show_form" class="flex flex-col gap-2">
      <EditorEntityRow
        v-for="item in equipment"
        :key="item.id"
        :title="item.name"
        :to="`/editor/coffee/equipment/${item.id}`"
        @delete="remove(item)"
      >
        <span v-if="item.type" class="truncate text-xs text-foreground/60">{{ item.type }}</span>
      </EditorEntityRow>

      <p v-if="!equipment.length" class="py-4 text-center text-sm text-foreground/40">
        {{ $t('coffee.editor.noEquipment') }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import EquipmentForm from '@/components/coffee/EquipmentForm.vue'
import EditorEntityRow from '@/components/editor/EditorEntityRow.vue'
import { deleteCoffeeEquipmentRaw, getCoffeeEquipmentRaw } from '@/services/coffee'
import { useToast } from '@/composables/useToast'

import type { BrewEquipment } from '@/interfaces/Coffee'

const { t } = useI18n()
const toast = useToast()

const equipment = ref<BrewEquipment[]>([])
const show_form = ref(false)

async function fetchEquipment(): Promise<void> {
  try {
    equipment.value = await getCoffeeEquipmentRaw()
  } catch {
    toast.error(t('editor.loadFailed'))
  }
}

async function onSaved(): Promise<void> {
  show_form.value = false
  await fetchEquipment()
}

async function remove(item: BrewEquipment): Promise<void> {
  if (!item.id || !confirm(t('coffee.editor.confirmDeleteEquipment', { name: item.name }))) return

  try {
    await deleteCoffeeEquipmentRaw(item.id)
    await fetchEquipment()
  } catch {
    toast.error(t('coffee.editor.failedDeleteEquipment'))
  }
}

onMounted(fetchEquipment)
</script>
