import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './assets/tailwind.css' // Importando o CSS do Tailwind

createApp(App).use(router).mount('#app')
