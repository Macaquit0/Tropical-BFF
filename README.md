# BFF Cognito

Este projeto é um Backend for Frontend (BFF) que utiliza o Amazon Cognito para autenticação e autorização.

## Requisitos

- Node.js v14 ou superior
- AWS CLI configurado
- Conta AWS com Amazon Cognito configurado

## Instalação

1. Clone o repositório:
    ```sh
    git clone https://github.com/seu-usuario/bff-cognito.git
    cd bff-cognito
    ```

2. Instale as dependências:
    ```sh
    npm install
    ```

## Configuração

1. Crie um arquivo `.env` na raiz do projeto e adicione as seguintes variáveis:
    ```env
    COGNITO_USER_POOL_ID=seu_user_pool_id
    COGNITO_CLIENT_ID=seu_client_id
    COGNITO_REGION=sua_regiao
    ```

## Uso

1. Inicie o servidor:
    ```sh
    npm start
    ```

2. Acesse `http://localhost:3000` no seu navegador.

## Estrutura do Projeto

- `src/` - Código fonte do projeto
- `src/routes/` - Definição das rotas da API
- `src/controllers/` - Lógica dos controladores
- `src/services/` - Serviços de integração com o Cognito

## Contribuição

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Faça push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request