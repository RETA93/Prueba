package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// MovimientoTipo tipo de movimiento
type MovimientoTipo string

const (
	MovimientoIN       MovimientoTipo = "IN"
	MovimientoOUT      MovimientoTipo = "OUT"
	MovimientoTRANSFER MovimientoTipo = "TRANSFER"
)

// Movimiento modelo básico
type Movimiento struct {
	ID            uuid.UUID      `json:"id"`
	ProductID     uuid.UUID      `json:"product_id"`
	SourceStoreID uuid.UUID      `json:"source_store_id"`
	TargetStoreID uuid.UUID      `json:"target_store_id"`
	Quantity      int            `json:"quantity"`
	Type          MovimientoTipo `json:"type"`
	Timestamp     time.Time      `json:"timestamp"`
}

// CrearMovimiento modelo para crear movimiento
type CrearMovimiento struct {
	ProductID     uuid.UUID      `json:"product_id" binding:"required"`
	SourceStoreID uuid.UUID      `json:"source_store_id" binding:"required"`
	TargetStoreID uuid.UUID      `json:"target_store_id" binding:"required"`
	Quantity      int            `json:"quantity" binding:"required,gt=0"`
	Type          MovimientoTipo `json:"type" binding:"required"`
}

// MovimientoDetalle modelo completo
type MovimientoDetalle struct {
	ID            uuid.UUID      `json:"id"`
	ProductID     uuid.UUID      `json:"product_id"`
	SourceStoreID uuid.UUID      `json:"source_store_id"`
	TargetStoreID uuid.UUID      `json:"target_store_id"`
	Quantity      int            `json:"quantity"`
	Type          MovimientoTipo `json:"type"`
	Timestamp     time.Time      `json:"timestamp"`
	Activo        bool           `json:"activo"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	// Campos adicionales para información relacionada
	ProductName     string `json:"product_name"`
	SourceStoreName string `json:"source_store_name"`
	TargetStoreName string `json:"target_store_name"`
}

type MovementHandler struct {
	db *sql.DB
}

func NewMovementHandler(db *sql.DB) *MovementHandler {
	return &MovementHandler{db: db}
}

// ListarMovimientos godoc
// @Summary      Listar movimientos
// @Description  Obtiene la lista de todos los movimientos de inventario
// @Tags         movimientos
// @Accept       json
// @Produce      json
// @Success      200  {array}   MovimientoDetalle
// @Failure      500  {object}  map[string]string
// @Router       /ListarMovimientos [get]
func (h *MovementHandler) ListarMovimientos(w http.ResponseWriter, r *http.Request) {
	query := `
        SELECT 
            m.id, m.productId, m.sourceStoreId, m.targetStoreId,
            m.quantity, m.type, m.timestamp, m.activo,
            m.created_at, m.updated_at,
            p.name as product_name,
            s1.name as source_store_name,
            s2.name as target_store_name
        FROM prueba.movimientos m
        JOIN catalogos.productos p ON m.productId = p.id
        JOIN catalogos.tiendas s1 ON m.sourceStoreId = s1.id
        JOIN catalogos.tiendas s2 ON m.targetStoreId = s2.id
        WHERE m.activo = true
        ORDER BY m.timestamp DESC`

	rows, err := h.db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var movimientos []MovimientoDetalle
	for rows.Next() {
		var m MovimientoDetalle
		err := rows.Scan(
			&m.ID, &m.ProductID, &m.SourceStoreID, &m.TargetStoreID,
			&m.Quantity, &m.Type, &m.Timestamp, &m.Activo,
			&m.CreatedAt, &m.UpdatedAt,
			&m.ProductName, &m.SourceStoreName, &m.TargetStoreName,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		movimientos = append(movimientos, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movimientos)
}

// CrearMovimiento godoc
// @Summary      Crear movimiento
// @Description  Registra un nuevo movimiento de inventario
// @Tags         movimientos
// @Accept       json
// @Produce      json
// @Param        movimiento body CrearMovimiento true "Datos del movimiento"
// @Success      201  {object}  MovimientoDetalle
// @Failure      400  {object}  map[string]string
// @Router       /CrearMovimiento [post]
func (h *MovementHandler) CrearMovimiento(w http.ResponseWriter, r *http.Request) {
	var mov CrearMovimiento
	if err := json.NewDecoder(r.Body).Decode(&mov); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	// Validar tipo de movimiento
	if mov.Type != MovimientoIN && mov.Type != MovimientoOUT && mov.Type != MovimientoTRANSFER {
		http.Error(w, "Tipo de movimiento inválido", http.StatusBadRequest)
		return
	}

	// Validar cantidad positiva
	if mov.Quantity <= 0 {
		http.Error(w, "La cantidad debe ser positiva", http.StatusBadRequest)
		return
	}

	// Verificar existencia de producto y tiendas
	var exists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM catalogos.productos WHERE id = $1)", mov.ProductID).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "Producto no encontrado", http.StatusBadRequest)
		return
	}

	// Iniciar transacción
	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Insertar movimiento
	var movimiento MovimientoDetalle
	err = tx.QueryRow(`
        INSERT INTO prueba.movimientos (
            id, productId, sourceStoreId, targetStoreId,
            quantity, type, timestamp
        ) VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)
        RETURNING id, productId, sourceStoreId, targetStoreId,
                  quantity, type, timestamp, activo, created_at, updated_at
    `, uuid.New(), mov.ProductID, mov.SourceStoreID, mov.TargetStoreID,
		mov.Quantity, mov.Type).Scan(
		&movimiento.ID, &movimiento.ProductID, &movimiento.SourceStoreID,
		&movimiento.TargetStoreID, &movimiento.Quantity, &movimiento.Type,
		&movimiento.Timestamp, &movimiento.Activo, &movimiento.CreatedAt,
		&movimiento.UpdatedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movimiento)
}

// ObtenerMovimiento godoc
// @Summary      Obtener movimiento
// @Description  Obtiene los detalles de un movimiento específico
// @Tags         movimientos
// @Accept       json
// @Produce      json
// @Param        id query string true "ID del movimiento"
// @Success      200  {object}  MovimientoDetalle
// @Failure      404  {object}  map[string]string
// @Router       /ObtenerMovimiento [get]
func (h *MovementHandler) ObtenerMovimiento(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	query := `
        SELECT 
            m.id, m.productId, m.sourceStoreId, m.targetStoreId,
            m.quantity, m.type, m.timestamp, m.activo,
            m.created_at, m.updated_at,
            p.name as product_name,
            s1.name as source_store_name,
            s2.name as target_store_name
        FROM prueba.movimientos m
        JOIN catalogos.productos p ON m.productId = p.id
        JOIN catalogos.tiendas s1 ON m.sourceStoreId = s1.id
        JOIN catalogos.tiendas s2 ON m.targetStoreId = s2.id
        WHERE m.id = $1`

	var mov MovimientoDetalle
	err = h.db.QueryRow(query, id).Scan(
		&mov.ID, &mov.ProductID, &mov.SourceStoreID, &mov.TargetStoreID,
		&mov.Quantity, &mov.Type, &mov.Timestamp, &mov.Activo,
		&mov.CreatedAt, &mov.UpdatedAt,
		&mov.ProductName, &mov.SourceStoreName, &mov.TargetStoreName,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Movimiento no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mov)
}
