<template>
  <div class="flex flex-col gap-4 bg-background rounded-xl p-4 border border-accent/20">
    <div class="flex items-center justify-between">
      <h3 class="text-lg font-bold">{{ isEdit ? $t('coffee.beanForm.titleEdit') : $t('coffee.beanForm.titleNew') }}</h3>
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
        <label class="text-sm font-semibold">{{ $t('coffee.beanForm.name') }}</label>
        <input
          v-model="form.name"
          class="bg-background border border-accent/40 rounded px-2 py-1 text-sm"
          :placeholder="$t('coffee.beanForm.namePlaceholder')"
        />
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-sm font-semibold">{{ $t('coffee.beanForm.origin') }}</label>
        <input
          v-model="form.origin"
          class="bg-background border border-accent/40 rounded px-2 py-1 text-sm"
          :placeholder="$t('coffee.beanForm.originPlaceholder')"
        />
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-sm font-semibold">{{ $t('coffee.beanForm.roaster') }}</label>
        <input
          v-model="form.roaster"
          class="bg-background border border-accent/40 rounded px-2 py-1 text-sm"
          :placeholder="$t('coffee.beanForm.roasterPlaceholder')"
        />
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-sm font-semibold">{{ $t('coffee.beanForm.roastDate') }}</label>
        <input
          type="date"
          v-model="form.roast_date"
          class="bg-background border border-accent/40 rounded px-2 py-1 text-sm"
        />
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-sm font-semibold">{{ $t('coffee.beanForm.variety') }}</label>
        <input
          v-model="form.variety"
          class="bg-background border border-accent/40 rounded px-2 py-1 text-sm"
          :placeholder="$t('coffee.beanForm.varietyPlaceholder')"
        />
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-sm font-semibold">{{ $t('coffee.beanForm.process') }}</label>
        <input
          v-model="form.process"
          class="bg-background border border-accent/40 rounded px-2 py-1 text-sm"
          :placeholder="$t('coffee.beanForm.processPlaceholder')"
        />
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-sm font-semibold">{{ $t('coffee.beanForm.altitude') }}</label>
        <input
          v-model="form.altitude"
          class="bg-background border border-accent/40 rounded px-2 py-1 text-sm"
          :placeholder="$t('coffee.beanForm.altitudePlaceholder')"
        />
      </div>
    </div>

    <div class="flex flex-col gap-1">
      <label class="text-sm font-semibold">{{ $t('coffee.beanForm.description') }}</label>
      <textarea
        v-model="form.description"
        rows="3"
        class="bg-background border border-accent/40 rounded px-2 py-1 text-sm resize-y"
        :placeholder="$t('coffee.beanForm.descriptionPlaceholder')"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { CoffeeBean } from '@/interfaces/Coffee'
import { upsertCoffeeBeanRaw } from '@/services/coffee'

const { t } = useI18n()

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
    alert(t('coffee.beanForm.failedSave'))
  } finally {
    saving.value = false
  }
}
</script>
