<template>
  <div class="flex flex-col gap-6">
    <header class="flex flex-col gap-2">
      <h1 class="text-2xl font-bold">{{ $t('editor.beans') }}<CursorBlink /></h1>
    </header>

    <div v-if="!show_form" class="flex justify-end">
      <button
        class="rounded bg-accent px-3 py-1 text-sm font-semibold text-background hover:brightness-110"
        @click="openNew"
      >
        {{ $t('coffee.editor.newBean') }}
      </button>
    </div>

    <BeanForm v-if="show_form" :bean="null" @saved="onSaved" @cancel="closeForm"></BeanForm>

    <div v-if="!show_form" class="flex flex-col gap-2">
      <EditorEntityRow
        v-for="bean in beans"
        :key="bean.id"
        :title="bean.name"
        :to="`/editor/coffee/beans/${bean.id}`"
        @delete="remove(bean)"
      >
        <span v-if="bean.roaster || bean.origin" class="truncate text-xs text-foreground/60">
          {{ [bean.roaster, bean.origin].filter(Boolean).join(' · ') }}
        </span>
      </EditorEntityRow>

      <p v-if="!beans.length" class="py-4 text-center text-sm text-foreground/40">
        {{ $t('coffee.editor.noBeans') }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import BeanForm from '@/components/coffee/BeanForm.vue'
import EditorEntityRow from '@/components/editor/EditorEntityRow.vue'
import { deleteCoffeeBeanRaw, getCoffeeBeansRaw } from '@/services/coffee'
import { useToast } from '@/composables/useToast'

import type { CoffeeBean } from '@/interfaces/Coffee'

const { t } = useI18n()
const toast = useToast()

const beans = ref<CoffeeBean[]>([])
const show_form = ref(false)

async function fetchBeans(): Promise<void> {
  try {
    beans.value = await getCoffeeBeansRaw()
  } catch {
    toast.error(t('editor.loadFailed'))
  }
}

function openNew(): void {
  show_form.value = true
}

function closeForm(): void {
  show_form.value = false
}

async function onSaved(): Promise<void> {
  closeForm()
  await fetchBeans()
}

async function remove(bean: CoffeeBean): Promise<void> {
  if (!bean.id || !confirm(t('coffee.editor.confirmDeleteBean', { name: bean.name }))) return

  try {
    await deleteCoffeeBeanRaw(bean.id)
    await fetchBeans()
  } catch {
    toast.error(t('coffee.editor.failedDeleteBean'))
  }
}

onMounted(fetchBeans)
</script>
