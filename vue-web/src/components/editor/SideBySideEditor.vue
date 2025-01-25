<template>
  <div class="flex gap-4 pb-4">
    <div class="flex flex-col gap-2 grow">
      <input class="text-3xl w-full" type="text" placeholder="Blog Post Title" v-model="editable_title" />
      <input class="text-2xl w-full" type="text" placeholder="Blog Post Subtitle" v-model="editable_subtitle" />
      <div class="flex gap-2">
        <span>Publish Date: </span>
        <input class="w-72 border border-dotted rounded-xl px-2" type="text" v-model="editable_publish_date" />
        <span>{{ dayjs(editable_publish_date).fromNow() }}</span>
      </div>
    </div>

    <div class="flex gap-2 mt-auto">
      <button class="p-4 border border-dotted rounded-full" title="New" @click="newPost">
        <Icon icon="mynaui:edit" height="2rem"></Icon>
      </button>
      <button class="p-4 border border-dotted rounded-full" title="Save" @click="savePost">
        <Icon icon="mynaui:save" height="2rem"></Icon>
      </button>
    </div>
  </div>

  <div class="flex gap-4">
    <div class="w-1/2">
      <textarea class="w-full h-full rounded-xl p-4 bg-accent shadow-inner shadow-background"
        placeholder="Markdown Contents Here" v-model="editable_content"></textarea>
    </div>

    <div class="w-1/2 max-h-lvh overflow-y-scroll">
      <MarkdownDisplay :md="editable_content"></MarkdownDisplay>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue';
import { getBlogPost } from '@/services/blogPost';

import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
dayjs.extend(relativeTime);

interface Props {
  uri: string;
}

const { uri } = defineProps<Props>();
const $emit = defineEmits(['new']);

const response = getBlogPost(uri);
const post = response?.data;

const editable_content = ref(post?.value?.content || '');
const editable_title = ref(post?.value?.title || '');
const editable_subtitle = ref(post?.value?.subtitle || '');
const editable_publish_date = ref(dayjs(post?.value?.published_at).toISOString() || dayjs().toISOString());

function newPost(): void {
  const confirmed = confirm('Are you sure to create new post? Any unsaved changes will be erased!');

  if (confirmed) $emit('new')
}

function savePost(): void {
  const original_post = post?.value;
  const content = {
    ...original_post,
    content: editable_content.value,
    title: editable_title.value,
    sub_title: editable_subtitle.value,
    updated_at: dayjs().toISOString(),
    published_at: editable_publish_date.value,
  };

  navigator.clipboard.writeText(JSON.stringify(content));
  alert('Post copied to clipboard.');
}

if (post) {
  watch(post, (newVal) => {
    editable_content.value = newVal?.content;
    editable_title.value = newVal?.title;
    editable_subtitle.value = newVal?.subtitle;
    editable_publish_date.value = dayjs(newVal?.published_at).toISOString() || dayjs().toISOString();
  });
}
</script>