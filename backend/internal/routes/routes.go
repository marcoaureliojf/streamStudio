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

	protectedRoutes := r.PathPrefix("/api").Subrouter()
	protectedRoutes.Use(middlewares.AuthMiddleware)
	protectedRoutes.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		user := middlewares.GetUserFromContext(r.Context())
		json.NewEncoder(w).Encode(user)
	}).Methods("GET")

	return r
}
