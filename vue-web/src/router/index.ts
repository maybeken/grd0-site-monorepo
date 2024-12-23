import { nextTick } from 'vue';
import { createRouter, createWebHistory } from 'vue-router';

import HomeView from '@/views/HomeView.vue';
import NotFound from '@/views/NotFound.vue';
import Editor from '@/views/Editor.vue';
import Gallery from '@/views/GalleryView.vue';

const DEFAULT_TITLE = import.meta.env.VITE_DEFAULT_TITLE;

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: {
        title: 'Home',
      },
      children: [
        {
          path: '/test',
          name: 'test',
          meta: {
            title: 'Test',
          },
          component: HomeView,
        }
      ],
    },
    {
      path: '/editor',
      name: 'Editor',
      meta: {
        title: 'Editor',
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
