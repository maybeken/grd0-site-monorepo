<template>
  <BlogCard v-if="content">
    <div class="p-4 lg:max-w-[66%] mx-auto">
      <div class="py-4">
        <h1 class="font-black">{{ content.title }}</h1>
        <h2 class="py-2 text-secondary">{{ content.subtitle }}</h2>
      </div>

      <div class="flex gap-4">
        <ProfileIcon class="w-24 md:w-16" :src="`${ASSET_URL}/profile/${content.author.email}.jpg`"></ProfileIcon>
        <div class="my-auto text-sm md:text-md">
          <p class="text-secondary">{{ content.author.display_name }}</p>
          <p class="text-secondary italic">Posted At: {{ dayjs(content.created_at).calendar() }} | Last Edited: {{ dayjs(content.updated_at).calendar() }}</p>
        </div>
      </div>

      <hr class="border-accent my-2" />
      
      <div class="content px-2 py-2 text-justify">
        <MarkdownDisplay :md="content.content"></MarkdownDisplay>
      </div>
    </div>
  </BlogCard>
  <div v-else>
    <NotFound></NotFound>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { useRoute } from 'vue-router';
import { listBlogPost } from '@/services/blogPost';

import dayjs from 'dayjs';
import calendar from 'dayjs/plugin/calendar';

import type { Ref } from 'vue';
import type { BlogPost } from '@/interfaces/Blog';

dayjs.extend(calendar);

const ASSET_URL = import.meta.env.VITE_ASSETS_URL;

const $route = useRoute();
const uri = $route.params.slug;

const blogPost = listBlogPost();
const content: Ref<BlogPost | undefined> = ref();

watch(
  blogPost,
  (newVal) => {
    content.value = newVal.find((item) => {
      return item.uri === uri;
    });
  },
)
</script>