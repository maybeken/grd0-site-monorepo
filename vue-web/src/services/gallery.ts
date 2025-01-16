import { useRequest } from 'alova';
import { dataInstance, assetsInstance } from './api';

import type { Ref } from 'vue';
import type { AssetFileList, GalleryDetail, GalleryCategory } from '@/interfaces/Gallery';

function listAssets(): Ref<AssetFileList> {
  try {
    const { data } = useRequest(assetsInstance.Get<AssetFileList>('/files.json'));

    return data;
  } catch(error: unknown) {
    throw error;
  }
}

function listGallery(): Ref<GalleryDetail[]> {
  try {
    const { data } = useRequest(dataInstance.Get<GalleryDetail[]>('/gallery'));

    return data;
  } catch(error: unknown) {
    throw error;
  }
}

function getGalleryDetail(path: string): Ref<GalleryDetail> | void {
  if (!path) return;
  
  try {
    const { data } = useRequest(dataInstance.Get<GalleryDetail>(`/gallery/${path}`));

    return data;
  } catch(error: unknown) {
    throw error;
  }
}

function getGalleryCategory(): Ref<GalleryCategory> {
  try {
    const { data } = useRequest(dataInstance.Get<GalleryCategory>('/gallery/category'));

    return data;
  } catch(error: unknown) {
    throw error;
  }
}

export { listAssets, listGallery, getGalleryDetail, getGalleryCategory };