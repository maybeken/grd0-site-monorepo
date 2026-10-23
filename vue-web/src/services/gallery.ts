import { useRequest } from 'alova/client'
import { adminInstance, dataInstance } from './api'

import type { Ref } from 'vue'
import type {
  GalleryDetail,
  GalleryCollection,
  GalleryCollectionDetail,
  GalleryDirectory,
  PaginatedResponse,
  Asset
} from '@/interfaces/Gallery'

const DETAIL_PAGE_SIZE = 200

function galleryDetailPath(collection: string): string {
  return `/gallery/${collection}`
}

function getGalleryDirectoriesRaw() {
  return adminInstance.Get<GalleryDirectory[]>('/gallery/directories')
}

function upsertGalleryCollectionRaw(collection: GalleryCollectionDetail) {
  return adminInstance.Put<GalleryCollectionDetail>('/gallery/collection', collection)
}

function deleteGalleryCollectionRaw(name: string) {
  return adminInstance.Delete(`/gallery/collection/${encodeURIComponent(name)}`)
}

function getGalleryDetailsRaw(collection: string, page = 1) {
  return dataInstance.Get<PaginatedResponse<GalleryDetail>>(
    `/v2/gallery/details/${galleryDetailPath(collection)}?page=${page}&page_size=${DETAIL_PAGE_SIZE}`
  )
}

function upsertGalleryDetailRaw(collection: string, detail: GalleryDetail) {
  return adminInstance.Put<GalleryDetail>(
    `/gallery/details/${galleryDetailPath(collection)}`,
    detail
  )
}

function deleteGalleryDetailRaw(collection: string, filename: string) {
  return adminInstance.Delete<unknown>(`/gallery/details/${galleryDetailPath(collection)}`, {
    filename
  })
}

function loadMoreAssets(collection: string, page: number) {
  return dataInstance.Get<PaginatedResponse<Asset>>(
    `/v2/gallery/${collection}?page=${page}&page_size=40`
  )
}

function getGalleryDetail(
  path?: string,
  page = 1
): { loading: Ref<boolean, boolean>; data: Ref<PaginatedResponse<GalleryDetail>> } | void {
  if (!path) return

  try {
    const { loading, data } = useRequest(
      dataInstance.Get<PaginatedResponse<GalleryDetail>>(
        `/v2/gallery/details/${path}?page=${page}&page_size=40`
      )
    )

    return { loading, data }
  } catch (error: unknown) {
    throw error
  }
}

function getGalleryCollectionRaw() {
  return dataInstance.Get<GalleryCollection>('/gallery/collection')
}

function getGalleryCollection(): { loading: Ref<boolean, boolean>; data: Ref<GalleryCollection> } {
  try {
    const { loading, data } = useRequest(getGalleryCollectionRaw())

    return { loading, data }
  } catch (error: unknown) {
    throw error
  }
}

export {
  galleryDetailPath,
  getGalleryDirectoriesRaw,
  upsertGalleryCollectionRaw,
  deleteGalleryCollectionRaw,
  getGalleryDetailsRaw,
  upsertGalleryDetailRaw,
  deleteGalleryDetailRaw,
  loadMoreAssets,
  getGalleryDetail,
  getGalleryCollectionRaw,
  getGalleryCollection
}
