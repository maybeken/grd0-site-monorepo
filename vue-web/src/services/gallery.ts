import { useRequest } from 'alova';
import { assetsInstance } from './api';

import type { Ref } from 'vue';
import type { AssetFileList } from '@/interfaces/Gallery';

function listAssets(): Ref<AssetFileList> {
  try {
    const { data } = useRequest(assetsInstance.Get<AssetFileList>('/files.json'));

    return data;
  } catch(error: unknown) {
    throw error;
  }
}

export { listAssets };