import { useRequest } from 'alova/client'
import { dataInstance, adminInstance } from './api'

import type { Ref } from 'vue'
import type { MapLocation } from '@/interfaces/TravelersMap'

function getPublishedMapLocationsRaw() {
  return dataInstance.Get<MapLocation[]>('/travel/map')
}

function getAllMapLocationsRaw() {
  return adminInstance.Get<MapLocation[]>('/travel/map/all')
}

function getMapLocation(): { loading: Ref<boolean, boolean>; data: Ref<MapLocation[]> } {
  try {
    const { loading, data } = useRequest(getPublishedMapLocationsRaw())

    return { loading, data }
  } catch (error: unknown) {
    throw error
  }
}

async function getMapLocationBySlug(slug: string): Promise<MapLocation | null> {
  const locations = await getAllMapLocationsRaw()

  if (!Array.isArray(locations)) {
    return null
  }

  return locations.find((location) => location.slug === slug) ?? null
}

function toKebabCase(title: string): string {
  return title
    .normalize('NFKD')
    .replace(/[\u0300-\u036f]/g, '')
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
}

function upsertMapLocation(location: MapLocation) {
  return adminInstance.Put<MapLocation>('/travel/map', location)
}

function deleteMapLocation(slug: string) {
  return adminInstance.Delete(`/travel/map/${encodeURIComponent(slug)}`)
}

export {
  getPublishedMapLocationsRaw,
  getAllMapLocationsRaw,
  getMapLocation,
  getMapLocationBySlug,
  toKebabCase,
  upsertMapLocation,
  deleteMapLocation
}
