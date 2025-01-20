<template>
  <div v-if="content || loading">
    <BlogCard class="py-4 px-16">
      <BlogDetail :content="content" :loading="loading"></BlogDetail>
    </BlogCard>
  </div>
  <div v-else>
    <NotFound></NotFound>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router';
import { getBlogPost } from '@/services/blogPost';

import dayjs from 'dayjs';
import calendar from 'dayjs/plugin/calendar';

dayjs.extend(calendar);

const $route = useRoute();
const uri = $route.params.slug;
const uri_sanitized = typeof uri === "string" ? uri : uri[0];

const response = getBlogPost(uri_sanitized);
const content = response?.data;
const loading = response?.loading;
</script>