package handlers

import (
	"database/sql"
	"encoding/json"
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
		// Si tiene ID en query params, es detalle, si no, es listado
		if id := r.URL.Query().Get("id"); id != "" {
			h.GetProductDetail(w, r)
		} else {
			h.GetProducts(w, r)
		}
	case http.MethodPost:
		h.CreateProduct(w, r)
	case http.MethodPut:
		h.UpdateProduct(w, r)
	case http.MethodPatch:
		// Si tiene el parámetro activate, es activación/desactivación
		if r.URL.Query().Get("activate") != "" {
			h.ToggleProductStatus(w, r)
		} else {
			http.Error(w, "Operación no válida", http.StatusBadRequest)
		}
	case http.MethodDelete:
		h.DeleteProduct(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// GetProducts godoc
// @Summary      Listar productos
// @Description  Obtiene la lista de todos los productos activos
// @Tags         productos
// @Accept       json
// @Produce      json
// @Success      200 {array}   Producto
// @Failure      500 {object}  map[string]string
// @Router       /ListarProductos [get]
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
        SELECT id, name, description, category, price, sku
        FROM catalogos.productos 
        WHERE activo = true
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

// CreateProduct godoc
// @Summary      Crear producto
// @Description  Crea un nuevo producto en el sistema
// @Tags         productos
// @Accept       json
// @Produce      json
// @Param        producto  body CrearProducto  true  "Datos del producto"
// @Success      201 {object}  ProductoDetalle
// @Failure      400 {object}  map[string]string
// @Failure      500 {object}  map[string]string
// @Router       /CrearProductos [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p CrearProducto
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	// Datos por default
	var ID uuid.UUID = uuid.New()
	var Activo bool = true
	var CreatedAt time.Time = time.Now()

	_, err := h.db.Exec(`
        INSERT INTO catalogos.productos (id, name, description, category, price, sku, activo, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `, ID, p.Name, p.Description, p.Category, p.Price, p.SKU, Activo, CreatedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Crear respuesta usando ProductoDetalle
	respuesta := ProductoDetalle{
		ID:          ID,
		Name:        p.Name,
		Description: p.Description,
		Category:    p.Category,
		Price:       p.Price,
		SKU:         p.SKU,
		Activo:      Activo,
		CreatedAt:   CreatedAt,
		UpdatedAt:   CreatedAt, // En la creación, UpdatedAt es igual a CreatedAt
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(respuesta)
}

// GetProductDetail godoc
// @Summary      Obtener detalle de producto
// @Description  Obtiene los detalles completos de un producto específico
// @Tags         productos
// @Accept       json
// @Produce      json
// @Param        id query string true "ID del producto"
// @Success      200  {object}  ProductoDetalle
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /ObtenerProductos [get]
func (h *ProductHandler) GetProductDetail(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "ID de producto inválido", http.StatusBadRequest)
		return
	}

	var producto ProductoDetalle
	err = h.db.QueryRow(`
        SELECT id, name, description, category, price, sku, activo, created_at, updated_at
        FROM catalogos.productos 
        WHERE id = $1
    `, id).Scan(
		&producto.ID, &producto.Name, &producto.Description, &producto.Category,
		&producto.Price, &producto.SKU, &producto.Activo, &producto.CreatedAt, &producto.UpdatedAt,
	)

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

// UpdateProduct godoc
// @Summary      Actualizar producto
// @Description  Actualiza los datos de un producto existente
// @Tags         productos
// @Accept       json
// @Produce      json
// @Param        id query string true "ID del producto"
// @Param        producto body ActualizarProducto true "Datos del producto"
// @Success      200  {object}  ProductoDetalle
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /ActualizarProductos [put]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "ID de producto inválido", http.StatusBadRequest)
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

	rows, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rows == 0 {
		http.Error(w, "Producto no encontrado o inactivo", http.StatusNotFound)
		return
	}

	// Obtener el producto actualizado
	var producto ProductoDetalle
	err = h.db.QueryRow(`
        SELECT id, name, description, category, price, sku, activo, created_at, updated_at
        FROM catalogos.productos 
        WHERE id = $1
    `, id).Scan(
		&producto.ID, &producto.Name, &producto.Description, &producto.Category,
		&producto.Price, &producto.SKU, &producto.Activo, &producto.CreatedAt, &producto.UpdatedAt,
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(producto)
}

// ToggleProductStatus godoc
// @Summary      Activar/Desactivar producto
// @Description  Cambia el estado activo/inactivo de un producto
// @Tags         productos
// @Accept       json
// @Produce      json
// @Param        id query string true "ID del producto"
// @Param        activate query boolean true "true para activar, false para desactivar"
// @Success      200  {object}  ProductoDetalle
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /ActivarDesactivarProductos [patch]
func (h *ProductHandler) ToggleProductStatus(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "ID de producto inválido", http.StatusBadRequest)
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

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}

	var producto ProductoDetalle
	err = h.db.QueryRow(`
        SELECT id, name, description, category, price, sku, activo, created_at, updated_at
        FROM catalogos.productos 
        WHERE id = $1
    `, id).Scan(
		&producto.ID, &producto.Name, &producto.Description, &producto.Category,
		&producto.Price, &producto.SKU, &producto.Activo, &producto.CreatedAt, &producto.UpdatedAt,
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(producto)
}

// DeleteProduct godoc
// @Summary      Eliminar producto
// @Description  Elimina permanentemente un producto
// @Tags         productos
// @Accept       json
// @Produce      json
// @Param        id query string true "ID del producto"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /EliminarProductos [delete]
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "ID de producto inválido", http.StatusBadRequest)
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

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
