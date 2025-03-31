<template>
  <div class="flex flex-col gap-16 h-fit">
    <div class="w-full">
      <PostTable :columns="columns" :actions="actions" :data="articles" @edit="(uri: string) => selected = uri"></PostTable>
    </div>

    <div>
      <SideBySideEditor :uri="selected" :key="selected" @new="selected = ''"></SideBySideEditor>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
dayjs.extend(relativeTime);

import { listBlogPostAdmin } from '@/services/blogPost';
import SideBySideEditor from '@/components/editor/SideBySideEditor.vue';

const response = listBlogPostAdmin();
const articles = response?.data;
const columns = {
  title: {
    display_name: "Title",
  },
  uri: {
    display_name: "URI",
    formatter: (uri: string) => `/blog/${uri}`
  },
  published_at: {
    display_name: "Published",
    formatter: (data: string) => {
      if (!data) return 'Unpublished';

      const publish_date = dayjs(data);
      const published = publish_date.isBefore(dayjs());

      if (published) {
        return `Published (${publish_date.fromNow()})`;
      }

      return `Scheduled (${publish_date.fromNow()})`;
    },
  }
};
const actions = [
  { name: 'edit', display_name: 'Edit', data_key: 'uri' },
  { name: 'delete', display_name: 'Delete', data_key: 'uri' },
]

const selected = ref("");
</script>

<style lang="postcss" scoped>
</style>