<template>
  <nav class="container mx-auto">
    <div class="flex w-full gap-4 h-16">
      <div class="content-center">
        <DropdownBtn :items="menuItems">
          <Icon icon="mynaui:menu" height="auto" />
        </DropdownBtn>
      </div>

      <div class="content-center mx-auto">
        <p class="text-base sm:text-xl lg:text-2xl text-secondary font-extrabold uppercase">
          <span class="hidden sm:inline">grd0.net</span
          ><span class="select-none hidden sm:inline"> | </span>{{ $t('nav.tagline') }}<CursorBlink />
        </p>
      </div>
      <div class="content-center"></div>
    </div>

    <div class="hidden md:flex w-full place-content-center gap-4">
      <MenuLinkBtn
        v-for="item in menuItems.filter((item) => !item.hidden && item.link.startsWith('/'))"
        :key="item.link"
        :display-name="$t(item.displayNameKey)"
        :link="item.link"
        :hidden="item.hidden"
        :onNavigate="navigateTo"
      ></MenuLinkBtn>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const $router = useRouter()

const menuItems = [
  { displayNameKey: 'nav.gallery', link: '/gallery' },
  { displayNameKey: 'nav.travelersMap', link: '/travel/map' },
  { displayNameKey: 'nav.beanSoup', link: '/coffee' },
  { displayNameKey: 'nav.blog', link: '/blog' },
  { displayNameKey: 'nav.listenToMe', link: '/music' },
  { displayNameKey: 'nav.instagram', link: 'https://www.instagram.com/maybe_ken/' },
  { displayNameKey: 'nav.youtube', link: 'https://www.youtube.com/@maybeken' },
  { displayNameKey: 'nav.linkedin', link: 'https://www.linkedin.com/in/maybeken/' },
  { displayNameKey: 'nav.buyMeCoffee', link: 'https://ko-fi.com/maybeken' },
  { displayNameKey: 'nav.cmsLogin', link: '/redirect/auth/login', hidden: true },
  { displayNameKey: 'nav.toolEditor', link: '/blog/editor', hidden: true },
  { displayNameKey: 'nav.toolTastingEditor', link: '/coffee/editor', hidden: true },
  { displayNameKey: 'nav.toolDotMatrix', link: '/tool/dotmatrix', hidden: true },
  { displayNameKey: 'nav.toolTldraw', link: '/tool/draw', hidden: true }
]

const navigateTo = (url: string) => {
  if (url[0] == '/') {
    $router.push(url)
  } else {
    window.open(url, '_blank')?.focus()
  }
}
</script>
