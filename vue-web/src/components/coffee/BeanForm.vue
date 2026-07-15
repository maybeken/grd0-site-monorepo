<template>
  <div class="flex flex-col gap-4 bg-background rounded-xl p-4 border border-accent/20">
    <div class="flex items-center justify-between">
      <h3 class="text-lg font-bold">{{ isEdit ? 'Edit Bean' : 'New Bean' }}</h3>
      <div class="flex gap-2">
        <button
          @click="$emit('cancel')"
          class="px-3 py-1 rounded border border-accent/40 text-sm hover:bg-accent/10"
        >
          Cancel
        </button>
        <button
          @click="save"
          :disabled="saving || !form.name.trim()"
          class="px-3 py-1 rounded bg-accent text-background text-sm font-semibold hover:brightness-110 disabled:opacity-50"
        >
          {{ saving ? 'Saving...' : 'Save' }}
        </button>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="flex flex-col gap-1">
        <label class="text-sm font-semibold">Name *</label>
        <input
          v-model="form.name"
          class="bg-background border border-accent/40 rounded px-2 py-1 text-sm"
          placeholder="e.g. Yirgacheffe Kochere"
        />
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-sm font-semibold">Origin</label>
        <input
          v-model="form.origin"
          class="bg-background border border-accent/40 rounded px-2 py-1 text-sm"
          placeholder="e.g. Ethiopia"
        />
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-sm font-semibold">Roaster</label>
        <input
          v-model="form.roaster"
          class="bg-background border border-accent/40 rounded px-2 py-1 text-sm"
          placeholder="e.g. Blue Bottle"
        />
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-sm font-semibold">Roast Date</label>
        <input
          type="date"
          v-model="form.roast_date"
          class="bg-background border border-accent/40 rounded px-2 py-1 text-sm"
        />
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-sm font-semibold">Variety</label>
        <input
          v-model="form.variety"
          class="bg-background border border-accent/40 rounded px-2 py-1 text-sm"
          placeholder="e.g. Heirloom"
        />
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-sm font-semibold">Process</label>
        <input
          v-model="form.process"
          class="bg-background border border-accent/40 rounded px-2 py-1 text-sm"
          placeholder="e.g. Washed"
        />
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-sm font-semibold">Altitude</label>
        <input
          v-model="form.altitude"
          class="bg-background border border-accent/40 rounded px-2 py-1 text-sm"
          placeholder="e.g. 1900-2100 masl"
        />
      </div>
    </div>

    <div class="flex flex-col gap-1">
      <label class="text-sm font-semibold">Description</label>
      <textarea
        v-model="form.description"
        rows="3"
        class="bg-background border border-accent/40 rounded px-2 py-1 text-sm resize-y"
        placeholder="Tasting notes, flavor profile..."
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed } from 'vue'
import type { CoffeeBean } from '@/interfaces/Coffee'
import { upsertCoffeeBeanRaw } from '@/services/coffee'

interface Props {
  bean?: CoffeeBean | null
}

const props = withDefaults(defineProps<Props>(), {
  bean: null
})

const emit = defineEmits<{
  saved: []
  cancel: []
}>()

const isEdit = computed(() => !!props.bean?.id)
const saving = ref(false)

function createEmptyForm(): CoffeeBean {
  return {
    name: '',
    origin: null,
    roaster: null,
    roast_date: null,
    variety: null,
    process: null,
    altitude: null,
    description: null
  }
}

const form = reactive<CoffeeBean>(createEmptyForm())

watch(
  () => props.bean,
  (b) => {
    if (b) {
      Object.assign(form, createEmptyForm(), b)
      if (b.roast_date) {
        form.roast_date = b.roast_date.split('T')[0] ?? null
      }
    } else {
      Object.assign(form, createEmptyForm())
    }
  },
  { immediate: true }
)

async function save() {
  if (!form.name.trim()) return
  saving.value = true
  try {
    await upsertCoffeeBeanRaw({
      ...form,
      roast_date: form.roast_date ? new Date(form.roast_date).toISOString() : null
    })
    emit('saved')
  } catch {
    alert('Failed to save bean')
  } finally {
    saving.value = false
  }
}
</script>
