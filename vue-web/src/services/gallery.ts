import { useRequest } from 'alova/client'
import { dataInstance } from './api'

import type { Ref } from 'vue'
import type { GalleryDetail, GalleryCollection, PaginatedResponse, Asset } from '@/interfaces/Gallery'

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

export { loadMoreAssets, getGalleryDetail, getGalleryCollectionRaw, getGalleryCollection }
