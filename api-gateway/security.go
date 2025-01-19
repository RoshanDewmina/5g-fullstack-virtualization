package api_gateway

import (
    "log"
    "net/http"
    "os"
    "strings"

    "github.com/golang-jwt/jwt/v4"
)

// Example minimal JWT validation middleware
func JWTAuth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
            return
        }

        tokenString := parts[1]
        secret := getJWTSecret()

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(secret), nil
        })
        if err != nil || !token.Valid {
            log.Println("JWT validation error:", err)
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // If we get here, token is valid
        next.ServeHTTP(w, r)
    }
}

func getJWTSecret() string {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        secret = "my-secret-key" // DO NOT use a hardcoded secret in real life
    }
    return secret
}
