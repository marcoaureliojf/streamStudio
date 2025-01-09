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

type PermissionHandler struct {
	db *gorm.DB
}

func NewPermissionHandler() *PermissionHandler {
	return &PermissionHandler{db: database.GetDB()}
}

type PermissionRegisterRequest struct {
	Name string `json:"name"`
}

type PermissionResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdatePermissionRequest struct {
	Name string `json:"name"`
}

func (h *PermissionHandler) Register(w http.ResponseWriter, r *http.Request) {
	user := middlewares.GetUserFromContext(r.Context())
	if user == nil {
		log.Println("Usuário não autenticado")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Não autorizado"})
		return
	}
	var request PermissionRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Erro ao decodificar o corpo da requisição:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Corpo da requisição inválido"})
		return
	}

	permission := models.Permission{
		Name: request.Name,
	}

	result := h.db.Create(&permission)
	if result.Error != nil {
		log.Println("Erro ao criar permissão:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao criar permissão"})
		return
	}

	response := PermissionResponse{
		ID:        permission.ID,
		Name:      permission.Name,
		CreatedAt: permission.CreatedAt,
		UpdatedAt: permission.UpdatedAt,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *PermissionHandler) GetPermissions(w http.ResponseWriter, r *http.Request) {
	user := middlewares.GetUserFromContext(r.Context())
	if user == nil {
		log.Println("Usuário não autenticado")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Não autorizado"})
		return
	}

	var permissions []models.Permission
	result := h.db.Find(&permissions)
	if result.Error != nil {
		log.Println("Erro ao buscar permissões:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao buscar permissões"})
		return
	}

	var response []PermissionResponse
	for _, permission := range permissions {
		response = append(response, PermissionResponse{
			ID:        permission.ID,
			Name:      permission.Name,
			CreatedAt: permission.CreatedAt,
			UpdatedAt: permission.UpdatedAt,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *PermissionHandler) GetPermission(w http.ResponseWriter, r *http.Request) {
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
		log.Println("ID da permissão inválido:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "ID da permissão inválido"})
		return
	}

	var permission models.Permission
	result := h.db.First(&permission, id)
	if result.Error != nil {
		log.Println("Erro ao buscar permissão:", result.Error)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Permissão não encontrada"})
		return
	}

	response := PermissionResponse{
		ID:        permission.ID,
		Name:      permission.Name,
		CreatedAt: permission.CreatedAt,
		UpdatedAt: permission.UpdatedAt,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *PermissionHandler) UpdatePermission(w http.ResponseWriter, r *http.Request) {
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
		log.Println("ID da permissão inválido:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "ID da permissão inválido"})
		return
	}

	var request UpdatePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Erro ao decodificar o corpo da requisição:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Corpo da requisição inválido"})
		return
	}

	var permission models.Permission
	result := h.db.First(&permission, id)
	if result.Error != nil {
		log.Println("Erro ao buscar permissão:", result.Error)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Permissão não encontrada"})
		return
	}

	permission.Name = request.Name
	result = h.db.Save(&permission)
	if result.Error != nil {
		log.Println("Erro ao atualizar permissão:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao atualizar permissão"})
		return
	}

	response := PermissionResponse{
		ID:        permission.ID,
		Name:      permission.Name,
		CreatedAt: permission.CreatedAt,
		UpdatedAt: permission.UpdatedAt,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *PermissionHandler) DeletePermission(w http.ResponseWriter, r *http.Request) {
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
		log.Println("ID da permissão inválido:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "ID da permissão inválido"})
		return
	}

	var permission models.Permission
	result := h.db.First(&permission, id)
	if result.Error != nil {
		log.Println("Erro ao buscar permissão:", result.Error)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Permissão não encontrada"})
		return
	}

	result = h.db.Delete(&permission)
	if result.Error != nil {
		log.Println("Erro ao excluir permissão:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao excluir permissão"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
