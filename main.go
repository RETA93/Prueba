package main

import (
	"database/sql"
	"fmt"
	_ "go-project/docs"
	"go-project/handlers"
	"go-project/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           API de Productos
// @version         1.0
// @description     API para gestión de productos
// @host            localhost:8080
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
	shopHandler := handlers.NewShopHandler(db)
	inventoryHandler := handlers.NewInventoryHandler(db)
	movementHandler := handlers.NewMovementHandler(db)
	// Configurar rutas
	r := mux.NewRouter() // Usamos mux.NewRouter()

	// Ruta para la documentación Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Rutas de la API Productos
	r.HandleFunc("/api/ListarProductos", productHandler.ListarProductos)
	r.HandleFunc("/api/CrearProducto", productHandler.CrearProducto)
	r.HandleFunc("/api/ObtenerProducto", productHandler.ObtenerProducto)
	r.HandleFunc("/api/ActualizarProducto", productHandler.ActualizarProducto)
	r.HandleFunc("/api/ActivarDesactivarProducto", productHandler.ToggleProductoEstado)
	r.HandleFunc("/api/EliminarProducto", productHandler.EliminarProducto)

	// Rutas de la API Tiendas
	r.HandleFunc("/api/ListarTiendas", shopHandler.ListarTiendas)
	r.HandleFunc("/api/CrearTiendas", shopHandler.CrearTienda)
	r.HandleFunc("/api/ObtenerTiendas", shopHandler.ObtenerTienda)
	r.HandleFunc("/api/ActualizarTiendas", shopHandler.ActualizarTienda)
	r.HandleFunc("/api/ActivarDesactivarTiendas", shopHandler.ToggleTiendaEstado)
	r.HandleFunc("/api/EliminarTiendas", shopHandler.EliminarTienda)

	// Rutas de la API Inventarios
	r.HandleFunc("/api/ListarInventarios", inventoryHandler.ListarInventarios)
	r.HandleFunc("/api/CrearInventario", inventoryHandler.CrearInventario)
	r.HandleFunc("/api/ObtenerInventario", inventoryHandler.ObtenerInventario)
	r.HandleFunc("/api/ActualizarInventario", inventoryHandler.ActualizarInventario)
	r.HandleFunc("/api/EliminarInventario", inventoryHandler.EliminarInventario)

	// Rutas de la API Movimientos
	r.HandleFunc("/api/ListarMovimientos", movementHandler.ListarMovimientos)
	r.HandleFunc("/api/CrearMovimiento", movementHandler.CrearMovimiento)
	r.HandleFunc("/api/ObtenerMovimiento", movementHandler.ObtenerMovimiento)
	// Rutas de la API Operacion
	r.HandleFunc("/api/stores/{id}/inventory", inventoryHandler.GetStoreInventory)
	r.HandleFunc("/api/inventory/transfer", inventoryHandler.TransferInventory)
	r.HandleFunc("/api/inventory/alerts", inventoryHandler.GetStockAlerts)

	// Aplicar middleware CORS
	handler := middleware.CORSMiddleware(r)

	// Iniciar servidor
	fmt.Println("Servidor iniciado en http://localhost:8080")
	fmt.Println("Documentación Swagger en http://localhost:8080/swagger/index.html")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
