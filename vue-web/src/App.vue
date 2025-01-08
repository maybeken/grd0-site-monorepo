<script setup lang="ts">
import { RouterView } from 'vue-router'
import SiteFooter from './components/SiteFooter.vue';
</script>

<template>
  <header class="sticky top-0 bg-background px-4 py-2 z-50">
    <Navbar></Navbar>
  </header>

  <main class="container mx-auto py-4">
    <RouterView />
  </main>

  <footer class="container mx-auto pb-8">
    <SiteFooter></SiteFooter>
  </footer>
</template>

<style type="postcss">
body {
  @apply bg-background text-foreground px-2 font-['Doto'];
}

h1 {
  @apply text-2xl md:text-4xl;
}

h2 {
  @apply text-xl md:text-2xl;
}

p {
  @apply text-sm md:text-base;
}
</style>
