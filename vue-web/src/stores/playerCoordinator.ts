import { ref } from 'vue'
import { defineStore } from 'pinia'

export const usePlayerCoordinator = defineStore('playerCoordinator', () => {
  const activePlayerId = ref<string | null>(null)

  const play = (id: string) => {
    activePlayerId.value = id
  }

  const pause = () => {
    activePlayerId.value = null
  }

  return { activePlayerId, play, pause }
})
