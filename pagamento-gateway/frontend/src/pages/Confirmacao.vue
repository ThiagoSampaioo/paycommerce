<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-100 to-white px-4">
    <div class="bg-white shadow-xl rounded-xl w-full max-w-md p-8 animate-fade-in">
      <h1 class="text-2xl font-bold text-center text-gray-800 mb-6">Confirmação de Pagamento</h1>

      <div v-if="loading" class="text-center text-gray-500">
        <svg class="mx-auto mb-2 animate-spin h-6 w-6 text-blue-600" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z" />
        </svg>
        Carregando pagamento...
      </div>

      <div v-else-if="erro" class="text-center">
        <p class="text-red-600 font-semibold">❌ Erro: {{ erro }}</p>
      </div>

      <div v-else>
        <div class="space-y-3">
          <div class="flex items-center justify-between border-b pb-2">
            <span class="text-gray-600 font-medium">🔖 Pedido:</span>
            <span class="font-mono text-gray-800">{{ pagamento.reference_id }}</span>
          </div>

          <div class="flex items-center justify-between border-b pb-2">
            <span class="text-gray-600 font-medium">💰 Valor:</span>
            <span class="text-gray-800 font-semibold text-lg">R$ {{ pagamento.amount.toFixed(2) }}</span>
          </div>

          <div class="flex items-center justify-between">
            <span class="text-gray-600 font-medium">📦 Status:</span>
            <span
              :class="[ 'font-semibold', pagamento.status === 'PAID' ? 'text-green-600' : 'text-yellow-600' ]"
            >
              {{ pagamento.status === 'PAID' ? 'Pago com sucesso' : 'Aguardando pagamento' }}
            </span>
          </div>
        </div>

        <div v-if="pagamento.status === 'PENDING'" class="mt-6">
          <button
            @click="confirmar"
            class="w-full bg-green-600 text-white font-semibold py-2 rounded-lg hover:bg-green-700 transition"
          >
            ✅ Confirmar Pagamento
          </button>
        </div>

        <div v-else class="text-center text-green-600 font-semibold mt-6 space-y-4">
  ✅ Pagamento já confirmado!

  <div v-if="redirectUrl">
    <button
      @click="voltarParaLoja"
      class="mt-4 bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition"
    >
      🛍️ Voltar para a loja
    </button>
  </div>
</div>

      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  data() {
    return {
      pagamento: null,
      loading: true,
      erro: null,
      apiKey: '',
      paymentId: null
    }
  },
  async mounted() {
    console.log('Query params:', this.$route.query)
    const amount = parseFloat(this.$route.query.amount)
    const referenceId = this.$route.query.referenceId
    this.apiKey = this.$route.query.apiKey
    const redirectUrl = this.$route.query.redirectUrl
  
    if (!this.apiKey || isNaN(amount) || !referenceId) {
      this.erro = 'Parâmetros inválidos na URL.'
      this.loading = false
      return
    }

    try {
      const { data } = await axios.post('http://localhost:8080/api/checkout', {
        amount,
        reference_id: referenceId
      }, {
        headers: {
          'x-api-key': this.apiKey
        }
      })

      this.paymentId = data.payment_id
      this.pagamento = {
        amount,
        reference_id: referenceId,
        status: 'PENDING'
      }
      this.loading = false
      this.redirectUrl = redirectUrl || null
    } catch (err) {
      console.error(err)
      this.erro = err.response?.data?.error || 'Erro ao criar pagamento.'
    } finally {
      this.loading = false
    }
  },
  methods: {
  async confirmar() {
    try {
      await axios.post(`http://localhost:8080/api/confirmacao/${this.paymentId}`, {}, {
        headers: {
          'x-api-key': this.apiKey
        }
      })

      alert('Pagamento confirmado com sucesso!')
      this.pagamento.status = 'PAID'
    } catch (err) {
      console.error(err)
      alert('Erro ao confirmar pagamento.')
    }
  },
  voltarParaLoja() {
    if (this.$route.query.redirectUrl) {
      window.location.href = this.$route.query.redirectUrl
    } else {
      alert('URL de retorno não encontrada.')
    }
  }
}


}

</script>

<style scoped>
.animate-fade-in {
  animation: fadeIn 0.4s ease-out;
}
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
