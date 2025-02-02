import { Tldraw } from 'tldraw'
import 'tldraw/tldraw.css'

function App() {
  return (
    <div style={{ position: 'fixed', inset: 0 }}>
      <Tldraw
        inferDarkMode={true}
        onMount={(editor) => {
          editor.updateInstanceState({ isGridMode: true })
        }}
      />
    </div>
  )
}

export default App
