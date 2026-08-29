import { useRequest } from 'alova/client'
import { dataInstance, adminInstance } from './api'

import type { Ref } from 'vue'
import type { FeedEntry, FeedUpsertInput } from '@/interfaces/Feed'

function listFeedRaw(limit?: number, before?: string) {
  let url = '/v2/feed?limit=' + (limit ?? 20)
  if (before) url += '&before=' + encodeURIComponent(before)
  return dataInstance.Get<FeedEntry[]>(url)
}

function listFeed(limit?: number, before?: string): {
  loading: Ref<boolean, boolean>
  data: Ref<FeedEntry[]>
} {
  try {
    const { loading, data } = useRequest(listFeedRaw(limit, before))
    return { loading, data }
  } catch (error: unknown) {
    throw error
  }
}

function getFeedBySlugRaw(slug: string) {
  return dataInstance.Get<FeedEntry>(`/feed/${slug}`)
}

function getFeedBySlug(slug: string): {
  loading: Ref<boolean, boolean>
  data: Ref<FeedEntry>
} | void {
  if (!slug) return
  try {
    const { loading, data } = useRequest(getFeedBySlugRaw(slug))
    return { loading, data }
  } catch (error: unknown) {
    throw error
  }
}

function upsertFeedRaw(slug: string, payload: FeedUpsertInput) {
  return adminInstance.Put<FeedEntry>(`/feed/${slug}`, payload)
}

function deleteFeedRaw(slug: string) {
  return adminInstance.Delete(`/feed/${slug}`)
}

export {
  listFeedRaw,
  listFeed,
  getFeedBySlugRaw,
  getFeedBySlug,
  upsertFeedRaw,
  deleteFeedRaw
}
