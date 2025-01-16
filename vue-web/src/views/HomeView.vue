<script setup lang="ts">
import { listBlogPost } from '@/services/blogPost';

const response = listBlogPost();
const blogPost = response.data;
const loading = response.loading;
</script>

<template>
  <div class="flex flex-col gap-8">
    <BlogSummary
      v-if="loading"
      :loading="loading"
    ></BlogSummary>
    <BlogSummary
      v-for="blog in blogPost"
      :key="blog.uri"
      :loading="loading"
      v-bind="blog"
    ></BlogSummary>
  </div>
</template>
