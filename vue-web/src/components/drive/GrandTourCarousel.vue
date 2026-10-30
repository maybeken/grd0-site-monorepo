<template>
  <section
    class="relative h-full w-full cursor-grab overflow-hidden"
    role="region"
    aria-roledescription="carousel"
    :aria-label="$t('grandTour.view.title')"
  >
    <div
      class="carousel-stage relative h-full w-full"
      @pointerdown="onPointerDown"
      @pointerup="onPointerUp"
      @pointercancel="onPointerCancel"
    >
      <div
        v-for="(vehicle, index) in vehicles"
        :key="vehicle.id ?? index"
        class="carousel-slide absolute inset-0 flex select-none items-center justify-center"
        :style="slideStyle(index)"
        :aria-hidden="index === current ? 'false' : 'true'"
        :data-active="index === current ? 'true' : 'false'"
        @click="onSlideClick(index)"
      >
        <div class="h-full w-[82%] max-w-5xl text-foreground">
          <VehicleSlide :vehicle="vehicle"></VehicleSlide>
        </div>
      </div>
    </div>

    <div
      class="pointer-events-none absolute inset-x-0 bottom-2 z-40 flex items-center justify-center gap-4 md:bottom-4"
    >
      <button
        type="button"
        class="pointer-events-auto rounded-full border border-dotted border-foreground/40 p-2 transition-colors hover:bg-accent/20 disabled:cursor-not-allowed disabled:opacity-40"
        :disabled="vehicles.length <= 1"
        :aria-label="$t('grandTour.view.previous')"
        @click="go(-1)"
      >
        <Icon icon="mynaui:arrow-left" height="1.25rem"></Icon>
      </button>
      <button
        type="button"
        class="pointer-events-auto rounded-full border border-dotted border-foreground/40 p-2 transition-colors hover:bg-accent/20 disabled:cursor-not-allowed disabled:opacity-40"
        :disabled="vehicles.length <= 1"
        :aria-label="$t('grandTour.view.next')"
        @click="go(1)"
      >
        <Icon icon="mynaui:arrow-right" height="1.25rem"></Icon>
      </button>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'

import VehicleSlide from '@/components/drive/VehicleSlide.vue'
import Icon from '@/components/shared/Icon.vue'

import type { Vehicle } from '@/interfaces/Vehicle'

const props = defineProps<{
  vehicles: Vehicle[]
}>()

const current = ref(0)

let pointer_start: { x: number; y: number } | null = null
let last_swipe = 0

function offsetFor(index: number): number {
  const count = props.vehicles.length

  if (count === 0) return 0

  let offset = index - current.value

  if (offset > count / 2) offset -= count
  if (offset < -count / 2) offset += count

  return offset
}

function slideStyle(index: number): Record<string, string> {
  const offset = offsetFor(index)
  const distance = Math.abs(offset)
  const hidden = distance > 2

  return {
    transform: `translateX(${offset * 55}%) translateZ(${-distance * 210}px) rotateY(${offset * -28}deg)`,
    opacity: hidden ? '0' : String(distance === 0 ? 1 : distance === 1 ? 0.45 : 0.22),
    zIndex: String(30 - distance * 10),
    visibility: hidden ? 'hidden' : 'visible',
    pointerEvents: hidden ? 'none' : 'auto'
  }
}

function goTo(index: number): void {
  const count = props.vehicles.length

  if (count === 0) return

  const next = ((index % count) + count) % count
  current.value = next
}

function go(delta: number): void {
  goTo(current.value + delta)
}

function onSlideClick(index: number): void {
  if (index === current.value) return
  if (Date.now() - last_swipe < 300) return

  goTo(index)
}

function onPointerDown(event: PointerEvent): void {
  if (event.pointerType === 'mouse' && event.button !== 0) return

  pointer_start = { x: event.clientX, y: event.clientY }
}

function onPointerUp(event: PointerEvent): void {
  if (!pointer_start) return

  const delta_x = event.clientX - pointer_start.x
  const delta_y = event.clientY - pointer_start.y

  pointer_start = null

  if (Math.abs(delta_x) < 50 || Math.abs(delta_x) <= Math.abs(delta_y)) return

  last_swipe = Date.now()
  go(delta_x < 0 ? 1 : -1)
}

function onPointerCancel(): void {
  pointer_start = null
}

function onKeydown(event: KeyboardEvent): void {
  const target = event.target as HTMLElement | null

  if (target?.closest('input, textarea, select, [contenteditable="true"]')) return

  if (event.key === 'ArrowRight') {
    event.preventDefault()
    go(1)
  } else if (event.key === 'ArrowLeft') {
    event.preventDefault()
    go(-1)
  }
}

onMounted(() => {
  window.addEventListener('keydown', onKeydown)
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', onKeydown)
})
</script>

<style scoped>
.carousel-stage {
  perspective: 1500px;
  perspective-origin: 50% 45%;
  -webkit-mask-image: linear-gradient(
    to right,
    transparent 0,
    black 12%,
    black 88%,
    transparent 100%
  );
  mask-image: linear-gradient(to right, transparent 0, black 12%, black 88%, transparent 100%);
}

.carousel-slide {
  transition:
    transform 600ms cubic-bezier(0.22, 1, 0.36, 1),
    opacity 450ms ease;
  will-change: transform, opacity;
}

@media (prefers-reduced-motion: reduce) {
  .carousel-slide {
    transition: none;
  }
}
</style>
