// src/composables/useAuth.js
import { ref, computed } from 'vue'

const user = ref(null)

function loadUser() {
  const token = localStorage.getItem('token')
  if (token) {
    try {
      const payload = JSON.parse(atob(token.split('.')[1]))
      user.value = {
        id: payload.user_id,
        role: payload.role,
        name: payload.name,
        email: payload.email
      }
    } catch (error) {
      console.error('Erro ao carregar usuário:', error)
      user.value = null
    }
  } else {
    user.value = null
  }
}

const isAdmin = computed(() => user.value?.role === 'admin')
const isCliente = computed(() => user.value?.role === 'cliente')

function logout() {
  localStorage.removeItem('token')
  user.value = null
  window.location.href = '/login'
}

export function useAuth() {
  return { user, loadUser, isAdmin, isCliente, logout }
}
