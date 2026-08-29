import type { BlogPost } from './Blog'
import type { Music } from './Music'

export interface FeedBlock {
  id: string
  feed_entry_id: string
  type: 'blog' | 'music'
  order: number
  blog?: BlogPost
  music?: Music
}

export interface FeedEntry {
  id: string
  slug: string
  title: string
  subtitle: string
  author: {
    email: string
    display_name: string
  }
  published_at: string | null
  created_at: string
  updated_at: string
  blocks: FeedBlock[]
  summary: string
  cover_image: string
}

export interface FeedBlockCreateData {
  title?: string
  content?: string
  v?: string
  artist?: string
}

export interface FeedBlockInput {
  type: 'blog' | 'music'
  order: number
  create_data?: FeedBlockCreateData
  ref?: string
}

export interface FeedUpsertInput {
  title: string
  subtitle: string
  published_at: string | null
  blocks: FeedBlockInput[]
}

export interface EditorBlock {
  id: string
  type: 'blog' | 'music'
  mode: 'create' | 'reference'
  order: number
  createData: FeedBlockCreateData
  refId?: string
  refType?: 'blog' | 'music'
  searchQuery?: string
  searchResults?: Array<{ id: string; label: string; ref: string }>
  refLabel?: string
}
