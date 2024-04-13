<template>
    <nav class="container mx-auto">
      <div class="flex w-full gap-4">
        <div class="content-center">
          <DropdownBtn :items="menuItems">
            <v-icon name="md-menu"></v-icon>
          </DropdownBtn>
        </div>
        <div class="content-center mx-auto">
          <p class="text-xl text-secondary font-extrabold font-mono uppercase">Where Everything Starts And Ends</p>
        </div>
        <div class="content-center"></div>
      </div>
    </nav>
</template>

<script setup lang="ts">
const menuItems = [
  { displayName: 'Home', link: '/' },
  { displayName: 'Gallery', link: '/gallery' },
  { displayName: 'Twitter', link: 'https://x.com/maybeken0' },
];
</script>