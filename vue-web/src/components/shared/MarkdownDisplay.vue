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

const isImageOnlyParagraph = (element: Element) => {
  if (element.tagName !== 'P') return false

  const children = Array.from(element.children)
  if (!children.length || !children.every((child) => child.tagName === 'IMG')) return false

  return Array.from(element.childNodes).every(
    (node) => node.nodeType !== Node.TEXT_NODE || !node.textContent?.trim()
  )
}

const mergeImageRuns = (html: string) => {
  const doc = new DOMParser().parseFromString(html, 'text/html')
  const elements = Array.from(doc.body.children)
  let run: Element[] = []

  const flush = () => {
    if (run.length > 1) {
      const grid = doc.createElement('div')
      grid.className = 'image-grid'

      for (const paragraph of run) {
        while (paragraph.firstChild) grid.appendChild(paragraph.firstChild)
      }

      run[0].replaceWith(grid)
      for (const paragraph of run.slice(1)) paragraph.remove()
    }

    run = []
  }

  for (const element of elements) {
    if (isImageOnlyParagraph(element)) run.push(element)
    else flush()
  }

  flush()

  return doc.body.innerHTML
}

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

  return mergeImageRuns(String(html))
}

onMounted(async () => {
  parsedContent.value = await parse2HTML($props.md || '')
})

watch(
  () => $props.md,
  async (newValue) => {
    parsedContent.value = await parse2HTML(newValue || '')
  }
)
</script>

<style lang="postcss"></style>
