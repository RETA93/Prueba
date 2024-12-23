package main

import (
	"database/sql"
	"fmt"
	_ "go-project/docs"
	"go-project/handlers"
	"go-project/middleware"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           API de Productos
// @version         1.0
// @description     API para gestión de productos
// @host            localhost:3000
// @BasePath        /api
// @schemes         http

func main() {
	// Configuración de la base de datos
	db, err := sql.Open("postgres", "host=postgres port=5432 user=root password=root dbname=root sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Verificar conexión
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Crear handlers
	productHandler := handlers.NewProductHandler(db)

	// Configurar rutas
	mux := http.NewServeMux()

	// Ruta para la documentación Swagger
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3000/swagger/doc.json"), //The url pointing to API definition
	))

	// Rutas de la API
	mux.HandleFunc("/api/ListarProductos", productHandler.HandleProducts)
	mux.HandleFunc("/api/CrearProductos", productHandler.HandleProducts)
	mux.HandleFunc("/api/ObtenerProductos", productHandler.HandleProducts)
	mux.HandleFunc("/api/ActualizarProductos", productHandler.HandleProducts)
	mux.HandleFunc("/api/ActivarDesactivarProductos", productHandler.HandleProducts)
	mux.HandleFunc("/api/EliminarProductos", productHandler.HandleProducts)

	// Aplicar middleware CORS
	handler := middleware.CORSMiddleware(mux)

	// Iniciar servidor
	fmt.Println("Servidor iniciado en http://localhost:3000")
	fmt.Println("Documentación Swagger en http://localhost:3000/swagger/index.html")
	log.Fatal(http.ListenAndServe(":3000", handler))
}
