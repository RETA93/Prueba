package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
	Alert_Type  string    `json:"alert_type"`
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
	vars := mux.Vars(r)   // Obtener los parámetros de la ruta
	storeID := vars["id"] // El parámetro 'id' es parte de la ruta

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
        WHERE t.Id = $1 AND i.activo = true
        ORDER BY i.quantity ASC`

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

	// Llamar a la función de la BD
	var result bool
	err := h.db.QueryRow(`
        SELECT transfer_inventory($1, $2, $3, $4)
    `, transfer.ProductID, transfer.SourceStoreID,
		transfer.TargetStoreID, transfer.Quantity).Scan(&result)

	if err != nil {
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
	query := `SELECT * FROM prueba.vw_inventory_alerts ORDER BY quantity ASC`

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
			&alert.Alert_Type,
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
