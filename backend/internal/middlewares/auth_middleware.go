package middlewares

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/marcoaureliojf/streamStudio/backend/internal/auth"
	"github.com/marcoaureliojf/streamStudio/backend/internal/config"
	"github.com/marcoaureliojf/streamStudio/backend/internal/database"
	"github.com/marcoaureliojf/streamStudio/backend/internal/database/models"
)

type contextKey string

const UserContextKey contextKey = "user"

type ErrorResponse struct {
	Message string `json:"message"`
}

func AuthMiddleware(next http.Handler) http.Handler {
	cfg := config.LoadConfig()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("Token não enviado")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "Token não enviado"})
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		claims, err := auth.ValidateToken(tokenString, cfg)
		if err != nil {
			log.Println("Token inválido:", err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "Token inválido"})
			return
		}

		var user models.User
		result := database.GetDB().First(&user, claims.UserID)
		if result.Error != nil {
			log.Println("Usuário não encontrado:", result.Error)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "Usuário não encontrado"})
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(ctx context.Context) *models.User {
	user, ok := ctx.Value(UserContextKey).(models.User)
	if !ok {
		return nil
	}
	return &user
}
