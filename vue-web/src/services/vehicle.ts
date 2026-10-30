import { useRequest } from 'alova/client'
import { dataInstance, adminInstance } from './api'

import type { Ref } from 'vue'
import type { Vehicle, VehicleUpload } from '@/interfaces/Vehicle'

function getVehiclesRaw() {
  return dataInstance.Get<Vehicle[]>('/vehicles')
}

async function getVehicleCollection(force = false): Promise<Vehicle[]> {
  const value = await getVehiclesRaw().send(force)

  if (!Array.isArray(value)) {
    throw new Error('Unexpected vehicle collection response')
  }

  return value
}

function getVehicles(): { loading: Ref<boolean, boolean>; data: Ref<Vehicle[]> } {
  try {
    const { loading, data } = useRequest(getVehiclesRaw())
    return { loading, data }
  } catch (error: unknown) {
    throw error
  }
}

async function getVehicleRaw(id: string): Promise<Vehicle | null> {
  const value = await dataInstance.Get<Vehicle>(`/vehicles/${id}`)

  return value && typeof value === 'object' && 'id' in value ? value : null
}

function upsertVehicleRaw(vehicle: Vehicle) {
  return adminInstance.Put<Vehicle>('/vehicles', vehicle)
}

function deleteVehicleRaw(id: string) {
  return adminInstance.Delete(`/vehicles/${id}`)
}

function requestVehicleAttachmentRaw(key: string) {
  return adminInstance.Put<VehicleUpload>(`/vehicles/attachment/${encodeURIComponent(key)}`)
}

async function uploadVehicleCover(upload: VehicleUpload, blob: Blob): Promise<void> {
  const response = await fetch(upload.url, {
    method: 'PUT',
    body: blob,
    headers: { 'Content-Type': blob.type || 'image/png' }
  })

  if (!response.ok) {
    throw new Error(`Upload failed with status ${response.status}`)
  }
}

export {
  getVehiclesRaw,
  getVehicleCollection,
  getVehicles,
  getVehicleRaw,
  upsertVehicleRaw,
  deleteVehicleRaw,
  requestVehicleAttachmentRaw,
  uploadVehicleCover
}
