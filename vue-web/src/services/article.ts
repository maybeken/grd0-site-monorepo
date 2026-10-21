import { useRequest } from 'alova/client'
import { dataInstance, adminInstance } from './api'

import type { Ref } from 'vue'
import type {
  Article,
  ArticleAttachment,
  ArticleListEntry,
  ArticlePayload,
  ArticleScene
} from '@/interfaces/Article'

function listArticlesRaw() {
  return dataInstance.Get<ArticleListEntry[]>('/article')
}

function listArticles(): { loading: Ref<boolean, boolean>; data: Ref<ArticleListEntry[]> } {
  try {
    const { loading, data } = useRequest(listArticlesRaw())

    return { loading, data }
  } catch (error: unknown) {
    throw error
  }
}

function getArticleRaw(uri: string) {
  return dataInstance.Get<Article>(`/article/${uri}`)
}

function getArticle(uri: string): { loading: Ref<boolean, boolean>; data: Ref<Article> } | void {
  if (!uri) return

  try {
    const { loading, data } = useRequest(getArticleRaw(uri))

    return { loading, data }
  } catch (error: unknown) {
    throw error
  }
}

function listArticlesAdminRaw() {
  return adminInstance.Get<ArticleListEntry[]>('/article/all')
}

function listArticlesAdmin(): { loading: Ref<boolean, boolean>; data: Ref<ArticleListEntry[]> } {
  try {
    const { loading, data } = useRequest(listArticlesAdminRaw())

    return { loading, data }
  } catch (error: unknown) {
    throw error
  }
}

function getArticleDraftRaw(uri: string) {
  return adminInstance.Get<Article>(`/article/${uri}/draft`)
}

function createArticleRaw(payload: ArticlePayload) {
  return adminInstance.Post<Article>('/article', payload)
}

function updateArticleRaw(uri: string, payload: ArticlePayload) {
  return adminInstance.Put<Article>(`/article/${uri}`, payload)
}

function unpublishArticleRaw(uri: string) {
  return adminInstance.Delete(`/article/${uri}`)
}

function purgeArticleRaw(uri: string) {
  return adminInstance.Delete(`/article/${uri}/purge`)
}

function getArticleSceneRaw(uri: string, scene_id: string) {
  return dataInstance.Get<unknown>(`/article/${uri}/scene/${scene_id}`)
}

function getArticleDraftSceneRaw(uri: string, scene_id: string) {
  return adminInstance.Get<unknown>(`/article/${uri}/draft/scene/${scene_id}`)
}

function putArticleSceneRaw(uri: string, scene_id: string, scene: object) {
  return adminInstance.Put<ArticleScene>(`/article/${uri}/scene/${scene_id}`, scene)
}

function deleteArticleSceneRaw(uri: string, scene_id: string) {
  return adminInstance.Delete(`/article/${uri}/scene/${scene_id}`)
}

function requestArticleAttachmentRaw(path: string) {
  return adminInstance.Put<ArticleAttachment>(`/article/attachment/${path}`)
}

export {
  listArticlesRaw,
  listArticles,
  getArticleRaw,
  getArticle,
  listArticlesAdminRaw,
  listArticlesAdmin,
  getArticleDraftRaw,
  createArticleRaw,
  updateArticleRaw,
  unpublishArticleRaw,
  purgeArticleRaw,
  getArticleSceneRaw,
  getArticleDraftSceneRaw,
  putArticleSceneRaw,
  deleteArticleSceneRaw,
  requestArticleAttachmentRaw
}
