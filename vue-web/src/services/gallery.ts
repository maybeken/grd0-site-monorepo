import { useRequest } from 'alova/client';
import { dataInstance } from './api';

import type { Ref } from 'vue';
import type { Asset, GalleryDetail, GalleryCategory } from '@/interfaces/Gallery';

function listAssets(category?: string): { loading: Ref<boolean, boolean>, data: Ref<Asset[]> } | void {
  if (!category) return;

  try {
    const { loading, data } = useRequest(dataInstance.Get<Asset[]>(`/gallery/${category}`));

    return { loading, data };
  } catch(error: unknown) {
    throw error;
  }
}

function getGalleryDetail(path?: string): { loading: Ref<boolean, boolean>, data: Ref<GalleryDetail[]> } | void {
  if (!path) return;
  
  try {
    const { loading, data } = useRequest(dataInstance.Get<GalleryDetail[]>(`/gallery/category/${path}`));

    return { loading, data };
  } catch(error: unknown) {
    throw error;
  }
}

function getGalleryCategoryRaw() {
  return dataInstance.Get<GalleryCategory>('/gallery/category');
}

function getGalleryCategory(): { loading: Ref<boolean, boolean>, data: Ref<GalleryCategory> } {
  try {
    const { loading, data } = useRequest(getGalleryCategoryRaw());

    return { loading, data };
  } catch(error: unknown) {
    throw error;
  }
}

export { listAssets, getGalleryDetail, getGalleryCategoryRaw, getGalleryCategory };