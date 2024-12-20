# Helpers
container_exists = $(shell docker ps -q -f name=$(1) | grep -q . && echo "1" || echo "0")

# Configurações
APP_NAME = bff-cognito
DOCKER_COMPOSE_FILE = docker-compose.yaml

# Comandos de setup
setup: tidy build setup-db
	@echo "Setup completo!"

tidy:
	@echo "Rodando go mod tidy..."
	go mod tidy

build:
	@echo "Construindo aplicação..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) build

setup-db:
	@echo "Configurando banco de dados..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) exec db sh -c "psql -U postgres -d bff -c 'CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, email VARCHAR(255), password VARCHAR(255));'"

# Comandos para rodar e parar containers
up:
	@echo "Subindo os serviços..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

down:
	@echo "Derrubando os serviços..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

restart: down up
	@echo "Reiniciando os serviços..."

logs:
	@echo "Mostrando logs..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

# Rodar aplicação local
run:
	@echo "Rodando aplicação localmente..."
	go run cmd/bff/main.go

# Criar novas migrações
create-migration:
	@echo "Criando nova migração: $(name)"
	go run cmd/migrator/main.go create $(name)

# Executar migrações
migrate:
	@echo "Executando migrações..."
	go run cmd/migrator/main.go up
# Help
help:
	@echo "Comandos disponíveis:"
	@echo "  make tidy              - Executa go mod tidy"
	@echo "  make build             - Constrói a aplicação com Docker"
	@echo "  make setup             - Configura o ambiente completo (tidy, build, setup-db)"
	@echo "  make up                - Sobe os serviços com Docker"
	@echo "  make down              - Derruba os serviços com Docker"
	@echo "  make restart           - Reinicia os serviços"
	@echo "  make logs              - Mostra os logs dos serviços"
	@echo "  make run               - Roda a aplicação localmente sem Docker"
	@echo "  make create-migration  - Cria uma nova migração (use com 'name=<nome>')"
	@echo "  make migrate           - Executa as migrações no banco de dados"
	@echo "  make setup-db          - Configura o banco de dados"
