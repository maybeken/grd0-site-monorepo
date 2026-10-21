import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'

import type { Ref } from 'vue'

interface Anchor {
  id: string
  top: number
  height: number
}

export function useScrollSync(
  editor: Ref<HTMLElement | null>,
  preview: Ref<HTMLElement | null>,
  enabled: Ref<boolean>
) {
  const active_block_id = ref('')

  let frame = 0
  let observer: ResizeObserver | null = null
  let detach: (() => void) | null = null

  function measure(pane: HTMLElement): Anchor[] {
    const pane_rect = pane.getBoundingClientRect()
    const anchors: Anchor[] = []

    pane.querySelectorAll<HTMLElement>('[data-block-id]').forEach((element) => {
      const rect = element.getBoundingClientRect()

      anchors.push({
        id: element.dataset.blockId || '',
        top: rect.top - pane_rect.top + pane.scrollTop,
        height: rect.height
      })
    })

    return anchors
  }

  function anchorAt(anchors: Anchor[], position: number): Anchor | null {
    if (anchors.length === 0) return null

    let candidate: Anchor | null = null

    for (const anchor of anchors) {
      if (anchor.top <= position + 1) {
        candidate = anchor
      } else {
        break
      }
    }

    return candidate || anchors[0] || null
  }

  function updateActive(): void {
    const source = editor.value
    if (!source) return

    const anchor = anchorAt(measure(source), source.scrollTop)
    if (anchor) active_block_id.value = anchor.id
  }

  function sync(force = false): void {
    if (!force && !enabled.value) return

    const source = editor.value
    const target = preview.value
    if (!source || !target) return

    const source_anchors = measure(source)
    const target_anchors = measure(target)
    if (source_anchors.length === 0 || target_anchors.length === 0) return

    const position = source.scrollTop
    const source_anchor = anchorAt(source_anchors, position)
    if (!source_anchor) return

    active_block_id.value = source_anchor.id

    const target_anchor =
      target_anchors.find((anchor) => anchor.id === source_anchor.id) ||
      target_anchors[target_anchors.length - 1]
    if (!target_anchor) return

    const ratio =
      source_anchor.height > 0
        ? Math.min(Math.max((position - source_anchor.top) / source_anchor.height, 0), 1)
        : 0
    const max = Math.max(target.scrollHeight - target.clientHeight, 0)
    const next = Math.min(target_anchor.top + ratio * target_anchor.height, max)

    target.scrollTop = next
  }

  function onEditorScroll(): void {
    if (frame) return

    frame = window.requestAnimationFrame(() => {
      frame = 0

      if (!enabled.value) {
        updateActive()
        return
      }

      sync()
    })
  }

  function realign(force = false): void {
    if (!force && !enabled.value) return

    window.requestAnimationFrame(() => sync(force))
  }

  function scrollPreviewToActive(): void {
    const target = preview.value
    if (!target) return

    const anchor = measure(target).find((item) => item.id === active_block_id.value)
    if (!anchor) return

    const max = Math.max(target.scrollHeight - target.clientHeight, 0)
    target.scrollTop = Math.min(Math.max(anchor.top, 0), max)
  }

  function reset(): void {
    active_block_id.value = ''

    if (frame) {
      window.cancelAnimationFrame(frame)
      frame = 0
    }

    if (preview.value) preview.value.scrollTop = 0
  }

  function attach(): void {
    detach?.()
    detach = null

    observer?.disconnect()
    observer = null

    const source = editor.value
    const target = preview.value
    if (!source || !target) return

    source.addEventListener('scroll', onEditorScroll, { passive: true })
    detach = () => source.removeEventListener('scroll', onEditorScroll)

    if (typeof ResizeObserver !== 'undefined') {
      observer = new ResizeObserver(() => realign())
      if (source.firstElementChild) observer.observe(source.firstElementChild)
      if (target.firstElementChild) observer.observe(target.firstElementChild)
    }
  }

  watch([editor, preview], attach, { immediate: true, flush: 'post' })

  onMounted(() => {
    nextTick(attach)
  })

  onBeforeUnmount(() => {
    detach?.()
    observer?.disconnect()

    if (frame) window.cancelAnimationFrame(frame)
  })

  return { active_block_id, realign, reset, scrollPreviewToActive }
}
