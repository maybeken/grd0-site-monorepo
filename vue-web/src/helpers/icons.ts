import { addCollection } from '@iconify/vue/offline'

import type { IconifyJSON } from '@iconify/vue/offline'

let pending: Promise<IconifyJSON> | null = null

async function fetchCollection(): Promise<IconifyJSON> {
  const response = await fetch(`${import.meta.env.BASE_URL}mynaui-icons.json`)

  if (!response.ok) {
    throw new Error(`Icon collection request failed with ${response.status}`)
  }

  const data: IconifyJSON = await response.json()

  addCollection(data)

  return data
}

export function loadIconCollection(): Promise<IconifyJSON> {
  if (!pending) {
    pending = fetchCollection().catch((error: unknown) => {
      pending = null
      throw error
    })
  }

  return pending
}

let pending_flags: Promise<IconifyJSON> | null = null
const flag_names = new Set<string>()
let flag_codes: string[] = []

async function fetchFlagCollection(): Promise<IconifyJSON> {
  const response = await fetch(`${import.meta.env.BASE_URL}circle-flags.json`)

  if (!response.ok) {
    throw new Error(`Flag collection request failed with ${response.status}`)
  }

  const data: IconifyJSON = await response.json()

  addCollection(data)

  for (const name of Object.keys(data.icons ?? {})) flag_names.add(name)
  for (const name of Object.keys(data.aliases ?? {})) flag_names.add(name)

  flag_codes = Array.from(flag_names)
    .filter((name) => /^[a-z]{2}$/.test(name))
    .sort()

  return data
}

export function hasFlag(code: string): boolean {
  return flag_names.has(code)
}

export function knownFlagCodes(): string[] {
  return flag_codes
}

export function loadFlagCollection(): Promise<IconifyJSON> {
  if (!pending_flags) {
    pending_flags = fetchFlagCollection().catch((error: unknown) => {
      pending_flags = null
      throw error
    })
  }

  return pending_flags
}
