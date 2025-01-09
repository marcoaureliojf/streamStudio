package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/marcoaureliojf/streamStudio/backend/internal/config"
	"github.com/marcoaureliojf/streamStudio/backend/internal/database/models"
)

type Claims struct {
	UserID uint `json:"userId"`
	jwt.StandardClaims
}

func GenerateToken(user models.User, cfg config.Config) (string, error) {
	claims := Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token válido por 24 horas
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		log.Println("Erro ao assinar o token:", err)
		return "", err
	}
	return signedToken, nil
}

func ValidateToken(tokenString string, cfg config.Config) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inválido: %v", token.Header["alg"])
		}
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		log.Println("Erro ao parsear token: ", err)
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		log.Println("Token inválido: ", err)
		return nil, fmt.Errorf("token inválido")
	}
	return claims, nil
}
