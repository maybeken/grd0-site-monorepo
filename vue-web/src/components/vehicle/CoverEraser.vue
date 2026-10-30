<template>
  <div class="flex flex-col gap-3">
    <p class="text-sm font-medium leading-relaxed text-foreground/60">
      {{ $t('grandTour.editor.cover.eraseHint') }}
    </p>

    <div
      ref="viewport"
      class="erase-viewport checker relative h-[42vh] touch-none select-none overflow-auto rounded border border-accent/20"
      :class="panning ? 'cursor-grabbing' : 'cursor-crosshair'"
      @pointerdown="onPointerDown"
      @pointermove="onPointerMove"
      @pointerup="onPointerUp"
      @pointercancel="onPointerUp"
      @pointerenter="update_cursor"
      @pointerleave="cursor_visible = false"
      @contextmenu.prevent
      @wheel="onWheel"
    >
      <div
        class="mx-auto"
        :style="{
          width: `${surface_size.width}px`,
          height: `${surface_size.height}px`,
          padding: `${surface_padding.y}px ${surface_padding.x}px`
        }"
      >
        <div class="relative" :style="{ width: `${display_size.width}px`, height: `${display_size.height}px` }">
          <canvas
            ref="canvas"
            class="block"
            :style="{ width: `${display_size.width}px`, height: `${display_size.height}px` }"
          ></canvas>

          <div
            v-if="cursor_visible"
            class="brush-cursor pointer-events-none absolute rounded-full border border-foreground shadow-[0_0_0_1px_rgba(0,0,0,0.55)]"
            :style="cursor_style"
          ></div>
        </div>
      </div>
    </div>

    <div class="flex flex-wrap items-center gap-3">
      <label class="flex items-center gap-2 text-sm font-medium text-foreground/60">
        {{ $t('grandTour.editor.cover.brushSize') }}
        <input
          v-model.number="brush"
          type="range"
          min="8"
          max="160"
          step="4"
          class="accent-foreground"
        />
      </label>

      <div class="flex items-center gap-1">
        <EditorButton
          variant="secondary"
          size="sm"
          :disabled="zoom <= 1"
          :label="$t('grandTour.editor.cover.zoomOut')"
          @click="zoom_out"
        >
          −
        </EditorButton>
        <span class="min-w-12 text-center text-sm font-medium text-foreground/60" aria-live="polite">
          {{ Math.round(zoom * 100) }}%
        </span>
        <EditorButton
          variant="secondary"
          size="sm"
          :disabled="zoom >= 6"
          :label="$t('grandTour.editor.cover.zoomIn')"
          @click="zoom_in"
        >
          +
        </EditorButton>
        <EditorButton
          variant="secondary"
          size="sm"
          :disabled="zoom === 1"
          :label="$t('grandTour.editor.cover.zoomFit')"
          @click="reset_zoom"
        >
          {{ $t('grandTour.editor.cover.zoomFit') }}
        </EditorButton>
      </div>

      <EditorButton variant="secondary" size="sm" :disabled="strokes.length === 0" @click="undo">
        {{ $t('grandTour.editor.cover.undoErase') }}
      </EditorButton>

      <EditorButton variant="secondary" size="sm" :disabled="strokes.length === 0" @click="clear">
        {{ $t('grandTour.editor.cover.clearErase') }}
      </EditorButton>
    </div>

    <div class="grid grid-cols-2 gap-3">
      <figure class="flex flex-col gap-1">
        <div
          class="flex h-24 items-center justify-center rounded border border-accent/20 bg-background"
        >
          <img
            :src="light_url"
            :alt="$t('grandTour.editor.cover.silhouetteLight')"
            class="h-full w-full object-contain"
          />
        </div>
        <figcaption class="text-sm font-medium text-foreground/60">
          {{ $t('grandTour.editor.cover.silhouetteLight') }}
        </figcaption>
      </figure>

      <figure class="flex flex-col gap-1">
        <div
          class="flex h-24 items-center justify-center rounded border border-accent/20 bg-foreground"
        >
          <img
            :src="dark_url"
            :alt="$t('grandTour.editor.cover.silhouetteDark')"
            class="h-full w-full object-contain"
          />
        </div>
        <figcaption class="text-sm font-medium text-foreground/60">
          {{ $t('grandTour.editor.cover.silhouetteDark') }}
        </figcaption>
      </figure>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'

import EditorButton from '@/components/editor/ui/EditorButton.vue'

import type { GeneratedCover } from '@/interfaces/Vehicle'

