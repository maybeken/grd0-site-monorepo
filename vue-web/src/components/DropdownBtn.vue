<template>
  <div class="relative z-20" @mouseleave="expanded = false">
    <button class="p-2 md:p-4 border-foreground border-dotted border rounded-lg" @click="expanded = !expanded"
      @mouseenter="expanded = true">
      <slot />
    </button>

    <div v-show="expanded"
      class="absolute z-10 bg-background motion-preset-slide-down motion-duration-500 bg-opacity-80 backdrop-blur-xs">
      <div class="flex flex-col border-foreground border-dotted border rounded-lg">
        <div class="p-2">
          <input type="text"
            class="w-full rounded-md bg-transparent p-1 text-center border-foreground border-dotted border"
            placeholder="Search..." v-model="searchText" @input="onSearch" />
        </div>

        <MenuLinkBtn v-for="item in filteredList" :key="item.link" :display-name="item.displayName" :link="item.link"
          :hidden="searchText == '' && item.hidden" :onNavigate="navigateTo"></MenuLinkBtn>

        <button v-if="filteredList.length === 0" class="w-48 p-2 hover:bg-shade rounded-lg disabled:bg-background"
          disabled>
          No Item Available
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import MenuLinkBtn from './MenuLinkBtn.vue';

interface Props {
  items: {
    displayName: string;
    link: string;
    hidden: boolean;
  }[];
}

const $props = defineProps<Props>();
const $router = useRouter();
const expanded = ref(false);
const searchText = ref('');
const filteredList = ref($props.items);

const navigateTo = (url: string) => {
  if (url[0] == '/') {
    $router.push(url);
  } else {
    window.open(url, '_blank')?.focus();
  }

  expanded.value = !expanded.value;
};

const onSearch = () => {
  // TODO: Support search from blog post, gallery collection, gallery image details
  filteredList.value = $props.items.filter((item) => { return item.displayName.toLowerCase().includes(searchText.value.toLowerCase()) });
};
</script>