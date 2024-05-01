<template>
  <div class="markdown" v-html="parsedContent"></div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';

import rehypeSanitize from 'rehype-sanitize';
import rehypeStringify from 'rehype-stringify';
import remarkGfm from 'remark-gfm';
import remarkParse from 'remark-parse';
import remarkRehype from 'remark-rehype';
import { unified } from 'unified';
import rehypeExternalLinks from 'rehype-external-links';
import remarkBreaks from 'remark-breaks';

interface Props {
  md: string,
};

const $props = defineProps<Props>();

const parsedContent = ref('');

const parse2HTML = async (content: string) => {
  const html = await unified()
  .use(remarkParse)
  .use(remarkGfm)
  .use(rehypeExternalLinks, { rel: ['nofollow'],target: '_blank' })
  .use(remarkBreaks)
  .use(remarkRehype)
  .use(rehypeSanitize)
  .use(rehypeStringify)
  .process(content);

  return String(html);
};

onMounted(async () => {
  parsedContent.value = await parse2HTML($props.md);
});

watch($props, async (oldVal, newVal) => {
  parsedContent.value = await parse2HTML(newVal.md);
})
</script>

<style type="postcss">
.markdown, .markdown * {
  h1 {
    @apply text-2xl font-bold;
  }

  h2 {
    @apply text-xl font-bold;
  }

  a {
    @apply hover:underline font-medium;
  }

  hr {
    @apply border-accent mt-2 mb-4;
  }

  p {
    @apply py-2;
  }
}
</style>