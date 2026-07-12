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
        <label class="text-sm font-semibold">Bean</label>
        <select v-model="form.bean.id" class="bg-background border border-accent/40 rounded px-2 py-1 text-sm">
          <option value="">Select bean...</option>
          <option v-for="b in beans" :key="b.id" :value="b.id">{{ b.name }}</option>
        </select>
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-sm font-semibold">Equipment</label>
        <select v-model="form.equipment.id" class="bg-background border border-accent/40 rounded px-2 py-1 text-sm">
          <option value="">Select equipment...</option>
          <option v-for="e in equipmentList" :key="e.id" :value="e.id">{{ e.name }}</option>
        </select>
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
          <IdkToggle v-model="idkFields.rating" />
        </div>
      </div>
    </div>

    <fieldset class="border border-accent/20 rounded p-3">
      <legend class="text-sm font-bold px-1">Brew Recipe</legend>
      <div class="grid grid-cols-2 md:grid-cols-3 gap-3">
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
            <span class="text-xs">{{ computedRatio != null ? '1:' + computedRatio.toFixed(1) : 'IDK' }}</span>
          </div>
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-xs">Brew Time</label>
          <div class="flex items-center gap-1">
            <BrewTimeInput
              v-model="form.brew_time"
              :disabled="idkFields.brew_time"
            />
            <IdkToggle v-model="idkFields.brew_time" />
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
            <IdkToggle v-model="(idkFields as any)[dim.key]" />
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
import { GRIND_SIZES, TASTE_DIMENSIONS } from '@/helpers/coffee'
import { upsertCoffeeTastingRaw } from '@/services/coffee'
import IdkToggle from './IdkToggle.vue'
import BrewTimeInput from './BrewTimeInput.vue'

interface Props {
  tasting?: TastingNote | null
  beans: CoffeeBean[]
  equipmentList: BrewEquipment[]
}

const props = withDefaults(defineProps<Props>(), {
  tasting: null,
})

const emit = defineEmits<{
  saved: []
  cancel: []
}>()

const isEdit = computed(() => !!props.tasting?.id)
const saving = ref(false)

const tasteDimensions = TASTE_DIMENSIONS.map((d) => ({
  key: `taste_${d}`,
  label: d,
}))

function createEmptyForm(): TastingNote {
  return {
    bean: { name: '' },
    equipment: { name: '' },
    tasted_at: new Date().toISOString(),
    pinned: false,
    grind_size: null,
    grind_setting: null,
    coffee_dose: null,
    water_in: null,
    coffee_out: null,
    ratio: null,
    brew_time: null,
    taste_fruity: null,
    taste_sour: null,
    taste_sweetness: null,
    taste_nutty: null,
    taste_spice: null,
    taste_floral: null,
    taste_green: null,
    overall_notes: null,
    rating: null,
  }
}

const form = reactive<TastingNote>(createEmptyForm())

const idkFields = reactive<Record<string, boolean>>({
  grind_size: false,
  grind_setting: false,
  coffee_dose: false,
  water_in: false,
  coffee_out: false,
  brew_time: false,
  rating: false,
  taste_fruity: false,
  taste_sour: false,
  taste_sweetness: false,
  taste_nutty: false,
  taste_spice: false,
  taste_floral: false,
  taste_green: false,
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

async function save() {
  saving.value = true
  try {
    await upsertCoffeeTastingRaw({ ...form })
    emit('saved')
  } catch (e) {
    alert('Failed to save tasting')
  } finally {
    saving.value = false
  }
}
</script>
