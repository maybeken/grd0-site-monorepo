const WEB_MERCATOR_RADIUS = 20037508.34
const MAX_RESOLUTION = 156543.03392

export function epsg4326toEpsg3857(coordinates: number[]): number[] {
  const latitude = coordinates[0]
  const longitude = coordinates[1]

  if (typeof latitude !== 'number' || typeof longitude !== 'number') {
    return [0, 0]
  }

  if (!Number.isFinite(latitude) || !Number.isFinite(longitude)) {
    return [0, 0]
  }

  const x = (longitude * WEB_MERCATOR_RADIUS) / 180
  let y = Math.log(Math.tan(((90 + latitude) * Math.PI) / 360)) / (Math.PI / 180)
  y = (y * WEB_MERCATOR_RADIUS) / 180

  return [x, y]
}

export function epsg3857toEpsg4326(coordinates: number[]): number[] {
  const x = coordinates[0]
  const y = coordinates[1]

  if (typeof x !== 'number' || typeof y !== 'number') {
    return [0, 0]
  }

  if (!Number.isFinite(x) || !Number.isFinite(y)) {
    return [0, 0]
  }

  const longitude = (x * 180) / WEB_MERCATOR_RADIUS
  const latitude =
    (Math.atan(Math.pow(Math.E, (y * 180 * (Math.PI / 180)) / WEB_MERCATOR_RADIUS)) * 360) /
      Math.PI -
    90

  return [latitude, longitude]
}

export function zoomToResolution(zoom: number): number {
  return MAX_RESOLUTION / Math.pow(2, zoom)
}

export function resolutionToZoom(resolution: number): number {
  return Math.log2(MAX_RESOLUTION / resolution)
}

export function formatZoom(resolution?: number | null): string {
  if (resolution === undefined || resolution === null || resolution <= 0) {
    return ''
  }

  return resolutionToZoom(resolution).toFixed(2)
}

export function parseZoom(value: string | number | null | undefined): number | null {
  if (value === null || value === undefined) {
    return null
  }

  if (typeof value === 'number') {
    return Number.isFinite(value) ? value : null
  }

  if (value.trim() === '') {
    return null
  }

  const zoom = Number(value)

  return Number.isFinite(zoom) ? zoom : null
}

export function isValidLatitude(value: number): boolean {
  return Number.isFinite(value) && value >= -90 && value <= 90
}

export function isValidLongitude(value: number): boolean {
  return Number.isFinite(value) && value >= -180 && value <= 180
}
