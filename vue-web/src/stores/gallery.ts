import { defineStore } from 'pinia';

export const useGalleryStore = defineStore('gallery', {
  state: () => {
    return { selected_collection: 'all' }; 
  },
})
