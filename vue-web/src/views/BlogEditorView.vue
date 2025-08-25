<template>
  <div class="flex flex-col gap-16 h-fit">
    <div class="w-full">
      <PostTable :columns="columns" :actions="actions" :data="articles" @edit="(uri: string) => selected = uri" @delete="unpublishPost"></PostTable>
    </div>

    <div>
      <SideBySideEditor :uri="selected" :key="selected" @new="(uri: string) => { selected = uri }"></SideBySideEditor>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
dayjs.extend(relativeTime);

import { deleteBlogPostRaw, listBlogPostAdmin } from '@/services/blogPost';

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

      if (publish_date <= dayjs("0001-01-01T00:00:00.000Z")) {
        return `Unpublished`;
      } else if (publish_date.isBefore(dayjs())) {
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

async function unpublishPost(uri: string) {
  const confirmed = confirm("Are you sure to unpublish the post?");

  if (confirmed) {
    try {
      await deleteBlogPostRaw(uri);
      alert("Done.");
    } catch(error) {
      alert("Failed.");
    }
  }
}
</script>

<style lang="postcss" scoped>
</style>