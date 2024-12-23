package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-project/utils"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Product es el modelo de producto para la documentación
// @Description Modelo de producto
type Producto struct {
	ID          uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string    `json:"name" example:"Laptop HP" binding:"required"`
	Description string    `json:"description" example:"Laptop HP con procesador Intel i5"`
	Category    string    `json:"category" example:"Electrónicos"`
	Price       float64   `json:"price" example:"999.99" binding:"required"`
	SKU         string    `json:"sku" example:"LAP-001" binding:"required"`
}

type CrearProducto struct {
	Name        string  `json:"name" example:"Laptop HP" binding:"required"`
	Description string  `json:"description" example:"Laptop HP con procesador Intel i5"`
	Category    string  `json:"category" example:"Electrónicos"`
	Price       float64 `json:"price" example:"999.99" binding:"required"`
	SKU         string  `json:"sku" example:"LAP-001" binding:"required"`
}
type ActualizarProducto struct {
	Name        string  `json:"name" example:"Laptop HP Actualizada"`
	Description string  `json:"description" example:"Laptop HP con procesador Intel i7"`
	Category    string  `json:"category" example:"Electrónicos"`
	Price       float64 `json:"price" example:"1299.99"`
	SKU         string  `json:"sku" example:"LAP-002"`
}

type ProductoDetalle struct {
	ID          uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string    `json:"name" example:"Laptop HP" binding:"required"`
	Description string    `json:"description" example:"Laptop HP con procesador Intel i5"`
	Category    string    `json:"category" example:"Electrónicos"`
	Price       float64   `json:"price" example:"999.99" binding:"required"`
	SKU         string    `json:"sku" example:"LAP-001" binding:"required"`
	Activo      bool      `json:"activo" example:"true"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type ProductHandler struct {
	db *sql.DB
}

func NewProductHandler(db *sql.DB) *ProductHandler {
	return &ProductHandler{db: db}
}

// HandleProducts maneja todas las peticiones relacionadas con productos
func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if id := r.URL.Query().Get("id"); id != "" {
			h.ObtenerProducto(w, r)
		} else {
			h.ListarProductos(w, r)
		}
	case http.MethodPost:
		h.CrearProducto(w, r)
	case http.MethodPut:
		h.ActualizarProducto(w, r)
	case http.MethodPatch:
		h.ToggleProductoEstado(w, r)
	case http.MethodDelete:
		h.EliminarProducto(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ListarProductos godoc
// @Summary      Listar productos
// @Description  Obtiene la lista de todos los productos activos
// @Tags         productos
// @Accept       json
// @Produce      json
// @Success      200  {array}   Producto
// @Failure      500  {object}  map[string]string
// @Router       /ListarProductos [get]
func (h *ProductHandler) ListarProductos(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
         SELECT id, name, description, category, price, sku
        FROM catalogos.productos 
        WHERE activo = true
        ORDER BY category, name  
    `)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var productos []Producto
	for rows.Next() {
		var p Producto
		err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.Category,
			&p.Price, &p.SKU,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		productos = append(productos, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(productos)
}

// CrearProducto godoc
// @Summary      Crear producto
// @Description  Crea un nuevo producto en el sistema
// @Tags         productos
// @Accept       json
// @Produce      json
// @Param        producto body CrearProducto true "Datos del producto"
// @Success      201  {object}  ProductoDetalle
// @Failure      400  {object}  map[string]string
// @Router       /CrearProductos [post]
func (h *ProductHandler) CrearProducto(w http.ResponseWriter, r *http.Request) {
	var p CrearProducto
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	err := utils.WithTransaction(h.db, func(tx *sql.Tx) error {
		// Verificar SKU único
		var exists bool
		err := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM catalogos.productos WHERE sku = $1)",
			p.SKU).Scan(&exists)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("SKU ya existe")
		}

		// Insertar producto
		_, err = tx.Exec(`
            INSERT INTO catalogos.productos (
                id, name, description, category, price, sku, activo, created_at
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        `, uuid.New(), p.Name, p.Description, p.Category, p.Price,
			p.SKU, true, time.Now())
		return err
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// ObtenerProducto godoc
// @Summary      Obtener producto por ID
// @Description  Obtiene los detalles de un producto específico
// @Tags         productos
// @Accept       json
// @Produce      json
// @Param        id query string true "ID del producto"
// @Success      200  {object}  ProductoDetalle
// @Failure      404  {object}  map[string]string
// @Router       /ObtenerProductos [get]
func (h *ProductHandler) ObtenerProducto(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var producto ProductoDetalle
	err = h.db.QueryRow(`
        SELECT id, name, description, category, price, sku, activo, created_at, updated_at
        FROM catalogos.productos 
        WHERE id = $1
    `, id).Scan(
		&producto.ID, &producto.Name, &producto.Description, &producto.Category,
		&producto.Price, &producto.SKU, &producto.Activo, &producto.CreatedAt, &producto.UpdatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(producto)
}

// ActualizarProducto godoc
// @Summary      Actualizar producto
// @Description  Actualiza los datos de un producto existente
// @Tags         productos
// @Accept       json
// @Produce      json
// @Param        id query string true "ID del producto"
// @Param        producto body ActualizarProducto true "Datos del producto"
// @Success      200  {object}  ProductoDetalle
// @Failure      404  {object}  map[string]string
// @Router       /ActualizarProductos [put]
func (h *ProductHandler) ActualizarProducto(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var p ActualizarProducto
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec(`
        UPDATE catalogos.productos 
        SET name = $1, description = $2, category = $3, price = $4, sku = $5, updated_at = CURRENT_TIMESTAMP
        WHERE id = $6 AND activo = true
    `, p.Name, p.Description, p.Category, p.Price, p.SKU, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "Producto no encontrado o inactivo", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// ToggleProductoEstado godoc
// @Summary      Activar/Desactivar producto
// @Description  Cambia el estado activo/inactivo de un producto
// @Tags         productos
// @Accept       json
// @Produce      json
// @Param        id query string true "ID del producto"
// @Param        activate query boolean true "true para activar, false para desactivar"
// @Success      200  {object}  ProductoDetalle
// @Failure      404  {object}  map[string]string
// @Router       /ActivarDesactivarProductos [patch]
func (h *ProductHandler) ToggleProductoEstado(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	activate := r.URL.Query().Get("activate") == "true"

	result, err := h.db.Exec(`
        UPDATE catalogos.productos 
        SET activo = $1, updated_at = CURRENT_TIMESTAMP
        WHERE id = $2
    `, activate, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// EliminarProducto godoc
// @Summary      Eliminar producto
// @Description  Elimina un producto del sistema
// @Tags         productos
// @Accept       json
// @Produce      json
// @Param        id query string true "ID del producto"
// @Success      204  "No Content"
// @Failure      404  {object}  map[string]string
// @Router       /EliminarProductos [delete]
func (h *ProductHandler) EliminarProducto(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec(`
        DELETE FROM catalogos.productos 
        WHERE id = $1
    `, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
