package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/marcoaureliojf/streamStudio/backend/internal/config"
	"github.com/marcoaureliojf/streamStudio/backend/internal/database"
	"github.com/marcoaureliojf/streamStudio/backend/internal/queue"
	"github.com/marcoaureliojf/streamStudio/backend/internal/routes"
	"github.com/rs/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	cfg := config.LoadConfig()
	database.Connect(cfg)

	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	queue.Init(cfg)

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	})

	// Configuração de rotas
	rs := routes.SetupStreamRoutes()
	r := routes.SetupRoutes()

	serverAddr := fmt.Sprintf(":%d", cfg.ServerPort)
	streamServerAddr := fmt.Sprintf(":%d", cfg.StreamServerPort)

	// Inicia o servidor de streaming em uma goroutine
	go func() {
		log.Printf("Serviço de streaming rodando na porta %s\n", streamServerAddr)
		if err := http.ListenAndServe(streamServerAddr, corsMiddleware.Handler(rs)); err != nil {
			log.Fatalf("Erro no serviço de streaming: %v", err)
		}
	}()

	// Inicia o servidor principal
	log.Printf("Servidor rodando na porta %s\n", serverAddr)
	if err := http.ListenAndServe(serverAddr, corsMiddleware.Handler(r)); err != nil {
		log.Fatalf("Erro no servidor principal: %v", err)
	}
}
