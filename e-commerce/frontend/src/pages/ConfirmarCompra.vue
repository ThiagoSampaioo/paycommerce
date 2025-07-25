<script setup>
import { ref, onMounted } from 'vue'
import { useCartStore } from '@/stores/cart'
import { useAuth } from '@/composables/useAuth'
import { useRouter } from 'vue-router'
import axios from 'axios'

const cart = useCartStore()
const { user } = useAuth()
const router = useRouter()

const loading = ref(false)
const erro = ref(null)
const providerApiKey = ref('')
const providerName = ref('')
const providerURLPayment = ref('')

const token = localStorage.getItem('token')

onMounted(async () => {
  try {
    const { data } = await axios.get('http://localhost:8084/api/providers/active', {
      headers: { Authorization: `Bearer ${token}` }
    })

    providerApiKey.value = data.api_key
    providerName.value = data.name
    providerURLPayment.value = data.payment_url
  } catch (err) {
    console.error('Erro ao carregar provedor de pagamento:', err)
    erro.value = 'Não foi possível carregar o provedor de pagamento.'
  }
})

const finalizarCompra = async () => {
  if (!user.value || providerApiKey.value === '') {
    erro.value = 'Você precisa estar logado com uma empresa válida.'
    return
  }

  const total = parseFloat(cart.totalGeral.toFixed(2))
  const referenceId = `REF-${Date.now()}`
  const items = cart.items.map(item => ({
    name: item.title,
    price: item.price,
    quantity: item.quantity
  }))

  try {
    loading.value = true

    await axios.post('http://localhost:8084/initiate-payment', {
      user_id: user.value.id,
      reference_id: referenceId,
      total,
      items,
      provider_api_key: providerApiKey.value
    })

    const redirectUrl = encodeURIComponent(window.location.origin + '/compras') // ou alguma rota da loja
const url = `${providerURLPayment.value}?amount=${total}&referenceId=${referenceId}&apiKey=${providerApiKey.value}&redirectUrl=${redirectUrl}`
window.location.href = url
    cart.limparCarrinho()
  } catch (err) {
    console.error(err)
    erro.value = 'Erro ao iniciar pagamento.'
  } finally {
    loading.value = false
  }
}
const cancelarCompra = () => {
  cart.limparCarrinho()
  router.push('/')
}

</script>

<template>
  <div class="min-h-screen bg-gray-100 px-6 py-10">
    <div class="max-w-4xl mx-auto">
      <h1 class="text-3xl font-bold text-gray-800 mb-2">Resumo do Pedido</h1>
      <p class="text-sm text-gray-500 mb-6">
        Pagando com: <span class="font-medium text-blue-700">{{ providerName || 'Carregando...' }}</span>
      </p>

      <div v-if="cart.items.length === 0" class="text-center py-20 text-gray-500">
  <svg xmlns="http://www.w3.org/2000/svg" class="mx-auto mb-4 h-16 w-16 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2 9h12l-2-9m-6 0v-4a1 1 0 112 0v4" />
  </svg>
  <p class="text-lg font-medium mb-4">Seu carrinho está vazio</p>
  <button
    @click="router.push('/')"
    class="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition"
  >
    Voltar à loja
  </button>
</div>


      <div v-else class="space-y-4">
        <div
          v-for="item in cart.items"
          :key="item.id"
          class="bg-white rounded-lg shadow-md p-4 flex items-center gap-5"
        >
          <img :src="item.image" alt="Produto" class="w-24 h-24 object-contain" />
          <div class="flex-1">
            <h2 class="text-lg font-semibold text-gray-800 ">{{ item.title }}</h2>
            <p class="text-sm text-gray-500">Quantidade: {{ item.quantity }}</p>
            <p class="text-md font-bold text-green-700">
              R$ {{ (item.price * item.quantity).toFixed(2) }}
            </p>
          </div>
          <button
            @click="cart.removerItem(item.id)"
            class="text-red-500 text-sm hover:underline"
          >
            Remover
          </button>
        </div>

        <!-- Total e Ações -->
        <div class="bg-white rounded-lg shadow-md p-6 mt-6">
          <div class="flex justify-between items-center mb-4">
            <span class="text-xl font-semibold text-gray-700">Total</span>
            <span class="text-2xl font-semibold text-gray-700">
              R$ {{ cart.totalGeral.toFixed(2) }}
            </span>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 mt-4">
            <button
              @click="router.push('/')"
              class="w-full py-2 px-4 border border-gray-400 rounded-md text-gray-700 hover:bg-gray-100 transition"
            >
              Adicionar Mais Itens
            </button>

            <button
              @click="cancelarCompra"
              class="w-full py-2 px-4 bg-red-500 text-white rounded-md hover:bg-red-600 transition"
            >
              Cancelar Compra
            </button>

            <button
              @click="finalizarCompra"
              :disabled="loading"
              class="w-full py-2 px-4 bg-green-600 text-white rounded-md hover:bg-green-700 transition disabled:opacity-50"
            >
              {{ loading ? 'Processando...' : 'Finalizar Compra' }}
            </button>
          </div>

          <p v-if="erro" class="text-red-600 mt-4 text-sm text-center">{{ erro }}</p>
        </div>
      </div>
    </div>
  </div>
</template>
