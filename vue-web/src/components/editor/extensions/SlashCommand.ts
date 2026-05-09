import { Extension } from '@tiptap/vue-3'
import Suggestion from '@tiptap/suggestion'

interface SlashCommandItem {
  title: string
  description: string
  command: (props: { editor: any; range: { from: number; to: number } }) => void
}

const commands: SlashCommandItem[] = [
  {
    title: 'Heading 1',
    description: 'Large heading',
    command: ({ editor, range }) => {
      editor.chain().focus().deleteRange(range).toggleHeading({ level: 1 }).run()
    },
  },
  {
    title: 'Heading 2',
    description: 'Medium heading',
    command: ({ editor, range }) => {
      editor.chain().focus().deleteRange(range).toggleHeading({ level: 2 }).run()
    },
  },
  {
    title: 'Heading 3',
    description: 'Small heading',
    command: ({ editor, range }) => {
      editor.chain().focus().deleteRange(range).toggleHeading({ level: 3 }).run()
    },
  },
  {
    title: 'Bullet List',
    description: 'Create a bullet list',
    command: ({ editor, range }) => {
      editor.chain().focus().deleteRange(range).toggleBulletList().run()
    },
  },
  {
    title: 'Ordered List',
    description: 'Create an ordered list',
    command: ({ editor, range }) => {
      editor.chain().focus().deleteRange(range).toggleOrderedList().run()
    },
  },
  {
    title: 'Blockquote',
    description: 'Create a blockquote',
    command: ({ editor, range }) => {
      editor.chain().focus().deleteRange(range).toggleBlockquote().run()
    },
  },
  {
    title: 'Code Block',
    description: 'Add a code block',
    command: ({ editor, range }) => {
      editor.chain().focus().deleteRange(range).toggleCodeBlock().run()
    },
  },
  {
    title: 'Horizontal Rule',
    description: 'Add a divider',
    command: ({ editor, range }) => {
      editor.chain().focus().deleteRange(range).setHorizontalRule().run()
    },
  },
]

function renderItems(
  popup: HTMLDivElement,
  items: SlashCommandItem[],
  selectedIndex: number,
  commandRef: ((item: SlashCommandItem) => void) | null
) {
  popup.innerHTML = ''

  items.forEach((item, index) => {
    const el = document.createElement('div')
    el.className =
      'flex flex-col px-3 py-2 rounded-lg cursor-pointer text-sm ' +
      (index === selectedIndex ? 'bg-gray-700 text-white' : 'text-gray-300 hover:bg-gray-800')

    const title = document.createElement('span')
    title.className = 'font-medium'
    title.textContent = item.title

    const desc = document.createElement('span')
    desc.className = 'text-xs text-gray-500'
    desc.textContent = item.description

    el.appendChild(title)
    el.appendChild(desc)

    el.addEventListener('click', () => {
      if (commandRef) {
        commandRef(item)
      }
    })

    popup.appendChild(el)
  })
}

function updatePosition(popup: HTMLDivElement, clientRect?: () => DOMRect) {
  if (!clientRect) return

  const rect = clientRect()
  if (!rect) return

  popup.style.left = `${rect.left + window.scrollX}px`
  popup.style.top = `${rect.bottom + window.scrollY + 4}px`
}

const SlashCommand = Extension.create({
  name: 'slashCommand',

  addOptions() {
    return {
      suggestion: {
        char: '/',
        command: ({ editor, range, props }: { editor: any; range: { from: number; to: number }; props: SlashCommandItem }) => {
          props.command({ editor, range })
        },
        items: ({ query }: { query: string }) => {
          return commands.filter((item) =>
            item.title.toLowerCase().startsWith(query.toLowerCase())
          )
        },
        render: () => {
          let popup: HTMLDivElement | null = null
          let items: SlashCommandItem[] = []
          let selectedIndex = 0
          let commandRef: ((item: SlashCommandItem) => void) | null = null

          return {
            onStart: (props: { items: SlashCommandItem[]; command: (item: SlashCommandItem) => void; editor: any; clientRect?: () => DOMRect }) => {
              items = props.items
              selectedIndex = 0
              commandRef = props.command

              popup = document.createElement('div')
              popup.className =
                'absolute z-50 bg-gray-900 border border-gray-700 rounded-xl p-1 shadow-xl min-w-[200px]'
              popup.style.position = 'absolute'

              renderItems(popup, items, selectedIndex, commandRef)
              document.body.appendChild(popup)
              updatePosition(popup, props.clientRect)
            },

            onUpdate: (props: { items: SlashCommandItem[]; command: (item: SlashCommandItem) => void; clientRect?: () => DOMRect }) => {
              items = props.items
              selectedIndex = 0
              commandRef = props.command

              if (popup) {
                renderItems(popup, items, selectedIndex, commandRef)
                updatePosition(popup, props.clientRect)
              }
            },

            onKeyDown: (props: { event: KeyboardEvent }) => {
              if (props.event.key === 'ArrowDown') {
                selectedIndex = (selectedIndex + 1) % items.length
                if (popup) renderItems(popup, items, selectedIndex, commandRef)
                return true
              }

              if (props.event.key === 'ArrowUp') {
                selectedIndex = (selectedIndex - 1 + items.length) % items.length
                if (popup) renderItems(popup, items, selectedIndex, commandRef)
                return true
              }

              if (props.event.key === 'Enter') {
                const selected = items[selectedIndex]
                if (commandRef && selected) {
                  commandRef(selected)
                }
                return true
              }

              return false
            },

            onExit: () => {
              if (popup) {
                document.body.removeChild(popup)
                popup = null
              }
              items = []
              selectedIndex = 0
              commandRef = null
            },
          }
        },
      },
    }
  },

  addProseMirrorPlugins() {
    return [
      Suggestion({
        editor: this.editor,
        ...this.options.suggestion,
      }),
    ]
  },
})

export { SlashCommand }
export default SlashCommand
