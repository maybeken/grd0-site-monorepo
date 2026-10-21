import { computed, ref } from 'vue'

import { errorStatus, isArticle } from '@/helpers/article'
import { createArticleRaw, getArticleDraftRaw, updateArticleRaw } from '@/services/article'

import type { Article, ArticleBlock } from '@/interfaces/Article'

export type DraftError = 'load' | 'create' | 'save' | ''

function cloneBlocks(blocks: ArticleBlock[]): ArticleBlock[] {
  return JSON.parse(JSON.stringify(blocks || [])) as ArticleBlock[]
}

export function useArticleDraft() {
  const article = ref<Article | null>(null)
  const uri = ref('')
  const title = ref('')
  const subtitle = ref('')
  const published_at = ref<string | null>(null)
  const blocks = ref<ArticleBlock[]>([])
  const loading = ref(false)
  const saving = ref(false)
  const error = ref<DraftError>('')
  const conflict = ref(false)
  const baseline = ref('')

  function currentSnapshot(): string {
    return JSON.stringify({
      uri: uri.value,
      title: title.value,
      subtitle: subtitle.value,
      published_at: published_at.value,
      blocks: blocks.value
    })
  }

  const dirty = computed(
    () => article.value !== null && currentSnapshot() !== baseline.value
  )

  function apply(value: Article): void {
    article.value = value
    uri.value = value.uri
    title.value = value.title
    subtitle.value = value.subtitle
    published_at.value = value.published_at
    blocks.value = cloneBlocks(value.blocks)
    error.value = ''
    conflict.value = false
    baseline.value = currentSnapshot()
  }

  function reset(): void {
    article.value = null
    uri.value = ''
    title.value = ''
    subtitle.value = ''
    published_at.value = null
    blocks.value = []
    loading.value = false
    saving.value = false
    error.value = ''
    conflict.value = false
    baseline.value = ''
  }

  async function open(target: string): Promise<boolean> {
    loading.value = true
    error.value = ''
    conflict.value = false

    try {
      const response = await getArticleDraftRaw(target).send()
      if (!isArticle(response)) {
        error.value = 'load'
        return false
      }

      apply(response)
      return true
    } catch {
      error.value = 'load'
      return false
    } finally {
      loading.value = false
    }
  }

  async function create(target: string): Promise<boolean> {
    loading.value = true
    error.value = ''
    conflict.value = false

    try {
      const response = await createArticleRaw(target ? { uri: target } : {}).send()
      if (!isArticle(response)) {
        error.value = 'create'
        return false
      }

      apply(response)
      return true
    } catch (failure) {
      conflict.value = errorStatus(failure) === 409
      if (!conflict.value) error.value = 'create'
      return false
    } finally {
      loading.value = false
    }
  }

  async function save(): Promise<boolean> {
    if (!article.value) return false

    const previous_uri = article.value.uri

    saving.value = true
    error.value = ''
    conflict.value = false

    try {
      const response = await updateArticleRaw(previous_uri, {
        uri: uri.value,
        title: title.value,
        subtitle: subtitle.value,
        published_at: published_at.value,
        blocks: blocks.value
      }).send()

      if (!isArticle(response)) {
        error.value = 'save'
        return false
      }

      apply(response)
      return true
    } catch (failure) {
      if (uri.value !== previous_uri && errorStatus(failure) === 409) {
        conflict.value = true
        uri.value = previous_uri
      } else {
        error.value = 'save'
      }
      return false
    } finally {
      saving.value = false
    }
  }

  return {
    article,
    uri,
    title,
    subtitle,
    published_at,
    blocks,
    loading,
    saving,
    error,
    conflict,
    dirty,
    open,
    create,
    save,
    reset
  }
}
