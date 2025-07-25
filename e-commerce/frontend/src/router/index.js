import { createRouter, createWebHistory } from 'vue-router'
import LoginPage from '../pages/LoginPage.vue'
import RegisterPage from '../pages/RegisterPage.vue'
import HomePage from '../pages/HomePage.vue'
import DashboardPage from '../pages/DashboardPage.vue'
import ConfirmarCompra from '../pages/ConfirmarCompra.vue'
import ProvidersPage from '@/pages/ProvidersPage.vue'
import ComprasPage from '@/pages/ComprasPage.vue' // Minhas compras
import ComprasAdmin from '@/pages/ComprasAdmin.vue' // Todas as compras

const routes = [
  {
    path: '/',
    name: 'home',
    component: HomePage
  },
  {
    path: '/login',
    name: 'login',
    component: LoginPage
  },
  {
    path: '/register',
    name: 'register',
    component: RegisterPage
  },
  {
    path: '/confirmar',
    name: 'confirmar',
    component: ConfirmarCompra,
    meta: { requiresAuth: true }
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    component: DashboardPage,
    meta: { requiresAuth: true }
  },
  {
    path: '/providers',
    name: 'providers',
    component: ProvidersPage,
    meta: { requiresAuth: true, role: 'admin' }
  },
  {
    path: '/compras',
    name: 'compras',
    component: ComprasPage, // Minhas compras
    meta: { requiresAuth: true }
  },
  {
    path: '/admin/compras',
    name: 'compras-admin',
    component: ComprasAdmin, // Todas as compras
    meta: { requiresAuth: true, role: 'admin' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Guardião de rotas
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const user = JSON.parse(localStorage.getItem('user') || '{}')
  const isAuthenticated = !!token

  if (to.meta.requiresAuth && !isAuthenticated) {
    return next({ name: 'login' })
  }

  // Evita acesso a login/register estando logado
  if ((to.name === 'login' || to.name === 'register') && isAuthenticated) {
    return next({ name: 'dashboard' })
  }

  // Verificação de role para rotas de admin
  if (to.meta.role && to.meta.role !== user.role) {
    return next({ name: 'dashboard' })
  }

  next()
})

export default router
