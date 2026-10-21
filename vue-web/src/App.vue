<script setup lang="ts">
import { RouterView } from 'vue-router'
import SiteFooter from './components/layout/SiteFooter.vue'
import ToastContainer from './components/layout/ToastContainer.vue'
import UpdateNotification from './components/layout/UpdateNotification.vue'
import ParticleBackground from './components/layout/ParticleBackground.vue'
</script>

<template>
  <ParticleBackground />

  <header class="sticky top-0 from-background from-90% bg-linear-to-b py-2 md:pb-4 px-2 md:px-0 z-50">
    <Navbar></Navbar>
  </header>

  <main class="container mx-auto py-4 px-2 md:px-0 min-h-lvh">
    <RouterView />
  </main>

  <footer class="container mx-auto px-2 md:px-0 pb-8">
    <SiteFooter></SiteFooter>
  </footer>

  <UpdateNotification />
  <ToastContainer />
</template>
