import { nextTick } from 'vue';
import { createRouter, createWebHistory } from 'vue-router';

import HomeView from '@/views/HomeView.vue';
import NotFound from '@/views/NotFoundView.vue';
import Editor from '@/views/BlogEditorView.vue';
import Gallery from '@/views/GalleryView.vue';
import Blog from '@/views/BlogView.vue';
import Map from '@/views/MapView.vue';

const DEFAULT_TITLE = import.meta.env.VITE_DEFAULT_TITLE;

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  scrollBehavior: (to, from, savedPosition) => {
    if (to.hash) {
      return {
        el: to.hash,
        behavior: 'smooth',
      }
    }
  },
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: {
        title: 'Home',
      },
    },
    {
      path: '/blog/:slug',
      name: 'blog',
      component: Blog,
      meta: {
        title: 'Blog',
      },
    },
    {
      path: '/blog/editor',
      name: 'Blog Editor',
      meta: {
        title: 'Blog Editor',
      },
      component: Editor,
    },
    {
      path: '/gallery',
      name: 'Gallery',
      meta: {
        title: 'Gallery',
      },
      component: Gallery,
    },
    {
      path: '/travel/map',
      name: 'map',
      meta: {
        title: 'Traveler\'s Map',
      },
      component: Map,
    },
    {
      path: '/:pathMatch(.*)',
      name: '404-not-found',
      meta: {
        title: 'Page Not Found',
      },
      component: NotFound
    },
  ]
});

router.afterEach((to, from) => {
  // Use next tick to handle router history correctly
  // see: https://github.com/vuejs/vue-router/issues/914#issuecomment-384477609
  nextTick(() => {
      document.title = to.meta.title ? `${DEFAULT_TITLE} | ${to.meta.title}` : DEFAULT_TITLE;
  });
});

export default router
