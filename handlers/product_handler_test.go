package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "host=postgres port=5432 user=root password=root dbname=root_test sslmode=disable")
	if err != nil {
		t.Fatalf("Error connecting to test database: %v", err)
	}
	return db
}

func TestListarProductos(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	handler := NewProductHandler(db)
	req := httptest.NewRequest("GET", "/api/ListarProductos", nil)
	w := httptest.NewRecorder()

	handler.ListarProductos(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var productos []Producto
	err := json.NewDecoder(w.Body).Decode(&productos)
	if err != nil {
		t.Errorf("Error decoding response: %v", err)
	}
}

func TestCrearProducto(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	handler := NewProductHandler(db)
	producto := CrearProducto{
		Name:        "Test Product",
		Description: "Test Description",
		Category:    "Test Category",
		Price:       99.99,
		SKU:         "TEST-001",
	}

	body, _ := json.Marshal(producto)
	req := httptest.NewRequest("POST", "/api/CrearProducto", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CrearProducto(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
}
