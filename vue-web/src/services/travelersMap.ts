import { useRequest } from 'alova/client'
import { dataInstance } from './api'

import type { Ref } from 'vue'
import type { MapLocation } from '@/interfaces/TravelersMap'

function getMapLocation(): { loading: Ref<boolean, boolean>; data: Ref<MapLocation[]> } {
  try {
    const { loading, data } = useRequest(dataInstance.Get<MapLocation[]>('/travel/map'))

    return { loading, data }
  } catch (error: unknown) {
    throw error
  }
}

export { getMapLocation }
