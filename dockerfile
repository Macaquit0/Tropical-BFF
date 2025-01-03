# Usar imagem base do Golang
FROM golang:1.23

# Criar diretório de trabalho
WORKDIR /app

# Copiar todos os arquivos do projeto
COPY . .

# Instalar dependências
RUN go mod tidy

# Compilar o binário da aplicação
RUN go build -o bff ./cmd/bff/main.go

# Expor a porta da aplicação
EXPOSE 8082

# Comando para rodar o servidor
CMD ["./bff"]
