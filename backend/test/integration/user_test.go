package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/marcoaureliojf/streamStudio/backend/internal/config"
	"github.com/marcoaureliojf/streamStudio/backend/internal/database"
	"github.com/marcoaureliojf/streamStudio/backend/internal/handlers"
	"github.com/marcoaureliojf/streamStudio/backend/internal/routes"
	"github.com/stretchr/testify/assert"
)

func setupTestDB() {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "planetbass")
	os.Setenv("DB_NAME", "streamstudio")
	os.Setenv("JWT_SECRET", "asdfasdfdfwerqwerytfghfgh")
	os.Setenv("SERVER_PORT", "8181")
	cfg := config.LoadConfig()
	database.Connect(cfg)

}

func TestIntegrationRegisterUser(t *testing.T) {
	setupTestDB()
	router := routes.SetupRoutes()
	userRegisterRequest := handlers.UserRegisterRequest{
		Name:     "Test User",
		Email:    "integration@test.com",
		Password: "testpassword",
		TeamID:   1,
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
	assert.Equal(t, userRegisterRequest.TeamID, responseBody.TeamID)
}

func TestIntegrationLoginUser(t *testing.T) {
	setupTestDB()
	router := routes.SetupRoutes()

	userRegisterRequest := handlers.UserRegisterRequest{
		Name:     "Test User",
		Email:    "integration2@test.com",
		Password: "testpassword",
		TeamID:   1,
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
	router := routes.SetupRoutes()

	userRegisterRequest := handlers.UserRegisterRequest{
		Name:     "Test User",
		Email:    "integration3@test.com",
		Password: "testpassword",
		TeamID:   1,
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
