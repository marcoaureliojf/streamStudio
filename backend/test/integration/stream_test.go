package integration

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/marcoaureliojf/streamStudio/backend/internal/handlers"
	"github.com/marcoaureliojf/streamStudio/backend/internal/routes"
	"github.com/marcoaureliojf/streamStudio/backend/test/integration/testhelper"
	"github.com/stretchr/testify/assert"
)

// func createTeam(db *gorm.DB) (uint, error) {
// 	team := models.Team{
// 		Name: fmt.Sprintf("Test Team %d", time.Now().UnixNano()),
// 	}
// 	result := db.Create(&team)
// 	if result.Error != nil {
// 		return 0, result.Error
// 	}
// 	return team.ID, nil
// }

// func setupTestDB() {

// 	cfg := config.Config{
// 		DBHost:           "localhost",
// 		DBPort:           5432,
// 		DBUser:           "postgres",
// 		DBPassword:       "postgres",
// 		DBName:           "streamstudio_test",
// 		JWTSecret:        "test-secret",
// 		ServerPort:       8080,
// 		StreamServerPort: 8081,
// 	}
// 	database.Connect(cfg)

// 	teamId, err := createTeam(database.GetDB())
// 	if err != nil {
// 		log.Fatal("Erro ao criar equipe de teste: ", err)
// 	}
// 	os.Setenv("TEST_TEAM_ID", strconv.FormatUint(uint64(teamId), 10))
// }

func TestIntegrationRegisterStream(t *testing.T) {
	testhelper.SetupTestDB()
	//setupTestDB()
	teamIdStr := os.Getenv("TEST_TEAM_ID")
	teamId, err := strconv.ParseUint(teamIdStr, 10, 64)
	if err != nil {
		log.Fatal("Erro ao converter teamId", err)
	}
	router := routes.SetupStreamRoutes()
	streamRegisterRequest := handlers.StreamRegisterRequest{
		Title:       "Live de Teste",
		Description: "Live teste",
		StartTime:   time.Now().Add(time.Hour),
		EndTime:     time.Now().Add(time.Hour * 2),
		TeamID:      uint(teamId),
	}

	requestBody, _ := json.Marshal(streamRegisterRequest)

	req, _ := http.NewRequest("POST", "/api/streams", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEsImV4cCI6MTcxNzU0OTg5NiwiaWF0IjoxNzE3NDYzNDk2fQ.5W7sR-y6P65a5u9n90_n6qL8wW74o1kZtE7G88y6nQk")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var responseBody handlers.StreamResponse
	json.NewDecoder(rr.Body).Decode(&responseBody)
	assert.Equal(t, streamRegisterRequest.Title, responseBody.Title)
	assert.Equal(t, streamRegisterRequest.Description, responseBody.Description)
	assert.Equal(t, streamRegisterRequest.TeamID, responseBody.TeamID)
}

func TestIntegrationRegisterSchedule(t *testing.T) {
	testhelper.SetupTestDB()
	//setupTestDB()
	teamIdStr := os.Getenv("TEST_TEAM_ID")
	teamId, err := strconv.ParseUint(teamIdStr, 10, 64)
	if err != nil {
		log.Fatal("Erro ao converter teamId", err)
	}
	router := routes.SetupStreamRoutes()

	streamRegisterRequest := handlers.StreamRegisterRequest{
		Title:       "Live de Teste",
		Description: "Live teste",
		StartTime:   time.Now().Add(time.Hour),
		EndTime:     time.Now().Add(time.Hour * 2),
		TeamID:      uint(teamId),
	}

	requestBodyStream, _ := json.Marshal(streamRegisterRequest)

	reqStream, _ := http.NewRequest("POST", "/api/streams", bytes.NewBuffer(requestBodyStream))
	reqStream.Header.Set("Content-Type", "application/json")
	reqStream.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEsImV4cCI6MTcxNzU0OTg5NiwiaWF0IjoxNzE3NDYzNDk2fQ.5W7sR-y6P65a5u9n90_n6qL8wW74o1kZtE7G88y6nQk")

	rrStream := httptest.NewRecorder()
	router.ServeHTTP(rrStream, reqStream)

	var responseBodyStream handlers.StreamResponse
	json.NewDecoder(rrStream.Body).Decode(&responseBodyStream)

	scheduleRegisterRequest := handlers.ScheduleRegisterRequest{
		StreamID:      responseBodyStream.ID,
		ScheduledTime: time.Now().Add(time.Hour * 3),
	}

	requestBodySchedule, _ := json.Marshal(scheduleRegisterRequest)

	reqSchedule, _ := http.NewRequest("POST", "/api/schedules", bytes.NewBuffer(requestBodySchedule))
	reqSchedule.Header.Set("Content-Type", "application/json")
	reqSchedule.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEsImV4cCI6MTcxNzU0OTg5NiwiaWF0IjoxNzE3NDYzNDk2fQ.5W7sR-y6P65a5u9n90_n6qL8wW74o1kZtE7G88y6nQk")
	rrSchedule := httptest.NewRecorder()
	router.ServeHTTP(rrSchedule, reqSchedule)

	assert.Equal(t, http.StatusCreated, rrSchedule.Code)
	var responseBodySchedule handlers.ScheduleResponse
	json.NewDecoder(rrSchedule.Body).Decode(&responseBodySchedule)
	assert.Equal(t, scheduleRegisterRequest.StreamID, responseBodySchedule.StreamID)
	assert.Equal(t, scheduleRegisterRequest.ScheduledTime.Format(time.RFC3339), responseBodySchedule.ScheduledTime.Format(time.RFC3339))

}
