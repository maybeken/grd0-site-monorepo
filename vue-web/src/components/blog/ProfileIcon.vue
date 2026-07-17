<template>
  <div
    class="border-solid border-accent border-4 rounded-full overflow-clip h-full w-full bg-secondary"
  >
    <CDNImage
      v-if="email"
      class="motion-preset-fade"
      :resolution="256"
      :quality="75"
      :uri="`/profile/${$props?.email}.jpg`"
    ></CDNImage>
  </div>
</template>

<script setup lang="ts">
interface Props {
  email: string
}

const { email } = defineProps<Props>()
</script>
