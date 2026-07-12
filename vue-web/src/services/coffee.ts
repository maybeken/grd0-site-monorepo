import { useRequest } from 'alova/client'
import { dataInstance, adminInstance } from './api'

import type { Ref } from 'vue'
import type { CoffeeBean, BrewEquipment, TastingNote } from '@/interfaces/Coffee'

function getCoffeeBeansRaw() {
  return dataInstance.Get<CoffeeBean[]>('/coffee/beans')
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
  upsertCoffeeBeanRaw,
  deleteCoffeeBeanRaw,
  getCoffeeEquipmentRaw,
  getCoffeeEquipment,
  upsertCoffeeEquipmentRaw,
  deleteCoffeeEquipmentRaw,
  getCoffeeTastingsRaw,
  getCoffeeTastings,
  upsertCoffeeTastingRaw,
  deleteCoffeeTastingRaw
}
