export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

export interface Asset {
  filename: string
  exif?: {
    datetime?: string
    shutter?: string
    fstop?: number
    iso?: number
    focal?: number
    equipment?: {
      camera?: string
      lens?: string
    }
  }
  collection?: string
}

export interface AssetFileList {
  [key: string]: Asset[]
}

export interface GalleryDetail {
  path: string
  filename: string
  tz_adjustment?: number
  description?: string
}

export interface GalleryCollection {
  [key: string]: {
    title?: string
    cover?: string
  }
}
