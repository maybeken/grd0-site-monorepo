import { nextTick } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import { formatCollectionName } from '@/helpers/collection'
import { getGalleryCollectionRaw } from '@/services/gallery'
import { getBlogPostRaw } from '@/services/blogPost'

const DEFAULT_TITLE = import.meta.env.VITE_DEFAULT_TITLE

declare module 'vue-router' {
  interface RouteMeta {
    title?: string
    url?: string
    requiresAuth?: boolean
  }
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  scrollBehavior: (to, from, savedPosition) => {
    if (to.hash) {
      return {
        el: to.hash,
        behavior: 'smooth'
      }
    }

    return { top: 0, behavior: 'smooth' }
  },
  routes: [
    {
      path: '/',
      name: 'home',
      redirect: '/gallery',
      meta: {
        title: 'Home'
      }
    },
    {
      path: '/blog',
      name: 'blog',
      component: () => import('@/views/HomeView.vue'),
      meta: {
        title: 'Blog'
      }
    },
    {
      path: '/blog/:slug',
      name: 'Blog',
      component: () => import('@/views/BlogView.vue'),
      meta: {
        title: 'Blog'
      }
    },
    {
      path: '/blog/editor',
      name: 'Blog Editor',
      meta: {
        title: 'Blog Editor',
        requiresAuth: true
      },
      component: () => import('@/views/BlogEditorView.vue')
    },
    {
      path: '/gallery',
      name: 'Gallery Index',
      meta: {
        title: 'Gallery'
      },
      children: [
        {
          name: 'Gallery',
          path: '/gallery/:collection',
          component: () => import('@/views/GalleryView.vue')
        }
      ],
      component: () => import('@/views/GalleryView.vue')
    },
    {
      path: '/travel/map',
      name: 'map',
      meta: {
        title: "Traveler's Map"
      },
      component: () => import('@/views/MapView.vue')
    },
    {
      path: '/music',
      name: 'music',
      meta: { title: 'Music Matters' },
      component: () => import('@/views/MusicView.vue')
    },
    {
      path: '/coffee',
      name: 'coffee',
      meta: { title: 'Coffee Taste Log' },
      component: () => import('@/views/CoffeeView.vue')
    },
    {
      path: '/coffee/editor',
      name: 'Coffee Tasting Editor',
      meta: {
        title: 'Coffee Tasting Editor',
        requiresAuth: true
      },
      component: () => import('@/views/CoffeeEditorView.vue'),
    },
    {
      path: '/redirect',
      component: () => import('@/views/RedirectView.vue'),
      meta: {
        title: 'Redirecting...'
      },
      beforeEnter(to, from) {
        if (!to.meta.url) return

        window.location.href = to.meta.url
      },
      children: [
        {
          path: 'auth/login',
          name: 'login',
          component: {},
          meta: {
            url: `${import.meta.env.VITE_API_URL}/auth/login`
          }
        },
        {
          path: 'auth/callback',
          name: 'callback',
          component: {},
          beforeEnter(to, from, next) {
            if (!to.query.token || typeof to.query.token !== 'string') {
              next({ name: 'login' })
              return
            }

            const token = to.query.token
            const expires_at = Number(to.query.expires_at)

            sessionStorage.setItem('jwt_token', token)
            sessionStorage.setItem('jwt_expires', String(expires_at))
            next({ path: '/' })
          }
        }
      ]
    },
    {
      path: '/:pathMatch(.*)',
      name: '404-not-found',
      meta: {
        title: 'Page Not Found'
      },
      component: () => import('@/views/NotFoundView.vue')
    },
    {
      path: '/tool',
      name: 'Tool',
      children: [
        {
          path: 'dotmatrix',
          name: 'Dot Matrix',
          meta: {
            title: 'Dot Matrix Generator'
          },
          component: () => import('@/views/DotMatrix.vue')
        },
        {
          path: 'draw',
          name: 'Tldraw',
          meta: {
            title: 'tldraw;'
          },
          component: () => import('@/views/TldrawView.vue')
        }
      ]
    }
  ]
})

router.afterEach((to, from) => {
  // Use next tick to handle router history correctly
  // see: https://github.com/vuejs/vue-router/issues/914#issuecomment-384477609
  nextTick(() => {
    document.title = to.meta.title ? `${DEFAULT_TITLE} | ${to.meta.title}` : DEFAULT_TITLE
  })
})

router.beforeEach(async (to, from, next) => {
  if (to.meta.requiresAuth) {
    const token = sessionStorage.getItem('jwt_token')
    const expires = sessionStorage.getItem('jwt_expires')

    if (!token || !expires || Date.now() >= Number(expires)) {
      next({ name: 'login' })
      return
    }
  }

  const name = to.name
  const params = to.params
  const title = to.meta.title

  if (name === 'Blog' && params.slug) {
    const slug = typeof params.slug === 'string' ? params.slug : params.slug[0]

    if (slug) {
      const response = await getBlogPostRaw(slug)
      const blog_title = response?.title || ''

      to.meta.title = `${title} - ${blog_title}`
    }
  } else if (name === 'Gallery' && params.collection) {
    const response = await getGalleryCollectionRaw()
    const collection =
      typeof params.collection === 'string' ? params.collection : params.collection[0]
    const album = formatCollectionName(response, collection)

    to.meta.title = `${title} - ${album}`
  }

  next()
})

export default router
