package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/marcoaureliojf/streamStudio/backend/internal/auth"
	"github.com/marcoaureliojf/streamStudio/backend/internal/config"
	"github.com/marcoaureliojf/streamStudio/backend/internal/database"
	"github.com/marcoaureliojf/streamStudio/backend/internal/database/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler() *UserHandler {
	return &UserHandler{db: database.GetDB()}
}

type UserRegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	TeamID   uint   `json:"teamId"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	TeamID    uint      `json:"teamId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var request UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Erro ao decodificar o corpo da requisição:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Corpo da requisição inválido"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Erro ao gerar hash da senha:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao processar a requisição"})
		return
	}
	user := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hashedPassword),
		TeamID:   request.TeamID,
	}

	result := h.db.Create(&user)
	if result.Error != nil {
		log.Println("Erro ao criar usuário:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao criar usuário"})
		return
	}

	response := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		TeamID:    user.TeamID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Erro ao decodificar o corpo da requisição:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Corpo da requisição inválido"})
		return
	}

	var user models.User
	result := h.db.Where("email = ?", request.Email).First(&user)
	if result.Error != nil {
		log.Println("Erro ao buscar usuário:", result.Error)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Credenciais inválidas"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		log.Println("Erro ao comparar senhas:", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Credenciais inválidas"})
		return
	}
	cfg := config.LoadConfig()
	token, err := auth.GenerateToken(user, cfg)
	if err != nil {
		log.Println("Erro ao gerar token:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao gerar token"})
		return
	}

	response := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		TeamID:    user.TeamID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	type LoginResponse struct {
		User  UserResponse `json:"user"`
		Token string       `json:"token"`
	}

	loginResponse := LoginResponse{
		User:  response,
		Token: token,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginResponse)
}
