package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"api/utils"
)

// Middleware pour vérifier le token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Token manquant", http.StatusUnauthorized)
			return
		}

		// Extraire le token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Vérifier le token
		claims, err := utils.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, "Token invalide", http.StatusUnauthorized)
			return
		}

		// Ajouter l'ID utilisateur dans le contexte
		r.Header.Set("UserID", strconv.Itoa(claims.UserID))
		next.ServeHTTP(w, r)
	})
}
