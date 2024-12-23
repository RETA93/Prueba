package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Tienda modelo básico de tienda
// @Description Modelo de tienda para la API
type Tienda struct {
	ID      uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name    string    `json:"name" example:"Tienda Central"`
	Address string    `json:"address" example:"Av. Principal 123"`
	Phone   string    `json:"phone" example:"555-0123"`
}

// CrearTienda modelo para crear una nueva tienda
type CrearTienda struct {
	Name    string `json:"name" example:"Tienda Central"`
	Address string `json:"address" example:"Av. Principal 123"`
	Phone   string `json:"phone" example:"555-0123"`
}

// TiendaDetalle modelo completo con campos de auditoría
type TiendaDetalle struct {
	ID        uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string    `json:"name" example:"Tienda Central"`
	Address   string    `json:"address" example:"Av. Principal 123"`
	Phone     string    `json:"phone" example:"555-0123"`
	Activo    bool      `json:"activo" example:"true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ShopHandler struct {
	db *sql.DB
}

func NewShopHandler(db *sql.DB) *ShopHandler {
	return &ShopHandler{db: db}
}

// ListarTiendas godoc
// @Summary      Listar todas las tiendas
// @Description  Obtiene la lista de todas las tiendas activas
// @Tags         tiendas
// @Accept       json
// @Produce      json
// @Success      200  {array}   Tienda
// @Failure      500  {object}  map[string]string
// @Router       /ListarTiendas [get]
func (h *ShopHandler) ListarTiendas(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	rows, err := h.db.Query(`
        SELECT id, name, address, phone
        FROM catalogos.tiendas 
        WHERE activo = true
    `)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tiendas []Tienda
	for rows.Next() {
		var t Tienda
		err := rows.Scan(&t.ID, &t.Name, &t.Address, &t.Phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tiendas = append(tiendas, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tiendas)
}

// CrearTienda godoc
// @Summary      Crear nueva tienda
// @Description  Crea una nueva tienda en el sistema
// @Tags         tiendas
// @Accept       json
// @Produce      json
// @Param        tienda body CrearTienda true "Datos de la tienda"
// @Success      201  {object}  TiendaDetalle
// @Failure      400  {object}  map[string]string
// @Router       /CrearTiendas [post]
func (h *ShopHandler) CrearTienda(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var t CrearTienda
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	var tiendaDetalle TiendaDetalle
	tiendaDetalle.ID = uuid.New()

	err := h.db.QueryRow(`
        INSERT INTO catalogos.tiendas (id, name, address, phone, activo, created_at, updated_at)
        VALUES ($1, $2, $3, $4, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        RETURNING id, name, address, phone, activo, created_at, updated_at
    `, tiendaDetalle.ID, t.Name, t.Address, t.Phone).Scan(
		&tiendaDetalle.ID, &tiendaDetalle.Name, &tiendaDetalle.Address,
		&tiendaDetalle.Phone, &tiendaDetalle.Activo, &tiendaDetalle.CreatedAt,
		&tiendaDetalle.UpdatedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tiendaDetalle)
}

// ObtenerTienda godoc
// @Summary      Obtener tienda por ID
// @Description  Obtiene los detalles de una tienda específica
// @Tags         tiendas
// @Accept       json
// @Produce      json
// @Param        id query string true "ID de la tienda"
// @Success      200  {object}  TiendaDetalle
// @Failure      404  {object}  map[string]string
// @Router       /ObtenerTiendas [get]
func (h *ShopHandler) ObtenerTienda(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var tienda TiendaDetalle
	err = h.db.QueryRow(`
        SELECT id, name, address, phone, activo, created_at, updated_at
        FROM catalogos.tiendas 
        WHERE id = $1
    `, id).Scan(
		&tienda.ID, &tienda.Name, &tienda.Address, &tienda.Phone,
		&tienda.Activo, &tienda.CreatedAt, &tienda.UpdatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Tienda no encontrada", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tienda)
}

// ActualizarTienda godoc
// @Summary      Actualizar tienda
// @Description  Actualiza los datos de una tienda existente
// @Tags         tiendas
// @Accept       json
// @Produce      json
// @Param        id query string true "ID de la tienda"
// @Param        tienda body CrearTienda true "Datos actualizados de la tienda"
// @Success      200  {object}  TiendaDetalle
// @Failure      404  {object}  map[string]string
// @Router       /ActualizarTiendas [put]
func (h *ShopHandler) ActualizarTienda(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var t CrearTienda
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	var tienda TiendaDetalle
	err = h.db.QueryRow(`
        UPDATE catalogos.tiendas 
        SET name = $1, address = $2, phone = $3, updated_at = CURRENT_TIMESTAMP
        WHERE id = $4 AND activo = true
        RETURNING id, name, address, phone, activo, created_at, updated_at
    `, t.Name, t.Address, t.Phone, id).Scan(
		&tienda.ID, &tienda.Name, &tienda.Address, &tienda.Phone,
		&tienda.Activo, &tienda.CreatedAt, &tienda.UpdatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Tienda no encontrada o inactiva", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tienda)
}

// ToggleTiendaEstado godoc
// @Summary      Activar/Desactivar tienda
// @Description  Cambia el estado activo/inactivo de una tienda
// @Tags         tiendas
// @Accept       json
// @Produce      json
// @Param        id query string true "ID de la tienda"
// @Param        activate query boolean true "true para activar, false para desactivar"
// @Success      200  {object}  TiendaDetalle
// @Failure      404  {object}  map[string]string
// @Router       /ActivarDesactivarTiendas [patch]
func (h *ShopHandler) ToggleTiendaEstado(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	activate := r.URL.Query().Get("activate") == "true"

	var tienda TiendaDetalle
	err = h.db.QueryRow(`
        UPDATE catalogos.tiendas 
        SET activo = $1, updated_at = CURRENT_TIMESTAMP
        WHERE id = $2
        RETURNING id, name, address, phone, activo, created_at, updated_at
    `, activate, id).Scan(
		&tienda.ID, &tienda.Name, &tienda.Address, &tienda.Phone,
		&tienda.Activo, &tienda.CreatedAt, &tienda.UpdatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Tienda no encontrada", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tienda)
}

// EliminarTienda godoc
// @Summary      Eliminar tienda
// @Description  Elimina una tienda del sistema
// @Tags         tiendas
// @Accept       json
// @Produce      json
// @Param        id query string true "ID de la tienda"
// @Success      204  "No Content"
// @Failure      404  {object}  map[string]string
// @Router       /EliminarTiendas [delete]
func (h *ShopHandler) EliminarTienda(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec(`
        DELETE FROM catalogos.tiendas 
        WHERE id = $1
    `, id)

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
		http.Error(w, "Tienda no encontrada", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
