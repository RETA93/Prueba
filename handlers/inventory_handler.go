package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Inventario modelo básico
// @Description Modelo de inventario
type Inventario struct {
	ID        uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ProductID uuid.UUID `json:"product_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	StoreID   uuid.UUID `json:"store_id" example:"550e8400-e29b-41d4-a716-446655440002"`
	Quantity  int       `json:"quantity" example:"100"`
	MinStock  int       `json:"min_stock" example:"10"`
}

// CrearInventario modelo para crear nuevo inventario
// @Description Modelo para crear nuevo inventario
type CrearInventario struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	StoreID   uuid.UUID `json:"store_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=0"`
	MinStock  int       `json:"min_stock" binding:"required,min=0"`
}

// InventarioDetalle modelo completo con campos de auditoría
// @Description Modelo detallado de inventario
type InventarioDetalle struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	StoreID   uuid.UUID `json:"store_id"`
	Quantity  int       `json:"quantity"`
	MinStock  int       `json:"min_stock"`
	Activo    bool      `json:"activo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Campos adicionales para mostrar información relacionada
	ProductName string `json:"product_name"`
	StoreName   string `json:"store_name"`
}

type InventoryHandler struct {
	db *sql.DB
}

func NewInventoryHandler(db *sql.DB) *InventoryHandler {
	return &InventoryHandler{db: db}
}

// ListarInventarios godoc
// @Summary      Listar inventarios
// @Description  Obtiene la lista de todos los inventarios activos
// @Tags         inventarios
// @Accept       json
// @Produce      json
// @Success      200  {array}   InventarioDetalle
// @Failure      500  {object}  map[string]string
// @Router       /ListarInventarios [get]
func (h *InventoryHandler) ListarInventarios(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	query := `
        SELECT 
            i.id, i.productId, i.storeId, i.quantity, i.minStock,
            i.activo, i.created_at, i.updated_at,
            p.name as product_name, t.name as store_name
        FROM prueba.inventarios i
        JOIN catalogos.productos p ON i.productId = p.id
        JOIN catalogos.tiendas t ON i.storeId = t.id
        WHERE i.activo = true`

	rows, err := h.db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var inventarios []InventarioDetalle
	for rows.Next() {
		var i InventarioDetalle
		err := rows.Scan(
			&i.ID, &i.ProductID, &i.StoreID, &i.Quantity, &i.MinStock,
			&i.Activo, &i.CreatedAt, &i.UpdatedAt,
			&i.ProductName, &i.StoreName,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		inventarios = append(inventarios, i)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inventarios)
}

// CrearInventario godoc
// @Summary      Crear inventario
// @Description  Crea un nuevo registro de inventario
// @Tags         inventarios
// @Accept       json
// @Produce      json
// @Param        inventario body CrearInventario true "Datos del inventario"
// @Success      201  {object}  InventarioDetalle
// @Failure      400  {object}  map[string]string
// @Router       /CrearInventario [post]
func (h *InventoryHandler) CrearInventario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var inv CrearInventario
	if err := json.NewDecoder(r.Body).Decode(&inv); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	if inv.Quantity < 0 || inv.MinStock < 0 {
		http.Error(w, "Cantidad y stock mínimo deben ser no negativos", http.StatusBadRequest)
		return
	}

	// Verificar que el producto y la tienda existen
	var exists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM catalogos.productos WHERE id = $1)", inv.ProductID).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "Producto no encontrado", http.StatusBadRequest)
		return
	}

	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM catalogos.tiendas WHERE id = $1)", inv.StoreID).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "Tienda no encontrada", http.StatusBadRequest)
		return
	}

	var inventario InventarioDetalle
	err = h.db.QueryRow(`
        INSERT INTO prueba.inventarios (id, productId, storeId, quantity, minStock, activo)
        VALUES ($1, $2, $3, $4, $5, true)
        RETURNING id, productId, storeId, quantity, minStock, activo, created_at, updated_at
    `, uuid.New(), inv.ProductID, inv.StoreID, inv.Quantity, inv.MinStock).Scan(
		&inventario.ID, &inventario.ProductID, &inventario.StoreID,
		&inventario.Quantity, &inventario.MinStock, &inventario.Activo,
		&inventario.CreatedAt, &inventario.UpdatedAt,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(inventario)
}

// ObtenerInventario godoc
// @Summary      Obtener inventario por ID
// @Description  Obtiene los detalles de un inventario específico
// @Tags         inventarios
// @Accept       json
// @Produce      json
// @Param        id query string true "ID del inventario"
// @Success      200  {object}  InventarioDetalle
// @Failure      404  {object}  map[string]string
// @Router       /ObtenerInventario [get]
func (h *InventoryHandler) ObtenerInventario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	query := `
        SELECT 
            i.id, i.productId, i.storeId, i.quantity, i.minStock,
            i.activo, i.created_at, i.updated_at,
            p.name as product_name, t.name as store_name
        FROM prueba.inventarios i
        JOIN catalogos.productos p ON i.productId = p.id
        JOIN catalogos.tiendas t ON i.storeId = t.id
        WHERE i.id = $1`

	var inv InventarioDetalle
	err = h.db.QueryRow(query, id).Scan(
		&inv.ID, &inv.ProductID, &inv.StoreID, &inv.Quantity, &inv.MinStock,
		&inv.Activo, &inv.CreatedAt, &inv.UpdatedAt,
		&inv.ProductName, &inv.StoreName,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Inventario no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inv)
}

// ActualizarInventario godoc
// @Summary      Actualizar inventario
// @Description  Actualiza un registro de inventario existente
// @Tags         inventarios
// @Accept       json
// @Produce      json
// @Param        id query string true "ID del inventario"
// @Param        inventario body CrearInventario true "Datos del inventario"
// @Success      200  {object}  InventarioDetalle
// @Failure      404  {object}  map[string]string
// @Router       /ActualizarInventario [put]
func (h *InventoryHandler) ActualizarInventario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var inv CrearInventario
	if err := json.NewDecoder(r.Body).Decode(&inv); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec(`
        UPDATE prueba.inventarios 
        SET quantity = $1, minStock = $2, updated_at = CURRENT_TIMESTAMP
        WHERE id = $3 AND activo = true
    `, inv.Quantity, inv.MinStock, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "Inventario no encontrado o inactivo", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// EliminarInventario godoc
// @Summary      Eliminar inventario
// @Description  Elimina un registro de inventario
// @Tags         inventarios
// @Accept       json
// @Produce      json
// @Param        id query string true "ID del inventario"
// @Success      204  "No Content"
// @Failure      404  {object}  map[string]string
// @Router       /EliminarInventario [delete]
func (h *InventoryHandler) EliminarInventario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec("DELETE FROM prueba.inventarios WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "Inventario no encontrado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
