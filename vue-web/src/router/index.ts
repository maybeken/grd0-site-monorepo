import { nextTick } from 'vue';
import { createRouter, createWebHistory } from 'vue-router';

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
      component: () => import('@/views/HomeView.vue'),
      meta: {
        title: 'Home',
      },
    },
    {
      path: '/blog/:slug',
      name: 'blog',
      component: () => import('@/views/BlogView.vue'),
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
      component: () => import('@/views/BlogEditorView.vue'),
    },
    {
      path: '/gallery',
      name: 'Gallery Index',
      meta: {
        title: 'Gallery',
      },
      children: [
        {
          name: 'Gallery',
          path: '/gallery/:category',
          component: () => import('@/views/GalleryView.vue'),
        }
      ],
      component: () => import('@/views/GalleryView.vue'),
    },
    {
      path: '/travel/map',
      name: 'map',
      meta: {
        title: 'Traveler\'s Map',
      },
      component: () => import('@/views/MapView.vue'),
    },
    {
      path: '/:pathMatch(.*)',
      name: '404-not-found',
      meta: {
        title: 'Page Not Found',
      },
      component: () => import('@/views/NotFoundView.vue'),
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
