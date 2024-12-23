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
	shopHandler := handlers.NewShopHandler(db)
	inventoryHandler := handlers.NewInventoryHandler(db)
	movementHandler := handlers.NewMovementHandler(db)
	// Configurar rutas
	mux := http.NewServeMux()

	// Ruta para la documentación Swagger
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3000/swagger/doc.json"), //The url pointing to API definition
	))

	// Rutas de la API Productos
	mux.HandleFunc("/api/ListarProductos", productHandler.ListarProductos)
	mux.HandleFunc("/api/CrearProducto", productHandler.CrearProducto)
	mux.HandleFunc("/api/ObtenerProducto", productHandler.ObtenerProducto)
	mux.HandleFunc("/api/ActualizarProducto", productHandler.ActualizarProducto)
	mux.HandleFunc("/api/ActivarDesactivarProducto", productHandler.ToggleProductoEstado)
	mux.HandleFunc("/api/EliminarProducto", productHandler.EliminarProducto)

	// Rutas de la API Tiendas
	mux.HandleFunc("/api/ListarTiendas", shopHandler.ListarTiendas)
	mux.HandleFunc("/api/CrearTiendas", shopHandler.CrearTienda)
	mux.HandleFunc("/api/ObtenerTiendas", shopHandler.ObtenerTienda)
	mux.HandleFunc("/api/ActualizarTiendas", shopHandler.ActualizarTienda)
	mux.HandleFunc("/api/ActivarDesactivarTiendas", shopHandler.ToggleTiendaEstado)
	mux.HandleFunc("/api/EliminarTiendas", shopHandler.EliminarTienda)

	// Rutas de la API Inventarios
	mux.HandleFunc("/api/ListarInventarios", inventoryHandler.ListarInventarios)
	mux.HandleFunc("/api/CrearInventario", inventoryHandler.CrearInventario)
	mux.HandleFunc("/api/ObtenerInventario", inventoryHandler.ObtenerInventario)
	mux.HandleFunc("/api/ActualizarInventario", inventoryHandler.ActualizarInventario)
	mux.HandleFunc("/api/EliminarInventario", inventoryHandler.EliminarInventario)

	// Rutas de la API Movimientos
	mux.HandleFunc("/api/ListarMovimientos", movementHandler.ListarMovimientos)
	mux.HandleFunc("/api/CrearMovimiento", movementHandler.CrearMovimiento)
	mux.HandleFunc("/api/ObtenerMovimiento", movementHandler.ObtenerMovimiento)

	mux.HandleFunc("/api/stores/{id}/inventory", inventoryHandler.GetStoreInventory)
	mux.HandleFunc("/api/inventory/transfer", inventoryHandler.TransferInventory)
	mux.HandleFunc("/api/inventory/alerts", inventoryHandler.GetStockAlerts)

	// Aplicar middleware CORS
	handler := middleware.CORSMiddleware(mux)

	// Iniciar servidor
	fmt.Println("Servidor iniciado en http://localhost:3000")
	fmt.Println("Documentación Swagger en http://localhost:3000/swagger/index.html")
	log.Fatal(http.ListenAndServe(":3000", handler))
}
