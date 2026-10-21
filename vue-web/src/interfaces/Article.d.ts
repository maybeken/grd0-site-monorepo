import type { PaginatedResponse } from './Gallery'

export interface ArticleBlockMarkdown {
  id: string
  type: 'markdown'
  md: string
}

export interface ArticleBlockImage {
  id: string
  type: 'image'
  collection: string
  filename: string
  caption?: string
}

export interface ArticleBlockMusic {
  id: string
  type: 'music'
  v: string
  title: string
  artist?: string
}

export interface ArticleBlockDrawing {
  id: string
  type: 'excalidraw'
  scene_path: string
}

export type ArticleBlock =
  | ArticleBlockMarkdown
  | ArticleBlockImage
  | ArticleBlockMusic
  | ArticleBlockDrawing

export interface ArticleAuthor {
  email: string
  display_name: string
}

export interface Article {
  uri: string
  author: ArticleAuthor
  title: string
  subtitle: string
  blocks: ArticleBlock[]
  summary: string
  cover: string
  created_at: string
  updated_at: string
  published_at: string | null
}

export interface ArticleListEntry {
  uri: string
  author: ArticleAuthor
  title: string
  subtitle: string
  summary: string
  cover: string
  created_at: string
  updated_at: string
  published_at: string | null
}

export type ArticleListPage = PaginatedResponse<ArticleListEntry>

export interface ArticlePayload {
  uri?: string
  title?: string
  subtitle?: string
  blocks?: ArticleBlock[]
  published_at?: string | null
}

export interface ArticleScene {
  scene_path: string
}

export interface ArticleAttachment {
  url: string
  key: string
}

export type ArticleBlockType = ArticleBlock['type']