interface Point {
  x: number
  y: number
}

interface Stroke {
  points: Point[]
  radius: number
}

const props = defineProps<{
  cover: GeneratedCover
}>()

const emit = defineEmits<{
  update: [cover: GeneratedCover]
}>()

const canvas = ref<HTMLCanvasElement | null>(null)
const viewport = ref<HTMLDivElement | null>(null)
const brush = ref(40)
const zoom = ref(1)
const strokes = ref<Stroke[]>([])
const light_url = ref(props.cover.silhouette_light_url)
const dark_url = ref(props.cover.silhouette_dark_url)
const viewport_size = ref({ width: 0, height: 0 })
const image_size = ref({ width: 0, height: 0 })
const cursor_x = ref(0)
const cursor_y = ref(0)
const cursor_visible = ref(false)
const panning = ref(false)

let pan_start: { x: number; y: number; left: number; top: number } | null = null

const fit_scale = computed(() => {
  const { width, height } = viewport_size.value
  const image = image_size.value

  if (!width || !height || !image.width || !image.height) return 1

  return Math.min(width / image.width, height / image.height, 1)
})

const surface_padding = computed(() => {
  if (zoom.value <= 1) return { x: 0, y: 0 }

  return {
    x: Math.round(viewport_size.value.width / 2),
    y: Math.round(viewport_size.value.height / 2)
  }
})

const display_size = computed(() => ({
  width: Math.max(1, Math.round(image_size.value.width * fit_scale.value * zoom.value)),
  height: Math.max(1, Math.round(image_size.value.height * fit_scale.value * zoom.value))
}))

const surface_size = computed(() => ({
  width: display_size.value.width + surface_padding.value.x * 2,
  height: display_size.value.height + surface_padding.value.y * 2
}))

const cursor_style = computed(() => {
  const scale = image_size.value.width ? display_size.value.width / image_size.value.width : 1
  const radius = Math.max(2, brush.value * scale)

  return {
    left: `${cursor_x.value - radius}px`,
    top: `${cursor_y.value - radius}px`,
    width: `${radius * 2}px`,
    height: `${radius * 2}px`
  }
})

function center_view(): void {
  void nextTick(() => {
    const element = viewport.value
    if (!element) return

    element.scrollLeft = Math.max(0, Math.round((element.scrollWidth - element.clientWidth) / 2))
    element.scrollTop = Math.max(0, Math.round((element.scrollHeight - element.clientHeight) / 2))
  })
}

function zoom_in(): void {
  zoom.value = Math.min(6, zoom.value + 0.5)
  center_view()
}

function zoom_out(): void {
  zoom.value = Math.max(1, zoom.value - 0.5)
  center_view()
}

function reset_zoom(): void {
  zoom.value = 1
  center_view()
}

function onWheel(event: WheelEvent): void {
  if (!event.ctrlKey && !event.metaKey) return

  event.preventDefault()

  const element = viewport.value
  const image = canvas.value
  if (!element || !image) return

  const rect_before = image.getBoundingClientRect()
  if (rect_before.width === 0 || rect_before.height === 0) return

  const client_x = event.clientX
  const client_y = event.clientY
  const ratio_x = (client_x - rect_before.left) / rect_before.width
  const ratio_y = (client_y - rect_before.top) / rect_before.height

  const next = Math.min(
    6,
    Math.max(1, Math.round((zoom.value - event.deltaY / 200) * 10) / 10)
  )

  if (next === zoom.value) return

  zoom.value = next

  void nextTick(() => {
    const current = viewport.value
    const canvas_element = canvas.value
    if (!current || !canvas_element) return

    const rect_after = canvas_element.getBoundingClientRect()

    current.scrollLeft += rect_after.left + ratio_x * rect_after.width - client_x
    current.scrollTop += rect_after.top + ratio_y * rect_after.height - client_y
  })
}

let mask: HTMLCanvasElement | null = null
let cutout_image: HTMLImageElement | null = null
let light_image: HTMLImageElement | null = null
let dark_image: HTMLImageElement | null = null
let derived_urls: string[] = []
let render_seq = 0
let drawing = false
let last_point: Point | null = null

function load_image(src: string): Promise<HTMLImageElement> {
  return new Promise((resolve, reject) => {
    const image = new Image()
    image.onload = () => resolve(image)
    image.onerror = () => reject(new Error('Could not load the generated cover.'))
    image.src = src
  })
}

function preview_context(): CanvasRenderingContext2D | null {
  return canvas.value?.getContext('2d') ?? null
}

