<template>
    <nav class="container mx-auto">
      <div class="flex w-full gap-4">
        <div class="content-center">
          <DropdownBtn :items="menuItems">
            <v-icon name="md-menu"></v-icon>
          </DropdownBtn>
        </div>
        <div class="content-center mx-auto">
          <p class="text-base sm:text-xl lg:text-2xl text-secondary font-extrabold uppercase"><span class="hidden sm:inline">grd0.net | </span>Where Everything Starts And Ends</p>
        </div>
        <div class="content-center"></div>
      </div>
    </nav>
</template>

<script setup lang="ts">
const menuItems = [
  { displayName: 'Home', link: '/' },
  { displayName: 'Gallery', link: '/gallery' },
  { displayName: 'YouTube', link: 'https://www.youtube.com/@maybeken' },
  { displayName: 'X', link: 'https://x.com/maybeken0' },
];
</script>