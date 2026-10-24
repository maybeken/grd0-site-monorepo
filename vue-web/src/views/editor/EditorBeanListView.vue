<template>
  <EditorPage>
    <EditorPageHeader :title="$t('editor.beans')">
      <template #actions>
        <EditorButton v-if="!show_form" variant="primary" @click="openNew">
          <Icon icon="mynaui:plus" height="1rem"></Icon>
          {{ $t('coffee.editor.newBean') }}
        </EditorButton>
      </template>
    </EditorPageHeader>

    <BeanForm v-if="show_form" :bean="null" @saved="onSaved" @cancel="closeForm"></BeanForm>

    <EditorError v-else-if="load_failed" :message="t('editor.loadFailed')">
      <template #actions>
        <EditorButton variant="secondary" @click="fetchBeans">
          {{ $t('common.retry') }}
        </EditorButton>
      </template>
    </EditorError>

    <EditorLoading v-else-if="loading"></EditorLoading>

    <div v-else class="flex flex-col gap-2">
      <EditorEntityRow
        v-for="bean in beans"
        :key="bean.id"
        :title="bean.name"
        :to="`/editor/coffee/beans/${bean.id}`"
        @delete="remove(bean)"
      >
        <span
          v-if="bean.roaster || bean.origin"
          class="truncate text-sm font-medium text-foreground/60"
        >
          {{ [bean.roaster, bean.origin].filter(Boolean).join($t('coffee.editor.dotSeparator')) }}
        </span>
      </EditorEntityRow>

      <EditorEmptyState v-if="!beans.length" :message="$t('coffee.editor.noBeans')"></EditorEmptyState>
    </div>
  </EditorPage>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import BeanForm from '@/components/coffee/BeanForm.vue'
import EditorEntityRow from '@/components/editor/EditorEntityRow.vue'
import EditorEmptyState from '@/components/editor/ui/EditorEmptyState.vue'
import EditorError from '@/components/editor/ui/EditorError.vue'
import EditorLoading from '@/components/editor/ui/EditorLoading.vue'
import EditorPage from '@/components/editor/ui/EditorPage.vue'
import EditorPageHeader from '@/components/editor/ui/EditorPageHeader.vue'
import EditorButton from '@/components/editor/ui/EditorButton.vue'
import Icon from '@/components/shared/Icon.vue'
import { deleteCoffeeBeanRaw, getCoffeeBeansRaw } from '@/services/coffee'
import { useToast } from '@/composables/useToast'

import type { CoffeeBean } from '@/interfaces/Coffee'

const { t } = useI18n()
const toast = useToast()

const beans = ref<CoffeeBean[]>([])
const show_form = ref(false)
const loading = ref(true)
const load_failed = ref(false)

async function fetchBeans(): Promise<void> {
  loading.value = true
  load_failed.value = false

  try {
    beans.value = await getCoffeeBeansRaw()
  } catch {
    beans.value = []
    load_failed.value = true
  } finally {
    loading.value = false
  }
}

function openNew(): void {
  show_form.value = true
}

function closeForm(): void {
  show_form.value = false
}

async function onSaved(): Promise<void> {
  show_form.value = false
  toast.success(t('common.done'))
  await fetchBeans()
}

async function remove(bean: CoffeeBean): Promise<void> {
  if (!bean.id || !confirm(t('common.confirmDeleteItem', { name: bean.name }))) return

  try {
    await deleteCoffeeBeanRaw(bean.id)
    toast.success(t('common.deleted'))
    await fetchBeans()
  } catch {
    toast.error(t('coffee.editor.failedDeleteBean'))
  }
}

onMounted(fetchBeans)
</script>
