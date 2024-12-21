package main

import (
	"encoding/json"
	"fmt"
	"go-project/config"
	"go-project/logger"
	"go-project/middleware"
	"net/http"

	_ "go-project/docs" // Importa la documentación generada

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// Estructura para un recurso de ejemplo
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Price       float32 `json:"price"`
	Sku         string  `json:"sku"`
}

// Datos simulados
var (
	users    = []User{{ID: 1, Name: "Alice", Age: 25}, {ID: 2, Name: "Bob", Age: 30}}
	products = []Product{}
)

// GET /api/product
func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// Función para manejar la solicitud de obtener usuarios
func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Función para manejar la solicitud de crear un nuevo usuario
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	user.ID = len(users) + 1
	users = append(users, user)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func main() {
	// Configuración y logger
	cfg := config.LoadConfig()
	log := logger.NewLogger(cfg.LogLevel)

	// Crear multiplexer
	mux := http.NewServeMux()

	// Rutas públicas
	mux.HandleFunc("/public", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message": "Public endpoint"}`))
	})
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getUsers(w, r)
		case "POST":
			createUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getProducts(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Ruta protegida
	protectedHandler := middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message": "Protected endpoint"}`))
	}))
	mux.Handle("/protected", protectedHandler)

	// Ruta de Swagger
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// Aplicar middlewares globales
	handlerWithCommonHeaders := middleware.CommonHeadersMiddleware(mux)

	// Iniciar servidor
	log.Info("Starting server at port 8080")
	err := http.ListenAndServe(":8080", handlerWithCommonHeaders)
	if err != nil {
		log.Fatal(fmt.Sprintf("Server failed: %s", err))
	}
}
