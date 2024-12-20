# Usar imagem base do Golang
FROM golang:1.20

# Criar diretório de trabalho
WORKDIR /app

# Copiar todos os arquivos do projeto
COPY . .

# Instalar dependências
RUN go mod tidy

# Compilar o binário da aplicação
RUN go build -o bff-cognito ./cmd/bff/main.go

# Expor a porta da aplicação
EXPOSE 8080

# Comando para rodar o servidor
CMD ["./bff-cognito"]
