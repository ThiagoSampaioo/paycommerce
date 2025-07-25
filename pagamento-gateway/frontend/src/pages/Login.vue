<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-r from-blue-100 to-indigo-100 px-4">
    <div class="bg-white shadow-2xl rounded-xl p-8 w-full max-w-md animate-fade-in">
      <h2 class="text-3xl font-bold text-center text-gray-800 mb-6">PAYCOMMERCE</h2>

      <form @submit.prevent="login" class="space-y-5">
        <div class="relative">
          <EnvelopeIcon class="w-5 h-5 absolute top-3 left-3 text-gray-400" />
          <input
            v-model="email"
            type="email"
            placeholder="Email"
            required
            class="pl-10 w-full py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:outline-none"
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
            <circle
              class="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              stroke-width="4"
            ></circle>
            <path
              class="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8v8z"
            ></path>
          </svg>
          <span>{{ carregando ? 'Entrando...' : 'Entrar' }}</span>
        </button>
      </form>

      <transition name="fade">
        <p
          v-if="mensagem"
          class="mt-4 text-center text-sm text-red-600 font-medium"
        >
          {{ mensagem }}
        </p>
      </transition>

      <p class="mt-6 text-sm text-center text-gray-600">
        Não tem uma conta?
        <router-link
          to="/registro"
          class="text-blue-600 hover:underline font-semibold"
        >
          Registrar
        </router-link>
      </p>
    </div>
  </div>
</template>

<script>
import { EnvelopeIcon, LockClosedIcon, EyeIcon, EyeSlashIcon } from '@heroicons/vue/24/solid'
import axios from 'axios'

export default {
  components: {
    EnvelopeIcon,
    LockClosedIcon,
    EyeIcon,
    EyeSlashIcon
  },
  data() {
    return {
      email: '',
      password: '',
      mensagem: '',
      carregando: false,
      verSenha: false
    }
  },
  methods: {
    async login() {
      this.carregando = true
      this.mensagem = ''
      try {
        const res = await axios.post('http://localhost:8080/login', {
          email: this.email,
          password: this.password
        })
        localStorage.setItem('token', res.data.token)
        this.$router.push('/empresa')
      } catch (err) {
        this.mensagem = err.response?.data?.error || 'Erro ao fazer login'
      } finally {
        this.carregando = false
      }
    }
  }
}
</script>

<style scoped>
/* animação suave */
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
