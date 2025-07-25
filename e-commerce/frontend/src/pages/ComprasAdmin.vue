<!-- src/pages/admin/ComprasAdmin.vue -->
<template>
  <div class="min-h-screen bg-gray-100 p-6">
    <div class="max-w-7xl mx-auto bg-white shadow rounded-xl p-6 space-y-4">
      <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <h1 class="text-2xl font-bold text-gray-800">Todas as Compras</h1>
        <router-link to="/compras" class="text-sm bg-gray-100 px-4 py-2 rounded hover:bg-gray-200 flex items-center gap-2">
          <ShoppingBag class="w-4 h-4 mr-1" />
          Ver Minhas Compras
        </router-link>
        <router-link to="/dashboard" class="text-sm bg-gray-100 px-4 py-2 rounded hover:bg-gray-200 flex items-center gap-2">
          <ArrowBigLeft class="w-4 h-4 mr-1" />
          Voltar
        </router-link>
      </div>

      <!-- Filtros -->
      <div class="flex flex-col md:flex-row md:items-center gap-4">
        <select v-model="filtroStatus" class="border px-4 py-2 rounded w-full md:w-64">
          <option value="">Todos os Status</option>
          <option value="AGUARDANDO_PAGAMENTO">Pendente</option>
          <option value="PAGO">Pago</option>
          <option value="CANCELADA">Cancelado</option>
        </select>
        <button @click="buscarCompras" class="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-700">Filtrar</button>
      </div>

      <!-- Lista -->
      <div v-if="loading" class="text-gray-500">Carregando...</div>
      <div v-else-if="compras.length === 0" class="text-gray-500">Nenhuma compra encontrada.</div>

      <div v-else class="overflow-x-auto">
        <table class="w-full table-auto text-sm">
          <thead class="bg-gray-100 text-left">
            <tr>
              <th class="px-4 py-2">#</th>
              <th class="px-4 py-2 hidden md:table-cell">Usuário</th>
              <th class="px-4 py-2 hidden md:table-cell">Referência</th>
              <th class="px-4 py-2 hidden md:table-cell">Valor</th>
              <th class="px-4 py-2">Status</th>
              <th class="px-4 py-2 hidden md:table-cell">Data</th>
              <th class="px-4 py-2 text-right">Ações</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="compra in compras" :key="compra.ID" class="border-b">
              <td class="px-4 py-2">{{ compra.ID }}</td>
              <td class="px-4 py-2 hidden md:table-cell">{{ compra.User?.name || '-' }}</td>
              <td class="px-4 py-2 hidden md:table-cell">{{ compra.reference_id }}</td>
              <td class="px-4 py-2 hidden md:table-cell">R$ {{ compra.amount }}</td>
              <td class="px-4 py-2"><span :class="badgeStatus(compra.status)">{{ compra.status }}</span></td>
              <td class="px-4 py-2 hidden md:table-cell">{{ formatarData(compra.CreatedAt) }}</td>
              <td class="px-4 py-2 text-right">
                <button @click="verDetalhes(compra)" class="text-blue-600 hover:underline"> detalhes</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Paginação -->
      <div class="flex justify-between items-center pt-4 border-t text-sm text-gray-500">
        <span>Total: {{ total }}</span>
        <div class="space-x-2">
          <button @click="pagina--" :disabled="pagina === 1" class="px-3 py-1 border rounded">Anterior</button>
          <button @click="pagina++" :disabled="pagina * limite >= total" class="px-3 py-1 border rounded">Próxima</button>
        </div>
      </div>
    </div>

    <!-- Modal -->
    <div v-if="compraSelecionada" class="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50">
      <div class="bg-white rounded-xl shadow p-6 w-full max-w-lg space-y-4">
        <h2 class="text-lg font-bold">Detalhes da Compra</h2>
        <p><strong>Usuário:</strong> {{ compraSelecionada.User?.name }}</p>
        <p><strong>Referência:</strong> {{ compraSelecionada.reference_id }}</p>
        <p><strong>Status:</strong> {{ compraSelecionada.status }}</p>
        <p><strong>Valor:</strong> R$ {{ compraSelecionada.amount }}</p>
        <p><strong>Data:</strong> {{ formatarData(compraSelecionada.CreatedAt) }}</p>
        <button @click="compraSelecionada = null" class="mt-4 bg-gray-100 px-4 py-2 rounded hover:bg-gray-200">Fechar</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'
import { ShoppingBag, ArrowBigLeft } from 'lucide-vue-next'

const compras = ref([])
const compraSelecionada = ref(null)
const filtroStatus = ref('')
const pagina = ref(1)
const limite = 10
const total = ref(0)
const loading = ref(true)

const badgeStatus = (status) => ({
  AGUARDANDO_PAGAMENTO: 'bg-yellow-100 text-yellow-800 px-2 py-1 rounded',
  PAGO: 'bg-green-100 text-green-800 px-2 py-1 rounded',
  CANCELADA: 'bg-red-100 text-red-800 px-2 py-1 rounded'
}[status] || '')

const formatarData = (data) =>
  new Date(data).toLocaleString('pt-BR', {
    day: '2-digit', month: '2-digit', year: 'numeric',
    hour: '2-digit', minute: '2-digit'
  })

const buscarCompras = async () => {
  loading.value = true
  const token = localStorage.getItem('token')
  try {
    const res = await axios.get('http://localhost:8084/api/sales/all', {
      headers: { Authorization: `Bearer ${token}` },
      params: {
        page: pagina.value,
        limit: limite,
        status: filtroStatus.value || undefined
      }
    })
    compras.value = res.data.data
    total.value = res.data.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

watch(pagina, buscarCompras)

const verDetalhes = (compra) => {
  compraSelecionada.value = compra
}



onMounted(buscarCompras)
</script>
