<script setup lang="ts">
import { RouterView } from 'vue-router'
import SiteFooter from './components/SiteFooter.vue'
</script>

<template>
  <header class="sticky top-0 from-background from-90% bg-linear-to-b py-2 md:pb-4 z-50">
    <Navbar></Navbar>
  </header>

  <main class="container mx-auto py-4 min-h-lvh">
    <RouterView />
  </main>

  <footer class="container mx-auto pb-8">
    <SiteFooter></SiteFooter>
  </footer>
</template>
