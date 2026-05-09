<template>
  <div class="flex gap-4 pb-4">
    <div class="flex flex-col gap-2 grow">
      <input
        class="text-3xl w-full rounded-full px-4"
        type="text"
        placeholder="Blog Post Title"
        v-model="editable_title"
        :disabled="loading"
      />
      <input
        class="text-2xl w-full rounded-full px-4"
        type="text"
        placeholder="Blog Post Subtitle"
        v-model="editable_subtitle"
        :disabled="loading"
      />
      <div class="flex gap-2 px-4">
        <span>Publish Date: </span>
        <input
          class="w-72 border border-dotted rounded-xl px-2"
          type="text"
          v-model="editable_publish_date"
          :disabled="loading"
        />
        <span>{{
          dayjs(editable_publish_date) <= dayjs('0001-01-01T00:00:00.000Z')
            ? 'Unpublished'
            : dayjs(editable_publish_date).fromNow()
        }}</span>
      </div>
    </div>

    <div class="flex gap-2 mt-auto">
      <button
        class="p-4 border border-dotted rounded-full"
        title="New"
        @click="newPost"
        :disabled="loading"
      >
        <Icon icon="mynaui:edit" height="2rem"></Icon>
      </button>
      <button
        class="p-4 border border-dotted rounded-full"
        title="Save"
        @click="savePost"
        :disabled="loading"
      >
        <Icon icon="mynaui:save" height="2rem"></Icon>
      </button>
    </div>
  </div>

  <div class="flex gap-4">
    <div class="w-1/2">
      <bubble-menu v-if="editor" :editor="editor" :tippy-options="{ duration: 150 }">
        <div
          class="flex gap-1 bg-gray-900 border border-gray-700 rounded-xl p-1 shadow-lg"
        >
          <button
            class="px-2 py-1 rounded-lg text-sm font-bold hover:bg-gray-700"
            :class="{ 'bg-gray-700 text-white': editor.isActive('bold') }"
            @click="editor.chain().focus().toggleBold().run()"
          >
            B
          </button>
          <button
            class="px-2 py-1 rounded-lg text-sm italic hover:bg-gray-700"
            :class="{ 'bg-gray-700 text-white': editor.isActive('italic') }"
            @click="editor.chain().focus().toggleItalic().run()"
          >
            I
          </button>
          <button
            class="px-2 py-1 rounded-lg text-sm line-through hover:bg-gray-700"
            :class="{ 'bg-gray-700 text-white': editor.isActive('strike') }"
            @click="editor.chain().focus().toggleStrike().run()"
          >
            S
          </button>
          <button
            class="px-2 py-1 rounded-lg text-sm font-mono hover:bg-gray-700"
            :class="{ 'bg-gray-700 text-white': editor.isActive('code') }"
            @click="editor.chain().focus().toggleCode().run()"
          >
            &lt;/&gt;
          </button>
          <button
            class="px-2 py-1 rounded-lg text-sm hover:bg-gray-700"
            :class="{ 'bg-gray-700 text-white': editor.isActive('link') }"
            @click="toggleLink"
          >
            <Icon icon="mynaui:link" height="1rem" />
          </button>
        </div>
      </bubble-menu>

      <editor-content
        :editor="editor"
        class="editor-wrapper w-full h-full rounded-xl overflow-hidden bg-gray-900 border border-gray-700 p-4 min-h-[400px] prose prose-invert max-w-none"
      />
    </div>

    <div class="w-1/2 max-h-lvh overflow-y-scroll">
      <MarkdownDisplay :md="editable_content"></MarkdownDisplay>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, watch, onBeforeUnmount } from 'vue'
import { useEditor, EditorContent } from '@tiptap/vue-3'
import { BubbleMenu } from '@tiptap/vue-3/menus'
import StarterKit from '@tiptap/starter-kit'
import Placeholder from '@tiptap/extension-placeholder'
import Image from '@tiptap/extension-image'
import Link from '@tiptap/extension-link'
import { Markdown } from 'tiptap-markdown'
import { SlashCommand } from './extensions/SlashCommand'
import { getBlogPost, upsertBlogPostRaw, requestBlogAttachmentRaw } from '@/services/blogPost'
import { useModal } from '@/composables/useModal'
import { useToast } from '@/composables/useToast'

