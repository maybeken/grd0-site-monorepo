<template>
  <div class="flex gap-4 h-lvh">
    <div class="flex flex-col gap-2 w-1/2">
      File:
      <div class="flex flex-col gap-2 bg-secondary text-black font-mono p-2 rounded-lg">
        <div class="flex gap-2">
          <input
            v-model="selectedArticleUri"
            type="text"
            class="grow path bg-secondary text-black font-mono p-2 rounded-lg border border-black"
          ></input>
          <button
            @click="saveToClipboard"
            class="path bg-secondary text-black font-mono p-2 rounded-lg border border-black"
          >Save</button>
          <button
            @click="reloadArticles"
            class="path bg-secondary text-black font-mono p-2 rounded-lg border border-black"
          >Reset</button>
        </div>
        
        <span class="cursor-pointer" @click="newArticle">&lt; New Article &gt;</span>
        <span v-for="post of articles" :key="post.uri" class="cursor-pointer" @click="selectArticle(post.uri)">/{{ post.uri }} [{{ post.title }}]</span>
      </div>

      Editor:
      <input
        v-model="selectedArticleTitle"
        @input="changeTitle(selectedArticleUri, selectedArticleTitle)"
        :disabled="!selectedArticleUri"
        type="text"
        class="text-2xl bg-transparent text-secondary font-mono p-2 rounded-lg border border-secondary"
      ></input>
      <input
        v-model="selectedArticleSubtitle"
        @input="changeSubtitle(selectedArticleUri, selectedArticleSubtitle)"
        :disabled="!selectedArticleUri"
        type="text"
        class="text-lg bg-transparent text-secondary font-mono p-2 rounded-lg border border-secondary"
      ></input>
      <textarea
        v-model="selectedArticleContent"
        :disabled="!selectedArticleUri"
        @input="changeContent(selectedArticleUri, selectedArticleContent)"
        class="bg-transparent text-secondary rounded-lg p-2 border border-secondary grow"
      ></textarea>
    </div>

    <div class="flex flex-col gap-2 w-1/2">
      Preview:

      <MarkdownDisplay class="overflow-y-scroll" :md="selectedArticleContent"></MarkdownDisplay>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import dayjs from 'dayjs';

import { listBlogPost } from '@/services/blogPost';

const response = listBlogPost();
const blogPost = response?.data;
const articles = response?.data;
const selectedArticleUri = ref('');
const selectedArticleTitle = ref('');
const selectedArticleSubtitle = ref('');
const selectedArticleContent = ref('');

const selectArticle = (uri: string) => {
  if (!articles.value) return;

  const article = articles.value.find((article) => article.uri === uri);

  if (!article) return;

  selectedArticleUri.value = uri;
  selectedArticleTitle.value = article.title;
  selectedArticleSubtitle.value = article.subtitle;
  selectedArticleContent.value = article.content;
};

const changeTitle = (uri: string, title: string) => {
  updateBlogPost(uri, { title });
};

const changeSubtitle = (uri: string, subtitle: string) => {
  updateBlogPost(uri, { subtitle });
};

const changeContent = (uri: string, content: string) => {
  updateBlogPost(uri, { content });
};

const updateBlogPost = (uri: string, payload: { title?: string, subtitle?: string, content?: string }) => {
  if (!uri) return;
  if (!articles.value) return;

  const articleIdx = articles.value.findIndex((article) => article.uri === uri);
  const article = articles.value[articleIdx];

  articles.value = articles.value.filter((article) => article.uri !== uri);
  articles.value = [...articles.value, {
    ...article,
    title: payload.title ?? article.title,
    subtitle: payload.subtitle ?? article.subtitle,
    content: payload.content ?? article.content,
    updated_at: dayjs().toISOString(),
  }];
};

const newArticle = () => {
  const uri = window.prompt("Enter the URI of the post:");

  if (!uri) return;
  if (!articles.value) return;

  articles.value = [...articles.value, {
    uri,
    author: {
      email: "ken.lam@grd0.net",
      display_name: "Ken Lam",
    },
    title: "",
    subtitle: "",
    content: "",
    created_at: dayjs().toISOString(),
    updated_at: dayjs().toISOString(),
  }];

  selectArticle(uri);
};

const saveToClipboard = () => {
  if (!articles.value) return;

  articles.value = articles.value.sort((a, b) => dayjs(a.created_at).unix() - dayjs(b.created_at).unix()).reverse();
  
  navigator.clipboard.writeText(JSON.stringify(articles.value));
  alert('Saved to clipboard!');
};

const reloadArticles = () => {
  const confirmation = confirm('Are you sure to reload from save?');

  if (!confirmation) return;

  articles.value = blogPost.value;
  resetValues();
};

const resetValues = () => {
  selectedArticleUri.value = '';
  selectedArticleTitle.value = '';
  selectedArticleSubtitle.value = '';
  selectedArticleContent.value = '';
};
</script>

<style type="postcss" scoped>
</style>