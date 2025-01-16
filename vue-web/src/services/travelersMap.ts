import { useRequest } from 'alova';
import { dataInstance } from './api';

import type { Ref } from 'vue';
import type { MapLocation } from '@/interfaces/TravelersMap';

function getMapLocation(): Ref<MapLocation[]> {
  try {
    const { data } = useRequest(dataInstance.Get<MapLocation[]>('/travel/map'));

    return data;
  } catch(error: unknown) {
    throw error;
  }
}

export { getMapLocation };