import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
dayjs.extend(relativeTime)

const modal = useModal()
const toast = useToast()

interface Props {
  uri: string
}

const { uri } = defineProps<Props>()
const $emit = defineEmits(['new'])

const response = getBlogPost(uri)
const post = response?.data
const loading = response?.loading || ref(false)

const editable_content = ref('')
const editable_title = ref('')
const editable_subtitle = ref('')
const editable_publish_date = ref(dayjs('0001-01-01T00:00:00Z').toISOString())

let serializeTimer: ReturnType<typeof setTimeout> | null = null

const editor = useEditor({
  extensions: [
    StarterKit.configure({
      heading: {
        levels: [1, 2, 3],
      },
    }),
    Placeholder.configure({
      placeholder: 'Markdown Contents Here — type / for commands',
    }),
    Image,
    Link.configure({
      openOnClick: false,
    }),
    Markdown.configure({
      html: false,
      breaks: true,
      linkify: true,
      transformPastedText: true,
      transformCopiedText: true,
    }),
    SlashCommand,
  ],
  onUpdate: ({ editor: ed }) => {
    if (serializeTimer) clearTimeout(serializeTimer)
    serializeTimer = setTimeout(() => {
      editable_content.value = (ed.storage as any).markdown.getMarkdown()
    }, 300)
  },
  editorProps: {
    attributes: {
      class: 'focus:outline-none',
    },
    handleDrop: (view, event: DragEvent, _slice, _moved) => {
      if (!event.dataTransfer || !event.dataTransfer.files.length) return false

      const files = event.dataTransfer.files
      const coordinates = view.posAtCoords({
        left: event.clientX,
        top: event.clientY,
      })

      if (!coordinates) return false

      const pos = coordinates.pos

      for (let i = 0; i < files.length; i++) {
        const file = files[i]
        if (!file) continue

        if (!file.type.startsWith('image/')) continue

        handleImageUpload(file, pos + i)
      }

      return true
    },
  },
})

async function toggleLink() {
  if (!editor.value) return

  const previousUrl = editor.value.getAttributes('link').href

  if (previousUrl) {
    editor.value.chain().focus().unsetLink().run()
    return
  }

  const url = await modal.prompt('Enter the URL:', 'Link', 'https://', 'https://')
  if (url) {
    editor.value.chain().focus().setLink({ href: url }).run()
  }
}

async function handleImageUpload(file: File, pos: number) {
  if (!editor.value) return

  const placeholderSrc =
    'data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100" viewBox="0 0 100 100"><rect fill="%23374151" width="100" height="100"/><text fill="%239CA3AF" font-size="12" x="50" y="55" text-anchor="middle">Uploading...</text></svg>'

  const { state, dispatch } = editor.value.view
  const imageType = state.schema.nodes.image
  if (!imageType) return

  const node = imageType.create({ src: placeholderSrc })
  const tr = state.tr.insert(pos, node)
  dispatch(tr)

  try {
    loading.value = true

    const res = await requestBlogAttachmentRaw(file.name)
    await s3PresignedUrlUpload(res.url, file)

    const resolvedPos = editor.value.state.doc.resolve(pos + 1)
    const imageNode = resolvedPos.nodeAfter

    if (imageNode && imageNode.type.name === 'image') {
      editor.value
        .chain()
        .setNodeSelection(pos + 1)
        .updateAttributes('image', { src: `/${res.key}` })
        .run()
    }

    toast.info('Image uploaded.')
  } catch {
    if (editor.value) {
      const resolvedPos = editor.value.state.doc.resolve(pos + 1)

      if (resolvedPos.nodeAfter && resolvedPos.nodeAfter.type.name === 'image') {
        editor.value
          .chain()
          .setNodeSelection(pos + 1)
          .deleteSelection()
          .run()
      }
    }

    toast.error('Failed to upload image.')
  } finally {
    loading.value = false
  }
}

async function s3PresignedUrlUpload(url: string, payload: File) {
  const options = {
    method: 'PUT',
    headers: {
      'Content-Type': payload.type,
    },
    body: payload,
  }

  const response = await fetch(url, options)

  if (response.status !== 200) {
    throw Error('S3 upload failed, check network tab for more info.')
  }
}

