<template>
  <button
    v-if="!hidden"
    class="relative w-48 p-2 hover:bg-shade rounded-lg"
    @click="onNavigate(link)"
  >
    {{ displayName }}
    <span v-if="link[0] !== '/'" class="absolute right-2">
      <Icon icon="mynaui:external-link" height="auto" />
    </span>
  </button>
</template>

<script setup lang="ts">
interface Props {
  displayName: string;
  link: string;
  hidden?: boolean;
  onNavigate: Function;
}

const {
  displayName,
  link,
  hidden,
} = defineProps<Props>();
</script>