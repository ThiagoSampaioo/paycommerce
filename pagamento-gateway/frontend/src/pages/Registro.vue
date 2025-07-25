<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-r from-indigo-100 to-blue-100 px-4">
    <div class="bg-white shadow-2xl rounded-xl p-8 w-full max-w-md animate-fade-in">
      <h2 class="text-3xl font-bold text-center mb-6 text-gray-800">PAYCOMMERCE</h2>

      <form @submit.prevent="registrar" class="space-y-5">
        <div class="relative">
          <BuildingOfficeIcon class="w-5 h-5 absolute top-3 left-3 text-gray-400" />
          <input
            v-model="name"
            type="text"
            placeholder="Nome da empresa"
            required
            class="pl-10 w-full py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div class="relative">
          <EnvelopeIcon class="w-5 h-5 absolute top-3 left-3 text-gray-400" />
          <input
            v-model="email"
            type="email"
            placeholder="Email corporativo"
            required
            class="pl-10 w-full py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div class="relative">
            <LockClosedIcon class="w-5 h-5 absolute top-3 left-3 text-gray-400" />
            <input
                :type="verSenha ? 'text' : 'password'"
                v-model="password"
                placeholder="Crie uma senha"
                required
                class="pl-10 pr-10 w-full py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <button
                type="button"
                class="absolute right-3 top-2.5 text-gray-400 hover:text-gray-600"
                @click="verSenha = !verSenha"
            >
                <EyeIcon v-if="!verSenha" class="w-5 h-5" />
                <EyeSlashIcon v-else class="w-5 h-5" />
            </button>
        </div>

        <button
          type="submit"
          :disabled="carregando"
          class="w-full flex items-center justify-center gap-2 bg-blue-600 text-white font-semibold py-2 rounded-lg hover:bg-blue-700 transition disabled:opacity-50"
        >
          <svg
            v-if="carregando"
            class="animate-spin h-5 w-5 text-white"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z" />
          </svg>
          <span>{{ carregando ? 'Registrando...' : 'Registrar' }}</span>
        </button>
      </form>

      <transition name="fade">
        <p
          v-if="mensagem"
          class="mt-4 text-center text-sm font-medium"
          :class="mensagem.includes('sucesso') ? 'text-green-600' : 'text-red-600'"
        >
          {{ mensagem }}
        </p>
      </transition>

      <p class="mt-6 text-sm text-center text-gray-600">
        Já tem uma conta?
        <router-link
          to="/login"
          class="text-blue-600 hover:underline font-semibold"
        >
          Entrar
        </router-link>
      </p>
    </div>
  </div>
</template>

<script>
import {
  EnvelopeIcon,
  LockClosedIcon,
  BuildingOfficeIcon,
  EyeIcon,
  EyeSlashIcon
} from '@heroicons/vue/24/solid'
import axios from 'axios'

export default {
  components: {
    EnvelopeIcon,
    LockClosedIcon,
    BuildingOfficeIcon, 
    EyeIcon,
    EyeSlashIcon
  },
  data() {
    return {
      name: '',
      email: '',
      password: '',
      mensagem: '',
      carregando: false,
        verSenha: false
    }
  },
  methods: {
    async registrar() {
      this.carregando = true
      this.mensagem = ''
      try {
        const res = await axios.post('http://localhost:8080/registrar', {
          name: this.name,
          email: this.email,
          password: this.password
        })
        this.mensagem = res.data.message || 'Registrado com sucesso!'
      } catch (err) {
        this.mensagem = err.response?.data?.error || 'Erro ao registrar'
      } finally {
        this.carregando = false
      }
    }
  }
}
</script>

<style scoped>
.animate-fade-in {
  animation: fadeIn 0.5s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
