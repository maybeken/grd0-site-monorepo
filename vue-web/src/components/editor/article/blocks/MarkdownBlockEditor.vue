<template>
  <div class="flex flex-col gap-2" @dragover.prevent="" @drop.prevent="onDrop">
    <div ref="editor" class="w-full"></div>

    <div class="flex items-center gap-2">
      <button
        class="rounded-full border border-dotted px-3 py-1 disabled:opacity-50"
        :disabled="uploading"
        @click="pickFile"
      >
        {{ $t('article.editor.uploadImage') }}
      </button>
      <span v-if="uploading" class="text-sm text-secondary">{{ $t('common.saving') }}</span>
      <input ref="fileInput" type="file" accept="image/*" class="hidden" @change="onPick" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref, useTemplateRef, watch } from 'vue'
// @ts-ignore
import OverType from 'overtype'

import { t } from '@/i18n'
import { useToast } from '@/composables/useToast'
import { requestArticleAttachmentRaw } from '@/services/article'

import type { ArticleBlockMarkdown } from '@/interfaces/Article'

interface Props {
  block: ArticleBlockMarkdown
}

const $props = defineProps<Props>()

const toast = useToast()

const editor = useTemplateRef('editor')
const fileInput = useTemplateRef('fileInput')
const uploading = ref(false)

let overtype: any = null
let syncing = false
let initializing = true

onMounted(() => {
  if (!editor.value) return

  const initial_md = $props.block.md || ''

  const [instance] = new OverType(editor.value, {
    theme: 'cave',
    toolbar: true,
    showStats: true,
    smartLists: true,
    autoResize: true,
    minHeight: '12rem',
    placeholder: t('article.editor.markdownPlaceholder'),
    onChange: (value: string) => {
      if (initializing || syncing) return

      syncing = true
      $props.block.md = value
      syncing = false
    }
  })

  overtype = instance
  overtype?.setValue(initial_md)
  $props.block.md = initial_md
  initializing = false
})

onUnmounted(() => {
  overtype?.destroy?.()
  overtype = null
})

watch(
  () => $props.block.md,
  (value) => {
    if (!overtype || syncing) return
    if (overtype.getValue() === (value || '')) return

    syncing = true
    overtype.setValue(value || '')
    syncing = false
  }
)

function pickFile(): void {
  fileInput.value?.click()
}

function onPick(event: Event): void {
  const input = event.target as HTMLInputElement
  if (input.files) uploadFiles(input.files)
  input.value = ''
}

function onDrop(event: DragEvent): void {
  const files = event.dataTransfer?.files
  if (files) uploadFiles(files)
}

async function uploadFiles(files: FileList): Promise<void> {
  for (const file of Array.from(files)) {
    if (!file.name) continue

    try {
      uploading.value = true

      const target = await requestArticleAttachmentRaw(file.name).send()
      if (!target?.url || !target?.key) throw new Error('Missing upload target')

      const response = await fetch(target.url, {
        method: 'PUT',
        headers: { 'Content-Type': file.type },
        body: file
      })

      if (!response.ok) throw new Error('Upload failed')

      const reference = `![](/${target.key})`
      if (overtype) {
        overtype.insertAtCursor(reference)
        $props.block.md = overtype.getValue()
      } else {
        $props.block.md = `${$props.block.md}${reference}`
      }
    } catch {
      toast.error(t('article.editor.uploadFailed'))
    } finally {
      uploading.value = false
    }
  }
}
</script>