function paint_mask(from: Point, to: Point): void {
  if (!mask) return

  const context = mask.getContext('2d')
  if (!context) return

  const radius = brush.value

  context.fillStyle = '#fff'
  context.strokeStyle = '#fff'
  context.lineWidth = radius * 2
  context.lineCap = 'round'
  context.lineJoin = 'round'
  context.beginPath()
  context.moveTo(from.x, from.y)
  context.lineTo(to.x, to.y)
  context.stroke()

  if (from.x === to.x && from.y === to.y) {
    context.beginPath()
    context.arc(to.x, to.y, radius, 0, Math.PI * 2)
    context.fill()
  }
}

function erase_preview(from: Point, to: Point): void {
  const context = preview_context()
  if (!context) return

  const radius = brush.value

  context.globalCompositeOperation = 'destination-out'
  context.strokeStyle = '#000'
  context.fillStyle = '#000'
  context.lineWidth = radius * 2
  context.lineCap = 'round'
  context.lineJoin = 'round'
  context.beginPath()
  context.moveTo(from.x, from.y)
  context.lineTo(to.x, to.y)
  context.stroke()

  if (from.x === to.x && from.y === to.y) {
    context.beginPath()
    context.arc(to.x, to.y, radius, 0, Math.PI * 2)
    context.fill()
  }

  context.globalCompositeOperation = 'source-over'
}

function replay_mask(): void {
  if (!mask) return

  const context = mask.getContext('2d')
  if (!context) return

  context.clearRect(0, 0, mask.width, mask.height)

  for (const stroke of strokes.value) {
    const [first, ...rest] = stroke.points

    if (!first) continue

    context.fillStyle = '#fff'
    context.strokeStyle = '#fff'
    context.lineWidth = stroke.radius * 2
    context.lineCap = 'round'
    context.lineJoin = 'round'

    if (rest.length === 0) {
      context.beginPath()
      context.arc(first.x, first.y, stroke.radius, 0, Math.PI * 2)
      context.fill()
      continue
    }

    context.beginPath()
    context.moveTo(first.x, first.y)
    for (const point of rest) context.lineTo(point.x, point.y)
    context.stroke()
  }
}

function redraw(): void {
  const element = canvas.value

  if (!element || !cutout_image || !mask) return

  if (element.width !== cutout_image.naturalWidth) element.width = cutout_image.naturalWidth
  if (element.height !== cutout_image.naturalHeight) element.height = cutout_image.naturalHeight

  const context = element.getContext('2d')
  if (!context) return

  context.globalCompositeOperation = 'source-over'
  context.clearRect(0, 0, element.width, element.height)
  context.drawImage(cutout_image, 0, 0)

  if (strokes.value.length > 0) {
    context.globalCompositeOperation = 'destination-out'
    context.drawImage(mask, 0, 0)
    context.globalCompositeOperation = 'source-over'
  }
}

function reset_urls(): void {
  for (const url of derived_urls) URL.revokeObjectURL(url)

  derived_urls = []
  light_url.value = props.cover.silhouette_light_url
  dark_url.value = props.cover.silhouette_dark_url
}

function composite(image: HTMLImageElement): Promise<Blob> {
  const buffer = document.createElement('canvas')
  buffer.width = image.naturalWidth
  buffer.height = image.naturalHeight

  const context = buffer.getContext('2d')

  if (!context) return Promise.reject(new Error('Canvas is unavailable.'))

  context.drawImage(image, 0, 0)

  if (mask) {
    context.globalCompositeOperation = 'destination-out'
    context.drawImage(mask, 0, 0)
    context.globalCompositeOperation = 'source-over'
  }

  return new Promise((resolve, reject) => {
    buffer.toBlob((blob) => {
      if (blob) resolve(blob)
      else reject(new Error('Could not encode the cover.'))
    }, 'image/png')
  })
}

async function rebuild(): Promise<void> {
  if (!cutout_image || !light_image || !dark_image) return

  if (strokes.value.length === 0) {
    reset_urls()
    emit('update', props.cover)
    return
  }

  const sequence = (render_seq += 1)
  const [cutout, light, dark] = await Promise.all([
    composite(cutout_image),
    composite(light_image),
    composite(dark_image)
  ])

  if (sequence !== render_seq) return

  const urls = [
    URL.createObjectURL(cutout),
    URL.createObjectURL(light),
    URL.createObjectURL(dark)
  ]

  for (const url of derived_urls) URL.revokeObjectURL(url)

  derived_urls = urls
  light_url.value = urls[1]!
  dark_url.value = urls[2]!

  emit('update', {
    cutout,
    silhouette_light: light,
    silhouette_dark: dark,
    cutout_url: urls[0]!,
    silhouette_light_url: urls[1]!,
    silhouette_dark_url: urls[2]!
  })
}

