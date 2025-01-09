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
	r := mux.NewRouter()
	streamHandler := handlers.NewStreamHandler()
	scheduleHandler := handlers.NewScheduleHandler()

	protectedRoutes := r.PathPrefix("/api").Subrouter()
	protectedRoutes.Use(middlewares.AuthMiddleware)
	protectedRoutes.HandleFunc("/streams", streamHandler.Register).Methods("POST")
	protectedRoutes.HandleFunc("/schedules", scheduleHandler.Register).Methods("POST")
	return r
}
