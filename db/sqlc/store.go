package db

import (
	"database/sql"
)

// Store provides all functions to execute do queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}


// Here is place for implement all transaction in future





// // execTx executes a function within a database transaction
// func (store *Store) execTx(c context.Context, fn func(*Queries) error) error {
// 	tx, err := store.db.BeginTx(c, nil)
// 	if err != nil {
// 		return err
// 	}
// 	q := New(tx)
// 	err = fn(q)
// 	if err != nil{
// 		if rbErr := tx.Rollback(); rbErr != nil {
// 			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr) 
// 		}
// 		return err
// 	}
// 	return tx.Commit()
// }