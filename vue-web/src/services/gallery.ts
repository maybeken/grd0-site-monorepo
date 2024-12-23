import { useRequest } from 'alova';
import { alovaInstance } from './api';

import type { Ref } from 'vue';
import type { AssetFileList } from '@/interfaces/Gallery';

function listAssets(): Ref<AssetFileList> {
  try {
    const { data } = useRequest(alovaInstance.Get<AssetFileList>('/files.json'));

    return data;
  } catch(error: unknown) {
    throw error;
  }
}

export { listAssets };