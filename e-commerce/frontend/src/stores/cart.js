import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'

export const useCartStore = defineStore('cart', () => {
  const items = ref([])

  // Carregar carrinho do localStorage (se existir)
  const savedCart = localStorage.getItem('cart')
  if (savedCart) {
    items.value = JSON.parse(savedCart)
  }

  // Persistir alterações do carrinho no localStorage
  watch(items, (newItems) => {
    localStorage.setItem('cart', JSON.stringify(newItems))
  }, { deep: true })

  const adicionarItem = (produto, quantidade = 1) => {
    const existente = items.value.find(item => item.id === produto.id)
    if (existente) {
      existente.quantity += quantidade
    } else {
      items.value.push({ ...produto, quantity: quantidade })
    }
  }

  const removerItem = (produtoId) => {
    items.value = items.value.filter(item => item.id !== produtoId)
  }

  const limparCarrinho = () => {
    items.value = []
  }

  const totalItens = computed(() =>
    items.value.reduce((total, item) => total + item.quantity, 0)
  )

  const totalValor = computed(() =>
    items.value.reduce((total, item) => total + item.quantity * item.price, 0)
  )

  const totalGeral = computed(() =>
    items.value.reduce((acc, item) => acc + item.price * item.quantity, 0)
  )

  return {
    items,
    adicionarItem,
    removerItem,
    limparCarrinho,
    totalItens,
    totalValor,
    totalGeral
  }
})
