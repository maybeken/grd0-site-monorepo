<template>
  <MarkdownDisplay :md="$props.block.md" :lazy_loading="true"></MarkdownDisplay>
</template>

<script setup lang="ts">
import type { ArticleBlockMarkdown } from '@/interfaces/Article'

interface Props {
  block: ArticleBlockMarkdown
}

const $props = defineProps<Props>()
</script>
