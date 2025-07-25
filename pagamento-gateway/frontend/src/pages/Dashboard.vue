<template>
  <div class="min-h-screen bg-gray-100 p-6">
    <div class="max-w-4xl mx-auto">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-2xl font-bold text-gray-800">Dashboard</h2>
       
        <button @click="logout" class="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600 transition">
          Sair
        </button>
      </div>

      <div class="bg-white shadow-md rounded-lg p-6">
         <div v-if="plano?.is_active" class="flex flex-col sm:flex-row flex-wrap gap-2 mb-6">
        <div>
              <button @click="abrirModalCallback" class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">
                Atualizar Callback
              </button>
            </div>

            <div>
              <button @click="abrirModalBancario" class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">
                Atualizar Dados Bancários
              </button>
            </div>

            <div>
              <button @click="mostrarModalSaque = true" class="bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700">
                Realizar Saque
              </button>
            </div>
          </div>
        <h3 class="text-xl font-semibold mb-4 text-gray-700">Informações</h3>

        <div v-if="plano" class="space-y-6">
          <div>
            <div  v-if="plano?.is_active" class="flex items-center gap-4">
                <label class="font-medium text-gray-600">Saldo disponível:</label>
                <p class="text-lg font-bold text-green-700">R$ {{ saldo.toFixed(2) }}</p>
              </div>
            <p><span class="font-medium text-gray-600">Nome:</span> {{ plano.name }}</p>
            <p><span class="font-medium text-gray-600">Email:</span> {{ plano.email }}</p>
            <p>
              <span class="font-medium text-gray-600 mr-1">Status:</span>
              <span class="inline-block px-2 py-1 rounded text-white text-sm" :class="plano?.is_active ? 'bg-green-500' : 'bg-yellow-500'">
                {{ plano.is_active ? 'Ativo' : 'Inativo' }}
              </span>
            </p>
            <div class="mb-4">
              <p><span class="font-medium text-gray-600">Rota de Callback:</span> {{ plano.callback_url || 'Não definida' }}</p>
              <p class="text-sm text-gray-500">Use esta URL para receber notificações de eventos.</p>
            </div>

            <p v-if="plano?.bank_name">
              <span class="font-medium text-gray-600">Dados Bancários:</span>
              <div class="mt-2">
                <p><strong>Banco:</strong> {{ plano.bank_name }} ({{ plano.bank_code }})</p>
                <p><strong>Agência:</strong> {{ plano.agency_account }}</p>
                <p><strong>Conta:</strong> {{ plano.number_account }}</p>
                <p><strong>Tipo de Conta:</strong> {{ plano.type_account === 'CONTA_CORRENTE' ? 'Conta Corrente' : 'Poupança' }}</p>
                <p><strong>Tipo de Chave Pix:</strong> {{ plano.type_key_pix }}</p>
                <p><strong>Chave Pix:</strong> {{ plano.key_pix || 'Não definida' }}</p>
              </div>
            </p>

            
          </div>

          <div v-if="plano?.is_active" class="space-y-6">
            <div>
              <label class="block text-sm text-gray-600 mb-1">API Key:</label>
              <div class="flex items-center gap-2 bg-gray-100 px-3 py-2 rounded-md">
                <span class="text-sm font-mono break-all">{{ plano.api_key }}</span>
                <button @click="copiarApiKey" class="ml-auto text-blue-600 hover:underline text-sm">Copiar</button>
              </div>
            </div>

            

            

            <!-- Modal de Saque -->
            <div v-if="mostrarModalSaque" class="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50">
              <div class="bg-white p-6 rounded shadow-md w-full max-w-sm">
                <h3 class="text-lg font-semibold text-gray-700 mb-4">Realizar Saque</h3>
                <input v-model="valorSaque" type="number" placeholder="Valor" class="w-full px-3 py-2 rounded border border-gray-300 mb-4" />
                <div class="flex justify-end gap-2">
                  <button @click="mostrarModalSaque = false" class="px-4 py-2 rounded border">Cancelar</button>
                  <button @click="sacar" class="bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700">Confirmar</button>
                </div>
              </div>
            </div>

          </div>
          <!-- Plano Inativo -->
          <div v-else>
            <div class="border rounded-lg p-4 bg-gray-50">
              <p class="text-lg font-semibold mb-2">Plano Profissional</p>
              <ul class="text-sm text-gray-600 mb-4 list-disc pl-5">
                <li>Acesso à API com chave exclusiva</li>
                <li>Suporte prioritário</li>
                <li>Integração com checkout</li>
              </ul>
              <p class="text-2xl font-bold text-gray-800 mb-4">R$ 99,90 / mês</p>

              <label class="block text-sm font-medium text-gray-700 mb-1">Forma de pagamento:</label>
              <select
                v-model="formaPagamento"
                class="w-full px-3 py-2 rounded border border-gray-300 mb-4"
              >
                <option value="pix">Pix</option>
                <option value="cartao">Cartão de Crédito</option>
                <option value="boleto">Boleto</option>
              </select>

              <button
                @click="pagar"
                :disabled="carregando"
                class="w-full bg-green-600 text-white font-semibold py-2 rounded hover:bg-green-700 transition disabled:opacity-50 flex items-center justify-center gap-2"
              >
                <svg v-if="carregando" class="animate-spin h-5 w-5 text-white" viewBox="0 0 24 24" fill="none">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
                </svg>
                <span>{{ carregando ? 'Processando...' : 'Pagar e Ativar' }}</span>
              </button>
            </div>
          </div>
        

          <!-- Modal Callback -->
          <div v-if="mostrarModalCallback" class="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50">
            <div class="bg-white p-6 rounded shadow-md w-full max-w-md">
              <h3 class="text-lg font-semibold text-gray-700 mb-4">Atualizar Callback</h3>
              <input v-model="callbackUrl" placeholder="https://seusite.com/callback" class="w-full px-3 py-2 rounded border border-gray-300 mb-4" />
              <div class="flex justify-end gap-2">
                <button @click="mostrarModalCallback = false" class="px-4 py-2 rounded border">Cancelar</button>
                <button @click="atualizarCallback" class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">Salvar</button>
              </div>
            </div>
          </div>

          <!-- Modal Dados Bancários -->
          <div v-if="mostrarModalBancario" class="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50">
            <div class="bg-white p-6 rounded shadow-md w-full max-w-xl max-h-[90vh] overflow-y-auto">
              <h3 class="text-lg font-semibold text-gray-700 mb-4">Atualizar Dados Bancários</h3>
              <div class="grid grid-cols-2 gap-4">
                <input v-model="banco.bank_code" placeholder="Código Banco" class="border px-2 py-1 rounded" />
                <input v-model="banco.bank_name" placeholder="Nome Banco" class="border px-2 py-1 rounded" />
                <input v-model="banco.agency_account" placeholder="Agência" class="border px-2 py-1 rounded" />
                <input v-model="banco.number_account" placeholder="Conta" class="border px-2 py-1 rounded" />
                <select v-model="banco.type_account" class="border px-2 py-1 rounded">
                  <option value="CONTA_CORRENTE">Conta Corrente</option>
                  <option value="POUPANCA">Poupança</option>
                </select>
                <select v-model="banco.type_key_pix" class="border px-2 py-1 rounded">
                  <option value="CPF">CPF</option>
                  <option value="CNPJ">CNPJ</option>
                  <option value="EMAIL">Email</option>
                  <option value="TELEFONE">Telefone</option>
                  <option value="CHAVE_ALEATORIA">Aleatória</option>
                </select>
                <input v-model="banco.chave_pix" placeholder="Chave Pix" class="border px-2 py-1 rounded col-span-2" />
              </div>
              <div class="flex justify-end gap-2 mt-4">
                <button @click="mostrarModalBancario = false" class="px-4 py-2 rounded border">Cancelar</button>
                <button @click="atualizarBancario" class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">Salvar</button>
              </div>
            </div>
          </div>

               <!--  Listagem de pagamentos -->
          <div v-if="plano?.is_active" class="mt-10">
            <h4 class="text-lg font-semibold text-gray-700 mb-2">Pagamentos Recebidos</h4>

            <div class="flex items-center gap-4 mb-4">
              <select v-model="filtroStatus" @change="buscarPagamentos" class="border px-3 py-2 rounded">
                <option value="">Todos os Status</option>
                <option value="AGUARDANDO_PAGAMENTO">Pendente</option>
                <option value="PAGO">Pago</option>
                <option value="CANCELADO">Cancelado</option>
              </select>
            </div>

            <div class="overflow-x-auto bg-white rounded shadow">
              <table class="table-auto w-full text-sm">
                <thead class="bg-gray-100 text-gray-700">
                  <tr>
                    <th class="px-4 py-2 text-left">Referência</th>
                    <th class="px-4 py-2 text-left">Valor</th>
                    <th class="px-4 py-2 text-left">Status</th>
                    <th class="px-4 py-2 text-left">Data</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="pag in pagamentos" :key="pag.id" class="border-t">
                    <td class="px-4 py-2 font-mono">{{ pag.reference_id }}</td>
                    <td class="px-4 py-2">R$ {{ pag.amount.toFixed(2) }}</td>
                    <td class="px-4 py-2">
                      <span
                        :class="{
                          'text-yellow-600': pag.status === 'AGUARDANDO_PAGAMENTO',
                          'text-green-600': pag.status === 'PAGO',
                          'text-red-600': pag.status === 'CANCELADO'
                        }"
                        class="font-semibold"
                      >
                        {{ pag.status }}
                      </span>
                    </td>
                    <td class="px-4 py-2">{{ new Date(pag.created_at).toLocaleString('pt-BR') }}</td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div class="mt-4 flex justify-between items-center text-sm text-gray-600">
              <span>Total: {{ totalPagamentos }}</span>
              <div class="space-x-2">
                <button @click="pagina--" :disabled="pagina === 1" class="px-3 py-1 border rounded">Anterior</button>
                <button @click="pagina++" :disabled="pagina * limite >= totalPagamentos" class="px-3 py-1 border rounded">Próxima</button>
              </div>
            </div>
          </div>
        </div>
        <div v-else class="text-gray-500">Carregando informações do plano...</div>
      </div>
    </div>
  </div>