if (post) {
  watch(
    post,
    (newVal) => {
      if (!newVal) return

      editable_title.value = newVal.title || ''
      editable_subtitle.value = newVal.subtitle || ''
      editable_publish_date.value = newVal.published_at
        ? dayjs(newVal.published_at).toISOString()
        : dayjs('0001-01-01T00:00:00Z').toISOString()

      if (editor.value && newVal.content !== undefined) {
        editor.value.commands.setContent(newVal.content)
        editable_content.value = newVal.content
      }
    },
    { immediate: true }
  )
}

async function newPost() {
  const confirmed = await modal.confirm(
    'Create a new post? Any unsaved changes will be erased.',
    'New Post'
  )

  if (!confirmed) return

  const newUri = await modal.prompt(
    'Enter the URI for the new post:',
    'New Post URI',
    'my-new-post',
  )

  if (newUri) {
    $emit('new', newUri)
  }
}

async function savePost() {
  const confirmed = await modal.confirm(
    'Are you sure to create/update the post?',
    'Save Post'
  )

  if (!confirmed || !uri) return

  if (serializeTimer) {
    clearTimeout(serializeTimer)
    editable_content.value = (editor.value?.storage as any).markdown.getMarkdown()
  }

  const original_post = post?.value
  const content = {
    ...original_post,
    uri,
    author: {
      email: 'ken.lam@grd0.net',
      display_name: 'Ken Lam',
    },
    content: editable_content.value,
    title: editable_title.value,
    subtitle: editable_subtitle.value,
    created_at: dayjs().toISOString(),
    updated_at: dayjs().toISOString(),
    published_at: editable_publish_date.value,
  }

  try {
    loading.value = true
    await upsertBlogPostRaw(uri, content)
    loading.value = false
    toast.success('Post saved successfully.')
  } catch {
    loading.value = false
    toast.error('Failed to save post.')
  }
}

onBeforeUnmount(() => {
  if (serializeTimer) clearTimeout(serializeTimer)
  editor.value?.destroy()
})
</script>

<style scoped>
.editor-wrapper :deep(.ProseMirror) {
  min-height: 400px;
  outline: none;
}

.editor-wrapper :deep(.ProseMirror p.is-editor-empty:first-child::before) {
  color: #6b7280;
  content: attr(data-placeholder);
  float: left;
  height: 0;
  pointer-events: none;
}

.editor-wrapper :deep(.ProseMirror) h1 {
  font-size: 2rem;
  font-weight: 700;
  margin-top: 1rem;
  margin-bottom: 0.5rem;
}

.editor-wrapper :deep(.ProseMirror) h2 {
  font-size: 1.5rem;
  font-weight: 600;
  margin-top: 0.75rem;
  margin-bottom: 0.5rem;
}

.editor-wrapper :deep(.ProseMirror) h3 {
  font-size: 1.25rem;
  font-weight: 600;
  margin-top: 0.5rem;
  margin-bottom: 0.25rem;
}

.editor-wrapper :deep(.ProseMirror) ul,
.editor-wrapper :deep(.ProseMirror) ol {
  padding-left: 1.5rem;
}

.editor-wrapper :deep(.ProseMirror) ul {
  list-style-type: disc;
}

.editor-wrapper :deep(.ProseMirror) ol {
  list-style-type: decimal;
}

.editor-wrapper :deep(.ProseMirror) blockquote {
  border-left: 3px solid #4b5563;
  padding-left: 1rem;
  color: #9ca3af;
}

.editor-wrapper :deep(.ProseMirror) pre {
  background: #1f2937;
  border-radius: 0.75rem;
  padding: 1rem;
  font-family: 'Fira Code', monospace;
  font-size: 0.875rem;
}

.editor-wrapper :deep(.ProseMirror) code {
  background: #374151;
  border-radius: 0.25rem;
  padding: 0.125rem 0.25rem;
  font-size: 0.875rem;
}

.editor-wrapper :deep(.ProseMirror) pre code {
  background: none;
  padding: 0;
}

.editor-wrapper :deep(.ProseMirror) img {
  max-width: 100%;
  border-radius: 0.5rem;
}
</style>
