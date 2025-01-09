package unit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/marcoaureliojf/streamStudio/backend/internal/config"
	"github.com/marcoaureliojf/streamStudio/backend/internal/database"
	"github.com/marcoaureliojf/streamStudio/backend/internal/handlers"
	"github.com/marcoaureliojf/streamStudio/backend/internal/routes"
	"github.com/stretchr/testify/assert"
)

func setupTestDB() {
	cfg := config.LoadConfig()
	database.Connect(cfg)

}

func TestRegisterUser(t *testing.T) {
	setupTestDB()
	// userHandler := handlers.NewUserHandler()
	router := routes.SetupRoutes()

	userRegisterRequest := handlers.UserRegisterRequest{
		Name:     "Test User",
		Email:    "test@test.com",
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

func TestLoginUser(t *testing.T) {
	setupTestDB()
	//userHandler := handlers.NewUserHandler()
	router := routes.SetupRoutes()

	userRegisterRequest := handlers.UserRegisterRequest{
		Name:     "Test User",
		Email:    "test2@test.com",
		Password: "testpassword",
		TeamID:   1,
	}

	requestBodyRegister, _ := json.Marshal(userRegisterRequest)

	reqRegister, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(requestBodyRegister))
	reqRegister.Header.Set("Content-Type", "application/json")

	rrRegister := httptest.NewRecorder()
	router.ServeHTTP(rrRegister, reqRegister)

	userLoginRequest := handlers.UserLoginRequest{
		Email:    "test2@test.com",
		Password: "testpassword",
	}
	requestBodyLogin, _ := json.Marshal(userLoginRequest)

	reqLogin, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBodyLogin))
	reqLogin.Header.Set("Content-Type", "application/json")

	rrLogin := httptest.NewRecorder()
	router.ServeHTTP(rrLogin, reqLogin)

	assert.Equal(t, http.StatusOK, rrLogin.Code)

	var responseBodyLogin handlers.UserResponse
	json.NewDecoder(rrLogin.Body).Decode(&responseBodyLogin)
	assert.Equal(t, userRegisterRequest.Email, responseBodyLogin.Email)

}

func TestLoginUserFail(t *testing.T) {
	setupTestDB()
	//userHandler := handlers.NewUserHandler()
	router := routes.SetupRoutes()

	userLoginRequest := handlers.UserLoginRequest{
		Email:    "fail@test.com",
		Password: "testpassword",
	}
	requestBodyLogin, _ := json.Marshal(userLoginRequest)

	reqLogin, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBodyLogin))
	reqLogin.Header.Set("Content-Type", "application/json")

	rrLogin := httptest.NewRecorder()
	router.ServeHTTP(rrLogin, reqLogin)

	assert.Equal(t, http.StatusUnauthorized, rrLogin.Code)
}
