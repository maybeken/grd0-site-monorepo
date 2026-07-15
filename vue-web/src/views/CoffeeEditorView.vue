<template>
  <div class="flex flex-col gap-6 max-w-4xl mx-auto">
    <div class="flex gap-2 border-b border-accent/20 pb-2">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        @click="activeTab = tab.key"
        :class="[
          'px-4 py-2 rounded-t text-sm font-semibold transition-colors',
          activeTab === tab.key
            ? 'bg-accent text-background'
            : 'text-foreground/60 hover:text-foreground hover:bg-accent/10'
        ]"
      >
        {{ tab.label }}
      </button>
    </div>

    <div v-if="activeTab === 'beans'">
      <div v-if="!showBeanForm" class="flex justify-end mb-4">
        <button
          @click="openNewBean"
          class="px-3 py-1 rounded bg-accent text-background text-sm font-semibold hover:brightness-110"
        >
          + New Bean
        </button>
      </div>
      <BeanForm
        v-if="showBeanForm"
        :bean="editingBean"
        @saved="onBeanSaved"
        @cancel="closeBeanForm"
      />
      <div v-if="!showBeanForm" class="flex flex-col gap-2">
        <div
          v-for="bean in beans"
          :key="bean.id"
          class="flex items-center justify-between bg-background rounded-lg p-3 border border-accent/20"
        >
          <div>
            <span class="font-semibold text-sm">{{ bean.name }}</span>
            <span v-if="bean.roaster" class="text-xs text-foreground/60 ml-2">{{
              bean.roaster
            }}</span>
            <span v-if="bean.origin" class="text-xs text-foreground/40 ml-2">{{
              bean.origin
            }}</span>
          </div>
          <div class="flex gap-2">
            <button
              @click="editBean(bean)"
              class="px-2 py-1 rounded text-xs border border-accent/40 hover:bg-accent/10"
            >
              Edit
            </button>
            <button
              @click="deleteBean(bean)"
              class="px-2 py-1 rounded text-xs border border-red-500/40 text-red-400 hover:bg-red-500/10"
            >
              Delete
            </button>
          </div>
        </div>
        <p v-if="!beans.length" class="text-sm text-foreground/40 text-center py-4">No beans yet</p>
      </div>
    </div>

    <div v-if="activeTab === 'equipment'">
      <div v-if="!showEquipmentForm" class="flex justify-end mb-4">
        <button
          @click="openNewEquipment"
          class="px-3 py-1 rounded bg-accent text-background text-sm font-semibold hover:brightness-110"
        >
          + New Equipment
        </button>
      </div>
      <EquipmentForm
        v-if="showEquipmentForm"
        :equipment="editingEquipment"
        @saved="onEquipmentSaved"
        @cancel="closeEquipmentForm"
      />
      <div v-if="!showEquipmentForm" class="flex flex-col gap-2">
        <div
          v-for="eq in equipment"
          :key="eq.id"
          class="flex items-center justify-between bg-background rounded-lg p-3 border border-accent/20"
        >
          <div>
            <span class="font-semibold text-sm">{{ eq.name }}</span>
            <span v-if="eq.type" class="text-xs text-foreground/60 ml-2">{{ eq.type }}</span>
          </div>
          <div class="flex gap-2">
            <button
              @click="editEquipment(eq)"
              class="px-2 py-1 rounded text-xs border border-accent/40 hover:bg-accent/10"
            >
              Edit
            </button>
            <button
              @click="deleteEquipment(eq)"
              class="px-2 py-1 rounded text-xs border border-red-500/40 text-red-400 hover:bg-red-500/10"
            >
              Delete
            </button>
          </div>
        </div>
        <p v-if="!equipment.length" class="text-sm text-foreground/40 text-center py-4">
          No equipment yet
        </p>
      </div>
    </div>

    <div v-if="activeTab === 'tasting'">
      <div v-if="!showTastingForm" class="flex justify-end mb-4">
        <button
          @click="openNewTasting"
          class="px-3 py-1 rounded bg-accent text-background text-sm font-semibold hover:brightness-110"
        >
          + New Tasting
        </button>
      </div>
      <TastingForm
        v-if="showTastingForm"
        :tasting="editingTasting"
        :beans="beans"
        :equipment-list="equipment"
        @saved="onTastingSaved"
        @cancel="closeTastingForm"
        @bean-added="onBeanAddedFromTasting"
        @equipment-added="onEquipmentAddedFromTasting"
      />
      <div v-if="!showTastingForm" class="flex flex-col gap-2">
        <div
          v-for="tasting in tastings"
          :key="tasting.id"
          class="flex items-center justify-between bg-background rounded-lg p-3 border border-accent/20"
        >
          <div>
            <span class="font-semibold text-sm">{{ tasting.bean?.name || 'Unknown bean' }}</span>
            <span v-if="tasting.equipment?.name" class="text-xs text-foreground/60 ml-2">{{
              tasting.equipment.name
            }}</span>
            <span v-if="tasting.tasted_at" class="text-xs text-foreground/40 ml-2">{{
              tasting.tasted_at.split('T')[0]
            }}</span>
            <span v-if="tasting.rating != null" class="text-xs text-foreground/60 ml-2">
              ★ {{ tasting.rating }}
            </span>
          </div>
          <div class="flex gap-2">
            <button
              @click="editTasting(tasting)"
              class="px-2 py-1 rounded text-xs border border-accent/40 hover:bg-accent/10"
            >
              Edit
            </button>
            <button
              @click="deleteTasting(tasting)"
              class="px-2 py-1 rounded text-xs border border-red-500/40 text-red-400 hover:bg-red-500/10"
            >
              Delete
            </button>
          </div>
        </div>
        <p v-if="!tastings.length" class="text-sm text-foreground/40 text-center py-4">
          No tastings yet
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { CoffeeBean, BrewEquipment, TastingNote } from '@/interfaces/Coffee'
import {
  getCoffeeBeansRaw,
  getCoffeeEquipmentRaw,
  getCoffeeTastingsRaw,
  deleteCoffeeBeanRaw,
  deleteCoffeeEquipmentRaw,
  deleteCoffeeTastingRaw
} from '@/services/coffee'
import BeanForm from '@/components/coffee/BeanForm.vue'
import EquipmentForm from '@/components/coffee/EquipmentForm.vue'
import TastingForm from '@/components/coffee/TastingForm.vue'

