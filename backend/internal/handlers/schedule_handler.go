// package handlers

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"strconv"
// 	"time"

// 	"github.com/gorilla/mux"
// 	"github.com/marcoaureliojf/streamStudio/backend/internal/config"
// 	"github.com/marcoaureliojf/streamStudio/backend/internal/database"
// 	"github.com/marcoaureliojf/streamStudio/backend/internal/database/models"
// 	"github.com/marcoaureliojf/streamStudio/backend/internal/middlewares"
// 	"github.com/marcoaureliojf/streamStudio/backend/internal/queue"
// 	"gorm.io/gorm"
// )

// type ScheduleHandler struct {
// 	db       *gorm.DB
// 	rabbitMQ *queue.RabbitMQ
// 	cfg      config.Config
// }

// func NewScheduleHandler() *ScheduleHandler {
// 	cfg := config.LoadConfig()
// 	rabbitMQ, err := queue.NewRabbitMQ(cfg)
// 	if err != nil {
// 		log.Fatal("Erro ao iniciar RabbitMQ: ", err)
// 	}
// 	return &ScheduleHandler{db: database.GetDB(), rabbitMQ: rabbitMQ, cfg: cfg}
// }

// type ScheduleRegisterRequest struct {
// 	StreamID      uint      `json:"streamId"`
// 	ScheduledTime time.Time `json:"scheduledTime"`
// }

// type ScheduleResponse struct {
// 	ID            uint      `json:"id"`
// 	StreamID      uint      `json:"streamId"`
// 	ScheduledTime time.Time `json:"scheduledTime"`
// 	CreatedAt     time.Time `json:"createdAt"`
// 	UpdatedAt     time.Time `json:"updatedAt"`
// }

// type UpdateScheduleRequest struct {
// 	StreamID      uint      `json:"streamId"`
// 	ScheduledTime time.Time `json:"scheduledTime"`
// }

// func (h *ScheduleHandler) Register(w http.ResponseWriter, r *http.Request) {
// 	user := middlewares.GetUserFromContext(r.Context())
// 	if user == nil {
// 		log.Println("Usuário não autenticado")
// 		w.WriteHeader(http.StatusUnauthorized)
// 		json.NewEncoder(w).Encode(ErrorResponse{Message: "Não autorizado"})
// 		return
// 	}

// 	var request ScheduleRegisterRequest
// 	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
// 		log.Println("Erro ao decodificar o corpo da requisição:", err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(ErrorResponse{Message: "Corpo da requisição inválido"})
// 		return
// 	}

// 	schedule := models.Schedule{
// 		StreamID:      request.StreamID,
// 		ScheduledTime: request.ScheduledTime,
// 	}

// 	result := h.db.Create(&schedule)
// 	if result.Error != nil {
// 		log.Println("Erro ao criar agendamento:", result.Error)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao criar agendamento"})
// 		return
// 	}

// 	message, err := json.Marshal(schedule)
// 	if err != nil {
// 		log.Println("Erro ao serializar mensagem", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao processar agendamento"})
// 		return
// 	}

// 	err = h.rabbitMQ.Publish(message)
// 	if err != nil {
// 		log.Println("Erro ao publicar mensagem na fila", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao processar agendamento"})
// 		return
// 	}

// 	response := ScheduleResponse{
// 		ID:            schedule.ID,
// 		StreamID:      schedule.StreamID,
// 		ScheduledTime: schedule.ScheduledTime,
// 		CreatedAt:     schedule.CreatedAt,
// 		UpdatedAt:     schedule.UpdatedAt,
// 	}
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(response)
// }

// func (h *ScheduleHandler) GetSchedules(w http.ResponseWriter, r *http.Request) {
// 	user := middlewares.GetUserFromContext(r.Context())
// 	if user == nil {
// 		log.Println("Usuário não autenticado")
// 		w.WriteHeader(http.StatusUnauthorized)
// 		json.NewEncoder(w).Encode(ErrorResponse{Message: "Não autorizado"})
// 		return
// 	}

// 	var schedules []models.Schedule
// 	result := h.db.Find(&schedules)
// 	if result.Error != nil {
// 		log.Println("Erro ao buscar agendamentos:", result.Error)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao buscar agendamentos"})
// 		return
// 	}

// 	var response []ScheduleResponse
// 	for _, schedule := range schedules {
// 		response = append(response, ScheduleResponse{
// 			ID:            schedule.ID,
// 			StreamID:      schedule.StreamID,
// 			ScheduledTime: schedule.ScheduledTime,
// 			CreatedAt:     schedule.CreatedAt,
// 			UpdatedAt:     schedule.UpdatedAt,
// 		})
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(response)
// }

