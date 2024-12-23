package models

import (
	"time"

	"github.com/google/uuid"
)

type Producto struct {
	// ID único del producto
	ID uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	// Nombre del producto
	Name string `json:"name" example:"Laptop"`
	// Descripción detallada del producto
	Description string `json:"description" example:"Laptop HP Pavilion con procesador Intel i5"`
	// Categoría del producto
	Category string `json:"category" example:"Electrónicos"`
	// Precio del producto
	Price float64 `json:"price" example:"12999.99"`
	// SKU único del producto
	SKU string `json:"sku" example:"LAP-HP-001"`
	// Indica si el producto está activo
	Activo bool `json:"activo" example:"true"`
	// Fecha de creación del registro
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	// Fecha de última actualización
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

type Tienda struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	Activo    bool      `json:"activo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Inventario struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	StoreID   uuid.UUID `json:"store_id"`
	Quantity  int       `json:"quantity"`
	MinStock  int       `json:"min_stock"`
	Activo    bool      `json:"activo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Movimiento struct {
	ID            uuid.UUID `json:"id"`
	ProductID     uuid.UUID `json:"product_id"`
	SourceStoreID uuid.UUID `json:"source_store_id"`
	TargetStoreID uuid.UUID `json:"target_store_id"`
	Quantity      int       `json:"quantity"`
	Timestamp     time.Time `json:"timestamp"`
	Type          string    `json:"type"`
	Activo        bool      `json:"activo"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
