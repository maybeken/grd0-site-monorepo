<template>
  <div class="flex gap-4 pb-4">
    <div class="flex flex-col gap-2 grow">
      <input class="text-3xl w-full rounded-full px-4" type="text" placeholder="Blog Post Title" v-model="editable_title" :disabled="loading" />
      <input class="text-2xl w-full rounded-full px-4" type="text" placeholder="Blog Post Subtitle" v-model="editable_subtitle" :disabled="loading" />
      <div class="flex gap-2 px-4">
        <span>Publish Date: </span>
        <input class="w-72 border border-dotted rounded-xl px-2" type="text" v-model="editable_publish_date" :disabled="loading" />
        <span>{{ dayjs(editable_publish_date) <= dayjs('0001-01-01T00:00:00.000Z') ? "Unpublished" : dayjs(editable_publish_date).fromNow() }}</span>
      </div>
    </div>

    <div class="flex gap-2 mt-auto">
      <button class="p-4 border border-dotted rounded-full" title="New" @click="newPost" :disabled="loading">
        <Icon icon="mynaui:edit" height="2rem"></Icon>
      </button>
      <button class="p-4 border border-dotted rounded-full" title="Save" @click="savePost" :disabled="loading">
        <Icon icon="mynaui:save" height="2rem"></Icon>
      </button>
    </div>
  </div>

  <div class="flex gap-4" @dragover.prevent="" @drop.prevent="fileUpload">
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
import { getBlogPost, upsertBlogPostRaw, requestBlogAttachmentRaw } from '@/services/blogPost';

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
const loading = response?.loading || ref(false);

const editable_content = ref(post?.value?.content || '');
const editable_title = ref(post?.value?.title || '');
const editable_subtitle = ref(post?.value?.subtitle || '');
const editable_publish_date = ref(post?.value ? dayjs(post?.value?.published_at).toISOString() : dayjs("0001-01-01T00:00:00Z").toISOString());

function newPost(): void {
  const confirmed = confirm('Are you sure to create new post? Any unsaved changes will be erased!');

  if (confirmed) $emit('new')
}

async function savePost() {
  const confirmed = confirm('Are you sure to create/update the post?');

  if (!confirmed || !uri) return

  const original_post = post?.value;
  const content = {
    ...original_post,
    uri,
    author: {
      email: "ken.lam@grd0.net",
      display_name: "Ken Lam",
    },
    content: editable_content.value,
    title: editable_title.value,
    subtitle: editable_subtitle.value,
    created_at: dayjs().toISOString(),
    updated_at: dayjs().toISOString(),
    published_at: editable_publish_date.value,
  };

  try {
    loading.value = true;
    await upsertBlogPostRaw(uri, content)
    loading.value = false;
    
    alert("Done.")
  } catch(error) {
    alert("Failed.")
  }
}

if (post) {
  watch(post, (newVal) => {
    editable_content.value = newVal?.content;
    editable_title.value = newVal?.title;
    editable_subtitle.value = newVal?.subtitle;
    editable_publish_date.value = dayjs(newVal?.published_at).toISOString() || dayjs().toISOString();
  });
}

async function fileUpload($event: DragEvent) {
  if (!$event || !$event.dataTransfer) return;

  if ($event.dataTransfer.items) {
    for (const item of $event.dataTransfer.items) {
      if (item.kind === 'file') {
        const file = item.getAsFile();

        if (!file?.name) return;
        const res = await requestBlogAttachmentRaw(file?.name);

        try {
          loading.value = true;
          await s3PresignedUrlUpload(res.url, file);

          editable_content.value += `\n![](/${res.key})\n`;
        } catch (error) {
          console.error(error);
        } finally {
          loading.value = false;
        }
      }
    }
  } else {
    for (const file of $event.dataTransfer.files) {
      if (!file?.name) return;
      const res = await requestBlogAttachmentRaw(file?.name);

      try {
        loading.value = true;
        await s3PresignedUrlUpload(res.url, file);

        editable_content.value += `\n![](/${res.key})\n`;
      } catch (error) {
        console.error(error);
      } finally {
        loading.value = false;
      }
    }
  }
}

async function s3PresignedUrlUpload(url: string, payload: File) {
  const options = {
    method: 'PUT',
    headers: {
      "Content-Type": payload.type,
    },
    body: payload,
  };

  try {
    const response = await fetch(url, options);

    if (response.status !== 200) {
      throw Error("S3 upload failed, check network tab for more info.");
    }
  } catch (error) {
    throw error;
  }
}
</script>