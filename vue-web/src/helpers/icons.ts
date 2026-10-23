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
