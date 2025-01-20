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
import inspectUrls from '@jsdevtools/rehype-url-inspector';
import lazyLoadPlugin from 'rehype-plugin-image-native-lazy-loading'

import { generateResizedSrc, generateSrcset } from '@/helpers/cdn';

interface Props {
  md: string,
};

const $props = defineProps<Props>();

const parsedContent = ref('');

const parse2HTML = async (content: string) => {
  const processor = unified()
  .use(remarkParse)
  .use(remarkGfm)
  .use(remarkRehype)
  .use(rehypeExternalLinks, { rel: ['nofollow'],target: '_blank' })
  .use(remarkBreaks)
  .use(inspectUrls, {
    inspectEach(url) {
      if (!url.url.includes('//') && url?.node?.properties?.src) {
        url.node.properties.src = generateResizedSrc(url.url, 1024);
        url.node.properties.srcset = generateSrcset(url.url, 75);
        url.node.properties.loading = 'lazy';
      }

      return url;
    },
    selectors: [
      "img[src]",
    ]
  })
  .use(rehypeSanitize)
  // @ts-expect-error: TS function overloading issue with upstream plugin
  .use(lazyLoadPlugin)
  .use(rehypeStringify);

  const html = await processor.process(content);

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
    @apply text-2xl font-bold mt-8;
  }

  h2 {
    @apply text-xl font-bold mt-8;
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
  
  blockquote {
    @apply bg-accent px-4 py-2 text-sm rounded-xl text-secondary text-center mb-4;
  }

  img {
    @apply rounded-xl w-full;
  }

  p > img {
    @apply grow object-cover w-[32%];
  }

  p:has(img) {
    @apply flex flex-wrap gap-2;
  }
}
</style>