</template>


<script>
import axios from 'axios'

export default {
  data() {
    return {
      plano: null,
      saldo: 0,
      carregando: false,
      formaPagamento: 'pix',
      callbackUrl: '',
      valorSaque: 0,
      mostrarModalCallback: false,
      mostrarModalBancario: false,
      mostrarModalSaque: false,
      banco: {
        bank_code: '',
        bank_name: '',
        agency_account: '',
        number_account: '',
        type_account: 'CONTA_CORRENTE',
        type_key_pix: 'CPF',
        chave_pix: ''
      },
      pagamentos: [],
      filtroStatus: '',
      pagina: 1,
      limite: 10,
      totalPagamentos: 0
    }
  },
  mounted() {
    this.verPlano()
    this.buscarSaldo()
    this.buscarPagamentos()
  },
  watch: {
    pagina() {
      this.buscarPagamentos()
    }
  },
  methods: {
      abrirModalBancario() {
      this.banco = {
        bank_code: this.plano.bank_code || '',
        bank_name: this.plano.bank_name || '',
        agency_account: this.plano.agency_account || '',
        number_account: this.plano.number_account || '',
        type_account: this.plano.type_account || 'CONTA_CORRENTE',
        type_key_pix: this.plano.type_key_pix || 'CPF',
        chave_pix: this.plano.key_pix || ''
      }
      this.mostrarModalBancario = true
    },
    abrirModalCallback() {
      this.callbackUrl = this.plano.callback_url || ''
      this.mostrarModalCallback = true
    },
    async verPlano() {
      try {
        const res = await axios.get('http://localhost:8080/empresa/plano', {
          headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
        })
        console.log('Plano:', this.plano)
        this.plano = res.data
      } catch {
        alert('Erro ao buscar plano')
      }
    },
    async buscarSaldo() {
      try {
        const res = await axios.get('http://localhost:8080/empresa/saldo', {
          headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
        })
        this.saldo = res.data.saldo
      } catch {
        this.saldo = 0
      }
    },
    async buscarPagamentos() {
      try {
        const res = await axios.get('http://localhost:8080/empresa/pagamentos', {
          headers: { Authorization: `Bearer ${localStorage.getItem('token')}` },
          params: {
            page: this.pagina,
            limit: this.limite,
            status: this.filtroStatus || undefined
          }
        })
        this.pagamentos = res.data.data
        this.totalPagamentos = res.data.total
      } catch (err) {
        console.error('Erro ao buscar pagamentos:', err)
        this.pagamentos = []
      }
    },
    async pagar() {
      this.carregando = true
      try {
        await new Promise(resolve => setTimeout(resolve, 2000))
        await axios.post('http://localhost:8080/empresa/ativar', {}, {
          headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
        })
        alert(`Pagamento via ${this.formaPagamento.toUpperCase()} confirmado!`)
        this.verPlano()
      } catch {
        alert('Erro no pagamento ou ativação')
      } finally {
        this.carregando = false
      }
    },
    async copiarApiKey() {
      navigator.clipboard.writeText(this.plano.api_key)
      alert('API Key copiada!')
    },
    async atualizarCallback() {
      try {
        await axios.put('http://localhost:8080/empresa/callback', {
          callback_url: this.callbackUrl
        }, {
          headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
        })
        alert('Callback URL atualizada!')
        this.mostrarModalCallback = false
      } catch {
        alert('Erro ao atualizar callback')
      }
    },
    async atualizarBancario() {
      try {
        await axios.put('http://localhost:8080/empresa/bancario', this.banco, {
          headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
        })
        alert('Dados bancários atualizados!')
        this.mostrarModalBancario = false
      } catch {
        alert('Erro ao atualizar dados bancários')
      }
    },
    async sacar() {
      try {
        const res = await axios.post('http://localhost:8080/empresa/saque', {
          amount: this.valorSaque
        }, {
          headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
        })
        alert('Saque realizado com sucesso!')
        this.valorSaque = 0
        this.buscarSaldo()
        this.mostrarModalSaque = false
      } catch (err) {
        alert(err.response?.data?.error || 'Erro ao sacar')
      }
    },
    logout() {
      localStorage.removeItem('token')
      this.$router.push('/login')
    }
  }
}
</script>

