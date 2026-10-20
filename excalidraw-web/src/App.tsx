import { useCallback, useEffect, useRef, useState } from 'react'
import {
  CaptureUpdateAction,
  Excalidraw,
  restore,
  serializeAsJSON
} from '@excalidraw/excalidraw'
import type {
  ExcalidrawImperativeAPI,
  ExcalidrawInitialDataState,
  ExcalidrawProps
} from '@excalidraw/excalidraw/types'
import type { ImportedDataState } from '@excalidraw/excalidraw/data/types'
import '@excalidraw/excalidraw/index.css'

const STORAGE_KEY = 'grd0.draw.scene'
const SAVE_DELAY_MS = 500

function isExcalidrawScene(data: unknown): data is ImportedDataState {
  return (
    typeof data === 'object' &&
    data !== null &&
    (data as { type?: unknown }).type === 'excalidraw'
  )
}

function readStoredScene(): ExcalidrawInitialDataState | null {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return null

    const data: unknown = JSON.parse(raw)
    if (!isExcalidrawScene(data)) return null

    const { elements, appState, files } = restore(data, null, null)
    return {
      elements,
      appState: { ...appState, theme: 'dark' },
      files,
      scrollToContent: true
    }
  } catch {
    return null
  }
}

function App() {
  const apiRef = useRef<ExcalidrawImperativeAPI | null>(null)
  const saveTimer = useRef<number | undefined>(undefined)
  const [initialData] = useState(readStoredScene)

  const post = useCallback((message: unknown) => {
    window.parent.postMessage(message, window.location.origin)
  }, [])

  const save = useCallback(() => {
    const api = apiRef.current
    if (!api) return

    const json = serializeAsJSON(
      api.getSceneElements(),
      api.getAppState(),
      api.getFiles(),
      'local'
    )
    post({ save: JSON.parse(json) })
  }, [post])

  const load = useCallback(
    (payload: string) => {
      const api = apiRef.current
      if (!api) return

      let data: unknown
      try {
        data = JSON.parse(payload)
      } catch {
        data = null
      }

      if (!isExcalidrawScene(data)) {
        const legacy = typeof data === 'object' && data !== null && 'store' in data
        api.setToast({
          message: legacy
            ? 'This is a legacy tldraw drawing, which is no longer supported.'
            : 'Could not read this drawing.',
          duration: 6000
        })
        post({ loadError: legacy ? 'legacy' : 'invalid' })
        return
      }

      const scene = restore(data, null, null)
      api.updateScene({
        elements: scene.elements,
        appState: { ...scene.appState, theme: 'dark' },
        captureUpdate: CaptureUpdateAction.IMMEDIATELY
      })
      api.addFiles(Object.values(scene.files))
      api.scrollToContent(scene.elements, { fitToContent: true })
    },
    [post]
  )

  const handleChange = useCallback<NonNullable<ExcalidrawProps['onChange']>>(
    (elements, appState, files) => {
      window.clearTimeout(saveTimer.current)
      saveTimer.current = window.setTimeout(() => {
        try {
          localStorage.setItem(
            STORAGE_KEY,
            serializeAsJSON(elements, appState, files, 'local')
          )
        } catch (error) {
          if (error instanceof DOMException && error.name === 'QuotaExceededError') {
            apiRef.current?.setToast({
              message: 'Autosave failed: this drawing is too large for local storage.',
              duration: 6000
            })
          }
        }
      }, SAVE_DELAY_MS)
    },
    []
  )

  useEffect(() => {
    const onKeyDown = (event: KeyboardEvent) => {
      if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === 's') {
        event.preventDefault()
        save()
      }
    }

    const onMessage = (event: MessageEvent) => {
      if (event.source !== window.parent) return

      const data = event.data as
        | { event?: string; func?: string; payload?: unknown }
        | null
      if (data?.event !== 'command') return

      if (data.func === 'save') {
        save()
      } else if (data.func === 'load' && typeof data.payload === 'string') {
        load(data.payload)
      }
    }

    window.addEventListener('keydown', onKeyDown)
    window.addEventListener('message', onMessage)
    return () => {
      window.removeEventListener('keydown', onKeyDown)
      window.removeEventListener('message', onMessage)
    }
  }, [save, load])

  useEffect(() => {
    const flush = () => {
      const api = apiRef.current
      if (!api) return

      window.clearTimeout(saveTimer.current)
      try {
        localStorage.setItem(
          STORAGE_KEY,
          serializeAsJSON(
            api.getSceneElements(),
            api.getAppState(),
            api.getFiles(),
            'local'
          )
        )
      } catch {
        // Ignore storage errors when flushing on unload.
      }
    }

    window.addEventListener('pagehide', flush)
    return () => {
      window.removeEventListener('pagehide', flush)
      window.clearTimeout(saveTimer.current)
    }
  }, [])

  return (
    <div style={{ position: 'fixed', inset: 0 }}>
      <Excalidraw
        theme="dark"
        gridModeEnabled
        autoFocus
        initialData={initialData}
        excalidrawAPI={(api) => {
          apiRef.current = api
        }}
        onChange={handleChange}
      />
    </div>
  )
}

export default App
