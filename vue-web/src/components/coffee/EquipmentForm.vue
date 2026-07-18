<template>
  <div class="flex flex-col gap-4 bg-background rounded-xl p-4 border border-accent/20">
    <div class="flex items-center justify-between">
      <h3 class="text-lg font-bold">{{ isEdit ? $t('coffee.equipmentForm.titleEdit') : $t('coffee.equipmentForm.titleNew') }}</h3>
      <div class="flex gap-2">
        <button
          @click="$emit('cancel')"
          class="px-3 py-1 rounded border border-accent/40 text-sm hover:bg-accent/10"
        >
          {{ $t('common.cancel') }}
        </button>
        <button
          @click="save"
          :disabled="saving || !form.name.trim()"
          class="px-3 py-1 rounded bg-accent text-background text-sm font-semibold hover:brightness-110 disabled:opacity-50"
        >
          {{ saving ? $t('common.saving') : $t('common.save') }}
        </button>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="flex flex-col gap-1">
        <label class="text-sm font-semibold">{{ $t('coffee.equipmentForm.name') }}</label>
        <input
          v-model="form.name"
          class="bg-background border border-accent/40 rounded px-2 py-1 text-sm"
          :placeholder="$t('coffee.equipmentForm.namePlaceholder')"
        />
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-sm font-semibold">{{ $t('coffee.equipmentForm.type') }}</label>
        <select
          v-model="form.type"
          class="bg-background border border-accent/40 rounded px-2 py-1 text-sm"
        >
          <option :value="null">{{ $t('common.selectType') }}</option>
          <option v-for="t in allTypes" :key="t.value" :value="t.value">{{ $t(t.label) }}</option>
        </select>
      </div>
    </div>

    <div class="flex flex-col gap-1">
      <label class="text-sm font-semibold">{{ $t('coffee.equipmentForm.description') }}</label>
      <textarea
        v-model="form.description"
        rows="3"
        class="bg-background border border-accent/40 rounded px-2 py-1 text-sm resize-y"
        :placeholder="$t('coffee.equipmentForm.descriptionPlaceholder')"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { BrewEquipment } from '@/interfaces/Coffee'
import { upsertCoffeeEquipmentRaw } from '@/services/coffee'
import { EQUIPMENT_TYPES } from '@/helpers/coffee'

const { t } = useI18n()

interface Props {
  equipment?: BrewEquipment | null
}

const props = withDefaults(defineProps<Props>(), {
  equipment: null
})

const emit = defineEmits<{
  saved: []
  cancel: []
}>()

const isEdit = computed(() => !!props.equipment?.id)
const saving = ref(false)

const allTypes = [...EQUIPMENT_TYPES]

function createEmptyForm(): BrewEquipment {
  return {
    name: '',
    type: null,
    description: null
  }
}

const form = reactive<BrewEquipment>(createEmptyForm())

watch(
  () => props.equipment,
  (e) => {
    if (e) {
      Object.assign(form, createEmptyForm(), e)
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
    await upsertCoffeeEquipmentRaw(form)
    emit('saved')
  } catch {
    alert(t('coffee.equipmentForm.failedSave'))
  } finally {
    saving.value = false
  }
}
</script>
