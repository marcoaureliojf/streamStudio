package database

import (
	"fmt"
	"log"

	"github.com/marcoaureliojf/streamStudio/backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbInstance *gorm.DB

func Connect(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Erro ao conectar ao banco de dados: %w", err)
	}

	log.Println("Conectado ao banco de dados com sucesso!")

	// err = db.AutoMigrate(&models.User{}, &models.Team{}, &models.Permission{}, &models.Stream{}, &models.Schedule{})
	// if err != nil {
	// 	log.Fatal("Erro ao realizar auto migração:", err)
	// }

	dbInstance = db
	return db, nil
}

func GetDB() *gorm.DB {
	return dbInstance
}
