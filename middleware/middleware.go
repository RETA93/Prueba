package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// CommonHeadersMiddleware agrega encabezados comunes a todas las respuestas
func CommonHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware valida el token JWT en el encabezado Authorization
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtener el token del encabezado Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Extraer el token del encabezado
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // No se encontró el prefijo Bearer
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		// Validar el token (reemplaza "your-secret-key" con tu clave secreta)
		_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Asegúrate de usar el mismo método de firma
			return []byte("your-secret-key"), nil
		})
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Continuar con el siguiente controlador
		next.ServeHTTP(w, r)
	})
}
