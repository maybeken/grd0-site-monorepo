import { Tldraw, getSnapshot, loadSnapshot } from 'tldraw'
import 'tldraw/tldraw.css'



function App() {
  return (
    <div style={{ position: 'fixed', inset: 0 }}>
      <Tldraw
        colorScheme='dark'
        onMount={(editor) => {
          window.addEventListener(
            'keydown',
            (event) => {
              if ((event.ctrlKey || event.metaKey) && event.key === 's') {
                event.preventDefault();
                
                const { document } = getSnapshot(editor.store)

                window.parent.postMessage({
                  save: document,
                }, origin)
              }
            },
            false,
          );

          window.addEventListener('message', ({data}) => {
            if (data.event === 'command') {
              if (data.func === 'save') {
                const { document } = getSnapshot(editor.store)

                window.parent.postMessage({
                  save: document,
                }, origin)
              } else if (data.func === 'load') {
                loadSnapshot(editor.store, JSON.parse(data.payload))
              }
            }
          });

          editor.updateInstanceState({ isGridMode: true })
          editor.store.listen(
            (update) => {
              if (!window.parent) return;
              window.parent.postMessage(update, origin)
            },
            { scope: 'document', source: 'user' }
          )
        }}
      />
    </div>
  )
}

export default App
