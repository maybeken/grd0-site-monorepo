<template>
  <div class="flex flex-col gap-16 h-fit">
    <div class="w-full">
      <PostTable
        :columns="columns"
        :actions="actions"
        :data="articles"
        @edit="(uri: string) => (selected = uri)"
        @delete="unpublishPost"
      ></PostTable>
    </div>

    <div>
      <SideBySideEditor
        :uri="selected"
        :key="selected"
        @new="
          (uri: string) => {
            selected = uri
          }
        "
      ></SideBySideEditor>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
dayjs.extend(relativeTime)

import { deleteBlogPostRaw, listBlogPostAdmin } from '@/services/blogPost'

const { t } = useI18n()

const response = listBlogPostAdmin()
const articles = response?.data
const columns = {
  title: {
    display_name: t('blog.editor.columnTitle')
  },
  uri: {
    display_name: t('blog.editor.columnUri'),
    formatter: (uri: string) => `/blog/${uri}`
  },
  published_at: {
    display_name: t('blog.editor.columnPublished'),
    formatter: (data: string) => {
      if (!data) return t('blog.editor.unpublished')

      const publish_date = dayjs(data)

      if (publish_date <= dayjs('0001-01-01T00:00:00.000Z')) {
        return t('blog.editor.unpublished')
      } else if (publish_date.isBefore(dayjs())) {
        return t('blog.editor.publishedAgo', { time: publish_date.fromNow() })
      }

      return t('blog.editor.scheduledAgo', { time: publish_date.fromNow() })
    }
  }
}
const actions = [
  { name: 'edit', display_name: t('common.edit'), data_key: 'uri' },
  { name: 'delete', display_name: t('common.delete'), data_key: 'uri' }
]

const selected = ref('')

async function unpublishPost(uri: string) {
  const confirmed = confirm(t('blog.editor.confirmUnpublish'))

  if (confirmed) {
    try {
      await deleteBlogPostRaw(uri)
      alert(t('common.done'))
    } catch (error) {
      alert(t('common.failed'))
    }
  }
}
</script>

<style lang="postcss" scoped></style>
