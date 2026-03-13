import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/text-to-image',
  },
  {
    path: '/text-to-image',
    name: 'TextToImage',
    component: () => import('@/views/TextToImage.vue'),
  },
  {
    path: '/image-edit',
    name: 'ImageEdit',
    component: () => import('@/views/ImageEdit.vue'),
  },
  {
    path: '/cross-border',
    name: 'CrossBorder',
    component: () => import('@/views/CrossBorder.vue'),
  },
  {
    path: '/food-business',
    name: 'FoodBusiness',
    component: () => import('@/views/FoodBusiness.vue'),
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { guest: true },
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { guest: true },
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/Profile.vue'),
    meta: { requiresAuth: true },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token')

  if (to.meta.requiresAuth && !token) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else if (to.meta.guest && token) {
    next({ name: 'TextToImage' })
  } else {
    next()
  }
})

export default router
