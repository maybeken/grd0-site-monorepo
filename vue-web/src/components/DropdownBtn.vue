<template>
  <div class="relative z-20" @mouseleave="expanded = false">
    <button
      class="p-2 md:p-4 border-accent border-solid border-[1px] rounded-lg"
      @click="expanded = !expanded" @mouseenter="expanded = true"
    >
      <slot />
    </button>

    <div v-show="expanded" class="absolute z-10 bg-background motion-preset-slide-down motion-duration-500">
      <div class="flex flex-col border-accent border-solid border-[1px] rounded-lg">
        <div class="p-2">
          <input 
            type="text" 
            class="w-full rounded-md bg-transparent p-1 text-center border-accent border-solid border-[1px]" 
            placeholder="Search..."
            v-model="searchText"
            @input="onSearch"
          />
        </div>

        <template v-for="item in filteredList" :key="item.link">
          <button
            class="w-48 p-2 hover:bg-shade rounded-lg"
            @click="navigateTo(item.link)"
          >
            {{ item.displayName }}
            <span v-if="item.link[0] !== '/'" class="absolute right-2">
              <v-icon name="md-openinnew-outlined"></v-icon>
            </span>
          </button>
        </template>

        <button
          v-if="filteredList.length === 0"
          class="w-48 p-2 hover:bg-shade rounded-lg disabled:bg-background"
          disabled
        >
          No Item Available
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';

interface Props {
  items: {
    displayName: string;
    link: string;
  }[];
}

const $props = defineProps<Props>();
const $router = useRouter();
const expanded = ref(false);
const searchText = ref('');
const filteredList = ref($props.items);

const navigateTo = (url: string) => {
  if(url[0] == '/') {
    $router.push(url);
  } else {
    window.open(url, '_blank')?.focus();
  }

  expanded.value = !expanded.value;
};

const onSearch = () => {
  // TODO: Support search from blog post, gallery category, gallery image details
  filteredList.value = $props.items.filter((item) => { return item.displayName.toLowerCase().includes(searchText.value.toLowerCase()) });
};
</script>