// func (h *ScheduleHandler) GetSchedule(w http.ResponseWriter, r *http.Request) {
// 	user := middlewares.GetUserFromContext(r.Context())
// 	if user == nil {
// 		log.Println("Usuário não autenticado")
// 		w.WriteHeader(http.StatusUnauthorized)
// 		json.NewEncoder(w).Encode(ErrorResponse{Message: "Não autorizado"})
// 		return
// 	}

// 	vars := mux.Vars(r)
// 	id, err := strconv.ParseUint(vars["id"], 10, 64)
// 	if err != nil {
// 		log.Println("ID do agendamento inválido:", err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(ErrorResponse{Message: "ID do agendamento inválido"})
// 		return
// 	}

// 	var schedule models.Schedule
// 	result := h.db.First(&schedule, id)
// 	if result.Error != nil {
// 		log.Println("Erro ao buscar agendamento:", result.Error)
// 		w.WriteHeader(http.StatusNotFound)
// 		json.NewEncoder(w).Encode(ErrorResponse{Message: "Agendamento não encontrado"})
// 		return
// 	}

// 	response := ScheduleResponse{
// 		ID:            schedule.ID,
// 		StreamID:      schedule.StreamID,
// 		ScheduledTime: schedule.ScheduledTime,
// 		CreatedAt:     schedule.CreatedAt,
// 		UpdatedAt:     schedule.UpdatedAt,
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(response)
// }

// func (h *ScheduleHandler) UpdateSchedule(w http.ResponseWriter, r *http.Request) {
// 	user := middlewares.GetUserFromContext(r.Context())
// 	if user == nil {
// 		log.Println("Usuário não autenticado")
// 		w.WriteHeader(http.StatusUnauthorized)
// 		json.NewEncoder(w).Encode(ErrorResponse{Message: "Não autorizado"})
// 		return
// 	}

// 	vars := mux.Vars(r)
// 	id, err := strconv.ParseUint(vars["id"], 10, 64)
// 	if err != nil {
// 		log.Println("ID do agendamento inválido:", err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(ErrorResponse{Message: "ID do agendamento inválido"})
// 		return
// 	}

// 	var request UpdateScheduleRequest
// 	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
// 		log.Println("Erro ao decodificar o corpo da requisição:", err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(ErrorResponse{Message: "Corpo da requisição inválido"})
// 		return
// 	}

// 	var schedule models.Schedule
// 	result := h.db.First(&schedule, id)
// 	if result.Error != nil {
// 		log.Println("Erro ao buscar agendamento:", result.Error)
// 		w.WriteHeader(http.StatusNotFound)
// 		json.NewEncoder(w).Encode(ErrorResponse{Message: "Agendamento não encontrado"})
// 		return
// 	}

// 	schedule.StreamID = request.StreamID
// 	schedule.ScheduledTime = request.ScheduledTime

// 	result = h.db.Save(&schedule)
// 	if result.Error != nil {
// 		log.Println("Erro ao atualizar agendamento:", result.Error)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao atualizar agendamento"})
// 		return
// 	}

// 	response := ScheduleResponse{
// 		ID:            schedule.ID,
// 		StreamID:      schedule.StreamID,
// 		ScheduledTime: schedule.ScheduledTime,
// 		CreatedAt:     schedule.CreatedAt,
// 		UpdatedAt:     schedule.UpdatedAt,
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(response)
// }

// func (h *ScheduleHandler) DeleteSchedule(w http.ResponseWriter, r *http.Request) {
// 	user := middlewares.GetUserFromContext(r.Context())
// 	if user == nil {
// 		log.Println("Usuário não autenticado")
// 		w.WriteHeader(http.StatusUnauthorized)
// 		json.NewEncoder(w).Encode(ErrorResponse{Message: "Não autorizado"})
// 		return
// 	}

// 	vars := mux.Vars(r)
// 	id, err := strconv.ParseUint(vars["id"], 10, 64)
// 	if err != nil {
// 		log.Println("ID do agendamento inválido:", err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(ErrorResponse{Message: "ID do agendamento inválido"})
// 		return
// 	}

// 	var schedule models.Schedule
// 	result := h.db.First(&schedule, id)
// 	if result.Error != nil {
// 		log.Println("Erro ao buscar agendamento:", result.Error)
// 		w.WriteHeader(http.StatusNotFound)
// 		json.NewEncoder(w).Encode(ErrorResponse{Message: "Agendamento não encontrado"})
// 		return
// 	}

// 	result = h.db.Delete(&schedule)
// 	if result.Error != nil {
// 		log.Println("Erro ao excluir agendamento:", result.Error)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao excluir agendamento"})
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }

