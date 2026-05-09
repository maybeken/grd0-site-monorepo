import '@tiptap/core'

declare module '@tiptap/core' {
  interface Storage {
    markdown: {
      options: Record<string, unknown>
      getMarkdown(): string
    }
  }
}
