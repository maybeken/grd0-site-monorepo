import { nextTick } from 'vue';
import { createRouter, createWebHistory } from 'vue-router';
import { formatCategoryName } from '@/helpers/category';
import { getGalleryCategory } from '@/services/gallery';
import { getBlogPost } from '@/services/blogPost';

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
      name: 'Blog',
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
      path: '/music',
      name: 'music',
      meta: {
        title: 'Music Matters',
      },
      component: () => import('@/views/MusicView.vue'),
    },
    {
      path: '/:pathMatch(.*)',
      name: '404-not-found',
      meta: {
        title: 'Page Not Found',
      },
      component: () => import('@/views/NotFoundView.vue'),
    },
    {
      path: '/tool',
      name: 'Tool',
      children: [
        {
          path: 'dotmatrix',
          name: 'Dot Matrix',
          meta: {
            title: 'Dot Matrix Generator',
          },
          component: () => import('@/views/DotMatrix.vue'),
        },
      ],
    }
  ]
});

router.afterEach((to, from) => {
  // Use next tick to handle router history correctly
  // see: https://github.com/vuejs/vue-router/issues/914#issuecomment-384477609
  nextTick(() => {
      document.title = to.meta.title ? `${DEFAULT_TITLE} | ${to.meta.title}` : DEFAULT_TITLE;
  });
});

router.beforeEach(async(to, from, next) => {
  const name = to.name;
  const params = to.params;
  const title = to.meta.title;

  if (name === "Blog" && params.slug) {
    const slug = typeof params.slug === 'string' ? params.slug : params.slug[0];
    const response = getBlogPost(slug);
    const blog_title = response?.data.value.title || "";

    to.meta.title = `${title} - ${blog_title}`
  } else if (name === "Gallery" && params.category) {
    const response = getGalleryCategory();
    const category = typeof params.category === 'string' ? params.category : params.category[0];
    const album = formatCategoryName(response?.data.value, category);
    
    to.meta.title = `${title} - ${album}`
  }

  next();
})


export default router
