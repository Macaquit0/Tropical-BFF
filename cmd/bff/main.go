package main

import (
	"log"
	"net/http"
	"os"

	"bff-cognito/internal/auth"
	"bff-cognito/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Carregar as variáveis de ambiente
	port := os.Getenv("APP_PORT")
	region := os.Getenv("COGNITO_REGION")
	clientID := os.Getenv("COGNITO_APP_CLIENT_ID")
	userPoolID := os.Getenv("COGNITO_USER_POOL_ID")
	jwtSecret := os.Getenv("JWT_SECRET")

	// Verificar se as variáveis estão definidas
	if port == "" || region == "" || clientID == "" || userPoolID == "" || jwtSecret == "" {
		log.Fatal("Configurações de ambiente incompletas")
	}

	// Inicializar o cliente Cognito
	cognitoClient := auth.NewCognitoClient(region, clientID)

	// Configurar o router
	router := chi.NewRouter()
	handlers.RegisterRoutes(router, cognitoClient)

	// Iniciar o servidor
	log.Printf("Servidor rodando na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
