<template>
  <EditorPage>
    <EditorPageHeader
      :title="bean?.name || $t('editor.beans')"
      :back-to="BEANS_ROUTE"
      :back-label="$t('editor.ui.backTo', { section: $t('editor.beans') })"
    >
      <template #actions>
        <template v-if="bean">
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

    <EditorError v-else-if="!bean" :message="$t('editor.notFound')">
      <template #actions>
        <EditorButton variant="secondary" :to="BEANS_ROUTE">
          {{ $t('editor.ui.backTo', { section: $t('editor.beans') }) }}
        </EditorButton>
      </template>
    </EditorError>

    <BeanForm
      v-else
      ref="bean_form"
      :bean="bean"
      :standalone="false"
      :form-id="FORM_ID"
      @saved="onSaved"
      @cancel="back"
      @update:saving="saving = $event"
      @update:can-save="can_save = $event"
    ></BeanForm>
  </EditorPage>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

import BeanForm from '@/components/coffee/BeanForm.vue'
import EditorButton from '@/components/editor/ui/EditorButton.vue'
import EditorError from '@/components/editor/ui/EditorError.vue'
import EditorLoading from '@/components/editor/ui/EditorLoading.vue'
import EditorPage from '@/components/editor/ui/EditorPage.vue'
import EditorPageHeader from '@/components/editor/ui/EditorPageHeader.vue'
import { useToast } from '@/composables/useToast'
import { getCoffeeBeanRaw } from '@/services/coffee'

import type { CoffeeBean } from '@/interfaces/Coffee'

const BEANS_ROUTE = '/editor/coffee/beans'
const FORM_ID = 'bean-form'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const toast = useToast()

const bean = ref<CoffeeBean | null>(null)
const bean_form = ref<InstanceType<typeof BeanForm> | null>(null)
const loading = ref(true)
const saving = ref(false)
const can_save = ref(false)

const id = computed(() => {
  const value = route.params.id

  return typeof value === 'string' ? value : ''
})

async function load(): Promise<void> {
  loading.value = true
  bean.value = null

  try {
    bean.value = await getCoffeeBeanRaw(id.value)
  } catch {
    bean.value = null
  } finally {
    loading.value = false
  }
}

function onSaved(): void {
  toast.success(t('common.done'))
}

function back(): void {
  router.push(BEANS_ROUTE)
}

function cancel(): void {
  bean_form.value?.requestCancel()
}

watch(id, load, { immediate: true })
</script>
