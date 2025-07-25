# 🛒 PayCommerce

PayCommerce é um projeto completo de e-commerce integrado a uma plataforma de pagamentos, desenvolvido com foco educacional e profissional. Simula o fluxo real de compras online com processamento externo de pagamento, redirecionamentos, confirmação por callback e gerenciamento de provedores.

## 🧱 Tecnologias Utilizadas

### E-commerce
- **Frontend**: Vue.js + Tailwind CSS (porta 5174)
- **Backend**: Golang (porta 8084)
- **Banco de Dados**: MySQL 8 (porta 3308)
- **Gerenciador de Banco**: Adminer (porta 8088)

### Plataforma de Pagamento
- **Frontend**: Vue.js + Tailwind CSS (porta 5173)
- **Backend**: Golang (porta 8080)
- **Banco de Dados**: MySQL 8 (porta 3307)
- **Gerenciador de Banco**: Adminer (porta 8081)

### Outros
- Docker + Docker Compose

---

## ⚙️ Estrutura de Pastas

paycommerce/
│
├── e-commerce/
│ ├── backend/
│ └── frontend/
│
└── pagamento-gateway/
├── backend/
└── frontend/

yaml


---

## 🚀 Como iniciar o projeto

### 1. Clone o repositório

```bash
git clone https://github.com/seu-usuario/paycommerce.git
cd paycommerce
2. Suba os serviços com Docker Compose
Para o e-commerce:


cd e-commerce
docker compose up -d

Para a plataforma de pagamento:


cd ../pagamento-gateway
docker compose up -d

Nota: Certifique-se de que as portas 3307, 3308, 8080, 8081, 8084, 8088, 5173 e 5174 estejam livres.

3. Inicie os backends e frontends manualmente
Cada backend e frontend deve ser iniciado individualmente. Exemplo:



# Em cada pasta backend (e-commerce e pagamento)
go run main.go

# Em cada pasta frontend
npm install
npm run dev

👤 Usuário Administrador
Ao iniciar o backend do e-commerce pela primeira vez, um usuário sistêmico será criado automaticamente:

Nome: Usuário Sistêmico

Email: admin@shop.com

Senha: 123456

Role: admin

🔁 Integração entre sistemas
Registro do provedor de pagamento no e-commerce
Acesse o painel admin do e-commerce e cadastre um provedor com os seguintes dados:

Nome: paycommerce

API Key: fornecida pela plataforma de pagamento

URL de Pagamento: http://localhost:5173/confirmacao

URL de Cancelamento: http://localhost:8080/api/callback/cancelamento

Configuração de callback na plataforma de pagamento
Cadastre no sistema da plataforma a seguinte URL de callback, para que o status do pagamento seja retornado ao e-commerce:

Callback de confirmação:
http://localhost:8084/callback-payment

📦 Fluxo de Compra
Cliente seleciona produtos e inicia checkout no e-commerce.

E-commerce registra a venda e redireciona o cliente para a plataforma de pagamento.

Cliente confirma o pagamento na tela da plataforma.

Plataforma chama a URL de callback do e-commerce com status da transação.

E-commerce atualiza o status da compra.

Caso necessário, o e-commerce pode cancelar uma venda, enviando requisição para:

http://localhost:8080/api/callback/cancelamento

📝 Considerações Finais
Este projeto foi desenvolvido com o objetivo de simular um ambiente profissional de integração entre um sistema de vendas e um gateway de pagamento, com foco em boas práticas, organização de código e arquitetura realista.

Contribuições e sugestões são bem-vindas!
