package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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
type UpdateStreamRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	TeamID      uint      `json:"teamId"`
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

func (h *StreamHandler) GetStreams(w http.ResponseWriter, r *http.Request) {
	user := middlewares.GetUserFromContext(r.Context())
	if user == nil {
		log.Println("Usuário não autenticado")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Não autorizado"})
		return
	}

	var streams []models.Stream
	result := h.db.Find(&streams)
	if result.Error != nil {
		log.Println("Erro ao buscar transmissões:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao buscar transmissões"})
		return
	}

	var response []StreamResponse
	for _, stream := range streams {
		response = append(response, StreamResponse{
			ID:          stream.ID,
			Title:       stream.Title,
			Description: stream.Description,
			StartTime:   stream.StartTime,
			EndTime:     stream.EndTime,
			TeamID:      stream.TeamID,
			CreatedAt:   stream.CreatedAt,
			UpdatedAt:   stream.UpdatedAt,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *StreamHandler) GetStream(w http.ResponseWriter, r *http.Request) {
	user := middlewares.GetUserFromContext(r.Context())
	if user == nil {
		log.Println("Usuário não autenticado")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Não autorizado"})
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		log.Println("ID da transmissão inválido:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "ID da transmissão inválido"})
		return
	}

	var stream models.Stream
	result := h.db.First(&stream, id)
	if result.Error != nil {
		log.Println("Erro ao buscar transmissão:", result.Error)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Transmissão não encontrada"})
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *StreamHandler) UpdateStream(w http.ResponseWriter, r *http.Request) {
	user := middlewares.GetUserFromContext(r.Context())
	if user == nil {
		log.Println("Usuário não autenticado")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Não autorizado"})
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		log.Println("ID da transmissão inválido:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "ID da transmissão inválido"})
		return
	}

	var request UpdateStreamRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Erro ao decodificar o corpo da requisição:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Corpo da requisição inválido"})
		return
	}

	var stream models.Stream
	result := h.db.First(&stream, id)
	if result.Error != nil {
		log.Println("Erro ao buscar transmissão:", result.Error)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Transmissão não encontrada"})
		return
	}

	stream.Title = request.Title
	stream.Description = request.Description
	stream.StartTime = request.StartTime
	stream.EndTime = request.EndTime
	stream.TeamID = request.TeamID

	result = h.db.Save(&stream)
	if result.Error != nil {
		log.Println("Erro ao atualizar transmissão:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao atualizar transmissão"})
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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *StreamHandler) DeleteStream(w http.ResponseWriter, r *http.Request) {
	user := middlewares.GetUserFromContext(r.Context())
	if user == nil {
		log.Println("Usuário não autenticado")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Não autorizado"})
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		log.Println("ID da transmissão inválido:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "ID da transmissão inválido"})
		return
	}

	var stream models.Stream
	result := h.db.First(&stream, id)
	if result.Error != nil {
		log.Println("Erro ao buscar transmissão:", result.Error)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Transmissão não encontrada"})
		return
	}

	result = h.db.Delete(&stream)
	if result.Error != nil {
		log.Println("Erro ao excluir transmissão:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao excluir transmissão"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
