package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/marcoaureliojf/streamStudio/backend/internal/config"
	"github.com/marcoaureliojf/streamStudio/backend/internal/database"
	"github.com/marcoaureliojf/streamStudio/backend/internal/database/models"
	"github.com/marcoaureliojf/streamStudio/backend/internal/handlers"
	"github.com/marcoaureliojf/streamStudio/backend/internal/routes"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func createTeam(db *gorm.DB) (uint, error) {
	team := models.Team{
		Name: fmt.Sprintf("Test Team %d", time.Now().UnixNano()),
	}
	result := db.Create(&team)
	if result.Error != nil {
		return 0, result.Error
	}
	return team.ID, nil
}

func setupTestDB() {

	cfg := config.Config{
		DBHost:     "localhost",
		DBPort:     5432,
		DBUser:     "postgres",
		DBPassword: "planetbass",
		DBName:     "streamstudio_test",
		JWTSecret:  "test-secret",
		ServerPort: 8080,
	}
	database.Connect(cfg)

	teamId, err := createTeam(database.GetDB())
	if err != nil {
		log.Fatal("Erro ao criar equipe de teste: ", err)
	}
	os.Setenv("TEST_TEAM_ID", strconv.FormatUint(uint64(teamId), 10))
}

func TestIntegrationRegisterUser(t *testing.T) {
	setupTestDB()

	teamIdStr := os.Getenv("TEST_TEAM_ID")
	teamId, err := strconv.ParseUint(teamIdStr, 10, 64)
	if err != nil {
		log.Fatal("Erro ao converter teamId", err)
	}

	router := routes.SetupRoutes()
	userRegisterRequest := handlers.UserRegisterRequest{
		Name:     "Test User",
		Email:    "integration@test.com",
		Password: "testpassword",
		TeamID:   uint(teamId),
	}

	requestBody, _ := json.Marshal(userRegisterRequest)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var responseBody handlers.UserResponse
	json.NewDecoder(rr.Body).Decode(&responseBody)
	assert.Equal(t, userRegisterRequest.Name, responseBody.Name)
	assert.Equal(t, userRegisterRequest.Email, responseBody.Email)
	assert.Equal(t, uint(teamId), responseBody.TeamID)
}

func TestIntegrationLoginUser(t *testing.T) {
	setupTestDB()
	teamIdStr := os.Getenv("TEST_TEAM_ID")
	teamId, err := strconv.ParseUint(teamIdStr, 10, 64)
	if err != nil {
		log.Fatal("Erro ao converter teamId", err)
	}
	router := routes.SetupRoutes()

	userRegisterRequest := handlers.UserRegisterRequest{
		Name:     "Test User",
		Email:    "integration2@test.com",
		Password: "testpassword",
		TeamID:   uint(teamId),
	}

	requestBodyRegister, _ := json.Marshal(userRegisterRequest)

	reqRegister, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(requestBodyRegister))
	reqRegister.Header.Set("Content-Type", "application/json")

	rrRegister := httptest.NewRecorder()
	router.ServeHTTP(rrRegister, reqRegister)

	userLoginRequest := handlers.UserLoginRequest{
		Email:    "integration2@test.com",
		Password: "testpassword",
	}
	requestBodyLogin, _ := json.Marshal(userLoginRequest)

	reqLogin, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBodyLogin))
	reqLogin.Header.Set("Content-Type", "application/json")

	rrLogin := httptest.NewRecorder()
	router.ServeHTTP(rrLogin, reqLogin)

	assert.Equal(t, http.StatusOK, rrLogin.Code)
	var responseBodyLogin struct {
		User  handlers.UserResponse `json:"user"`
		Token string                `json:"token"`
	}
	json.NewDecoder(rrLogin.Body).Decode(&responseBodyLogin)
	assert.Equal(t, userRegisterRequest.Email, responseBodyLogin.User.Email)

}

func TestIntegrationLoginUserFail(t *testing.T) {
	setupTestDB()
	router := routes.SetupRoutes()

	userLoginRequest := handlers.UserLoginRequest{
		Email:    "integrationFail@test.com",
		Password: "testpassword",
	}
	requestBodyLogin, _ := json.Marshal(userLoginRequest)

	reqLogin, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBodyLogin))
	reqLogin.Header.Set("Content-Type", "application/json")

	rrLogin := httptest.NewRecorder()
	router.ServeHTTP(rrLogin, reqLogin)

	assert.Equal(t, http.StatusUnauthorized, rrLogin.Code)
}

func TestIntegrationProtectedEndpoint(t *testing.T) {
	setupTestDB()
	teamIdStr := os.Getenv("TEST_TEAM_ID")
	teamId, err := strconv.ParseUint(teamIdStr, 10, 64)
	if err != nil {
		log.Fatal("Erro ao converter teamId", err)
	}
	router := routes.SetupRoutes()

	userRegisterRequest := handlers.UserRegisterRequest{
		Name:     "Test User",
		Email:    "integration3@test.com",
		Password: "testpassword",
		TeamID:   uint(teamId),
	}

	requestBodyRegister, _ := json.Marshal(userRegisterRequest)

	reqRegister, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(requestBodyRegister))
	reqRegister.Header.Set("Content-Type", "application/json")

	rrRegister := httptest.NewRecorder()
	router.ServeHTTP(rrRegister, reqRegister)

	userLoginRequest := handlers.UserLoginRequest{
		Email:    "integration3@test.com",
		Password: "testpassword",
	}
	requestBodyLogin, _ := json.Marshal(userLoginRequest)

	reqLogin, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBodyLogin))
	reqLogin.Header.Set("Content-Type", "application/json")

	rrLogin := httptest.NewRecorder()
	router.ServeHTTP(rrLogin, reqLogin)

	var responseBodyLogin struct {
		User  handlers.UserResponse `json:"user"`
		Token string                `json:"token"`
	}
	json.NewDecoder(rrLogin.Body).Decode(&responseBodyLogin)

	reqProtected, _ := http.NewRequest("GET", "/api/test", nil)
	reqProtected.Header.Set("Authorization", "Bearer "+responseBodyLogin.Token)
	rrProtected := httptest.NewRecorder()
	router.ServeHTTP(rrProtected, reqProtected)

	assert.Equal(t, http.StatusOK, rrProtected.Code)
	var responseBodyProtected handlers.UserResponse
	json.NewDecoder(rrProtected.Body).Decode(&responseBodyProtected)
	assert.Equal(t, userRegisterRequest.Email, responseBodyProtected.Email)
}

func TestIntegrationProtectedEndpointFail(t *testing.T) {
	setupTestDB()
	router := routes.SetupRoutes()

	reqProtected, _ := http.NewRequest("GET", "/api/test", nil)
	rrProtected := httptest.NewRecorder()
	router.ServeHTTP(rrProtected, reqProtected)
	assert.Equal(t, http.StatusUnauthorized, rrProtected.Code)
}
