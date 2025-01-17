package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/marcoaureliojf/streamStudio/backend/internal/handlers"
	"github.com/marcoaureliojf/streamStudio/backend/internal/middlewares"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	userHandler := handlers.NewUserHandler()
	r.HandleFunc("/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")

	teamHandler := handlers.NewTeamHandler()
	permissionHandler := handlers.NewPermissionHandler()

	protectedRoutes := r.PathPrefix("/api").Subrouter()
	protectedRoutes.Use(middlewares.AuthMiddleware)
	protectedRoutes.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		user := middlewares.GetUserFromContext(r.Context())
		json.NewEncoder(w).Encode(user)
	}).Methods("GET")
	protectedRoutes.HandleFunc("/teams", teamHandler.Register).Methods("POST")
	protectedRoutes.HandleFunc("/teams", teamHandler.GetTeams).Methods("GET")
	protectedRoutes.HandleFunc("/teams/{id}", teamHandler.GetTeam).Methods("GET")
	protectedRoutes.HandleFunc("/teams/{id}", teamHandler.UpdateTeam).Methods("PUT")
	protectedRoutes.HandleFunc("/teams/{id}", teamHandler.DeleteTeam).Methods("DELETE")

	protectedRoutes.HandleFunc("/permissions", permissionHandler.Register).Methods("POST")
	protectedRoutes.HandleFunc("/permissions", permissionHandler.GetPermissions).Methods("GET")
	protectedRoutes.HandleFunc("/permissions/{id}", permissionHandler.GetPermission).Methods("GET")
	protectedRoutes.HandleFunc("/permissions/{id}", permissionHandler.UpdatePermission).Methods("PUT")
	protectedRoutes.HandleFunc("/permissions/{id}", permissionHandler.DeletePermission).Methods("DELETE")

	return r
}

func SetupStreamRoutes() *mux.Router {
	sr := mux.NewRouter()
	streamHandler := handlers.NewStreamHandler()
	scheduleHandler := handlers.NewScheduleHandler()
	signalingHandler := handlers.NewSignalingHandler()

	protectedRoutes := sr.PathPrefix("/api").Subrouter()
	protectedRoutes.Use(middlewares.AuthMiddleware)
	protectedRoutes.HandleFunc("/offer", signalingHandler.Offer).Methods("POST")
	protectedRoutes.HandleFunc("/icecandidate", signalingHandler.IceCandidate).Methods("POST")
	protectedRoutes.HandleFunc("/streams", streamHandler.Register).Methods("POST")
	protectedRoutes.HandleFunc("/streams", streamHandler.GetStreams).Methods("GET")
	protectedRoutes.HandleFunc("/streams/{id}", streamHandler.GetStream).Methods("GET")
	protectedRoutes.HandleFunc("/streams/{id}", streamHandler.UpdateStream).Methods("PUT")
	protectedRoutes.HandleFunc("/streams/{id}", streamHandler.DeleteStream).Methods("DELETE")

	protectedRoutes.HandleFunc("/schedules", scheduleHandler.Register).Methods("POST")
	protectedRoutes.HandleFunc("/schedules", scheduleHandler.GetSchedules).Methods("GET")
	protectedRoutes.HandleFunc("/schedules/{id}", scheduleHandler.GetSchedule).Methods("GET")
	protectedRoutes.HandleFunc("/schedules/{id}", scheduleHandler.UpdateSchedule).Methods("PUT")
	protectedRoutes.HandleFunc("/schedules/{id}", scheduleHandler.DeleteSchedule).Methods("DELETE")

	return sr
}
