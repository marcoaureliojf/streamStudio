package database

import (
	"fmt"
	"log"

	"github.com/marcoaureliojf/streamStudio/backend/internal/config"
	"github.com/marcoaureliojf/streamStudio/backend/internal/database/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg config.Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}

	log.Println("Conectado ao banco de dados com sucesso!")

	err = db.AutoMigrate(&models.User{}, &models.Team{}, &models.Permission{})
	if err != nil {
		log.Fatal("Erro ao realizar auto migração:", err)
	}

	DB = db
}

func GetDB() *gorm.DB {
	return DB
}
