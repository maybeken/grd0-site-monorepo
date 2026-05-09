import { reactive } from 'vue'

export interface ModalState {
  isOpen: boolean
  type: 'confirm' | 'prompt'
  title: string
  message: string
  inputValue: string
  inputPlaceholder: string
  resolve: ((value: boolean | string | null) => void) | null
}

const state = reactive<ModalState>({
  isOpen: false,
  type: 'confirm',
  title: '',
  message: '',
  inputValue: '',
  inputPlaceholder: '',
  resolve: null,
})

export function useModal() {
  const confirm = (message: string, title = 'Confirm'): Promise<boolean> => {
    return new Promise<boolean>((resolve) => {
      state.type = 'confirm'
      state.title = title
      state.message = message
      state.resolve = (value) => resolve(value as boolean)
      state.isOpen = true
    })
  }

  const prompt = (
    message: string,
    title = 'Input',
    placeholder = '',
    defaultValue = ''
  ): Promise<string | null> => {
    return new Promise<string | null>((resolve) => {
      state.type = 'prompt'
      state.title = title
      state.message = message
      state.inputValue = defaultValue
      state.inputPlaceholder = placeholder
      state.resolve = (value) => resolve(value as string | null)
      state.isOpen = true
    })
  }

  const open = () => {
    state.isOpen = true
  }

  const close = () => {
    if (state.resolve) {
      state.resolve(false)
      state.resolve = null
    }
    state.isOpen = false
  }

  const submit = () => {
    if (state.resolve) {
      if (state.type === 'prompt') {
        state.resolve(state.inputValue)
      } else {
        state.resolve(true)
      }
      state.resolve = null
    }
    state.isOpen = false
  }

  return {
    state,
    confirm,
    prompt,
    open,
    close,
    submit,
  }
}
