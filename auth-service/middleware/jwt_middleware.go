package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Ključevi za context da ne bi došlo do konflikta
type contextKey string

const (
	ContextUserIDKey contextKey = "userID"
	ContextRoleKey   contextKey = "role"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// JWTAuth middleware validira JWT token i ubacuje userID i role u request context
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Mora početi sa "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Parsiranje tokena sa claims
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Dohvati user_id i role iz claims
		userID, okUser := claims["user_id"].(string)
		role, okRole := claims["role"].(string)
		fmt.Printf("✅ Token claims -> user_id: %s, role: %s\n", userID, role)

		if !okUser || !okRole || userID == "" || role == "" {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// LOGOVANJE
		fmt.Printf("✅ Token claims -> user_id: %s, role: %s\n", userID, role)

		// Ubaci u context
		ctx := context.WithValue(r.Context(), ContextUserIDKey, userID)
		ctx = context.WithValue(ctx, ContextRoleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
