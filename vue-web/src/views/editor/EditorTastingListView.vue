<template>
  <div class="flex flex-col gap-6">
    <header class="flex flex-col gap-2">
      <h1 class="text-2xl font-bold">{{ $t('editor.tastings') }}<CursorBlink /></h1>
    </header>

    <div v-if="!show_form" class="flex justify-end">
      <button
        class="rounded bg-accent px-3 py-1 text-sm font-semibold text-background hover:brightness-110"
        @click="show_form = true"
      >
        {{ $t('coffee.editor.newTasting') }}
      </button>
    </div>

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

    <div v-if="!show_form" class="flex flex-col gap-2">
      <EditorEntityRow
        v-for="tasting in tastings"
        :key="tasting.id"
        :title="tasting.bean?.name || $t('coffee.editor.unknownBean')"
        :to="`/editor/coffee/tastings/${tasting.id}`"
        @delete="remove(tasting)"
      >
        <span class="truncate text-xs text-foreground/60">
          {{ [tasting.equipment?.name, tasting.tasted_at?.split('T')[0]].filter(Boolean).join(' · ') }}
          <template v-if="tasting.rating != null"> · ★ {{ tasting.rating }}</template>
        </span>
      </EditorEntityRow>

      <p v-if="!tastings.length" class="py-4 text-center text-sm text-foreground/40">
        {{ $t('coffee.editor.noTastingsYet') }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import TastingForm from '@/components/coffee/TastingForm.vue'
import EditorEntityRow from '@/components/editor/EditorEntityRow.vue'
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
  try {
    await Promise.all([fetchBeans(), fetchEquipment(), fetchTastings()])
  } catch {
    toast.error(t('editor.loadFailed'))
  }
}

async function onSaved(): Promise<void> {
  show_form.value = false
  await fetchAll()
}

async function remove(tasting: TastingNote): Promise<void> {
  const label =
    tasting.bean?.name || tasting.tasted_at?.split('T')[0] || t('coffee.editor.unknownBean')

  if (!tasting.id || !confirm(t('coffee.editor.confirmDeleteTasting', { label }))) return

  try {
    await deleteCoffeeTastingRaw(tasting.id)
    await fetchTastings()
  } catch {
    toast.error(t('coffee.editor.failedDeleteTasting'))
  }
}

onMounted(fetchAll)
</script>
