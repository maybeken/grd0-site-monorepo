import { useRequest } from 'alova/client'
import { dataInstance } from './api'

import type { Ref } from 'vue'
import type { Music } from '@/interfaces/Music'

function getMusic(): { loading: Ref<boolean, boolean>; data: Ref<Music[]> } {
  try {
    const { loading, data } = useRequest(dataInstance.Get<Music[]>('/music'))

    return { loading, data }
  } catch (error: unknown) {
    throw error
  }
}

export { getMusic }
