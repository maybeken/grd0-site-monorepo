<script setup lang="ts">
import { listBlogPost } from '@/services/blogPost'

const response = listBlogPost()
const blogPost = response.data
const loading = response.loading
</script>

<template>
  <div class="flex flex-col gap-8">
    <template v-if="loading">
      <BlogSummary v-for="i in 5" :loading="loading"></BlogSummary>
    </template>
    <BlogSummary
      v-for="blog in blogPost"
      :key="blog.uri"
      :loading="loading"
      v-bind="blog"
    ></BlogSummary>
  </div>
</template>
