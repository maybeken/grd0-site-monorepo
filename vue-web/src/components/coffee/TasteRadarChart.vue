<template>
  <div v-if="hasData" class="w-full aspect-square max-h-full">
    <Radar :data="chartData" :options="chartOptions" />
  </div>
  <div v-else class="flex items-center justify-center w-full aspect-square text-sm opacity-50">
    No taste data
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Radar } from 'vue-chartjs'
import {
  Chart as ChartJS,
  RadialLinearScale,
  PointElement,
  LineElement,
  Filler,
  Tooltip,
  Legend
} from 'chart.js'

ChartJS.register(RadialLinearScale, PointElement, LineElement, Filler, Tooltip, Legend)

interface Props {
  tasteFruity?: number | null
  tasteSour?: number | null
  tasteFermented?: number | null
  tasteSweetness?: number | null
  tasteNutty?: number | null
  tasteSpice?: number | null
  tasteFloral?: number | null
  tasteGreen?: number | null
  tasteTobacco?: number | null
  tasteBitter?: number | null
}

const props = defineProps<Props>()

const labels = ['Nutty', 'Spice', 'Fruity', 'Sour', 'Fermented', 'Sweetness', 'Floral', 'Green', 'Tobacco', 'Bitter']

const values = computed(() => [
  props.tasteNutty,
  props.tasteSpice,
  props.tasteFruity,
  props.tasteSour,
  props.tasteFermented,
  props.tasteSweetness,
  props.tasteFloral,
  props.tasteGreen,
  props.tasteTobacco,
  props.tasteBitter,
])

const hasData = computed(() => values.value.some((v) => v != null))

const chartData = computed(() => ({
  labels,
  datasets: [
    {
      data: values.value.map((v) => v ?? 0),
      backgroundColor: 'rgba(99, 102, 241, 0.2)',
      borderColor: 'rgba(99, 102, 241, 0.8)',
      borderWidth: 2,
      pointBackgroundColor: 'rgba(99, 102, 241, 1)',
      pointRadius: 3
    }
  ]
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: true,
  scales: {
    r: {
      min: -1,
      max: 10,
      ticks: {
        stepSize: 2,
        display: false
      },
      pointLabels: {
        font: { size: 10 }
      },
      grid: {
        color: 'rgba(255, 255, 255, 0.1)'
      },
      angleLines: {
        color: 'rgba(255, 255, 255, 0.1)'
      }
    }
  },
  plugins: {
    legend: { display: false },
    tooltip: { enabled: true }
  }
}
</script>
