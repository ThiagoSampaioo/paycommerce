<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-100 to-white flex items-center justify-center p-6">
    <div
      v-if="user"
      class="w-full max-w-6xl bg-white rounded-2xl shadow-xl p-8 animate-fade-in space-y-6"
    >
      <!-- Cabeçalho -->
      <div>
        <h1 class="text-3xl font-bold text-gray-800 mb-1">Olá, {{ user.name }}</h1>
        <p class="text-sm text-gray-500">Bem-vindo de volta à sua área administrativa</p>
      </div>

      <!-- Informações do Usuário -->
      <div class="border-t pt-4 space-y-3 text-gray-700">
        <div class="flex items-center gap-2">
          <Mail class="w-5 h-5 text-gray-500" />
          <span class="font-medium text-gray-600">Email:</span>
          <span>{{ user.email }}</span>
        </div>
        <div class="flex items-center gap-2">
          <Shield class="w-5 h-5 text-gray-500" />
          <span class="font-medium text-gray-600">Perfil:</span>
          <span>
            {{ user.role === 'admin' ? 'Administrador do sistema' : 'Cliente da plataforma' }}
          </span>
        </div>
      </div>

      <!-- Ações -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 pt-6 border-t">
        <RouterLink
          to="/providers"
          class="btn-icon bg-blue-600 hover:bg-blue-700 text-white"
        >
          <Settings class="w-5 h-5 mr-2" />
          Provedor de Pagamento
        </RouterLink>

        <RouterLink
          to="/admin/compras"
          class="btn-icon bg-purple-600 hover:bg-purple-700 text-white"
        >
          <ShoppingBag class="w-5 h-5 mr-2" />
          Vendas
        </RouterLink>

        <RouterLink
          to="/"
          class="btn-icon bg-green-600 hover:bg-green-700 text-white"
        >
          <ShoppingBag class="w-5 h-5 mr-2" />
          Voltar à Loja
        </RouterLink>

        <button
          @click="logout"
          class="btn-icon bg-gray-100 text-red-600 hover:bg-red-600 hover:text-white col-span-full sm:col-span-2 lg:col-span-1"
        >
          <LogOut class="w-5 h-5 mr-2" />
          Sair da Conta
        </button>
      </div>
    </div>

    <div v-else class="text-gray-500 animate-pulse text-center">Carregando...</div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Mail, Shield, Settings, ShoppingBag, LogOut } from 'lucide-vue-next'
import { useAuth } from '@/composables/useAuth'

const router = useRouter()
const { user, loadUser, logout: doLogout } = useAuth()

const logout = () => {
  doLogout()
}

// Garante que loadUser será chamado após montagem
onMounted(() => {
  loadUser()

  // Se não tiver token, redireciona
  if (!localStorage.getItem('token')) {
    router.push('/login')
  }
})
</script>

<style scoped>
@keyframes fade-in {
  from {
    opacity: 0;
    transform: translateY(12px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-fade-in {
  animation: fade-in 0.4s ease-out;
}

.btn-icon {
  @apply w-full flex items-center justify-center px-4 py-3 rounded-xl font-medium transition;
}
</style>
