import type { Article, ArticleBlock, ArticleBlockType } from '@/interfaces/Article'

export function createArticleBlock(type: ArticleBlockType): ArticleBlock {
  const id = crypto.randomUUID()

  switch (type) {
    case 'markdown':
      return { id, type, md: '' }
    case 'image':
      return { id, type, collection: '', filename: '' }
    case 'music':
      return { id, type, v: '', title: '' }
    case 'excalidraw':
      return { id, type, scene_path: '' }
  }
}

export function isArticle(value: unknown): value is Article {
  if (!value || typeof value !== 'object') return false

  const article = value as Article
  return typeof article.uri === 'string' && Array.isArray(article.blocks)
}

export function articleSceneId(scene_path: string): string | null {
  if (!scene_path) return null

  const filename = scene_path.split('/').pop() || ''
  const match = filename.match(/^([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})\.json$/i)

  return match ? match[1]! : null
}

export function errorStatus(error: unknown): number | undefined {
  if (!error || typeof error !== 'object') return undefined

  const candidate = error as { status?: unknown; response?: { status?: unknown } }
  const status = candidate.status ?? candidate.response?.status

  return typeof status === 'number' ? status : undefined
}
