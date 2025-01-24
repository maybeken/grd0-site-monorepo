<script setup lang="ts">
import { RouterView } from 'vue-router'
import SiteFooter from './components/SiteFooter.vue';
</script>

<template>
  <header class="sticky top-0 bg-background px-4 py-2 z-50">
    <Navbar></Navbar>
  </header>

  <main class="container mx-auto py-4 min-h-lvh">
    <RouterView />
  </main>

  <footer class="container mx-auto pb-8">
    <SiteFooter></SiteFooter>
  </footer>
</template>