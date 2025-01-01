export interface Asset {
  filename: string,
  exif?: {
    datetime?: string,
    shutter?: string,
    fstop?: number,
    iso?: number,
    focal?: number,
    equipment?: {
      camera?: string,
      lens?: string,
    }
  },
  category?: string,
};
  
export interface AssetFileList {
  [key: string]: Asset[],
}

export interface GalleryDetail {
  path: string,
  filename: string,
  tz_adjustment?: number,
  description?: string,
}

export interface GalleryCategory {
  [key: string]: string,
}