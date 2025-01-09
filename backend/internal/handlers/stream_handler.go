package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/marcoaureliojf/streamStudio/backend/internal/database"
	"github.com/marcoaureliojf/streamStudio/backend/internal/database/models"
	"github.com/marcoaureliojf/streamStudio/backend/internal/middlewares"
	"gorm.io/gorm"
)

type StreamHandler struct {
	db *gorm.DB
}

func NewStreamHandler() *StreamHandler {
	return &StreamHandler{db: database.GetDB()}
}

type StreamRegisterRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	TeamID      uint      `json:"teamId"`
}
type StreamResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	TeamID      uint      `json:"teamId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (h *StreamHandler) Register(w http.ResponseWriter, r *http.Request) {
	user := middlewares.GetUserFromContext(r.Context())
	if user == nil {
		log.Println("Usuário não autenticado")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Não autorizado"})
		return
	}

	var request StreamRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Erro ao decodificar o corpo da requisição:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Corpo da requisição inválido"})
		return
	}

	stream := models.Stream{
		Title:       request.Title,
		Description: request.Description,
		StartTime:   request.StartTime,
		EndTime:     request.EndTime,
		TeamID:      request.TeamID,
	}

	result := h.db.Create(&stream)
	if result.Error != nil {
		log.Println("Erro ao criar transmissão:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao criar transmissão"})
		return
	}
	response := StreamResponse{
		ID:          stream.ID,
		Title:       stream.Title,
		Description: stream.Description,
		StartTime:   stream.StartTime,
		EndTime:     stream.EndTime,
		TeamID:      stream.TeamID,
		CreatedAt:   stream.CreatedAt,
		UpdatedAt:   stream.UpdatedAt,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
