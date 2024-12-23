package utils

import "database/sql"

// TxFn representa una función que se ejecutará dentro de una transacción
type TxFn func(*sql.Tx) error

// WithTransaction ejecuta una función dentro de una transacción
func WithTransaction(db *sql.DB, fn TxFn) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic después del rollback
		}
	}()

	err = fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
