import { useRequest } from 'alova/client';
import { dataInstance } from './api';

import type { Ref } from 'vue';
import type { Asset, GalleryDetail, GalleryCollection } from '@/interfaces/Gallery';

function listAssets(collection?: string): { loading: Ref<boolean, boolean>, data: Ref<Asset[]> } | void {
  if (!collection) return;

  try {
    const { loading, data } = useRequest(dataInstance.Get<Asset[]>(`/gallery/${collection}`));

    return { loading, data };
  } catch(error: unknown) {
    throw error;
  }
}

function getGalleryDetail(path?: string): { loading: Ref<boolean, boolean>, data: Ref<GalleryDetail[]> } | void {
  if (!path) return;
  
  try {
    const { loading, data } = useRequest(dataInstance.Get<GalleryDetail[]>(`/gallery/collection/${path}`));

    return { loading, data };
  } catch(error: unknown) {
    throw error;
  }
}

function getGalleryCollectionRaw() {
  return dataInstance.Get<GalleryCollection>('/gallery/collection');
}

function getGalleryCollection(): { loading: Ref<boolean, boolean>, data: Ref<GalleryCollection> } {
  try {
    const { loading, data } = useRequest(getGalleryCollectionRaw());

    return { loading, data };
  } catch(error: unknown) {
    throw error;
  }
}

export { listAssets, getGalleryDetail, getGalleryCollectionRaw, getGalleryCollection };