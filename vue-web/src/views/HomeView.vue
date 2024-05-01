<script setup lang="ts">
import { listBlogPost } from '@/services/blogPost';

const blogPost = listBlogPost();
</script>

<template>
  <div class="flex flex-col gap-8">
    <BlogSummary
      v-for="blog in blogPost"
      :key="blog.uri"
      v-bind="blog"
    ></BlogSummary>
  </div>
</template>
