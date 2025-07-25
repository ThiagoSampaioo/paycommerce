import { createRouter, createWebHistory } from 'vue-router'
import Login from '../pages/Login.vue'
import Registro from '../pages/Registro.vue'
import Dashboard from '../pages/Dashboard.vue'
import Confirmacao from '../pages/Confirmacao.vue'

const routes = [
  { path: '/', redirect: '/login' },
  { path: '/registro', component: Registro },
  { path: '/login', component: Login },
  {
    path: '/empresa',
    component: Dashboard,
    meta: { requiresAuth: true } // Protegida por login
  },
  { path: '/confirmacao', component: Confirmacao } // público
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Proteção de rotas com JWT
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
