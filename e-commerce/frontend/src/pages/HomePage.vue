<template>
  <div class="min-h-screen flex flex-col bg-gray-50 text-gray-800">
    <!-- Header -->
    <header class="bg-white shadow px-6 py-4 flex justify-between items-center border-b border-gray-200">
      <RouterLink to="/" class="text-2xl font-bold text-blue-700 hover:text-orange-700 transition hidden sm:block">
        SHOP
      </RouterLink>

      <div class="flex items-center gap-3">
        <template v-if="user">
          <span class="hidden sm:block text-sm text-gray-600">Olá, {{ user.name }}</span>

          <RouterLink v-if="isCliente" to="/compras" class="btn btn-info text-sm">Minhas Compras</RouterLink>
          <RouterLink v-if="isAdmin" to="/dashboard" class="btn btn-success text-sm">Dashboard</RouterLink>
          <RouterLink to="/confirmar" class="btn btn-info text-sm flex items-center gap-1">
            <ShoppingCart class="w-4 h-4" />
            Carrinho ({{ cart.totalItens }})
          </RouterLink>
          <button @click="logout" class="btn btn-outline text-sm flex items-center gap-1">
            <LogOut class="w-4 h-4" />
            Sair
          </button>
        </template>

        <template v-else>
          <RouterLink to="/login" class="btn btn-primary">Login</RouterLink>
          <RouterLink to="/register" class="btn btn-outline">Cadastrar</RouterLink>
        </template>
      </div>
    </header>

    <!-- Conteúdo -->
    <main class="flex-grow px-6 py-10">
      <h2 class="text-3xl font-bold text-gray-900 mb-10 text-center">Produtos em Destaque</h2>

      <div v-if="loading" class="text-center text-gray-500">Carregando produtos...</div>
      <div v-else-if="erro" class="text-center text-red-600">Erro ao carregar produtos.</div>

      <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        <div
          v-for="produto in produtos"
          :key="produto.id"
          class="bg-white rounded-xl shadow hover:shadow-lg transition p-4 flex flex-col"
        >
          <img
            :src="produto.image"
            alt="Produto"
            class="w-full h-44 object-contain mb-4 transition-transform duration-300 hover:scale-105"
          />
          <h3 class="font-semibold text-gray-800 text-base mb-1 line-clamp-2">{{ produto.title }}</h3>
          <p class="text-xs text-gray-500 mb-1 capitalize">{{ produto.category }}</p>
          <p class="text-lg font-bold text-blue-800 mb-3">R$ {{ produto.price.toFixed(2) }}</p>

          <div class="mt-auto flex flex-col gap-2">
            <button @click="cart.adicionarItem(produto, 1)" class="btn btn-outline text-sm">
              Adicionar ao Carrinho
            </button>
            <button @click="comprarAgora(produto)" class="btn btn-primary text-sm">
              Comprar Agora
            </button>
          </div>
        </div>
      </div>
    </main>

    <!-- Rodapé -->
    <footer class="bg-white text-center py-4 text-sm text-gray-500 border-t border-gray-200">
      &copy; 2025 PayCommerce. Todos os direitos reservados.
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuth } from '@/composables/useAuth'
import { useCartStore } from '@/stores/cart'
import { RouterLink, useRouter } from 'vue-router'
import { ShoppingCart, LogOut } from 'lucide-vue-next'

const produtos = ref([])
const loading = ref(true)
const erro = ref(false)

const { user, loadUser, isAdmin, isCliente, logout } = useAuth()
const cart = useCartStore()
const router = useRouter()

onMounted(async () => {
  loadUser()
  try {
    const res = await fetch('https://fakestoreapi.com/products')
    produtos.value = await res.json()
  } catch {
    erro.value = true
  } finally {
    loading.value = false
  }
})

const comprarAgora = (produto) => {
  cart.limparCarrinho()
  cart.adicionarItem(produto, 1)
  router.push('/confirmar')
}
</script>

<style scoped>
/* Botões base */
.btn {
  @apply px-4 py-2 rounded-full font-medium transition duration-200 shadow-sm text-center;
}

/* Botão Laranja - Ação principal */
.btn-primary {
  background-color: #016ac0;
  color: white;
}
.btn-primary:hover {
  background-color: #03508f;
}

/* Botão Azul - Carrinho ou info */
.btn-info {
  background-color: #9333ea;
  color: white;
}
.btn-info:hover {
  background-color: #5b16c9;
}

/* Botão Verde - Dashboard/Admin */
.btn-success {
  background-color: #2563eb;
  color: white;
}
.btn-success:hover {
  background-color: #0540c0;
}

/* Botão Cinza com borda - Alternativo/Outline */
.btn-outline {
  background-color: transparent;
  border: 1px solid #6B7280;
  color: #374151;
}
.btn-outline:hover {
  background-color: #F3F4F6;
}

/* Estilo para texto com limite de linhas */
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
 