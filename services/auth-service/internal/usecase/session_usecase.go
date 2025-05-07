package usecase

import (
	"BikeStoreGolang/services/auth-service/internal/domain"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type SessionUsecase struct {
	jwtSecret string
}

func NewSessionUsecase(jwtSecret string) *SessionUsecase {
	return &SessionUsecase{jwtSecret: jwtSecret}
}

// GenerateToken создает JWT токен для пользователя
func (s *SessionUsecase) GenerateToken(userID string, role domain.Role) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Токен на 24 часа
	})
	return token.SignedString([]byte(s.jwtSecret))
}

// ValidateToken проверяет JWT токен
func (s *SessionUsecase) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}