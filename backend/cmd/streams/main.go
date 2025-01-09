package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/marcoaureliojf/streamStudio/backend/internal/config"
	"github.com/marcoaureliojf/streamStudio/backend/internal/database"
	"github.com/marcoaureliojf/streamStudio/backend/internal/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	cfg := config.LoadConfig()
	database.Connect(cfg)

	r := routes.SetupStreamRoutes()

	streamServerAddr := fmt.Sprintf(":%d", cfg.StreamServerPort)

	log.Printf("Serviço de streams rodando na porta %s\n", streamServerAddr)
	log.Fatal(http.ListenAndServe(streamServerAddr, r))
}
