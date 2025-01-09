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

type ScheduleHandler struct {
	db *gorm.DB
}

func NewScheduleHandler() *ScheduleHandler {
	return &ScheduleHandler{db: database.GetDB()}
}

type ScheduleRegisterRequest struct {
	StreamID      uint      `json:"streamId"`
	ScheduledTime time.Time `json:"scheduledTime"`
}

type ScheduleResponse struct {
	ID            uint      `json:"id"`
	StreamID      uint      `json:"streamId"`
	ScheduledTime time.Time `json:"scheduledTime"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func (h *ScheduleHandler) Register(w http.ResponseWriter, r *http.Request) {
	user := middlewares.GetUserFromContext(r.Context())
	if user == nil {
		log.Println("Usuário não autenticado")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Não autorizado"})
		return
	}

	var request ScheduleRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Erro ao decodificar o corpo da requisição:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Corpo da requisição inválido"})
		return
	}

	schedule := models.Schedule{
		StreamID:      request.StreamID,
		ScheduledTime: request.ScheduledTime,
	}

	result := h.db.Create(&schedule)
	if result.Error != nil {
		log.Println("Erro ao criar agendamento:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao criar agendamento"})
		return
	}

	response := ScheduleResponse{
		ID:            schedule.ID,
		StreamID:      schedule.StreamID,
		ScheduledTime: schedule.ScheduledTime,
		CreatedAt:     schedule.CreatedAt,
		UpdatedAt:     schedule.UpdatedAt,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
