<template>
  <div class="flex flex-col gap-2" aria-busy="true">
    <span class="sr-only">{{ resolved_label }}</span>
    <Skeleton v-for="row in rows" :key="row" h="sm" w="full" :loading="true"></Skeleton>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import Skeleton from '@/components/shared/Skeleton.vue'

const props = withDefaults(
  defineProps<{
    rows?: number
    label?: string
  }>(),
  {
    rows: 3,
    label: ''
  }
)

const { t } = useI18n()
const resolved_label = computed(() => props.label || t('common.loading'))
</script>
