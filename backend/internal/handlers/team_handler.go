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

type TeamHandler struct {
	db *gorm.DB
}

func NewTeamHandler() *TeamHandler {
	return &TeamHandler{db: database.GetDB()}
}

type TeamRegisterRequest struct {
	Name string `json:"name"`
}

type TeamResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateTeamRequest struct {
	Name string `json:"name"`
}

func (h *TeamHandler) Register(w http.ResponseWriter, r *http.Request) {
	user := middlewares.GetUserFromContext(r.Context())
	if user == nil {
		log.Println("Usuário não autenticado")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Não autorizado"})
		return
	}
	var request TeamRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Erro ao decodificar o corpo da requisição:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Corpo da requisição inválido"})
		return
	}

	team := models.Team{
		Name: request.Name,
	}

	result := h.db.Create(&team)
	if result.Error != nil {
		log.Println("Erro ao criar equipe:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao criar equipe"})
		return
	}

	response := TeamResponse{
		ID:        team.ID,
		Name:      team.Name,
		CreatedAt: team.CreatedAt,
		UpdatedAt: team.UpdatedAt,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *TeamHandler) GetTeams(w http.ResponseWriter, r *http.Request) {
	user := middlewares.GetUserFromContext(r.Context())
	if user == nil {
		log.Println("Usuário não autenticado")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Não autorizado"})
		return
	}

	var teams []models.Team
	result := h.db.Find(&teams)
	if result.Error != nil {
		log.Println("Erro ao buscar equipes:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao buscar equipes"})
		return
	}

	var response []TeamResponse
	for _, team := range teams {
		response = append(response, TeamResponse{
			ID:        team.ID,
			Name:      team.Name,
			CreatedAt: team.CreatedAt,
			UpdatedAt: team.UpdatedAt,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
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
		log.Println("ID da equipe inválido:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "ID da equipe inválido"})
		return
	}

	var team models.Team
	result := h.db.First(&team, id)
	if result.Error != nil {
		log.Println("Erro ao buscar equipe:", result.Error)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Equipe não encontrada"})
		return
	}

	response := TeamResponse{
		ID:        team.ID,
		Name:      team.Name,
		CreatedAt: team.CreatedAt,
		UpdatedAt: team.UpdatedAt,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *TeamHandler) UpdateTeam(w http.ResponseWriter, r *http.Request) {
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
		log.Println("ID da equipe inválido:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "ID da equipe inválido"})
		return
	}

	var request UpdateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Erro ao decodificar o corpo da requisição:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Corpo da requisição inválido"})
		return
	}

	var team models.Team
	result := h.db.First(&team, id)
	if result.Error != nil {
		log.Println("Erro ao buscar equipe:", result.Error)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Equipe não encontrada"})
		return
	}

	team.Name = request.Name
	result = h.db.Save(&team)
	if result.Error != nil {
		log.Println("Erro ao atualizar equipe:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao atualizar equipe"})
		return
	}

	response := TeamResponse{
		ID:        team.ID,
		Name:      team.Name,
		CreatedAt: team.CreatedAt,
		UpdatedAt: team.UpdatedAt,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *TeamHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
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
		log.Println("ID da equipe inválido:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "ID da equipe inválido"})
		return
	}

	var team models.Team
	result := h.db.First(&team, id)
	if result.Error != nil {
		log.Println("Erro ao buscar equipe:", result.Error)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Equipe não encontrada"})
		return
	}

	result = h.db.Delete(&team)
	if result.Error != nil {
		log.Println("Erro ao excluir equipe:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao excluir equipe"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
