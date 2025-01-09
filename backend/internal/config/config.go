package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost           string
	DBPort           int
	DBUser           string
	DBPassword       string
	DBName           string
	JWTSecret        string
	ServerPort       int
	StreamServerPort int
	RabbitMQHost     string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env:", err)
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal("Erro ao converter a porta do banco de dados:", err)
	}

	serverPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatal("Erro ao converter a porta do servidor:", err)
	}
	streamServerPort, err := strconv.Atoi(os.Getenv("STREAM_SERVER_PORT"))
	if err != nil {
		log.Fatal("Erro ao converter a porta do servidor:", err)
	}

	return Config{
		DBHost:           os.Getenv("DB_HOST"),
		DBPort:           dbPort,
		DBUser:           os.Getenv("DB_USER"),
		DBPassword:       os.Getenv("DB_PASSWORD"),
		DBName:           os.Getenv("DB_NAME"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		ServerPort:       serverPort,
		StreamServerPort: streamServerPort,
		RabbitMQHost:     os.Getenv("DB_HOST"),
	}
}