function point_from(event: PointerEvent): Point | null {
  const element = canvas.value
  if (!element) return null

  const rect = element.getBoundingClientRect()
  if (rect.width === 0 || rect.height === 0) return null

  return {
    x: ((event.clientX - rect.left) * element.width) / rect.width,
    y: ((event.clientY - rect.top) * element.height) / rect.height
  }
}

function update_cursor(event: PointerEvent): void {
  const element = canvas.value
  if (!element) return

  const rect = element.getBoundingClientRect()
  if (rect.width === 0 || rect.height === 0) return

  cursor_x.value = event.clientX - rect.left
  cursor_y.value = event.clientY - rect.top
  cursor_visible.value = true
}

function onPointerDown(event: PointerEvent): void {
  if (event.button === 2) {
    const element = viewport.value
    if (!element) return

    event.preventDefault()
    panning.value = true
    pan_start = {
      x: event.clientX,
      y: event.clientY,
      left: element.scrollLeft,
      top: element.scrollTop
    }
    element.setPointerCapture(event.pointerId)
    cursor_visible.value = false
    return
  }

  if (event.button !== 0 || !mask) return

  const point = point_from(event)
  if (!point) return

  event.preventDefault()
  update_cursor(event)
  viewport.value?.setPointerCapture(event.pointerId)

  drawing = true
  last_point = point
  strokes.value = [...strokes.value, { points: [point], radius: brush.value }]

  paint_mask(point, point)
  erase_preview(point, point)
}

function onPointerMove(event: PointerEvent): void {
  if (panning.value && pan_start) {
    const element = viewport.value

    if (element) {
      element.scrollLeft = pan_start.left - (event.clientX - pan_start.x)
      element.scrollTop = pan_start.top - (event.clientY - pan_start.y)
    }

    return
  }

  update_cursor(event)

  if (!drawing || !last_point) return

  const point = point_from(event)
  if (!point) return

  strokes.value.at(-1)?.points.push(point)

  paint_mask(last_point, point)
  erase_preview(last_point, point)

  last_point = point
}

function onPointerUp(event: PointerEvent): void {
  if (panning.value) {
    panning.value = false
    pan_start = null
    viewport.value?.releasePointerCapture(event.pointerId)
    return
  }

  if (!drawing) return

  drawing = false
  last_point = null
  viewport.value?.releasePointerCapture(event.pointerId)

  void rebuild()
}

function undo(): void {
  strokes.value = strokes.value.slice(0, -1)
  replay_mask()
  redraw()
  void rebuild()
}

function clear(): void {
  strokes.value = []
  replay_mask()
  redraw()
  void rebuild()
}

async function prepare(): Promise<void> {
  strokes.value = []
  reset_urls()

  const [cutout, light, dark] = await Promise.all([
    load_image(props.cover.cutout_url),
    load_image(props.cover.silhouette_light_url),
    load_image(props.cover.silhouette_dark_url)
  ])

  cutout_image = cutout
  light_image = light
  dark_image = dark
  image_size.value = { width: cutout.naturalWidth, height: cutout.naturalHeight }

  const buffer = document.createElement('canvas')
  buffer.width = cutout.naturalWidth
  buffer.height = cutout.naturalHeight
  mask = buffer

  redraw()
}

watch(() => props.cover, () => void prepare(), { immediate: true })

let observer: ResizeObserver | null = null

onMounted(() => {
  if (!viewport.value) return

  observer = new ResizeObserver((entries) => {
    const entry = entries[0]
    if (!entry) return

    viewport_size.value = {
      width: entry.contentRect.width,
      height: entry.contentRect.height
    }
  })

  observer.observe(viewport.value)
})

onBeforeUnmount(() => {
  observer?.disconnect()
  observer = null

  for (const url of derived_urls) URL.revokeObjectURL(url)
  derived_urls = []
})
</script>

<style scoped>
.checker {
  background-image: conic-gradient(
    from 90deg at 50% 50%,
    var(--color-shade) 25%,
    var(--color-background) 0 50%,
    var(--color-shade) 0 75%,
    var(--color-background) 0
  );
  background-size: 16px 16px;
}
</style>
