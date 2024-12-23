package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

// Estructuras necesarias
type CrearProducto struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	SKU         string  `json:"sku"`
}

type CrearTienda struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type CrearInventario struct {
	ProductID uuid.UUID `json:"product_id"`
	StoreID   uuid.UUID `json:"store_id"`
	Quantity  int       `json:"quantity"`
	MinStock  int       `json:"min_stock"`
}

type StockTransfer struct {
	ProductID     uuid.UUID `json:"product_id"`
	SourceStoreID uuid.UUID `json:"source_store_id"`
	TargetStoreID uuid.UUID `json:"target_store_id"`
	Quantity      int       `json:"quantity"`
}

type ProductoDetalle struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`
	SKU         string    `json:"sku"`
}

type TiendaDetalle struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Address string    `json:"address"`
	Phone   string    `json:"phone"`
}

type InventarioDetalle struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	StoreID   uuid.UUID `json:"store_id"`
	Quantity  int       `json:"quantity"`
}

func TestTransferInventoryFlow(t *testing.T) {
	// Crear producto
	producto := CrearProducto{
		Name:        "Laptop HP",
		Description: "Laptop HP con procesador Intel i5",
		Category:    "Electrónicos",
		Price:       999.99,
		SKU:         "LAP-101",
	}
	productoID := crearProducto(t, producto)

	// Crear tiendas
	tienda1 := CrearTienda{
		Name:    "Tienda Central",
		Address: "Av. Principal 123",
		Phone:   "555-0123",
	}
	tienda2 := CrearTienda{
		Name:    "Sucursal Norte",
		Address: "Blvd. Manuel Ávila Camacho 2000",
		Phone:   "555-0124",
	}
	tienda1ID := crearTienda(t, tienda1)
	tienda2ID := crearTienda(t, tienda2)

	// Crear inventario inicial
	inventario := CrearInventario{
		ProductID: productoID,
		StoreID:   tienda1ID,
		Quantity:  100,
		MinStock:  10,
	}
	crearInventario(t, inventario)

	// Realizar transferencia
	transfer := StockTransfer{
		ProductID:     productoID,
		SourceStoreID: tienda1ID,
		TargetStoreID: tienda2ID,
		Quantity:      50,
	}

	// Verificar resultado
	resp := realizarTransferencia(t, transfer)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}

	// Verificar inventario final
	verificarInventario(t, tienda1ID, productoID, 50)
	verificarInventario(t, tienda2ID, productoID, 50)
}

func crearProducto(t *testing.T, producto CrearProducto) uuid.UUID {
	cuerpo, err := json.Marshal(producto)
	if err != nil {
		t.Fatalf("Error al convertir producto a JSON: %v", err)
	}

	resp, err := http.Post("http://localhost:3000/api/CrearProducto",
		"application/json", bytes.NewBuffer(cuerpo))
	if err != nil {
		t.Fatalf("Error al crear producto: %v", err)
	}
	defer resp.Body.Close()

	// Depuración de respuesta
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error al leer respuesta: %v", err)
	}
	t.Logf("Respuesta del servidor: %s", string(bodyBytes))

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Status code inesperado: %d", resp.StatusCode)
	}

	var respuesta ProductoDetalle
	if err := json.Unmarshal(bodyBytes, &respuesta); err != nil {
		t.Fatalf("Error al decodificar respuesta: %v. Body: %s", err, string(bodyBytes))
	}
	return respuesta.ID
}

func crearTienda(t *testing.T, tienda CrearTienda) uuid.UUID {
	cuerpo, err := json.Marshal(tienda)
	if err != nil {
		t.Fatalf("Error al convertir tienda a JSON: %v", err)
	}

	resp, err := http.Post("http://localhost:3000/api/CrearTiendas",
		"application/json", bytes.NewBuffer(cuerpo))
	if err != nil {
		t.Fatalf("Error al crear tienda: %v", err)
	}
	defer resp.Body.Close()

	var respuesta TiendaDetalle
	err = json.NewDecoder(resp.Body).Decode(&respuesta)
	if err != nil {
		t.Fatalf("Error al decodificar respuesta: %v", err)
	}
	return respuesta.ID
}

func crearInventario(t *testing.T, inv CrearInventario) {
	cuerpo, err := json.Marshal(inv)
	if err != nil {
		t.Fatalf("Error al convertir inventario a JSON: %v", err)
	}

	resp, err := http.Post("http://localhost:3000/api/CrearInventario",
		"application/json", bytes.NewBuffer(cuerpo))
	if err != nil {
		t.Fatalf("Error al crear inventario: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Se esperaba estado 201, se obtuvo %d", resp.StatusCode)
	}
}

func realizarTransferencia(t *testing.T, transferencia StockTransfer) *http.Response {
	cuerpo, err := json.Marshal(transferencia)
	if err != nil {
		t.Fatalf("Error al convertir transferencia a JSON: %v", err)
	}

	resp, err := http.Post("http://localhost:3000/api/inventory/transfer",
		"application/json", bytes.NewBuffer(cuerpo))
	if err != nil {
		t.Fatalf("Error al realizar transferencia: %v", err)
	}
	return resp
}
func verificarInventario(t *testing.T, tiendaID, productoID uuid.UUID, cantidadEsperada int) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:3000/api/stores/%s/inventory", tiendaID))
	if err != nil {
		t.Fatalf("Error al obtener inventario: %v", err)
	}
	defer resp.Body.Close()

	var inventario []InventarioDetalle
	err = json.NewDecoder(resp.Body).Decode(&inventario)
	if err != nil {
		t.Fatalf("Error al decodificar inventario: %v", err)
	}

	for _, inv := range inventario {
		if inv.ProductID == productoID && inv.Quantity != cantidadEsperada {
			t.Errorf("Se esperaba cantidad %d, se obtuvo %d", cantidadEsperada, inv.Quantity)
		}
	}
}