package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/marcoaureliojf/streamStudio/backend/internal/config"
	"github.com/marcoaureliojf/streamStudio/backend/internal/database"
	"github.com/marcoaureliojf/streamStudio/backend/internal/database/models"
	"github.com/marcoaureliojf/streamStudio/backend/internal/middlewares"
	"github.com/marcoaureliojf/streamStudio/backend/internal/queue"
	"gorm.io/gorm"
)

type ScheduleHandler struct {
	db       *gorm.DB
	rabbitMQ *queue.RabbitMQ
	cfg      config.Config
}

func NewScheduleHandler() *ScheduleHandler {
	cfg := config.LoadConfig()
	rabbitMQ, err := queue.NewRabbitMQ(cfg)
	if err != nil {
		log.Fatal("Erro ao iniciar RabbitMQ: ", err)
	}
	return &ScheduleHandler{db: database.GetDB(), rabbitMQ: rabbitMQ, cfg: cfg}
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

type UpdateScheduleRequest struct {
	StreamID      uint      `json:"streamId"`
	ScheduledTime time.Time `json:"scheduledTime"`
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

	message, err := json.Marshal(schedule)
	if err != nil {
		log.Println("Erro ao serializar mensagem", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao processar agendamento"})
		return
	}

	err = h.rabbitMQ.Publish(message)
	if err != nil {
		log.Println("Erro ao publicar mensagem na fila", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao processar agendamento"})
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

func (h *ScheduleHandler) GetSchedules(w http.ResponseWriter, r *http.Request) {
	user := middlewares.GetUserFromContext(r.Context())
	if user == nil {
		log.Println("Usuário não autenticado")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Não autorizado"})
		return
	}

	var schedules []models.Schedule
	result := h.db.Find(&schedules)
	if result.Error != nil {
		log.Println("Erro ao buscar agendamentos:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao buscar agendamentos"})
		return
	}

	var response []ScheduleResponse
	for _, schedule := range schedules {
		response = append(response, ScheduleResponse{
			ID:            schedule.ID,
			StreamID:      schedule.StreamID,
			ScheduledTime: schedule.ScheduledTime,
			CreatedAt:     schedule.CreatedAt,
			UpdatedAt:     schedule.UpdatedAt,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
func (h *ScheduleHandler) GetSchedule(w http.ResponseWriter, r *http.Request) {
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
		log.Println("ID do agendamento inválido:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "ID do agendamento inválido"})
		return
	}

	var schedule models.Schedule
	result := h.db.First(&schedule, id)
	if result.Error != nil {
		log.Println("Erro ao buscar agendamento:", result.Error)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Agendamento não encontrado"})
		return
	}

	response := ScheduleResponse{
		ID:            schedule.ID,
		StreamID:      schedule.StreamID,
		ScheduledTime: schedule.ScheduledTime,
		CreatedAt:     schedule.CreatedAt,
		UpdatedAt:     schedule.UpdatedAt,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *ScheduleHandler) UpdateSchedule(w http.ResponseWriter, r *http.Request) {
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
		log.Println("ID do agendamento inválido:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "ID do agendamento inválido"})
		return
	}

	var request UpdateScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Erro ao decodificar o corpo da requisição:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Corpo da requisição inválido"})
		return
	}

	var schedule models.Schedule
	result := h.db.First(&schedule, id)
	if result.Error != nil {
		log.Println("Erro ao buscar agendamento:", result.Error)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Agendamento não encontrado"})
		return
	}

	schedule.StreamID = request.StreamID
	schedule.ScheduledTime = request.ScheduledTime

	result = h.db.Save(&schedule)
	if result.Error != nil {
		log.Println("Erro ao atualizar agendamento:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao atualizar agendamento"})
		return
	}

	response := ScheduleResponse{
		ID:            schedule.ID,
		StreamID:      schedule.StreamID,
		ScheduledTime: schedule.ScheduledTime,
		CreatedAt:     schedule.CreatedAt,
		UpdatedAt:     schedule.UpdatedAt,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *ScheduleHandler) DeleteSchedule(w http.ResponseWriter, r *http.Request) {
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
		log.Println("ID do agendamento inválido:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "ID do agendamento inválido"})
		return
	}

	var schedule models.Schedule
	result := h.db.First(&schedule, id)
	if result.Error != nil {
		log.Println("Erro ao buscar agendamento:", result.Error)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Agendamento não encontrado"})
		return
	}

	result = h.db.Delete(&schedule)
	if result.Error != nil {
		log.Println("Erro ao excluir agendamento:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao excluir agendamento"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
