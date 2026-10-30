import { env, pipeline, RawImage } from '@huggingface/transformers'

type IncomingMessage =
  | {
      type: 'segment'
      id: number
      model_host: string
      model_id: string
      runtime_host: string
      blob: Blob
    }

type Segmenter = (image: RawImage) => Promise<RawImage>

const scope = self as unknown as {
  postMessage: (message: unknown) => void
  onmessage: ((event: MessageEvent<IncomingMessage>) => void) | null
}

const MAX_EDGE = 2000
const MATTE_THRESHOLD = 128
const MIN_COMPONENT_RATIO = 0.0005
const LIGHT: readonly [number, number, number] = [237, 237, 237]
const DARK: readonly [number, number, number] = [10, 10, 10]

let segmenter: Segmenter | null = null
let loading: Promise<Segmenter> | null = null

function post(message: unknown): void {
  scope.postMessage(message)
}

async function select_device(): Promise<'webgpu' | 'wasm' | null> {
  if (typeof WebAssembly === 'undefined') return null

  const gpu = (navigator as unknown as { gpu?: { requestAdapter(): Promise<unknown | null> } }).gpu

  if (gpu) {
    try {
      const adapter = await gpu.requestAdapter()
      if (adapter) return 'webgpu'
    } catch {
      return 'wasm'
    }
  }

  return 'wasm'
}

function needs_legacy_wasm(): boolean {
  const agent = navigator.userAgent
  const is_safari = /Safari/.test(agent) && !/Chrome|Chromium|Edg|OPR|Firefox|Android/.test(agent)

  if (!is_safari) return false

  const version = Number(/Version\/(\d+)/.exec(agent)?.[1] ?? 0)

  return version < 26 && !('gpu' in navigator)
}

async function load_segmenter(
  model_host: string,
  model_id: string,
  runtime_host: string
): Promise<Segmenter> {
  if (segmenter) return segmenter
  if (loading) return loading

  loading = (async () => {
    const device = await select_device()

    if (!device) {
      post({ type: 'unsupported' })
      throw new Error('WebAssembly is unavailable in this browser.')
    }

    env.allowLocalModels = false
    env.allowRemoteModels = true
    env.useBrowserCache = true
    env.remoteHost = model_host
    env.remotePathTemplate = '{model}'

    if (env.backends.onnx.wasm) {
      const suffix = needs_legacy_wasm() ? '' : '.asyncify'

      env.backends.onnx.wasm.wasmPaths = {
        mjs: `${runtime_host}ort-wasm-simd-threaded${suffix}.mjs`,
        wasm: `${runtime_host}ort-wasm-simd-threaded${suffix}.wasm`
      }
    }

    const instance = await pipeline('background-removal', model_id, {
      device,
      dtype: 'q8',
      progress_callback: (event: { status?: string; progress?: number }) => {
        if (event.status === 'progress' && typeof event.progress === 'number') {
          post({ type: 'progress', progress: Math.round(event.progress) })
        }
      }
    })

    segmenter = instance as unknown as Segmenter
    post({ type: 'ready', device })

    return segmenter
  })()

  try {
    return await loading
  } finally {
    loading = null
  }
}

function extract_matte(cutout: RawImage): Uint8Array {
  const rgba = cutout.rgba()
  const matte = new Uint8Array(rgba.width * rgba.height)

  for (let i = 0; i < matte.length; i += 1) {
    matte[i] = rgba.data[i * 4 + 3] ?? 0
  }

  return matte
}

function binarize(matte: Uint8Array): Uint8Array {
  const binary = new Uint8Array(matte.length)

  for (let i = 0; i < matte.length; i += 1) {
    binary[i] = (matte[i] ?? 0) >= MATTE_THRESHOLD ? 1 : 0
  }

  return binary
}

