<template>
  <div class="flex flex-col gap-4 bg-background rounded-xl p-4 border border-accent/20">
    <div class="flex items-center justify-between">
      <h3 class="text-lg font-bold">{{ isEdit ? 'Edit Tasting' : 'New Tasting' }}</h3>
      <div class="flex gap-2">
        <button
          @click="$emit('cancel')"
          class="px-3 py-1 rounded border border-accent/40 text-sm hover:bg-accent/10"
        >
          Cancel
        </button>
        <button
          @click="save"
          :disabled="saving"
          class="px-3 py-1 rounded bg-accent text-background text-sm font-semibold hover:brightness-110 disabled:opacity-50"
        >
          {{ saving ? 'Saving...' : 'Save' }}
        </button>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="flex flex-col gap-1">
        <div class="flex items-center justify-between">
          <label class="text-sm font-semibold">Bean</label>
          <button
            v-if="!showQuickAddBean"
            type="button"
            @click="showQuickAddBean = true"
            class="text-xs text-accent hover:underline"
          >
            + new
          </button>
        </div>
        <select
          v-model="form.bean.id"
          class="bg-background border border-accent/40 rounded px-2 py-1 text-sm"
        >
          <option value="">Select bean...</option>
          <option v-for="b in beans" :key="b.id" :value="b.id">
            {{ b.roaster }} {{ b.name }} ({{ b.roast_date?.split('T')[0]?.replace(/-/g, '/') }})
          </option>
        </select>
        <div
          v-if="showQuickAddBean"
          class="flex flex-col gap-2 p-2 rounded border border-accent/20 bg-accent/5"
        >
          <input
            v-model="quickBean.name"
            placeholder="Bean name *"
            class="bg-background border border-accent/40 rounded px-2 py-1 text-xs"
          />
          <div class="grid grid-cols-2 gap-2">
            <input
              v-model="quickBean.roaster"
              placeholder="Roaster"
              class="bg-background border border-accent/40 rounded px-2 py-1 text-xs"
            />
            <input
              v-model="quickBean.origin"
              placeholder="Origin"
              class="bg-background border border-accent/40 rounded px-2 py-1 text-xs"
            />
          </div>
          <div class="flex gap-2 justify-end">
            <button
              type="button"
              @click="cancelQuickBean"
              class="px-2 py-0.5 rounded text-xs border border-accent/40 hover:bg-accent/10"
            >
              Cancel
            </button>
            <button
              type="button"
              :disabled="quickBeanSaving || !quickBean.name.trim()"
              @click="saveQuickBean"
              class="px-2 py-0.5 rounded text-xs bg-accent text-background font-semibold hover:brightness-110 disabled:opacity-50"
            >
              {{ quickBeanSaving ? '...' : 'Add' }}
            </button>
          </div>
        </div>
      </div>

      <div class="flex flex-col gap-1">
        <div class="flex items-center justify-between">
          <label class="text-sm font-semibold">Equipment</label>
          <button
            v-if="!showQuickAddEquipment"
            type="button"
            @click="showQuickAddEquipment = true"
            class="text-xs text-accent hover:underline"
          >
            + new
          </button>
        </div>
        <select
          v-model="form.equipment.id"
          class="bg-background border border-accent/40 rounded px-2 py-1 text-sm"
        >
          <option value="">Select equipment...</option>
          <option
            v-for="e in (equipmentList || []).filter((e) => e.type !== 'Grinder')"
            :key="e.id"
            :value="e.id"
          >
            {{ e.name }}
          </option>
        </select>
        <div
          v-if="showQuickAddEquipment"
          class="flex flex-col gap-2 p-2 rounded border border-accent/20 bg-accent/5"
        >
          <input
            v-model="quickEquipment.name"
            placeholder="Equipment name *"
            class="bg-background border border-accent/40 rounded px-2 py-1 text-xs"
          />
          <select
            v-model="quickEquipment.type"
            class="bg-background border border-accent/40 rounded px-2 py-1 text-xs"
          >
            <option :value="null">Select type...</option>
            <option v-for="t in equipmentTypes" :key="t" :value="t">{{ t }}</option>
          </select>
          <div class="flex gap-2 justify-end">
            <button
              type="button"
              @click="cancelQuickEquipment"
              class="px-2 py-0.5 rounded text-xs border border-accent/40 hover:bg-accent/10"
            >
              Cancel
            </button>
            <button
              type="button"
              :disabled="quickEquipmentSaving || !quickEquipment.name.trim()"
              @click="saveQuickEquipment"
              class="px-2 py-0.5 rounded text-xs bg-accent text-background font-semibold hover:brightness-110 disabled:opacity-50"
            >
              {{ quickEquipmentSaving ? '...' : 'Add' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="flex flex-col gap-1">
        <label class="text-sm font-semibold">Date</label>
        <input
          type="date"
          v-model="form.tasted_at"
          class="bg-background border border-accent/40 rounded px-2 py-1 text-sm"
        />
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-sm font-semibold">Rating (1-10)</label>
        <div class="flex items-center gap-2">
          <input
            type="number"
            min="1"
            max="10"
            v-model.number="form.rating"
            :disabled="idkFields.rating"
            class="bg-background border border-accent/40 rounded px-2 py-1 text-sm w-20 disabled:opacity-50"
          />
        </div>
      </div>
    </div>

    <fieldset class="border border-accent/20 rounded p-3">
      <legend class="text-sm font-bold px-1">Brew Recipe</legend>

      <div class="grid grid-cols-2 md:grid-cols-3 gap-3">
        <div class="flex flex-col gap-1">
          <label class="text-xs">Grinder</label>
          <div class="flex items-center gap-1">
            <select
              v-model="form.grinder!.id"
              :disabled="idkFields.grinder"
              class="bg-background border border-accent/40 rounded px-1 py-0.5 text-xs flex-1 disabled:opacity-50"
            >
              <option value="">Select equipment...</option>
              <option
                v-for="e in (equipmentList || []).filter((e) => e.type === 'Grinder')"
                :key="e.id"
                :value="e.id"
              >
                {{ e.name }}
              </option>
            </select>
            <IdkToggle v-model="idkFields.grinder" />
          </div>
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-xs">Grind Size</label>
          <div class="flex items-center gap-1">
            <select
              v-model="form.grind_size"
              :disabled="idkFields.grind_size"
              class="bg-background border border-accent/40 rounded px-1 py-0.5 text-xs flex-1 disabled:opacity-50"
            >
              <option value="">Select...</option>
              <option v-for="g in GRIND_SIZES" :key="g" :value="g">{{ g }}</option>
            </select>
            <IdkToggle v-model="idkFields.grind_size" />
          </div>
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-xs">Grind Setting</label>
          <div class="flex items-center gap-1">
            <input
              type="number"
              step="0.1"
              v-model.number="form.grind_setting"
              :disabled="idkFields.grind_setting"
              class="bg-background border border-accent/40 rounded px-1 py-0.5 text-xs w-16 disabled:opacity-50"
            />
            <IdkToggle v-model="idkFields.grind_setting" />
          </div>
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-xs">Coffee Dose (g)</label>
          <div class="flex items-center gap-1">
            <input
              type="number"
              v-model.number="form.coffee_dose"
              :disabled="idkFields.coffee_dose"
              class="bg-background border border-accent/40 rounded px-1 py-0.5 text-xs w-16 disabled:opacity-50"
            />
            <IdkToggle v-model="idkFields.coffee_dose" />
          </div>
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-xs">Water In (ml)</label>
          <div class="flex items-center gap-1">
            <input
              type="number"
              v-model.number="form.water_in"
              :disabled="idkFields.water_in"
              class="bg-background border border-accent/40 rounded px-1 py-0.5 text-xs w-16 disabled:opacity-50"
            />
            <IdkToggle v-model="idkFields.water_in" />
          </div>
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-xs">Coffee Out (ml)</label>
          <div class="flex items-center gap-1">
            <input
              type="number"
              v-model.number="form.coffee_out"
              :disabled="idkFields.coffee_out"
              class="bg-background border border-accent/40 rounded px-1 py-0.5 text-xs w-16 disabled:opacity-50"
            />
            <IdkToggle v-model="idkFields.coffee_out" />
          </div>
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-xs">Ratio</label>
          <div class="flex items-center gap-1">
            <span class="text-xs">{{
              computedRatio != null ? '1:' + computedRatio.toFixed(1) : 'IDK'
            }}</span>
          </div>
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-xs">Brew Time</label>
          <div class="flex items-center gap-1">
            <BrewTimeInput v-model="form.brew_time" :disabled="idkFields.brew_time" />
            <IdkToggle v-model="idkFields.brew_time" />
          </div>
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-xs">Water Temp (°C)</label>
          <div class="flex items-center gap-1">
            <input
              type="number"
              v-model.number="form.water_temperature"
              :disabled="idkFields.water_temperature"
              class="bg-background border border-accent/40 rounded px-1 py-0.5 text-xs w-16 disabled:opacity-50"
            />
            <IdkToggle v-model="idkFields.water_temperature" />
          </div>
        </div>
      </div>
    </fieldset>

    <fieldset class="border border-accent/20 rounded p-3">
      <legend class="text-sm font-bold px-1">Taste Profile (0-10)</legend>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
        <div v-for="dim in tasteDimensions" :key="dim.key" class="flex flex-col gap-1">
          <label class="text-xs capitalize">{{ dim.label }}</label>
          <div class="flex items-center gap-1">
            <input
              type="number"
              min="0"
              max="10"
              v-model.number="(form as any)[dim.key]"
              :disabled="(idkFields as any)[dim.key]"
              class="bg-background border border-accent/40 rounded px-1 py-0.5 text-xs w-14 disabled:opacity-50"
            />
          </div>
        </div>
      </div>
    </fieldset>

    <div class="flex flex-col gap-1">
      <label class="text-sm font-semibold">Notes</label>
      <textarea
        v-model="form.overall_notes"
        rows="3"
        class="bg-background border border-accent/40 rounded px-2 py-1 text-sm resize-y"
      />
    </div>

    <div class="flex items-center gap-2">
      <input type="checkbox" v-model="form.pinned" class="accent-accent" />
      <label class="text-sm">Pin this tasting</label>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, reactive } from 'vue'
import type { CoffeeBean, BrewEquipment, TastingNote } from '@/interfaces/Coffee'
import { GRIND_SIZES, TASTE_DIMENSIONS, EQUIPMENT_TYPES } from '@/helpers/coffee'
import {
  upsertCoffeeTastingRaw,
  upsertCoffeeBeanRaw,
  upsertCoffeeEquipmentRaw
} from '@/services/coffee'
import IdkToggle from './IdkToggle.vue'
import BrewTimeInput from './BrewTimeInput.vue'

interface Props {
  tasting?: TastingNote | null
  beans: CoffeeBean[]
  equipmentList: BrewEquipment[]
}

const props = withDefaults(defineProps<Props>(), {
  tasting: null
})

const emit = defineEmits<{
  saved: []
  cancel: []
  beanAdded: [bean: CoffeeBean]
  equipmentAdded: [equipment: BrewEquipment]
}>()

const isEdit = computed(() => !!props.tasting?.id)
const saving = ref(false)

const showQuickAddBean = ref(false)
const quickBeanSaving = ref(false)
const quickBean = reactive({
  name: '',
  roaster: null as string | null,
  origin: null as string | null
})

const showQuickAddEquipment = ref(false)
const quickEquipmentSaving = ref(false)
const quickEquipment = reactive({
  name: '',
  type: null as string | null
})

const equipmentTypes = EQUIPMENT_TYPES

const tasteDimensions = TASTE_DIMENSIONS.map((d) => ({
  key: `taste_${d}`,
  label: d
}))

function createEmptyForm(): TastingNote {
  return {
    bean: { name: '' },
    equipment: { name: '' },
    tasted_at: new Date().toISOString().split('T')[0]!,
    pinned: false,
    grinder: { name: '' },
    grind_size: null,
    grind_setting: null,
    coffee_dose: null,
    water_in: null,
    coffee_out: null,
    ratio: null,
    brew_time: null,
    water_temperature: null,
    overall_notes: null,
    rating: null,
    ...Object.fromEntries(
      TASTE_DIMENSIONS.map((d) => [`taste_${d}`, null])
    ),
  }
}

const form = reactive<TastingNote>(createEmptyForm())

const idkFields = reactive<Record<string, boolean>>({
  grider: false,
  grind_size: false,
  grind_setting: false,
  coffee_dose: false,
  water_in: false,
  coffee_out: false,
  brew_time: false,
  water_temperature: false,
  rating: false
})

const computedRatio = computed(() => {
  if (form.water_in != null && form.coffee_dose != null && form.coffee_dose > 0) {
    return form.water_in / form.coffee_dose
  }
  return null
})

watch(computedRatio, (val) => {
  form.ratio = val
})

watch(
  () => props.tasting,
  (t) => {
    if (t) {
      Object.assign(form, createEmptyForm(), t)
      form.tasted_at = t.tasted_at || new Date().toISOString().substring(0, 10)
      for (const key of Object.keys(idkFields)) {
        ;(idkFields as any)[key] = (t as any)[key] == null
      }
    } else {
      Object.assign(form, createEmptyForm())
      for (const key of Object.keys(idkFields)) {
        ;(idkFields as any)[key] = false
      }
    }
  },
  { immediate: true }
)

for (const key of Object.keys(idkFields)) {
  watch(
    () => (idkFields as any)[key],
    (idk) => {
      if (idk) {
        ;(form as any)[key] = null
      }
    }
  )
}

function cancelQuickBean() {
  showQuickAddBean.value = false
  quickBean.name = ''
  quickBean.roaster = null
  quickBean.origin = null
}

async function saveQuickBean() {
  if (!quickBean.name.trim()) return
  quickBeanSaving.value = true
  try {
    const newBean = await upsertCoffeeBeanRaw({
      name: quickBean.name.trim(),
      roaster: quickBean.roaster?.trim() || null,
      origin: quickBean.origin?.trim() || null
    })
    emit('beanAdded', newBean)
    form.bean = newBean
    cancelQuickBean()
  } catch {
    alert('Failed to add bean')
  } finally {
    quickBeanSaving.value = false
  }
}

function cancelQuickEquipment() {
  showQuickAddEquipment.value = false
  quickEquipment.name = ''
  quickEquipment.type = null
}

async function saveQuickEquipment() {
  if (!quickEquipment.name.trim()) return
  quickEquipmentSaving.value = true
  try {
    const newEquipment = await upsertCoffeeEquipmentRaw({
      name: quickEquipment.name.trim(),
      type: quickEquipment.type || null
    })
    emit('equipmentAdded', newEquipment)
    form.equipment = newEquipment
    cancelQuickEquipment()
  } catch {
    alert('Failed to add equipment')
  } finally {
    quickEquipmentSaving.value = false
  }
}

async function save() {
  saving.value = true
  try {
    await upsertCoffeeTastingRaw({
      ...form,
      tasted_at: new Date(form.tasted_at).toISOString()
    })
    emit('saved')
  } catch (e) {
    alert('Failed to save tasting')
  } finally {
    saving.value = false
  }
}
</script>
