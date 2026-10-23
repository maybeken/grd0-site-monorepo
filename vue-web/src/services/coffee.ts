import { useRequest } from 'alova/client'
import { dataInstance, adminInstance } from './api'

import type { Ref } from 'vue'
import type { CoffeeBean, BrewEquipment, TastingNote } from '@/interfaces/Coffee'

function getCoffeeBeansRaw() {
  return dataInstance.Get<CoffeeBean[]>('/coffee/beans')
}

async function getCoffeeBeanRaw(id: string): Promise<CoffeeBean | null> {
  const value = await dataInstance.Get<CoffeeBean>(`/coffee/bean/${id}`)

  return value && typeof value === 'object' && 'id' in value ? value : null
}

function getCoffeeBeans(): { loading: Ref<boolean, boolean>; data: Ref<CoffeeBean[]> } {
  try {
    const { loading, data } = useRequest(getCoffeeBeansRaw())
    return { loading, data }
  } catch (error: unknown) {
    throw error
  }
}

function upsertCoffeeBeanRaw(bean: CoffeeBean) {
  return adminInstance.Put<CoffeeBean>('/coffee/bean', bean)
}

function deleteCoffeeBeanRaw(id: string) {
  return adminInstance.Delete(`/coffee/bean/${id}`)
}

function getCoffeeEquipmentRaw() {
  return dataInstance.Get<BrewEquipment[]>('/coffee/equipment')
}

async function getCoffeeEquipmentItemRaw(id: string): Promise<BrewEquipment | null> {
  const value = await dataInstance.Get<BrewEquipment>(`/coffee/equipment/${id}`)

  return value && typeof value === 'object' && 'id' in value ? value : null
}

function getCoffeeEquipment(): { loading: Ref<boolean, boolean>; data: Ref<BrewEquipment[]> } {
  try {
    const { loading, data } = useRequest(getCoffeeEquipmentRaw())
    return { loading, data }
  } catch (error: unknown) {
    throw error
  }
}

function upsertCoffeeEquipmentRaw(equipment: BrewEquipment) {
  return adminInstance.Put<BrewEquipment>('/coffee/equipment', equipment)
}

function deleteCoffeeEquipmentRaw(id: string) {
  return adminInstance.Delete(`/coffee/equipment/${id}`)
}

function getCoffeeTastingsRaw() {
  return dataInstance.Get<TastingNote[]>('/coffee/tastings')
}

async function getCoffeeTastingRaw(id: string): Promise<TastingNote | null> {
  const value = await dataInstance.Get<TastingNote>(`/coffee/tasting/${id}`)

  return value && typeof value === 'object' && 'id' in value ? value : null
}

function getCoffeeTastings(): { loading: Ref<boolean, boolean>; data: Ref<TastingNote[]> } {
  try {
    const { loading, data } = useRequest(getCoffeeTastingsRaw())
    return { loading, data }
  } catch (error: unknown) {
    throw error
  }
}

function upsertCoffeeTastingRaw(tasting: TastingNote) {
  return adminInstance.Put<TastingNote>('/coffee/tasting', tasting)
}

function deleteCoffeeTastingRaw(id: string) {
  return adminInstance.Delete(`/coffee/tasting/${id}`)
}

export {
  getCoffeeBeansRaw,
  getCoffeeBeans,
  getCoffeeBeanRaw,
  upsertCoffeeBeanRaw,
  deleteCoffeeBeanRaw,
  getCoffeeEquipmentRaw,
  getCoffeeEquipment,
  getCoffeeEquipmentItemRaw,
  upsertCoffeeEquipmentRaw,
  deleteCoffeeEquipmentRaw,
  getCoffeeTastingsRaw,
  getCoffeeTastings,
  getCoffeeTastingRaw,
  upsertCoffeeTastingRaw,
  deleteCoffeeTastingRaw
}
