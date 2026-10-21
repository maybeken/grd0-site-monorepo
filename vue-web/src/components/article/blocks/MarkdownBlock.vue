<template>
  <MarkdownDisplay :md="$props.block.md" :lazy_loading="!$props.eager"></MarkdownDisplay>
</template>

<script setup lang="ts">
import type { ArticleBlockMarkdown } from '@/interfaces/Article'

interface Props {
  block: ArticleBlockMarkdown
  eager?: boolean
}

const $props = defineProps<Props>()
</script>
