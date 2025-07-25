<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-100 to-white p-6">
    <div class="max-w-4xl mx-auto space-y-8 animate-fade-in">
      <!-- Topo -->
      <div class="flex justify-between items-center">
        <h2 class="text-3xl font-bold text-gray-700">Provedor de Pagamentos</h2>
        <RouterLink to="/dashboard" class="btn-outline text-sm">
          <ArrowLeft class="w-4 h-4 mr-1" />
          Voltar
        </RouterLink>
      </div>

      <!-- Formulário de cadastro -->
      <form
        v-if="!providers.length"
        @submit.prevent="criarProvider"
        class="bg-white p-6 rounded-xl shadow space-y-6 border border-gray-200"
      >
        <h3 class="text-xl font-semibold text-gray-800">Cadastrar novo provedor</h3>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="form-label">Nome</label>
            <input v-model="novo.name" type="text" class="form-input" required />
          </div>
          <div>
            <label class="form-label">API Key</label>
            <input v-model="novo.api_key" type="text" class="form-input" required />
          </div>
          <div class="md:col-span-2">
            <label class="form-label">URL de Pagamento</label>
            <input v-model="novo.payment_url" type="url" class="form-input" required />
          </div>
          <div class="md:col-span-2">
            <label class="form-label">URL de Cancelamento</label>
            <input v-model="novo.cancellation_url" type="url" class="form-input" required />
          </div>
          <div class="flex items-center gap-2 md:col-span-2 mt-2">
            <input v-model="novo.active" type="checkbox" id="active" class="w-4 h-4" />
            <label for="active" class="text-sm text-gray-700">Marcar como ativo</label>
          </div>
        </div>

        <button type="submit" class="btn-primary">
          <Save class="w-4 h-4 mr-2" />
          Cadastrar Provedor
        </button>
      </form>

      <!-- Lista de Provedores -->
      <div class="bg-white p-6 rounded-xl shadow border border-gray-200">
        <h3 class="text-xl font-semibold text-gray-800 mb-4">Provedor Cadastrado</h3>
        <ul class="divide-y divide-gray-200">
          <li
            v-for="p in providers"
            :key="p.ID"
            class="py-4 flex flex-col sm:flex-row sm:justify-between sm:items-center gap-3"
          >
            <div class="text-gray-700 space-y-1">
              <p class="font-semibold text-lg">Nome: {{ p.name }}</p>
              <p class="font-semibold text-black-600">URL de Pagamento: </p>
              <p class="text-sm break-all flex items-center gap-1 text-gray-600">{{ p.payment_url }}</p>
              <p class="font-semibold text-black-600">URL de Cancelamento: </p>
              <p class="text-sm break-all flex items-center gap-1 text-gray-500">{{ p.cancellation_url }}</p>
              <p class="font-semibold text-black-600">API Key: </p>
              <p class="text-xs break-all flex items-center gap-1 text-gray-500">{{ p.api_key }}</p>
            </div>

            <div class="flex flex-col items-end gap-2">
              <span
                :class="p.active ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'"
                class="text-xs font-semibold px-3 py-1 rounded-full"
              >
                {{ p.active ? 'Ativo' : 'Inativo' }}
              </span>
              <div class="flex gap-2 text-sm">
                <button @click="abrirEdicao(p)" class="btn-icon bg-blue-600 hover:bg-blue-700 text-white flex items-center gap-1 p-2 rounded">
                  <Edit class="w-4 h-4" /> Editar
                </button>
                <button v-if="p.active" @click="desativarProvider(p)" class="btn-icon bg-yellow-600 hover:bg-yellow-700 text-white flex items-center gap-1 p-2 rounded">
                  <Ban class="w-4 h-4" /> Desativar
                </button>
                <button @click="excluirProvider(p.ID)" class="btn-icon bg-red-600 hover:bg-red-700 text-white flex items-center gap-1 p-2 rounded">
                  <Trash2 class="w-4 h-4" /> Excluir
                </button>
              </div>
            </div>
          </li>
        </ul>
      </div>
    </div>

    <!-- Modal de Edição -->
    <div v-if="providerEditando" class="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50">
      <div class="bg-white p-6 rounded-xl shadow max-w-lg w-full space-y-4">
        <h3 class="text-xl font-bold text-gray-800">Editar Provedor</h3>

        <div class="space-y-3">
          <div>
            <label class="form-label">Nome</label>
            <input v-model="providerEditando.name" type="text" class="form-input" />
          </div>
          <div>
            <label class="form-label">API Key</label>
            <input v-model="providerEditando.api_key" type="text" class="form-input" />
          </div>
          <div>
            <label class="form-label">URL de Pagamento</label>
            <input v-model="providerEditando.payment_url" type="url" class="form-input" />
          </div>
          <div>
            <label class="form-label">URL de Cancelamento</label>
            <input v-model="providerEditando.cancellation_url" type="url" class="form-input" />
          </div>
          <div class="flex items-center gap-2">
            <input v-model="providerEditando.active" type="checkbox" id="edit-active" class="w-4 h-4" />
            <label for="edit-active" class="text-sm">Ativo</label>
          </div>
        </div>

        <div class="flex justify-end gap-2 pt-4">
          <button @click="providerEditando = null" class="btn-outline">
            <X class="w-4 h-4 mr-1" /> Cancelar
          </button>
          <button @click="editarProvider" class="btn-primary">
            <Save class="w-4 h-4 mr-1" /> Salvar
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { RouterLink } from 'vue-router'
import {
  ArrowLeft,
  Edit,
  Trash2,
  Save,
  Ban,
  Key,
  Globe,
  X
} from 'lucide-vue-next'

