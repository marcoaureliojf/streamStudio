package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/marcoaureliojf/streamStudio/backend/internal/config"
	"github.com/marcoaureliojf/streamStudio/backend/internal/database"
	"github.com/marcoaureliojf/streamStudio/backend/internal/queue"
	"github.com/marcoaureliojf/streamStudio/backend/internal/routes"
)

func main() {
	cfg := config.LoadConfig()
	database.Connect(cfg)

	rabbitMQ, err := queue.NewRabbitMQ(cfg)
	if err != nil {
		log.Fatal("Erro ao iniciar RabbitMQ: ", err)
	}
	defer rabbitMQ.Close()

	r := routes.SetupStreamRoutes()

	serverAddr := fmt.Sprintf(":%d", cfg.StreamServerPort)

	log.Printf("Serviço de streams rodando na porta %s\n", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, r))
}
