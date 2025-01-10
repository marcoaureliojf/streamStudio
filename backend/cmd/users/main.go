package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/marcoaureliojf/streamStudio/backend/internal/config"
	"github.com/marcoaureliojf/streamStudio/backend/internal/database"
	"github.com/marcoaureliojf/streamStudio/backend/internal/routes"
	"github.com/rs/cors"
)

func main() {
	cfg := config.LoadConfig()
	database.Connect(cfg)

	r := routes.SetupRoutes()

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Permite requisições do frontend
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	})

	serverAddr := fmt.Sprintf(":%d", cfg.ServerPort)

	log.Printf("Servidor rodando na porta %s\n", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, corsMiddleware.Handler(r)))
}
