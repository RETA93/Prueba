package models

import (
	"database/sql"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

// ---------------------------------------------------------------------------------------------------------------------------
func (r *Repository) CreateProduct(p *Producto) error {
	query := `
        INSERT INTO catalogos.productos (id, name, description, category, price, sku)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id`

	return r.db.QueryRow(
		query,
		p.ID, p.Name, p.Description, p.Category, p.Price, p.SKU,
	).Scan(&p.ID)
}

// ---------------------------------------------------------------------------------------------------------------------------
func (r *Repository) UpdateProduct(p *Producto) error {
	query := `
        UPDATE catalogos.productos 
        SET name = $2, description = $3, category = $4, price = $5, sku = $6, updated_at = CURRENT_TIMESTAMP
        WHERE id = $1 AND activo = true`

	result, err := r.db.Exec(query, p.ID, p.Name, p.Description, p.Category, p.Price, p.SKU)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------------
func (r *Repository) DeleteProduct(id uuid.UUID) error {
	query := `
        UPDATE catalogos.productos 
        SET activo = false, updated_at = CURRENT_TIMESTAMP
        WHERE id = $1 AND activo = true`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------------
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAllProductos() ([]Producto, error) {
	query := `
        SELECT id, name, description, category, price, sku, 
               activo, created_at, updated_at
        FROM catalogos.productos 
        WHERE activo = true`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productos []Producto
	for rows.Next() {
		var p Producto
		err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.Category,
			&p.Price, &p.SKU, &p.Activo, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		productos = append(productos, p)
	}
	return productos, nil
}

func (r *Repository) GetAllTiendas() ([]Tienda, error) {
	query := `
        SELECT id, name, address, phone, activo, created_at, updated_at
        FROM catalogos.tiendas 
        WHERE activo = true`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tiendas []Tienda
	for rows.Next() {
		var t Tienda
		err := rows.Scan(
			&t.ID, &t.Name, &t.Address, &t.Phone,
			&t.Activo, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tiendas = append(tiendas, t)
	}
	return tiendas, nil
}

func (r *Repository) GetInventarioByTienda(tiendaID uuid.UUID) ([]Inventario, error) {
	query := `
        SELECT i.id, i.product_id, i.store_id, i.quantity, i.min_stock,
               i.activo, i.created_at, i.updated_at
        FROM prueba.inventarios i
        WHERE i.store_id = $1 AND i.activo = true`

	rows, err := r.db.Query(query, tiendaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventarios []Inventario
	for rows.Next() {
		var i Inventario
		err := rows.Scan(
			&i.ID, &i.ProductID, &i.StoreID, &i.Quantity,
			&i.MinStock, &i.Activo, &i.CreatedAt, &i.UpdatedAt)
		if err != nil {
			return nil, err
		}
		inventarios = append(inventarios, i)
	}
	return inventarios, nil
}