const tabs = [
  { key: 'beans', label: 'Beans' },
  { key: 'equipment', label: 'Equipment' },
  { key: 'tasting', label: 'Tasting' }
] as const

const activeTab = ref<(typeof tabs)[number]['key']>('beans')

const beans = ref<CoffeeBean[]>([])
const equipment = ref<BrewEquipment[]>([])
const tastings = ref<TastingNote[]>([])

const showBeanForm = ref(false)
const editingBean = ref<CoffeeBean | null>(null)

const showEquipmentForm = ref(false)
const editingEquipment = ref<BrewEquipment | null>(null)

const showTastingForm = ref(false)
const editingTasting = ref<TastingNote | null>(null)

async function fetchBeans() {
  beans.value = await getCoffeeBeansRaw()
}

async function fetchEquipment() {
  equipment.value = await getCoffeeEquipmentRaw()
}

async function fetchTastings() {
  tastings.value = await getCoffeeTastingsRaw()
}

fetchBeans()
fetchEquipment()
fetchTastings()

function openNewBean() {
  editingBean.value = null
  showBeanForm.value = true
}

function editBean(bean: CoffeeBean) {
  editingBean.value = bean
  showBeanForm.value = true
}

function closeBeanForm() {
  showBeanForm.value = false
  editingBean.value = null
}

async function onBeanSaved() {
  closeBeanForm()
  await fetchBeans()
}

async function deleteBean(bean: CoffeeBean) {
  if (!bean.id || !confirm(`Delete "${bean.name}"?`)) return
  try {
    await deleteCoffeeBeanRaw(bean.id)
    await fetchBeans()
  } catch {
    alert('Failed to delete bean')
  }
}

function openNewEquipment() {
  editingEquipment.value = null
  showEquipmentForm.value = true
}

function editEquipment(eq: BrewEquipment) {
  editingEquipment.value = eq
  showEquipmentForm.value = true
}

function closeEquipmentForm() {
  showEquipmentForm.value = false
  editingEquipment.value = null
}

async function onEquipmentSaved() {
  closeEquipmentForm()
  await fetchEquipment()
}

async function deleteEquipment(eq: BrewEquipment) {
  if (!eq.id || !confirm(`Delete "${eq.name}"?`)) return
  try {
    await deleteCoffeeEquipmentRaw(eq.id)
    await fetchEquipment()
  } catch {
    alert('Failed to delete equipment')
  }
}

function openNewTasting() {
  editingTasting.value = null
  showTastingForm.value = true
}

function editTasting(tasting: TastingNote) {
  editingTasting.value = tasting
  showTastingForm.value = true
}

function closeTastingForm() {
  showTastingForm.value = false
  editingTasting.value = null
}

async function onTastingSaved() {
  closeTastingForm()
  await fetchTastings()
}

async function deleteTasting(tasting: TastingNote) {
  const label = tasting.bean?.name || tasting.tasted_at?.split('T')[0] || 'this tasting'
  if (!tasting.id || !confirm(`Delete tasting "${label}"?`)) return
  try {
    await deleteCoffeeTastingRaw(tasting.id)
    await fetchTastings()
  } catch {
    alert('Failed to delete tasting')
  }
}

async function onBeanAddedFromTasting() {
  await fetchBeans()
}

async function onEquipmentAddedFromTasting() {
  await fetchEquipment()
}
</script>
