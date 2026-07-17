<template>
  <div class="markdown" v-html="parsedContent"></div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'

import rehypeSanitize from 'rehype-sanitize'
import rehypeStringify from 'rehype-stringify'
import remarkGfm from 'remark-gfm'
import remarkParse from 'remark-parse'
import remarkRehype from 'remark-rehype'
import { unified } from 'unified'
import rehypeExternalLinks from 'rehype-external-links'
import remarkBreaks from 'remark-breaks'
import inspectUrls from '@jsdevtools/rehype-url-inspector'
import lazyLoadPlugin from 'rehype-plugin-image-native-lazy-loading'

import { generateResizedSrc, generateSrcset } from '@/helpers/cdn'

interface Props {
  md: string
  lazy_loading: boolean
}

const $props = defineProps<Props>()

const parsedContent = ref('')

const parse2HTML = async (content: string) => {
  const processor = unified()
    .use(remarkParse)
    .use(remarkGfm)
    .use(remarkRehype)
    .use(rehypeExternalLinks, { rel: ['nofollow'], target: '_blank' })
    .use(remarkBreaks)
    .use(inspectUrls, {
      inspectEach(url) {
        if (!url.url.includes('//') && url?.node?.properties?.src) {
          url.node.properties.src = generateResizedSrc(url.url, 1024)
          url.node.properties.srcset = generateSrcset(url.url, 75)
        }

        return url
      },
      selectors: ['img[src]']
    })
    .use(rehypeSanitize)

  if ($props.lazy_loading != false) {
    processor
      // @ts-expect-error: TS function overloading issue with upstream plugin
      .use(lazyLoadPlugin)
  }

  processor.use(rehypeStringify)

  const html = processor.processSync(content)

  return String(html)
}

onMounted(async () => {
  parsedContent.value = await parse2HTML($props.md)
})

watch($props, async (oldVal, newVal) => {
  parsedContent.value = await parse2HTML(newVal.md)
})
</script>

<style lang="postcss"></style>