const token = localStorage.getItem('token')
const providers = ref([])
const novo = ref({
  name: '',
  api_key: '',
  payment_url: '',
  cancellation_url: '',
  active: false
})
const providerEditando = ref(null)

const carregarProviders = async () => {
  try {
    const res = await axios.get('http://localhost:8084/api/providers', {
      headers: { Authorization: `Bearer ${token}` }
    })
    providers.value = res.data
  } catch (err) {
    console.error('Erro ao buscar provedores:', err)
  }
}

const criarProvider = async () => {
  try {
    await axios.post('http://localhost:8084/api/providers', novo.value, {
      headers: { Authorization: `Bearer ${token}` }
    })
    novo.value = { name: '', api_key: '', payment_url: '', cancellation_url: '', active: false }
    await carregarProviders()
  } catch (err) {
    alert('Erro ao criar provedor')
    console.error(err)
  }
}

const abrirEdicao = (p) => {
  providerEditando.value = { ...p }
}

const editarProvider = async () => {
  try {
    await axios.put(`http://localhost:8084/api/providers/${providerEditando.value.ID}`, providerEditando.value, {
      headers: { Authorization: `Bearer ${token}` }
    })
    providerEditando.value = null
    await carregarProviders()
  } catch (err) {
    alert('Erro ao editar provedor')
    console.error(err)
  }
}

const desativarProvider = async (p) => {
  if (!confirm(`Deseja desativar o provedor ${p.name}?`)) return

  try {
    await axios.post(`http://localhost:8084/api/providers/${p.ID}/deactivate`, {}, {
      headers: { Authorization: `Bearer ${token}` }
    })
    await carregarProviders()
  } catch (err) {
    alert('Erro ao desativar provedor')
    console.error(err)
  }
}

const excluirProvider = async (id) => {
  if (!confirm('Deseja realmente excluir este provedor?')) return

  try {
    await axios.delete(`http://localhost:8084/api/providers/${id}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    await carregarProviders()
  } catch (err) {
    alert('Erro ao excluir provedor')
    console.error(err)
  }
}

onMounted(carregarProviders)
</script>

<style scoped>
.animate-fade-in {
  animation: fade-in 0.4s ease-out;
}
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

.form-label {
  @apply block text-sm font-medium text-gray-600 mb-1;
}

.form-input {
  @apply w-full border border-gray-300 px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500;
}

.btn-primary {
  @apply inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition;
}

.btn-outline {
  @apply inline-flex items-center px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-100 transition;
}
</style>
