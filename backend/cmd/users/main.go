package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/marcoaureliojf/streamStudio/backend/internal/config"
	"github.com/marcoaureliojf/streamStudio/backend/internal/database"
	"github.com/marcoaureliojf/streamStudio/backend/internal/routes"
)

func main() {
	cfg := config.LoadConfig()
	database.Connect(cfg)

	r := routes.SetupRoutes()

	serverAddr := fmt.Sprintf(":%d", cfg.ServerPort)

	log.Printf("Servidor rodando na porta %s\n", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, r))
}
