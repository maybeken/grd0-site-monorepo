<template>
    <nav class="container mx-auto">
      <div class="flex w-full gap-4 h-16">
        <div class="content-center">
          <DropdownBtn :items="menuItems">
            <Icon icon="mynaui:menu" height="auto" />
          </DropdownBtn>
        </div>
        
        <div class="content-center mx-auto">
          <p class="text-base sm:text-xl lg:text-2xl text-secondary font-extrabold uppercase"><span class="hidden sm:inline">grd0.net | </span>Where Everything Starts And Ends</p>
        </div>
        <div class="content-center"></div>
      </div>

      <div class="hidden md:flex w-full place-content-center gap-4">
        <MenuLinkBtn v-for="item in menuItems.filter(item => !item.hidden && item.link.startsWith('/'))" :key="item.link" :display-name="item.displayName" :link="item.link"
          :hidden="item.hidden" :onNavigate="navigateTo"></MenuLinkBtn>
      </div>
    </nav>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';

const $router = useRouter();

const menuItems = [
  { displayName: 'Home', link: '/' },
  { displayName: 'Gallery', link: '/gallery' },
  { displayName: 'Traveler\'s Map', link: '/travel/map' },
  { displayName: 'LISTEN TO ME!', link: '/music' },
  { displayName: 'YouTube', link: 'https://www.youtube.com/@maybeken' },
  { displayName: 'X', link: 'https://x.com/maybeken0' },
  { displayName: 'Tool - Editor', link: '/blog/editor', hidden: true },
  { displayName: 'Tool - Dot Matrix', link: '/tool/dotmatrix', hidden: true },
];

const navigateTo = (url: string) => {
  if (url[0] == '/') {
    $router.push(url);
  } else {
    window.open(url, '_blank')?.focus();
  }
};
</script>