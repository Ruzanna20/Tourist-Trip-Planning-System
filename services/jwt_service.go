package services

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secretKey []byte
	expiry    time.Duration
}

func NewJWTService(jwtSecret string, expiryMinutesStr string) *JWTService {
	return &JWTService{
		secretKey: []byte(jwtSecret),
		expiry:    24 * time.Hour,
	}
}

func (s *JWTService) GenerateToken(userID int) (string, error) {
	claims := CustomClaims{
		UserID: userID,

		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("could not sign token:%w", err)
	}
	return tokenStr, nil
}

func (s *JWTService) ValidateToken(tokenStr string) (*CustomClaims, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexcepted signing method:%v", t.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return claims, nil

}

func (s *JWTService) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := bearerToken[1]

		claims, err := s.ValidateToken(tokenString)
		if err != nil {
			MSG := fmt.Sprintf("Invalid or expired token: %v", err)
			if validationErr, ok := err.(*jwt.ValidationError); ok && validationErr.Errors == jwt.ValidationErrorExpired {
				MSG = "Token has expired"
			}
			http.Error(w, MSG, http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User-ID", fmt.Sprintf("%d", claims.UserID))
		next(w, r)
	}
}

func (s *JWTService) GenerateRefreshToken(userID int) (string, error) {
	refreshExpiry := 7 * 24 * time.Hour

	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("could not sign token:%w", err)
	}
	return tokenStr, nil

}
