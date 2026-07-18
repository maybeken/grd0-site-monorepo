<template>
  <div class="flex gap-4 mx-auto w-full pb-8">
    <div class="flex gap-2">
      <label>{{ $t('tool.dotMatrix.finest') }}</label>
      <input class="bg-inherit appearance-none text-center" type="number" v-model="finest" />
      <label>{{ $t('tool.dotMatrix.gap') }}</label>
      <input class="bg-inherit appearance-none text-center" type="number" v-model="gap" />
    </div>
    <div>
      <button class="px-4 py-px border rounded-lg" @click="getSvg">{{ $t('common.copy') }}</button>
    </div>
  </div>
  <div class="flex mx-auto bg-white w-64">
    <svg
      id="svg"
      xmlns="http://www.w3.org/2000/svg"
      :viewBox="`${-(finest * size + gap * size * 2)} ${-(finest * size + gap * size * 2)} ${(finest * size + gap * size * 2) * 2} ${(finest * size + gap * size * 2) * 2}`"
      fill="black"
    >
      <g v-for="x in Array.from({ length: 2 * finest + 1 }, (_, i) => i - finest)" :key="x">
        <circle
          v-for="y in Array.from({ length: 2 * finest + 1 }, (_, i) => i - finest)"
          :cx="x * (size * 2 + gap)"
          :cy="y * (size * 2 + gap)"
          :r="size"
          :key="y"
        />
      </g>
    </svg>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const size = ref(8)
const gap = ref(8)
const finest = ref(1)

function getSvg() {
  // @ts-expect-error
  navigator.clipboard.writeText(document.getElementById('svg').outerHTML)
  alert(t('common.copiedToClipboard'))
}
</script>
