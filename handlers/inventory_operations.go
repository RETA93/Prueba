package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// StockTransfer modelo para transferencia de stock
type StockTransfer struct {
	ProductID     uuid.UUID `json:"product_id" binding:"required"`
	SourceStoreID uuid.UUID `json:"source_store_id" binding:"required"`
	TargetStoreID uuid.UUID `json:"target_store_id" binding:"required"`
	Quantity      int       `json:"quantity" binding:"required,gt=0"`
}

// StockAlert modelo para alertas de stock
type StockAlert struct {
	ProductID   uuid.UUID `json:"product_id"`
	StoreID     uuid.UUID `json:"store_id"`
	ProductName string    `json:"product_name"`
	StoreName   string    `json:"store_name"`
	Quantity    int       `json:"current_quantity"`
	MinStock    int       `json:"min_stock"`
}

// GetStoreInventory godoc
// @Summary      Listar inventario por tienda
// @Description  Obtiene el inventario completo de una tienda específica
// @Tags         inventario
// @Accept       json
// @Produce      json
// @Param        id path string true "ID de la tienda"
// @Success      200  {array}   InventarioDetalle
// @Failure      404  {object}  map[string]string
// @Router       /stores/{id}/inventory [get]
func (h *InventoryHandler) GetStoreInventory(w http.ResponseWriter, r *http.Request) {
	// Obtener ID de la tienda de la URL
	storeID := r.URL.Query().Get("id")
	if storeID == "" {
		http.Error(w, "ID de tienda requerido", http.StatusBadRequest)
		return
	}

	storeUUID, err := uuid.Parse(storeID)
	if err != nil {
		http.Error(w, "ID de tienda inválido", http.StatusBadRequest)
		return
	}

	query := `
        SELECT 
            i.id, i.productId, i.storeId, i.quantity, i.minStock,
            i.activo, i.created_at, i.updated_at,
            p.name as product_name,
            t.name as store_name
        FROM prueba.inventarios i
        JOIN catalogos.productos p ON i.productId = p.id
        JOIN catalogos.tiendas t ON i.storeId = t.id
        WHERE i.storeId = $1 AND i.activo = true`

	rows, err := h.db.Query(query, storeUUID)
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

// TransferInventory godoc
// @Summary      Transferir productos entre tiendas
// @Description  Realiza una transferencia de productos entre tiendas con validación de stock
// @Tags         inventario
// @Accept       json
// @Produce      json
// @Param        transfer body StockTransfer true "Datos de la transferencia"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /inventory/transfer [post]
func (h *InventoryHandler) TransferInventory(w http.ResponseWriter, r *http.Request) {
	var transfer StockTransfer
	if err := json.NewDecoder(r.Body).Decode(&transfer); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	// Iniciar transacción
	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Verificar stock disponible
	var currentStock int
	err = tx.QueryRow(`
        SELECT quantity 
        FROM prueba.inventarios 
        WHERE productId = $1 AND storeId = $2 AND activo = true
        FOR UPDATE`,
		transfer.ProductID, transfer.SourceStoreID).Scan(&currentStock)

	if err == sql.ErrNoRows {
		http.Error(w, "No hay inventario disponible", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if currentStock < transfer.Quantity {
		http.Error(w, "Stock insuficiente", http.StatusBadRequest)
		return
	}

	// Reducir stock en tienda origen
	_, err = tx.Exec(`
        UPDATE prueba.inventarios
        SET quantity = quantity - $1, updated_at = CURRENT_TIMESTAMP
        WHERE productId = $2 AND storeId = $3`,
		transfer.Quantity, transfer.ProductID, transfer.SourceStoreID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Aumentar stock en tienda destino
	_, err = tx.Exec(`
        INSERT INTO prueba.inventarios (id, productId, storeId, quantity, minStock, activo)
        VALUES ($1, $2, $3, $4, 10, true)
        ON CONFLICT (productId, storeId) DO UPDATE
        SET quantity = prueba.inventarios.quantity + $4,
            updated_at = CURRENT_TIMESTAMP`,
		uuid.New(), transfer.ProductID, transfer.TargetStoreID, transfer.Quantity)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Registrar el movimiento
	_, err = tx.Exec(`
        INSERT INTO prueba.movimientos (
            id, productId, sourceStoreId, targetStoreId,
            quantity, type, timestamp, activo
        ) VALUES ($1, $2, $3, $4, $5, 'TRANSFER', CURRENT_TIMESTAMP, true)`,
		uuid.New(), transfer.ProductID, transfer.SourceStoreID,
		transfer.TargetStoreID, transfer.Quantity)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Transferencia completada exitosamente",
	})
}

// GetStockAlerts godoc
// @Summary      Listar alertas de stock bajo
// @Description  Obtiene una lista de productos que están por debajo del stock mínimo
// @Tags         inventario
// @Accept       json
// @Produce      json
// @Success      200  {array}   StockAlert
// @Router       /inventory/alerts [get]
func (h *InventoryHandler) GetStockAlerts(w http.ResponseWriter, r *http.Request) {
	query := `
        SELECT 
            i.productId,
            i.storeId,
            p.name as product_name,
            t.name as store_name,
            i.quantity,
            i.minStock
        FROM prueba.inventarios i
        JOIN catalogos.productos p ON i.productId = p.id
        JOIN catalogos.tiendas t ON i.storeId = t.id
        WHERE i.activo = true AND i.quantity <= i.minStock
        ORDER BY i.quantity ASC`

	rows, err := h.db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var alerts []StockAlert
	for rows.Next() {
		var alert StockAlert
		err := rows.Scan(
			&alert.ProductID,
			&alert.StoreID,
			&alert.ProductName,
			&alert.StoreName,
			&alert.Quantity,
			&alert.MinStock,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		alerts = append(alerts, alert)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}