function fill_holes(binary: Uint8Array, width: number, height: number): void {
  const reached = new Uint8Array(binary.length)
  const stack = new Int32Array(binary.length)
  let top = 0

  const seed = (index: number): void => {
    if (binary[index] === 0 && reached[index] === 0) {
      reached[index] = 1
      stack[top] = index
      top += 1
    }
  }

  for (let x = 0; x < width; x += 1) {
    seed(x)
    seed((height - 1) * width + x)
  }

  for (let y = 0; y < height; y += 1) {
    seed(y * width)
    seed(y * width + width - 1)
  }

  while (top > 0) {
    top -= 1
    const index = stack[top] ?? 0
    const x = index % width
    const y = (index - x) / width

    if (x > 0) seed(index - 1)
    if (x < width - 1) seed(index + 1)
    if (y > 0) seed(index - width)
    if (y < height - 1) seed(index + width)
  }

  for (let i = 0; i < binary.length; i += 1) {
    if (binary[i] === 0 && reached[i] === 0) binary[i] = 1
  }
}

function remove_specks(binary: Uint8Array, width: number, height: number): void {
  const labels = new Int32Array(binary.length)
  const stack = new Int32Array(binary.length)
  const areas: number[] = [0]
  const min_area = Math.max(256, Math.floor(width * height * MIN_COMPONENT_RATIO))

  let label = 0
  let largest_label = 0
  let largest_area = 0

  for (let start = 0; start < binary.length; start += 1) {
    if (binary[start] === 0 || labels[start] !== 0) continue

    label += 1
    let top = 0
    let area = 0
    stack[top] = start
    top += 1
    labels[start] = label

    while (top > 0) {
      top -= 1
      const index = stack[top] ?? 0
      const x = index % width
      const y = (index - x) / width
      area += 1

      const visit = (next: number): void => {
        if (binary[next] === 1 && labels[next] === 0) {
          labels[next] = label
          stack[top] = next
          top += 1
        }
      }

      if (x > 0) visit(index - 1)
      if (x < width - 1) visit(index + 1)
      if (y > 0) visit(index - width)
      if (y < height - 1) visit(index + width)
    }

    areas[label] = area

    if (area > largest_area) {
      largest_area = area
      largest_label = label
    }
  }

  for (let i = 0; i < binary.length; i += 1) {
    if (binary[i] === 0) continue

    const own = labels[i] ?? 0
    const area = areas[own] ?? 0

    if (own !== largest_label && area < min_area) binary[i] = 0
  }
}

function silhouette(
  binary: Uint8Array,
  width: number,
  height: number,
  color: readonly [number, number, number]
): RawImage {
  const data = new Uint8ClampedArray(width * height * 4)

  for (let i = 0; i < binary.length; i += 1) {
    if (binary[i] === 0) continue

    const offset = i * 4
    data[offset] = color[0]
    data[offset + 1] = color[1]
    data[offset + 2] = color[2]
    data[offset + 3] = 255
  }

  return new RawImage(data, width, height, 4)
}

async function segment(message: Extract<IncomingMessage, { type: 'segment' }>): Promise<void> {
  try {
    const instance = await load_segmenter(
      message.model_host,
      message.model_id,
      message.runtime_host
    )

    let source = await RawImage.fromBlob(message.blob)
    const longest = Math.max(source.width, source.height)

    if (longest > MAX_EDGE) {
      const scale = MAX_EDGE / longest
      source = await source.resize(
        Math.max(1, Math.round(source.width * scale)),
        Math.max(1, Math.round(source.height * scale))
      )
    }

    const cutout = await instance(source)
    const matte = extract_matte(cutout)
    const binary = binarize(matte)

    fill_holes(binary, cutout.width, cutout.height)
    remove_specks(binary, cutout.width, cutout.height)

    const light = silhouette(binary, cutout.width, cutout.height, LIGHT)
    const dark = silhouette(binary, cutout.width, cutout.height, DARK)

    const [cutout_blob, light_blob, dark_blob] = await Promise.all([
      cutout.toBlob('image/png'),
      light.toBlob('image/png'),
      dark.toBlob('image/png')
    ])

    post({
      type: 'result',
      id: message.id,
      cutout: cutout_blob,
      silhouette_light: light_blob,
      silhouette_dark: dark_blob,
      width: cutout.width,
      height: cutout.height
    })
  } catch (error) {
    post({
      type: 'error',
      id: message.id,
      message: error instanceof Error ? error.message : 'Segmentation failed.'
    })
  }
}

scope.onmessage = (event: MessageEvent<IncomingMessage>) => {
  if (event.data.type === 'segment') {
    void segment(event.data)
  }
}
