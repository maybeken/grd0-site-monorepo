import { hasFlag, knownFlagCodes } from '@/helpers/icons'

const ASSET_URL = import.meta.env.VITE_ASSETS_URL

export const CUTOUT_MODEL_HOST =
  import.meta.env.VITE_CUTOUT_MODEL_HOST || `${ASSET_URL}/models/`
export const CUTOUT_MODEL_ID = import.meta.env.VITE_CUTOUT_MODEL_ID || 'ormbg'

export interface CountryOption {
  code: string
  name: string
}

let country_options: CountryOption[] | null = null

function display_name(code: string): string {
  if (typeof Intl.DisplayNames === 'undefined') return code

  try {
    const names = new Intl.DisplayNames(['en'], { type: 'region' })

    return names.of(code) ?? code
  } catch {
    return code
  }
}

export function countryOptions(): CountryOption[] {
  if (country_options) return country_options

  const codes = knownFlagCodes()

  if (codes.length === 0) return []

  country_options = codes
    .map((code) => {
      const upper = code.toUpperCase()

      return { code: upper, name: display_name(upper) }
    })
    .filter((option) => option.name !== option.code)
    .sort((left, right) => left.name.localeCompare(right.name))

  return country_options
}

export const BODY_TYPES = [
  { value: 'sedan', label: 'grandTour.bodyTypes.sedan' },
  { value: 'hatchback', label: 'grandTour.bodyTypes.hatchback' },
  { value: 'suv', label: 'grandTour.bodyTypes.suv' },
  { value: 'van', label: 'grandTour.bodyTypes.van' },
  { value: 'pickup', label: 'grandTour.bodyTypes.pickup' },
  { value: 'motorbike', label: 'grandTour.bodyTypes.motorbike' },
  { value: 'other', label: 'grandTour.bodyTypes.other' }
] as const

export function assetUrl(key: string): string {
  return `${ASSET_URL}/${key.replace(/^\/+/, '')}`
}

export function flagIcon(country_code?: string | null): string | null {
  if (!country_code) return null

  const code = country_code.trim().toLowerCase()

  if (!/^[a-z]{2}$/.test(code)) return null

  const name = `circle-flags:${code}`

  return hasFlag(code) ? name : null
}
