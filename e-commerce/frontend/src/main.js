import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import './assets/tailwind.css' // Importando o CSS do Tailwind

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')
