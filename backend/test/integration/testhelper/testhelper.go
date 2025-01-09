package testhelper

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/marcoaureliojf/streamStudio/backend/internal/config"
	"github.com/marcoaureliojf/streamStudio/backend/internal/database"
	"github.com/marcoaureliojf/streamStudio/backend/internal/database/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func Connect(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func CreateTeam(db *gorm.DB) (uint, error) {
	team := models.Team{
		Name: fmt.Sprintf("Test Team %d", time.Now().UnixNano()),
	}
	result := db.Create(&team)
	return team.ID, result.Error
}

func SetupTestDB() *gorm.DB {
	cfg := config.Config{
		DBHost:     "localhost",
		DBPort:     5432,
		DBUser:     "postgres",
		DBPassword: "planetbass",
		DBName:     "streamstudio_test",
		JWTSecret:  "test-secret",
		ServerPort: 8080,
	}

	var err error
	testDB, err = database.Connect(cfg)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados de teste: %v", err)
	}

	// Clean existing data
	testDB.Exec("DELETE FROM users")
	testDB.Exec("DELETE FROM teams")

	teamId, err := CreateTeam(testDB)
	if err != nil {
		log.Fatal("Error creating test team: ", err)
	}
	os.Setenv("TEST_TEAM_ID", strconv.FormatUint(uint64(teamId), 10))

	return testDB
}

func GetTestDB() *gorm.DB {
	if testDB == nil {
		return SetupTestDB()
	}
	return testDB
}

func CleanupTestDB() {
	if testDB != nil {
		sqlDB, err := testDB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}
